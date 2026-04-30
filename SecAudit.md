# Security Audit – Heinen (Zahnquiz)

**Audit-Datum:** 2026-04-30
**Auditor:** Claude (Code-Review)
**Umfang:** gesamter Go-Server-Code (`*.go`), Frontend-JS in `frontend.go`, Deploy-Skript, Konfiguration

> Dieser Bericht klassifiziert Befunde nach **HOCH / MITTEL / NIEDRIG / INFO**.
> Hochriskante Befunde sollten zeitnah behoben werden, bevor die Anwendung öffentlich exponiert wird.

---

## Zusammenfassung

| Schwere | Anzahl |
|---------|--------|
| HOCH    | 3 |
| MITTEL  | 9 |
| NIEDRIG | 8 |
| INFO    | 4 |

Die Anwendung verfügt bereits über solide Grundlagen (bcrypt, CSRF-Double-Submit, parametrisierte SQL-Queries, Rate-Limiting, Security-Header). Die schwerwiegendsten Probleme liegen in einem **Path-Traversal in der Sound-Auslieferung**, in **schwachen WebSocket-Origin-Prüfungen** und in der **Preisgabe von Lobby-Passwörtern im State-Broadcast**.

---

## HOCH

### H1 – Path Traversal in `handleSoundFile` (`handlers.go:491-521`)

```go
name := strings.TrimPrefix(r.URL.Path, "/sounds/")
...
for _, p := range soundTypes {
    if strings.HasPrefix(name, p) { ok = true; break }
}
...
fp := filepath.Join("sounds", name)
...
http.ServeFile(w, r, fp)
```

Die Validierung verwendet `strings.HasPrefix`, prüft jedoch nicht auf Pfadtrennzeichen. Eine Anfrage an
`/sounds/intro_sound/../../etc/passwd` erfüllt `HasPrefix(name, "intro_sound")` und wird nach `filepath.Join("sounds", "intro_sound/../../etc/passwd")` aufgelöst – `http.ServeFile` liefert dann Dateien außerhalb des `sounds/`-Verzeichnisses aus.

**Empfehlung:** Nach `filepath.Clean(name)` prüfen, dass das Resultat keine `..`-Komponenten enthält und mit dem Sound-Verzeichnis beginnt; alternativ die Datei-Endung sowie den exakten Dateinamen aus einer Whitelist (`soundTypes[i] + ext`) prüfen, bevor `ServeFile` aufgerufen wird.

---

### H2 – Permissive WebSocket-Origin-Prüfung (`ws.go:16-34`)

```go
var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
    origin := r.Header.Get("Origin")
    if origin == "" {
        return true // Non-browser clients don't send Origin
    }
    ...
}}
```

Browser senden bei WebSocket-Handshakes **immer** einen Origin-Header. Das Akzeptieren leerer Origins erlaubt:
- Cross-Site-WebSocket-Hijacking via Tools, die den Header unterdrücken.
- Verbindungen aus Nicht-Browser-Clients ohne jegliche Origin-Bindung – kombiniert mit `SameSite=Lax` an der Session-Cookie kann ein Angreifer-Server einen authentifizierten Nutzer dazu bringen, einen WS-Tunnel zu öffnen, der dann auch admin-relevante Spielaktionen (`create_game`, `kick_player`, `transfer_host`) ausführen kann.

WebSockets unterliegen **nicht** der Same-Origin-Policy auf Cookie-Ebene wie XHR; ein striktes Origin-Whitelisting ist hier die einzige Verteidigung.

**Empfehlung:** Leere Origin verwerfen oder zusätzlich pro WebSocket-Aktion einen kurzlebigen, an die Session gebundenen Token mitsenden, der vor dem ersten privilegierten Befehl validiert wird.

---

### H3 – Lobby-Passwort wird an alle Spieler im State-Broadcast verschickt (`game.go:654-712` + `ws.go`)

`buildState()` serialisiert `g.Settings` direkt in den State, der per `broadcastState()` an alle verbundenen Sockets geht. `GameSettings.LobbyPassword` ist ein Klartextfeld und wird damit jedem Spieler in der Lobby zugänglich gemacht – auch über die Browser-Konsole / DevTools.

Konsequenz: Wer einmal beigetreten ist, kann das Passwort weitergeben oder bei Re-Use gegen andere Konten verwenden.

**Empfehlung:** `LobbyPassword` aus dem broadcasteten State entfernen (z. B. eigenes `clientSettings` ohne dieses Feld) oder server-seitig nur ein Bool `hasPassword` exponieren.

---

## MITTEL

### M1 – Kein CSRF-Schutz / kein Methoden-Check auf `/api/logout` (`handlers.go:67-82`)
Der Endpoint löscht die Session per Cookie ohne `r.Method`-Check und ohne CSRF-Token. Ein eingebettetes `<img src="/api/logout">` o. ä. zwingt Nutzer zum Logout. Geringer Schaden, aber unnötig.

**Empfehlung:** Methode auf `POST` beschränken und CSRF-Token verifizieren.

---

### M2 – API-Schlüssel im Klartext in `settings`-Tabelle (`db.go`, `handlers.go:241-340`, `ai.go`)
`OPENAI_API_KEY` / `ANTHROPIC_API_KEY` werden – wenn nicht aus ENV gelesen – in der SQLite-DB im Klartext gespeichert. Bei Backups, DB-Diebstahl oder fehlerhaften Berechtigungen liegen produktive API-Schlüssel offen.

**Empfehlung:** Bevorzugt ENV-Variablen verwenden; falls DB-Speicherung notwendig, mit einem aus dem Master-Secret abgeleiteten symmetrischen Schlüssel verschlüsseln. Mindestens DB-Datei mit `chmod 600` ablegen und Backup-Ausschluss prüfen.

---

### M3 – Geteilter Rate-Limit-State zwischen Login und Lobby-Passwort (`ratelimit.go`)
`CheckLoginRate` und `CheckLobbyPasswordRate` schreiben in dieselbe Map `ipLimits`. Ein Brute-Force-Versuch gegen ein Lobby-Passwort kann den Login-Counter zurücksetzen und umgekehrt. Außerdem inkrementiert `CheckLoginRate` den Counter selbst dann, wenn das Konto erfolgreich angemeldet wurde, weil der Pfad zuerst `CheckLoginRate` aufruft, bevor `RecordSuccess` greift.

**Empfehlung:** Getrennte Maps (`ipLoginLimits`, `ipLobbyPwLimits`) und die Off-by-One-Grenze (`> 5` statt `>= 5` bei initialem `count=1`) prüfen.

---

### M4 – Unbegrenztes Wachstum der Rate-Limit-Maps (`ratelimit.go`)
`ipLimits` / `accountLimits` werden niemals gekürzt. Ein Angreifer kann durch Anmeldeversuche von vielen IP-Adressen / unter vielen Account-Namen unbegrenzt RAM allokieren (DoS).

**Empfehlung:** Periodische Bereinigung von Einträgen, deren `lastAt` älter als 30 Min. ist, oder Verwendung einer LRU-Cache-Bibliothek mit Obergrenze.

---

### M5 – `json.NewDecoder(r.Body).Decode(...)` ohne Größenlimit
In sämtlichen JSON-Handlern (`handleLogin`, `handleUsers`, `handleChangePassword`, `handleAIConfig`, …) wird der Request-Body ohne `http.MaxBytesReader` gelesen. Ein Client kann beliebig große JSON-Objekte schicken und den Server zur Allokation zwingen. (Nur bei `handleSounds` ist das Limit korrekt gesetzt.)

**Empfehlung:** Pro Handler `r.Body = http.MaxBytesReader(w, r.Body, 64*1024)` o. ä. setzen.

---

### M6 – Tutorial-Markdown wird mittels regex-basiertem Parser per `innerHTML` injiziert (`frontend.go:358, 364, 374`)
`tutorialHtml=markdownToHtml(d.content||'')` und anschließend `home-tutorial.innerHTML=tutorialHtml`. Der Parser entfernt keine HTML-Tags. Eine `tutorial.md` mit `<script>…</script>` wird ausgeführt.

Aktuell ist `tutorial.md` lokal vom Admin gepflegt – aber bei einem System-Update / fehlerhaften Upload-Pfad oder einer späteren Editor-Funktion wird das schnell zur XSS-Quelle.

**Empfehlung:** Server-seitig `bluemonday`/`goldmark` mit Sanitizer rendern und das resultierende sichere HTML ausliefern; oder im Frontend DOMPurify einsetzen, bevor `innerHTML` gesetzt wird.

---

### M7 – CSP erlaubt `'unsafe-inline'` für Scripts (`main.go:19`)
```
script-src 'self' 'unsafe-inline' https://cdnjs.cloudflare.com;
```
`'unsafe-inline'` deaktiviert den größten Teil der CSP-XSS-Härtung. Das gesamte Frontend ist eine inline-`<script>`-Sektion – ein Refactor wäre invasiv, sollte aber langfristig erfolgen (Nonces / Hashes oder externe `.js`-Datei).

Zusätzlich: `https://cdnjs.cloudflare.com` ist erlaubt, in `frontend.go` (geprüfter Auszug) finde ich keinen aktiven Verbrauch. Wenn nicht benötigt, entfernen.

---

### M8 – Log-Injection via Benutzername / Topic / Lobby-Name (`logging.go:49-64`)
`logEvent` formatiert Felder ohne Escaping. Ein Angreifer kann z. B. bei der Anmeldung als Benutzer `admin\n[2026-01-01 00:00:00] [INFO] IP=- USER=victim ACTION=LOGIN_OK details=` einen gefälschten Logeintrag erzeugen, der den späteren Log-Viewer (Regex-basiert) verwirrt und Audit-Spuren manipuliert.

**Empfehlung:** Newlines / Steuerzeichen aus IP, User, Action, Details vor dem Schreiben strippen oder JSON-Logging verwenden.

---

### M9 – AI-Provider-Klartextantworten in Fehlermeldungen an Clients (`ai.go:117-132`, `handlers.go:1086`)
Bei Parse- oder „No-Q“-Fehlern wird die Roh-Antwort des AI-Anbieters (bis zu 800 Zeichen) zurück an den Frontend-Client geschickt. Das kann Modelldetails, Prompt-Leakage oder Provider-Fehlermeldungen enthalten. Im Free-Tier eines API-Anbieters wäre das möglicherweise unkritisch, in produktiven Setups ist eine Trennung „User-Meldung / interner Log“ aber sinnvoll.

**Empfehlung:** Generische Fehlermeldung an den Client, vollständige Antwort nur ins Server-Log.

---

## NIEDRIG

### N1 – Self-XSS durch fehlendes JS-String-Escaping bei Benutzernamen (`frontend.go:402-403`)
`esc(s)` gibt nur HTML-escapten Text zurück (kein `\'` / `\"`). In `renderUsers` wird `esc(u.username)` jedoch in ein Single-Quoted-JS-String-Argument eingebettet: `onclick="resetUserPassword('+u.id+',\''+esc(u.username)+'\')"`. Ein Username wie `bob');alert(1)//` führt zu JS-Injection beim nächsten Admin-Aufruf.
Da nur Admins Nutzer anlegen, ist das **Self-XSS / Stored-Privilege-Eskalation gering**, aber unnötig.

**Empfehlung:** `data-`-Attribute statt onclick-String-Konkatenation oder eine zusätzliche `escapeJsString`-Funktion.

---

### N2 – Anonymes Massen-Reporting möglich (`handlers.go:1095-1127`)
`/api/questions/report` erfordert nur ein CSRF-Token (kein Login, kein Rate-Limit). Ein Angreifer kann beliebige Fragen beliebig oft melden und so den Admin-Workflow verstopfen.

**Empfehlung:** Pro IP/Frage Rate-Limit, ggf. Login-Pflicht oder zumindest Captcha bei hohem Aufkommen.

---

### N3 – Nicht-konstante String-Vergleiche bei Lobby-Passwort und CSRF (`ws.go:130`, `csrf.go:45`)
`p.Password != g.Settings.LobbyPassword` und `strings.EqualFold(...)` sind nicht zeitkonstant. Bei aktiv exploitable Timing-Side-Channels (z. B. lokales LAN) theoretisch gegen den Passwort-Vergleich nutzbar; CSRF-Token mit 256 Bit ist praktisch unangreifbar.

**Empfehlung:** `crypto/subtle.ConstantTimeCompare` für Passwort und CSRF-Token verwenden.

---

### N4 – Session wird bei Privilegienänderung nicht rotiert (`handlers.go:206-237`)
Wenn ein User per `PUT /api/users` zum Admin gemacht oder degradiert wird, bleiben bestehende Session-Cookies gültig. Ein abgelaufener Browser-Tab kann mit alter Rolle weiter agieren bis zum nächsten `/api/me`-Refresh.

**Empfehlung:** `deleteUserSessions(targetUser)` auch beim Rollenwechsel aufrufen.

---

### N5 – Logout-Cookie wird nicht überall ungültig gemacht
Beim Passwort-Reset durch den Admin werden Sessions für den Zielnutzer gelöscht (`deleteUserSessions`), aber **nicht** die `heinen_csrf`-Cookies invalidiert. Auch `cleanExpiredSessions` löscht keine Reconnect-Tokens vor Ablauf bei Rollenänderung.

**Empfehlung:** Bei sensitiven Vorgängen alle relevanten Cookies entwerten.

---

### N6 – `getIP` parst IPv6 nicht korrekt (`logging.go:76-88`)
`strings.Split(r.RemoteAddr, ":")[0]` ergibt für `[::1]:1234` den String `"["`. Damit greift die "trusted proxy"-Erkennung für IPv6-Localhost nicht; Forwarded-Header werden in IPv6-Setups niemals akzeptiert. Zusätzlich werden IPv6-IPs in Logs verstümmelt protokolliert.

**Empfehlung:** `net.SplitHostPort(r.RemoteAddr)` verwenden.

---

### N7 – HSTS / weitere Security-Header fehlen (`main.go:17-26`)
`Strict-Transport-Security` ist nicht gesetzt. `Permissions-Policy` und `Cross-Origin-Opener-Policy` ebenfalls nicht. Im HTTPS-Reverse-Proxy-Setup kann der Proxy das übernehmen, sollte aber mindestens an einer Stelle erzwungen werden.

---

### N8 – `Game-ID` und `Player-ID` mit reduzierter Entropie (`game.go:124-141`, `ws.go:98,136`)
`uuid.New().String()[:8]` liefert nur 32 Bit Entropie für die Game-ID, `[:12]` für Player-IDs (48 Bit). Für Reconnect-Tokens wird zwar `crypto/rand` (256 Bit) genutzt – die Player-ID wird aber als Identitätsnachweis in `submitAnswer`, `kick_player`, `transfer_host` etc. behandelt: Wer eine 12-Hex-ID errät, kann unter dieser Identität spielen. Praktisch durch begrenzte Lobby-Größe schwierig, aber unnötig schwach.

**Empfehlung:** Vollständige UUIDs oder `crypto/rand`-Tokens (≥128 Bit) verwenden.

---

## INFO / Hinweise

### I1 – Markdown-„Parser“ ist regex-basiert (`frontend.go:364`)
Beliebige Markdown-Eigenheiten (Tabellen, Code-Blöcke mit Backticks, verschachtelte Listen) werden falsch oder unsicher gerendert. Siehe M6.

### I2 – `LobbyPassword` und Game-State werden nicht persistent gespeichert
Beim Neustart des Servers gehen alle aktiven Lobbys und Reconnect-States im RAM verloren – das ist eher Funktionalität als Sicherheit, sollte aber dokumentiert werden.

### I3 – `deploy.sh` enthält Server-Hostnamen und SSH-Pfad
Datei steht in `.gitignore`, ist also lokal. Stelle sicher, dass die Datei nie versehentlich commitet wird (Pre-Commit-Hook oder `git update-index --assume-unchanged`).

### I4 – `heinen.db` und `heinen.log` werden nicht ignoriert
Im Repo-Root liegen `heinen.db` und `heinen.log` (`git status` zeigt sie als „untracked“). Beide enthalten potenziell Anmeldedaten/Hashes bzw. IP-Adressen. `.gitignore` sollte sie explizit aufnehmen, damit niemand versehentlich `git add -A` macht.

---

## Positive Beobachtungen

- **bcrypt** wird für neue Passwörter konsistent verwendet; Legacy-SHA-256-Hashes werden bei Login on-the-fly migriert.
- **Bootstrap-Pflicht** für Admin-Account via ENV (`HEINEN_ADMIN_PASSWORD`) verhindert Standard-Credentials.
- Alle SQL-Statements verwenden **Parameter-Bindings** – keine SQL-Injection-Pfade gefunden.
- **CSRF-Double-Submit** ist auf allen state-ändernden HTTP-Endpoints (außer `/api/logout`) konsistent erzwungen.
- **Reconnect-Tokens** sind 256 Bit `crypto/rand` und werden bei Verwendung gelöscht (Single-Use).
- **`http.Server` Timeouts** (`ReadHeaderTimeout`, `ReadTimeout`, `WriteTimeout`, `IdleTimeout`) sind gesetzt – Slowloris-Schutz vorhanden.
- **Account-Lockout** nach 5 Fehlversuchen mit 15-Min-Cooldown.
- **Sound-Upload** validiert Datei-Endung (.mp3/.wav) und 20 MB Größenlimit.
- **Game-History-IDs** werden vor DB-Lookup syntaktisch validiert (`isValidHistoryID`).

---

## Priorisierte To-Do-Liste

1. **(H1)** Path-Traversal in `handleSoundFile` schließen.
2. **(H2)** WebSocket-Origin nur explizit whitelisten; leere Origin verwerfen.
3. **(H3)** `LobbyPassword` aus dem State-Broadcast entfernen.
4. **(M5)** Globaler `MaxBytesReader` für JSON-Bodies.
5. **(M3, M4)** Rate-Limit-Logik trennen + Bounded-Cleanup.
6. **(M1)** Logout: Methodenprüfung + CSRF.
7. **(M8)** Newlines aus Log-Feldern strippen.
8. **(M2)** API-Keys verschlüsselt oder ausschließlich aus ENV.
9. **(M6, M7)** Frontend-Sanitizer + CSP härten.
10. Niedrige/Info-Punkte als Backlog.
