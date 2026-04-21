package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
	origin := r.Header.Get("Origin")
	if origin == "" {
		return true // Non-browser clients don't send Origin
	}
	host := r.Host
	if fh := r.Header.Get("X-Forwarded-Host"); fh != "" {
		host = fh
	}
	if fproto := r.Header.Get("X-Forwarded-Proto"); fproto != "" {
		return origin == fproto+"://"+host
	}
	return origin == "http://"+host || origin == "https://"+host
}}

func wsError(conn *websocket.Conn, msg string) {
	d, _ := json.Marshal(map[string]interface{}{"type": "error", "payload": map[string]string{"message": msg}})
	conn.WriteMessage(websocket.TextMessage, d)
}

func handleWS(w http.ResponseWriter, r *http.Request) {
	// Verify authentication before upgrading WebSocket
	authUser := getSessionUser(r)
	if authUser == "" {
		http.Error(w, "Unauthorized", 401)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()
	wsIP := getIP(r)
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})
	pingDone := make(chan struct{})
	go func() {
		tk := time.NewTicker(25 * time.Second)
		defer tk.Stop()
		for {
			select {
			case <-tk.C:
				if conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(10*time.Second)) != nil {
					return
				}
			case <-pingDone:
				return
			}
		}
	}()
	defer close(pingDone)

	var playerID, curInvite string
	for {
		_, mb, err := conn.ReadMessage()
		if err != nil {
			break
		}
		var msg struct {
			Type    string          `json:"type"`
			Payload json.RawMessage `json:"payload"`
		}
		if json.Unmarshal(mb, &msg) != nil {
			continue
		}
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))

		switch msg.Type {
		case "ping":
			continue

		case "create_game":
			var p struct{ Name string }
			json.Unmarshal(msg.Payload, &p)
			playerID = uuid.New().String()[:12]
			g := newGame(playerID, p.Name, authUser)
			curInvite = g.InviteCode
			gamesMu.Lock()
			games[curInvite] = g
			gamesMu.Unlock()
			reconnectToken := issueReconnectToken(playerID, g.ID)
			g.connMu.Lock()
			g.connections[playerID] = conn
			g.connMu.Unlock()
			d, _ := json.Marshal(map[string]interface{}{"type": "joined", "payload": map[string]string{"playerId": playerID, "gameId": g.ID, "inviteCode": curInvite, "reconnectToken": reconnectToken}})
			conn.WriteMessage(websocket.TextMessage, d)
			g.broadcastState()
			logInfo(wsIP, authUser, "GAME_CREATE", fmt.Sprintf("invite=%s lobby=%s", curInvite, g.Settings.LobbyName))

		case "join_game":
			var p struct{ Name, InviteCode, Password string }
			json.Unmarshal(msg.Payload, &p)
			code := strings.ToLower(p.InviteCode)
			gamesMu.Lock()
			g, ok := games[code]
			gamesMu.Unlock()
			if !ok {
				wsError(conn, "Spiel nicht gefunden")
				continue
			}
			if g.Settings.LobbyMode == LobbyPassword && p.Password != g.Settings.LobbyPassword {
				logWarn(wsIP, p.Name, "LOBBY_AUTH_FAIL", "invite="+code)
				wsError(conn, "Falsches Passwort")
				continue
			}
			playerID = uuid.New().String()[:12]
			curInvite = code
			if !g.addPlayer(playerID, p.Name) {
				wsError(conn, "Beitritt nicht möglich")
				continue
			}
			reconnectToken := issueReconnectToken(playerID, g.ID)
			g.connMu.Lock()
			g.connections[playerID] = conn
			g.connMu.Unlock()
			d, _ := json.Marshal(map[string]interface{}{"type": "joined", "payload": map[string]string{"playerId": playerID, "gameId": g.ID, "inviteCode": curInvite, "reconnectToken": reconnectToken}})
			conn.WriteMessage(websocket.TextMessage, d)
			g.broadcastState()

		case "reconnect":
			var p struct{ ReconnectToken string }
			json.Unmarshal(msg.Payload, &p)
			playerID, gameID, ok := validateReconnectToken(p.ReconnectToken)
			if !ok {
				wsError(conn, "Ungültiges Reconnect-Token")
				continue
			}
			gamesMu.Lock()
			g, ok := games[curInvite]
			if !ok {
				for code, gg := range games {
					if gg.ID == gameID {
						g = gg
						curInvite = code
						ok = true
						break
					}
				}
			}
			gamesMu.Unlock()
			if !ok {
				wsError(conn, "Spiel nicht gefunden")
				continue
			}
			if ok := g.reconnectPlayerByID(playerID, conn); ok {
				d, _ := json.Marshal(map[string]interface{}{"type": "reconnected", "payload": map[string]string{"playerId": playerID, "gameId": g.ID, "inviteCode": curInvite}})
				conn.WriteMessage(websocket.TextMessage, d)
				g.broadcastState()
			} else {
				wsError(conn, "Spieler nicht gefunden")
			}

		case "update_settings":
			gamesMu.Lock()
			g, ok := games[curInvite]
			gamesMu.Unlock()
			if !ok || playerID != g.HostID {
				continue
			}
			var s struct {
				Topic, Difficulty, StartDifficulty, Mode, LobbyName, LobbyMode, LobbyPassword string
				NumQuestions, TimePerQ, NumOptions, NumTeeth                                  int
				ShowTutorial, WebSearch, PlayIntro                                            *bool
			}
			json.Unmarshal(msg.Payload, &s)
			g.mu.Lock()
			if s.Topic != "" {
				g.Settings.Topic = s.Topic
			}
			if s.Difficulty != "" {
				g.Settings.Difficulty = s.Difficulty
			}
			if s.StartDifficulty != "" {
				g.Settings.StartDifficulty = s.StartDifficulty
			}
			if s.NumQuestions > 0 {
				g.Settings.NumQuestions = s.NumQuestions
			}
			if s.TimePerQ > 0 {
				g.Settings.TimePerQ = s.TimePerQ
			}
			if s.NumOptions >= 2 {
				g.Settings.NumOptions = s.NumOptions
			}
			if s.NumTeeth > 0 {
				g.Settings.NumTeeth = s.NumTeeth
				if g.Phase == PhaseLobby {
					for _, p := range g.Players {
						p.Teeth = s.NumTeeth
						p.MaxTeeth = s.NumTeeth
					}
				}
			}
			if s.Mode != "" {
				m := GameMode(s.Mode)
				if m == ModeSingleplayer && len(g.Players) > 1 {
					m = ModeClassic
				}
				if m == ModeClassic || m == ModeElimination || m == ModeBattleRoyale || m == ModeSingleplayer {
					g.Settings.Mode = m
				}
			}
			if s.ShowTutorial != nil {
				g.Settings.ShowTutorial = *s.ShowTutorial
			}
			if s.LobbyName != "" {
				g.Settings.LobbyName = s.LobbyName
			}
			if s.LobbyMode != "" {
				lm := LobbyMode(s.LobbyMode)
				if lm == LobbyInvite || lm == LobbyPassword || lm == LobbyOpen {
					g.Settings.LobbyMode = lm
				}
			}
			if s.LobbyPassword != "" {
				g.Settings.LobbyPassword = s.LobbyPassword
			}
			if s.WebSearch != nil {
				g.Settings.WebSearch = *s.WebSearch
			}
			if s.PlayIntro != nil {
				g.Settings.PlayIntro = *s.PlayIntro
			}
			g.mu.Unlock()
			g.broadcastState()

		case "kick_player":
			gamesMu.Lock()
			g, ok := games[curInvite]
			gamesMu.Unlock()
			if !ok || playerID != g.HostID {
				continue
			}
			var p struct {
				PlayerID string `json:"playerId"`
			}
			json.Unmarshal(msg.Payload, &p)
			if p.PlayerID == g.HostID {
				continue
			}
			g.sendTo(p.PlayerID, "kicked", map[string]string{"message": "Du wurdest entfernt"})
			g.connMu.Lock()
			if c, ok := g.connections[p.PlayerID]; ok {
				c.Close()
				delete(g.connections, p.PlayerID)
			}
			g.connMu.Unlock()
			g.removePlayer(p.PlayerID)
			g.broadcastState()
			go cleanupIfEmpty(curInvite)

		case "transfer_host":
			gamesMu.Lock()
			g, ok := games[curInvite]
			gamesMu.Unlock()
			if !ok || playerID != g.HostID {
				continue
			}
			var p struct {
				PlayerID string `json:"playerId"`
			}
			json.Unmarshal(msg.Payload, &p)
			g.mu.Lock()
			if _, exists := g.Players[p.PlayerID]; exists && p.PlayerID != g.HostID {
				g.HostID = p.PlayerID
				logInfo(wsIP, g.HostUser, "HOST_TRANSFER", fmt.Sprintf("invite=%s new_host=%s", curInvite, p.PlayerID))
			}
			g.mu.Unlock()
			g.broadcastState()

		case "start_game":
			gamesMu.Lock()
			g, ok := games[curInvite]
			gamesMu.Unlock()
			if !ok || playerID != g.HostID {
				continue
			}
			go func() {
				if err := g.startGame(); err != nil {
					// setError already handles AI errors — only fall back to lobby for non-AI errors
					g.mu.Lock()
					if g.Phase != PhaseError {
						g.Phase = PhaseLobby
						g.mu.Unlock()
						g.broadcastState()
					} else {
						g.mu.Unlock()
					}
				}
			}()

		case "skip_tutorial":
			gamesMu.Lock()
			g, ok := games[curInvite]
			gamesMu.Unlock()
			if !ok || playerID != g.HostID {
				continue
			}
			g.skipTutorial()

		case "end_game":
			gamesMu.Lock()
			g, ok := games[curInvite]
			gamesMu.Unlock()
			if !ok || playerID != g.HostID {
				continue
			}
			g.forceEnd()

		case "answer":
			gamesMu.Lock()
			g, ok := games[curInvite]
			gamesMu.Unlock()
			if !ok {
				continue
			}
			var p struct {
				Answer int `json:"answer"`
			}
			json.Unmarshal(msg.Payload, &p)
			g.submitAnswer(playerID, p.Answer)

		case "play_again":
			gamesMu.Lock()
			g, ok := games[curInvite]
			gamesMu.Unlock()
			if !ok || playerID != g.HostID {
				continue
			}
			g.mu.Lock()
			g.Phase = PhaseLobby
			g.CurrentQ = 0
			g.Questions = nil
			g.prefetched = nil
			g.prefetchFailed = false
			g.prefetchErr = ""
			g.ErrorMsg = ""
			for _, p := range g.Players {
				p.Teeth = g.Settings.NumTeeth
				p.MaxTeeth = g.Settings.NumTeeth
				p.Alive = true
				p.Answer = -1
				p.Answered = false
				p.JustLost = false
				p.JustDied = false
				p.LostThisGame = false
				p.Eliminated = false
			}
			g.mu.Unlock()
			g.broadcastState()

		case "leave_game":
			// Any player can leave from error phase — return to lobby
			gamesMu.Lock()
			g, ok := games[curInvite]
			gamesMu.Unlock()
			if !ok {
				continue
			}
			g.mu.Lock()
			if g.Phase == PhaseError || g.Phase == PhaseEnd {
				g.Phase = PhaseLobby
				g.CurrentQ = 0
				g.Questions = nil
				g.prefetched = nil
				g.prefetchFailed = false
				g.prefetchErr = ""
				g.ErrorMsg = ""
				for _, p := range g.Players {
					p.Teeth = g.Settings.NumTeeth
					p.MaxTeeth = g.Settings.NumTeeth
					p.Alive = true
					p.Answer = -1
					p.Answered = false
					p.JustLost = false
					p.JustDied = false
					p.LostThisGame = false
					p.Eliminated = false
				}
			}
			g.mu.Unlock()
			g.broadcastState()
		}
	}

	// Disconnect cleanup
	if curInvite != "" {
		gamesMu.Lock()
		g, ok := games[curInvite]
		gamesMu.Unlock()
		if ok {
			g.connMu.Lock()
			delete(g.connections, playerID)
			g.connMu.Unlock()
			g.mu.Lock()
			if p, ok := g.Players[playerID]; ok {
				p.Connected = false
			}
			if g.Phase == PhaseLobby {
				delete(g.Players, playerID)
				no := make([]string, 0)
				for _, pid := range g.PlayerOrder {
					if pid != playerID {
						no = append(no, pid)
					}
				}
				g.PlayerOrder = no
				if playerID == g.HostID && len(g.PlayerOrder) > 0 {
					g.HostID = g.PlayerOrder[0]
					logInfo("system", "system", "HOST_TRANSFER", fmt.Sprintf("invite=%s new_host=%s (auto)", curInvite, g.HostID))
				}
			}
			g.mu.Unlock()
			g.broadcastState()
			go cleanupIfEmpty(curInvite)
		}
	}
}
