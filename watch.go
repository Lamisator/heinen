package main

import (
	"crypto/subtle"
	"encoding/json"

	"github.com/gorilla/websocket"
)

// ─── Zuschauer*innen ohne Platz ─────────────────────────────────────────────────
//
// Wer eine Lobby nur ansehen will (Link /w/CODE, „Nur zuschauen"), bekommt eine
// WebSocket-Verbindung ohne Player: keine Zähne, keine Aktionen, nur den
// öffentlichen State (buildStateFor("")). Das geht auch während eines laufenden
// Spiels. Zuschauer*innen lösen nie eine KI-Generierung aus.

const (
	maxWatchers      = 50 // je Lobby
	maxWatchersPerIP = 5  // je Lobby und IP, damit eine Person die 50 nicht allein belegt
)

// passwordOK vergleicht das Lobby-Passwort in konstanter Zeit.
func passwordOK(given, want string) bool {
	return subtle.ConstantTimeCompare([]byte(given), []byte(want)) == 1
}

// addWatcher meldet conn als Zuschauer*in an. Bei Erfolg gehen hello und der
// öffentliche State noch unter den Locks hinaus, damit kein Broadcast
// dazwischen auf dieselbe Verbindung schreibt. Rückgabe: "" oder Fehlercode.
func (g *Game) addWatcher(conn *websocket.Conn, ip string, hello []byte) string {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.Settings.NoWatch {
		return "no_watch"
	}
	g.connMu.Lock()
	defer g.connMu.Unlock()
	if g.watchers == nil {
		g.watchers = make(map[*websocket.Conn]string)
	}
	if len(g.watchers) >= maxWatchers {
		return "watch_full"
	}
	n := 0
	for _, wip := range g.watchers {
		if wip == ip {
			n++
		}
	}
	if n >= maxWatchersPerIP {
		return "watch_full"
	}
	g.watchers[conn] = ip
	g.watchCount.Store(int32(len(g.watchers)))
	st, _ := json.Marshal(map[string]interface{}{"type": "state", "payload": g.buildStateFor("")})
	conn.WriteMessage(websocket.TextMessage, hello)
	conn.WriteMessage(websocket.TextMessage, st)
	return ""
}

// removeWatcher meldet conn ab; true, wenn sie zugeschaut hat.
func (g *Game) removeWatcher(conn *websocket.Conn) bool {
	g.connMu.Lock()
	defer g.connMu.Unlock()
	return g.removeWatcherLocked(conn)
}

func (g *Game) removeWatcherLocked(conn *websocket.Conn) bool {
	if _, ok := g.watchers[conn]; !ok {
		return false
	}
	delete(g.watchers, conn)
	g.watchCount.Store(int32(len(g.watchers)))
	return true
}

// closeWatchers schickt allen Zuschauer*innen „unwatched" und schließt ihre
// Verbindungen (Zuschauen abgeschaltet, Lobby geschlossen oder umgezogen).
func (g *Game) closeWatchers(msg string) int {
	d, _ := json.Marshal(map[string]interface{}{"type": "unwatched", "payload": map[string]string{"message": msg}})
	g.connMu.Lock()
	defer g.connMu.Unlock()
	n := len(g.watchers)
	for c := range g.watchers {
		c.WriteMessage(websocket.TextMessage, d)
		c.Close()
	}
	g.watchers = make(map[*websocket.Conn]string)
	g.watchCount.Store(0)
	return n
}

// broadcastWatchers schickt nur die neue Zahl – ein voller State würde bei der
// Host*in gerade getippte Einstellungen überschreiben.
func (g *Game) broadcastWatchers() {
	g.broadcastMsg("watchers", map[string]int{"count": int(g.watchCount.Load())})
}

func wsErrorCode(conn *websocket.Conn, code, msg string) {
	d, _ := json.Marshal(map[string]interface{}{"type": "error", "payload": map[string]string{"message": msg, "code": code}})
	conn.WriteMessage(websocket.TextMessage, d)
}

var watchErrText = map[string]string{
	"no_watch":   "In dieser Lobby ist Zuschauen abgeschaltet.",
	"watch_full": "Hier schauen schon zu viele zu.",
}
