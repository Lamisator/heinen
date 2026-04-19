# 🦷 HEINEN – Das Zahnquiz

Ein kompetitives, rundenbasiertes Multiplayer-Quizspiel mit Eliminationsmechanik.

## Schnellstart

```bash
# Kompilieren (CGO nötig für SQLite)
CGO_ENABLED=1 go build -o heinen .

# Starten
./heinen
```

Browser öffnen: **http://localhost:8671**

**Standard-Login:** `admin` / `admin`

## Abhängigkeiten

- Go 1.21+
- CGO (für `go-sqlite3`) – auf Linux ggf. `apt install gcc` nötig
- `github.com/gorilla/websocket`
- `github.com/google/uuid`
- `github.com/mattn/go-sqlite3`

## Authentifizierung

- **Login erforderlich** zum Erstellen von Spielen
- Jede*r eingeloggte User (Admin oder nicht) kann Spiele hosten
- Standard-Account `admin`/`admin` – bitte sofort ändern
- Alle User können ihr **Passwort selbst ändern** (Dashboard → "Passwort ändern")
- Beitreten per **Einladungslink** erfordert keinen Login

## Admin Control Panel

Nur für Admin-User sichtbar:

- **OpenAI API-Key** hinterlegen, ändern und testen
- **Benutzer\*innen verwalten**: anlegen, löschen, zu Admins befördern/degradieren
- **Schutzregel**: Es muss immer mindestens ein\*e Admin existieren – der\*die letzte Admin kann weder gelöscht noch degradiert werden

## Spielmodi

### Klassisch (Standard)
- Feste Anzahl an Fragen (konfigurierbar)
- Nach allen Fragen gewinnt, wer die meisten Zähne hat
- Bei **Gleichstand** gibt es mehrere Gewinner\*innen
- Frühes Ende wenn nur 1 Spieler\*in übrig

### Elimination
- Es werden so lange Fragen gestellt, bis nur noch **ein\*e Spieler\*in** Zähne hat
- Fragen werden in Batches von 20 generiert
- Wenn alle 20 aufgebraucht sind, wird kurz pausiert ("Neue Fragen werden nachgeladen...") und 20 weitere generiert

## Spielablauf

1. **Einloggen** → Dashboard → "Neues Spiel erstellen"
2. **Konfigurieren**: Thema, Modus, Fragen, Zähne, Zeit, Optionen, Intro-Dauer, Startsound
3. **Einladungslink** an Mitspieler\*innen senden – kein Login nötig
4. **Spiel starten** → HEINEN-Intro mit Slogan und optionalem Sound
5. **Fragen beantworten** – falsche/fehlende Antwort = Zahnverlust
6. **Spielende** → Ergebnis mit Gewinner\*innen
7. **"Nochmal spielen"** → Zurück zur Lobby mit allen Spieler\*innen, Host kann Einstellungen anpassen

## Feedback-Mechanik

- **"Du wurdest geheint!"** – erscheint mittig auf den Zähnen, bleibt bis zur nächsten Frage sichtbar
- **"Es hat sich ausgeheint!"** – erscheint bei Elimination mittig auf den Zähnen, bleibt bis Spielende sichtbar

## Einstellungen

| Parameter | Standard | Bereich | Hinweis |
|---|---|---|---|
| Thema | Allgemeinwissen | Freitext | Basis für KI-Fragengenerierung |
| Spielmodus | Klassisch | Klassisch / Elimination | |
| Anzahl Fragen | 10 | 1–50 | Nur im klassischen Modus |
| Zeit pro Frage | 20 Sek. | 5–120 | |
| Antwortmöglichkeiten | 4 | 2–4 | |
| Anzahl Zähne | 5 | 1–20 | |
| Intro-Dauer | 4 Sek. | 1–30 | Wie lange der HEINEN-Splash bleibt |
| Startsound | – | .mp3 / .wav | Wird beim Intro abgespielt |

## Fragengenerierung

- **Mit OpenAI API-Key:** Fragen werden live per `gpt-4o-mini` generiert
- **Ohne Key:** Platzhalter-Fragen (zum Testen)
- Key wird im Admin Panel konfiguriert und in SQLite gespeichert

## Architektur

```
main.go       – Server, Game-Logik, Auth, API-Handler, WebSocket
frontend.go   – Eingebettetes HTML/CSS/JS Frontend
go.mod/go.sum – Dependencies
heinen.db     – SQLite-Datenbank (wird automatisch erstellt)
```

- **Backend:** Go mit WebSocket (gorilla/websocket)
- **Datenbank:** SQLite (Benutzer\*innen, API-Key)
- **Frontend:** Eingebettetes SPA
- **Port:** 8671
