package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

// runningGame baut eine Lobby mitten in einer Frage: host hat richtig (2),
// p2 falsch (1) geantwortet, zwei Fragen stehen noch aus.
func runningGame() *Game {
	g := newGame("host", "Hosti", "tester")
	g.Players["p2"] = &Player{ID: "p2", Name: "Zweite", Teeth: 5, MaxTeeth: 5, Alive: true, Answer: -1, Connected: true}
	g.PlayerOrder = append(g.PlayerOrder, "p2")
	g.Settings.LobbyMode = LobbyPassword
	g.Settings.LobbyPassword = "geheimes-pw"
	g.Questions = []Question{
		{Text: "Aktuelle Frage?", Options: []string{"a", "b", "c", "d"}, Correct: 2},
		{Text: "KOMMENDE-FRAGE-EINS", Options: []string{"x", "y"}, Correct: 0},
		{Text: "KOMMENDE-FRAGE-ZWEI", Options: []string{"x", "y"}, Correct: 1},
	}
	g.Phase = PhaseQuestion
	g.TimerEnd = time.Now().Add(20 * time.Second)
	g.Players["host"].Answer, g.Players["host"].Answered = 2, true
	g.Players["p2"].Answer, g.Players["p2"].Answered = 1, true
	return g
}

func answersOf(t *testing.T, st map[string]interface{}) map[string]float64 {
	t.Helper()
	b, _ := json.Marshal(st)
	var v struct {
		Players []struct {
			ID     string  `json:"id"`
			Answer float64 `json:"answer"`
		} `json:"players"`
	}
	json.Unmarshal(b, &v)
	m := map[string]float64{}
	for _, p := range v.Players {
		m[p.ID] = p.Answer
	}
	return m
}

func TestSpectatorViewHidesPrivateData(t *testing.T) {
	g := runningGame()
	g.mu.Lock()
	pub := g.buildStateFor("")
	own := g.buildStateFor("p2")
	g.mu.Unlock()
	b, _ := json.Marshal(pub)
	js := string(b)
	for _, leak := range []string{"correct\"", "correctAnswer", "KOMMENDE-FRAGE", "geheimes-pw", "results"} {
		if strings.Contains(js, leak) {
			t.Errorf("Zuschauer-State enthält %q: %s", leak, js)
		}
	}
	for id, a := range answersOf(t, pub) {
		if a != -1 {
			t.Errorf("Zuschauer*in sieht Antwort %v von %s vor der Auflösung", a, id)
		}
	}
	// Spieler*innen sehen nur die eigene Antwort – also nicht mehr als Zuschauer*innen über andere
	oa := answersOf(t, own)
	if oa["p2"] != 1 || oa["host"] != -1 {
		t.Errorf("Spieler*innen-Sicht falsch: %v", oa)
	}
	// Nach der Auflösung ist die richtige Antwort öffentlich
	g.evaluateRoundForTest()
	g.mu.Lock()
	after, _ := json.Marshal(g.buildStateFor(""))
	g.mu.Unlock()
	if !strings.Contains(string(after), "correctAnswer") || strings.Contains(string(after), "KOMMENDE-FRAGE") {
		t.Errorf("State nach Auflösung unerwartet: %s", after)
	}
}

// evaluateRoundForTest wertet ohne Timer-Goroutine aus.
func (g *Game) evaluateRoundForTest() {
	g.mu.Lock()
	g.Phase = PhaseResults
	g.CurrentQ = 1
	g.mu.Unlock()
}

type wsClient struct {
	t *testing.T
	c *websocket.Conn
}

func dialWS(t *testing.T, srv *httptest.Server) *wsClient {
	t.Helper()
	u := "ws" + strings.TrimPrefix(srv.URL, "http") + "/ws"
	h := http.Header{"Origin": {srv.URL}}
	c, _, err := websocket.DefaultDialer.Dial(u, h)
	if err != nil {
		t.Fatal(err)
	}
	return &wsClient{t, c}
}

func (w *wsClient) send(typ string, payload interface{}) {
	d, _ := json.Marshal(map[string]interface{}{"type": typ, "payload": payload})
	if err := w.c.WriteMessage(websocket.TextMessage, d); err != nil {
		w.t.Fatal(err)
	}
}

// next liest bis zur nächsten Nachricht vom Typ typ (andere werden übersprungen).
func (w *wsClient) next(typ string) map[string]interface{} {
	w.t.Helper()
	w.c.SetReadDeadline(time.Now().Add(3 * time.Second))
	for {
		_, b, err := w.c.ReadMessage()
		if err != nil {
			w.t.Fatalf("auf %q gewartet: %v", typ, err)
		}
		var m struct {
			Type    string                 `json:"type"`
			Payload map[string]interface{} `json:"payload"`
		}
		json.Unmarshal(b, &m)
		if m.Type == typ {
			return m.Payload
		}
	}
}

func TestWatcherCannotActAndIsClosed(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(handleWS))
	defer srv.Close()
	g := runningGame()
	g.Settings.LobbyMode = LobbyInvite
	g.Phase = PhaseEnd
	g.CurrentQ = 1
	gamesMu.Lock()
	games[g.InviteCode] = g
	gamesMu.Unlock()
	defer func() { gamesMu.Lock(); delete(games, g.InviteCode); gamesMu.Unlock() }()

	w := dialWS(t, srv)
	defer w.c.Close()
	w.send("watch_game", map[string]string{"InviteCode": g.InviteCode})
	w.next("watching")
	st := w.next("state")
	if st["watchers"].(float64) != 1 || st["phase"] != "end" {
		t.Fatalf("State für Zuschauer*in: %v", st)
	}
	for _, a := range []string{"leave_game", "play_again", "start_game", "end_game", "skip_tutorial", "set_spectator", "kick_player", "transfer_host", "delegate_settings", "move_game", "reconnect", "create_game"} {
		w.send(a, map[string]interface{}{"playerId": "p2", "spectator": true, "Name": "x", "game": "y"})
	}
	w.send("update_settings", map[string]interface{}{"Topic": "Gekapert", "NoWatch": true})
	w.send("answer", map[string]int{"answer": 1})
	// Platz nehmen geht nur in der Lobby
	w.send("join_game", map[string]string{"Name": "Gast", "InviteCode": g.InviteCode})
	if e := w.next("error"); e["code"] != "running" {
		t.Fatalf("join_game im laufenden Spiel: %v", e)
	}
	g.mu.Lock()
	if g.Phase != PhaseEnd || g.Settings.Topic == "Gekapert" || g.Settings.NoWatch || g.HostID != "host" || len(g.Players) != 2 || g.Players["p2"].Spectator {
		t.Errorf("Zuschauer*in konnte etwas ändern: phase=%s topic=%s host=%s players=%d", g.Phase, g.Settings.Topic, g.HostID, len(g.Players))
	}
	g.mu.Unlock()

	// höchstens maxWatchersPerIP je IP
	var extra []*wsClient
	for i := 1; i < maxWatchersPerIP; i++ {
		x := dialWS(t, srv)
		defer x.c.Close()
		x.send("watch_game", map[string]string{"InviteCode": g.InviteCode})
		x.next("watching")
		extra = append(extra, x)
	}
	x := dialWS(t, srv)
	defer x.c.Close()
	x.send("watch_game", map[string]string{"InviteCode": g.InviteCode})
	if e := x.next("error"); e["code"] != "watch_full" {
		t.Fatalf("erwartet watch_full: %v", e)
	}

	// Host schaltet Zuschauen ab: alle Streams werden geschlossen
	g.mu.Lock()
	g.Settings.NoWatch = true
	g.mu.Unlock()
	g.closeWatchers("aus")
	if m := w.next("unwatched"); m["message"] != "aus" {
		t.Fatalf("unwatched: %v", m)
	}
	for _, e := range extra {
		e.next("unwatched")
	}
	y := dialWS(t, srv)
	defer y.c.Close()
	y.send("watch_game", map[string]string{"InviteCode": g.InviteCode})
	if e := y.next("error"); e["code"] != "no_watch" {
		t.Fatalf("erwartet no_watch: %v", e)
	}
	if n := g.watchCount.Load(); n != 0 {
		t.Errorf("watchCount = %d", n)
	}
}

func TestWatchNeedsLobbyPassword(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(handleWS))
	defer srv.Close()
	g := runningGame()
	gamesMu.Lock()
	games[g.InviteCode] = g
	gamesMu.Unlock()
	defer func() { gamesMu.Lock(); delete(games, g.InviteCode); gamesMu.Unlock() }()
	w := dialWS(t, srv)
	defer w.c.Close()
	w.send("watch_game", map[string]string{"InviteCode": g.InviteCode, "Password": ""})
	if e := w.next("error"); e["code"] != "password" {
		t.Fatalf("ohne Passwort: %v", e)
	}
	w.send("watch_game", map[string]string{"InviteCode": g.InviteCode, "Password": "geheimes-pw"})
	w.next("watching")
	st := w.next("state")
	b, _ := json.Marshal(st)
	if strings.Contains(string(b), "geheimes-pw") || strings.Contains(string(b), "KOMMENDE") || strings.Contains(string(b), "correctAnswer") {
		t.Fatalf("Zuschauer-State leakt: %s", b)
	}
	for id, a := range answersOf(t, st) {
		if a != -1 {
			t.Errorf("Antwort von %s sichtbar: %v", id, a)
		}
	}
}
