# HEINEN – Das Zahnquiz

Kompetitives, rundenbasiertes Multiplayer-Quizspiel mit KI-generierten Fragen und Eliminationsmechanik. Keine Installation für Mitspieler nötig – Beitreten per Link.

## Schnellstart

```bash
# CGO ist für SQLite erforderlich
CGO_ENABLED=1 go build -o heinen .
./heinen
```

Browser öffnen: **http://localhost:8671**

Standard-Login: `admin` / `admin` — bitte sofort im Dashboard ändern.

## Abhängigkeiten

- Go 1.21+
- GCC / CGO (für `go-sqlite3`) — auf Linux: `apt install gcc`
- `github.com/gorilla/websocket`
- `github.com/google/uuid`
- `github.com/mattn/go-sqlite3`

## Spielprinzip

Alle Spieler*innen beantworten gleichzeitig dieselbe Frage. Falsche oder fehlende Antwort = ein Zahn weniger. Wer alle Zähne verliert, scheidet aus.

## Spielmodi

| Modus | Beschreibung |
|---|---|
| **Klassisch** | Feste Fragenzahl; am Ende gewinnt, wer die meisten Zähne hat. Bei Gleichstand mehrere Gewinner*innen. |
| **Elimination** | Fragen bis nur noch eine Person Zähne hat. Fragen werden in 20er-Batches nachgeladen. |
| **KFO Battle Royale** | Wie Elimination, aber der Schwierigkeitsgrad eskaliert nach jedem 20er-Batch automatisch. |
| **Singleplayer** | Alleine spielen; Ziel ist ein möglichst hoher Score (Anzahl beantworteter Fragen). |

## Lobby-Zugangsmodi

- **Einladungslink** — Beitreten ohne Login, Link teilen reicht
- **Passwort** — Lobby mit eigenem Passwort schützen
- **Offen** — Jeder kann über die öffentliche Lobby-Liste beitreten

## Spielablauf

1. Einloggen → Dashboard → „Neues Spiel erstellen"
2. Thema, Modus, Fragenzahl, Zähne, Zeit, Schwierigkeit, Optionen konfigurieren
3. Einladungslink an Mitspieler*innen senden
4. Spiel starten → optionales HEINEN-Intro mit Splash und Sound
5. Fragen beantworten — Timer läuft für alle gleichzeitig
6. Nach jeder Frage: Auswertung mit richtige Antwort + wer Zähne verloren hat
7. Spielende → Ergebnisanzeige; „Nochmal spielen" kehrt zur Lobby zurück

## KI-Fragengenerierung

Fragen werden live per KI-API generiert. Unterstützte Anbieter:

| Anbieter | Standard-Modell | Web-Search |
|---|---|---|
| **OpenAI** | `gpt-5.4-mini` | optional (`gpt-4o-mini-search-preview` o. ä.) |
| **Anthropic** | konfigurierbar | – |

API-Key und Modell werden im **Admin Panel** hinterlegt und in SQLite gespeichert. Ohne konfigurierten Key startet das Spiel nicht.

Web Search kann pro Spiel aktiviert werden — nützlich für aktuelle Themen und Ereignisse.

## Schwierigkeitsgrade

| Stufe | Beschreibung |
|---|---|
| `leicht` | Allgemeinwissen, das die meisten kennen |
| `mittel` | Man muss nachdenken |
| `schwer` | Nur mit gutem Wissen |
| `extrem` | Echtes Expertenwissen |

Im Battle-Royale- und Singleplayer-Modus startet das Spiel auf der konfigurierten Anfangsschwierigkeit und steigert sich automatisch alle 20 Fragen.

## Einstellungen (pro Spiel)

| Parameter | Standard | Bereich |
|---|---|---|
| Thema | Allgemeinwissen | Freitext |
| Spielmodus | Klassisch | Klassisch / Elimination / Battle Royale / Singleplayer |
| Anfangsschwierigkeit | leicht | leicht / mittel / schwer / extrem |
| Anzahl Fragen | 10 | 1–50 (nur Klassisch) |
| Zeit pro Frage | 20 Sek. | 5–120 |
| Antwortmöglichkeiten | 4 | 2–4 |
| Anzahl Zähne | 5 | 1–20 |
| Web Search | aus | an / aus |
| Tutorial anzeigen | an | an / aus |
| Intro abspielen | an | an / aus |

## Sound-System

Folgende Sounds sind konfigurierbar (`.mp3` oder `.wav`, Upload im Admin Panel):

| Slot | Zeitpunkt |
|---|---|
| `intro_sound` | Beim HEINEN-Splash-Intro |
| `background_sound` | Durchgehende Hintergrundmusik |
| `question_sound` | Beim Einblenden einer neuen Frage |
| `answer_sound` | Wenn jemand antwortet |
| `wrong_sound` | Bei falscher Antwort |
| `timeout_sound` | Wenn die Zeit abläuft |
| `hurry_sound` | In den letzten Sekunden |
| `allwrong_sound` | Wenn alle falsch lagen |
| `allcorrect_sound` | Wenn alle richtig lagen |

Lautstärke je Slot ist im Admin Panel individuell regelbar.

## Admin Panel

Nur für Admin-User sichtbar:

- **KI-Anbieter** konfigurieren (OpenAI / Anthropic), Modell wählen, Key testen
- **Sounds** hochladen und Lautstärken einstellen
- **Benutzer*innen** anlegen, löschen, zu Admins befördern/degradieren
- **Intro-Dauer** global festlegen (1–30 Sek.)
- **Logs** einsehen und exportieren
- Es muss immer mindestens ein*e Admin existieren — der*die letzte Admin kann weder gelöscht noch degradiert werden

## Authentifizierung

- Login erforderlich zum Erstellen von Spielen
- Alle eingeloggten User können Spiele hosten
- Passwort-Hashing: SHA-256 mit 16-Byte-Salt; Legacy-Plaintext-Passwörter werden bei erstem Login automatisch migriert
- Beitreten per Einladungslink erfordert keinen Account

## Architektur

```
main.go       — HTTP-Server, Routing, Session-Cleanup
game.go       — Spiellogik, Phasen, Timer, Broadcast
handlers.go   — API-Handler (REST + WebSocket)
ws.go         — WebSocket-Verbindungsverwaltung
ai.go         — KI-Fragengenerierung (OpenAI / Anthropic)
auth.go       — Passwort-Hashing und Authentifizierung
db.go         — SQLite-Initialisierung, Settings, Sound-Paths
session.go    — Session-Verwaltung
logging.go    — Strukturiertes Logging
frontend.go   — Eingebettetes SPA (HTML/CSS/JS)
heinen.db     — SQLite-Datenbank (wird automatisch erstellt)
```

- **Backend:** Go, WebSocket (`gorilla/websocket`)
- **Datenbank:** SQLite mit WAL-Modus
- **Frontend:** Eingebettetes SPA, keine Build-Pipeline nötig
- **Port:** 8671
