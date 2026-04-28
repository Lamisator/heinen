package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// ─── Types ──────────────────────────────────────────────────────────────────────

type GamePhase string

const (
	PhaseTutorial GamePhase = "tutorial"
	PhaseIntro    GamePhase = "intro"
	PhaseLobby    GamePhase = "lobby"
	PhaseLoading  GamePhase = "loading"
	PhaseQuestion GamePhase = "question"
	PhaseResults  GamePhase = "results"
	PhaseRefill   GamePhase = "refill"
	PhaseEnd      GamePhase = "end"
	PhaseError    GamePhase = "error"
)

type GameMode string

const (
	ModeClassic      GameMode = "classic"
	ModeElimination  GameMode = "elimination"
	ModeBattleRoyale GameMode = "kfo_battle_royale"
	ModeSingleplayer GameMode = "kfo_singleplayer"
)

type LobbyMode string

const (
	LobbyInvite   LobbyMode = "invite"
	LobbyPassword LobbyMode = "password"
	LobbyOpen     LobbyMode = "open"
)

type Player struct {
	ID, Name                                                                 string
	Teeth, MaxTeeth, Answer                                                  int
	Alive, Answered, JustLost, JustDied, LostThisGame, Eliminated, Connected bool
}

type Question struct {
	Text    string   `json:"text"`
	Options []string `json:"options"`
	Correct int      `json:"correct"`
}

type GameSettings struct {
	Topic           string    `json:"topic"`
	Difficulty      string    `json:"difficulty"`
	StartDifficulty string    `json:"startDifficulty"`
	NumQuestions    int       `json:"numQuestions"`
	TimePerQ        int       `json:"timePerQuestion"`
	NumOptions      int       `json:"numOptions"`
	NumTeeth        int       `json:"numTeeth"`
	Mode            GameMode  `json:"mode"`
	ShowTutorial    bool      `json:"showTutorial"`
	LobbyName       string    `json:"lobbyName"`
	LobbyMode       LobbyMode `json:"lobbyMode"`
	LobbyPassword   string    `json:"lobbyPassword"`
	WebSearch       bool      `json:"webSearch"`
	PlayIntro       bool      `json:"playIntro"`
}

type Game struct {
	mu                               sync.Mutex
	ID, InviteCode, HostID, HostUser string
	DelegatedTo                      string
	Phase                            GamePhase
	ErrorMsg                         string
	Settings                         GameSettings
	Players                          map[string]*Player
	PlayerOrder                      []string
	Questions                        []Question
	CurrentQ                         int
	TimerEnd                         time.Time
	timerCancel                      chan struct{}
	SomeoneLost                      bool
	AllWrong                         bool
	AllCorrect                       bool
	prefetchMu                       sync.Mutex
	prefetched                       []Question
	prefetching                      bool
	prefetchFailed                   bool
	prefetchErr                      string
	connections                      map[string]*websocket.Conn
	connMu                           sync.Mutex
}

// ─── Lobby Names ────────────────────────────────────────────────────────────────

var lobbyNames = []string{"Krone", "Amalgam", "Zahnschmelz", "Kiefer", "Brücke", "Wurzel", "Prothese", "Zahnstein", "Fluorid", "Plombe", "Implantat", "Retainer", "Backenzahn", "Inlay", "Zement", "Dentin", "Pulpa", "Gingivitis", "Brackets", "Aligner", "Abformung", "Okklusion", "Munddusche", "Fissur"}

func randomLobbyName() string {
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(lobbyNames))))
	return lobbyNames[int(n.Int64())]
}

// ─── Global Game Registry ───────────────────────────────────────────────────────

var (
	games   = make(map[string]*Game)
	gamesMu sync.Mutex
)

// ─── Game Methods ───────────────────────────────────────────────────────────────

func newGame(hostID, hostName, hostUser string) *Game {
	return &Game{
		ID:         uuid.New().String()[:8],
		InviteCode: strings.ToLower(uuid.New().String()[:8]),
		HostID:     hostID,
		HostUser:   hostUser,
		Phase:      PhaseLobby,
		Settings: GameSettings{
			Topic: "Allgemeinwissen", Difficulty: "mittel", StartDifficulty: "leicht",
			NumQuestions: 10, TimePerQ: 20, NumOptions: 4, NumTeeth: 5,
			Mode: ModeClassic, ShowTutorial: true,
			LobbyName: randomLobbyName(), LobbyMode: LobbyInvite,
			WebSearch: false, PlayIntro: true,
		},
		Players:     map[string]*Player{hostID: {ID: hostID, Name: hostName, Teeth: 5, MaxTeeth: 5, Alive: true, Answer: -1, Connected: true}},
		PlayerOrder: []string{hostID},
		connections: make(map[string]*websocket.Conn),
	}
}

func (g *Game) addPlayer(id, name string) bool {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.Phase != PhaseLobby {
		return false
	}
	if p, ok := g.Players[id]; ok {
		p.Connected = true
		p.Name = name
		return true
	}
	g.Players[id] = &Player{ID: id, Name: name, Teeth: g.Settings.NumTeeth, MaxTeeth: g.Settings.NumTeeth, Alive: true, Answer: -1, Connected: true}
	g.PlayerOrder = append(g.PlayerOrder, id)
	if len(g.Players) > 1 && g.Settings.Mode == ModeSingleplayer {
		g.Settings.Mode = ModeClassic
	}
	return true
}

func (g *Game) reconnectPlayerByID(playerID string, conn *websocket.Conn) bool {
	g.mu.Lock()
	defer g.mu.Unlock()
	p, ok := g.Players[playerID]
	if !ok || p.Connected {
		return false
	}
	p.Connected = true
	g.connMu.Lock()
	g.connections[playerID] = conn
	g.connMu.Unlock()
	return true
}

func (g *Game) removePlayer(id string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	delete(g.Players, id)
	o := make([]string, 0)
	for _, p := range g.PlayerOrder {
		if p != id {
			o = append(o, p)
		}
	}
	g.PlayerOrder = o
}

func cleanupIfEmpty(inviteCode string) {
	gamesMu.Lock()
	defer gamesMu.Unlock()
	g, ok := games[inviteCode]
	if !ok {
		return
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	for _, p := range g.Players {
		if p.Connected {
			return
		}
	}
	delete(games, inviteCode)
	logInfo("system", "system", "GAME_CLEANUP", fmt.Sprintf("invite=%s name=%s (no players)", inviteCode, g.Settings.LobbyName))
}

func getIntroDelay() int {
	v := getSetting("intro_delay")
	if v == "" {
		return 4
	}
	var d int
	fmt.Sscanf(v, "%d", &d)
	if d < 1 {
		d = 1
	}
	return d
}

func (g *Game) isEndless() bool {
	return g.Settings.Mode == ModeElimination || g.Settings.Mode == ModeBattleRoyale || g.Settings.Mode == ModeSingleplayer
}

func (g *Game) currentDifficulty() string {
	if g.Settings.Mode == ModeBattleRoyale || g.Settings.Mode == ModeSingleplayer {
		return escalatedDifficulty(g.Settings.StartDifficulty, g.CurrentQ/20)
	}
	return g.Settings.Difficulty
}

func (g *Game) prevQTexts() []string {
	t := make([]string, 0, len(g.Questions))
	for _, q := range g.Questions {
		t = append(t, q.Text)
	}
	return t
}

func (g *Game) countActive() int {
	n := 0
	for _, p := range g.Players {
		if p.Alive && p.Connected {
			n++
		}
	}
	return n
}

func (g *Game) startGame() error {
	g.mu.Lock()
	if g.Phase != PhaseLobby && g.Phase != PhaseLoading {
		g.mu.Unlock()
		return fmt.Errorf("not in lobby")
	}
	g.Phase = PhaseLoading
	for id, p := range g.Players {
		if !p.Connected {
			delete(g.Players, id)
			no := make([]string, 0)
			for _, pid := range g.PlayerOrder {
				if pid != id {
					no = append(no, pid)
				}
			}
			g.PlayerOrder = no
			continue
		}
		p.Teeth = g.Settings.NumTeeth
		p.MaxTeeth = g.Settings.NumTeeth
		p.Alive = true
		p.LostThisGame = false
		p.Eliminated = false
		p.JustLost = false
		p.JustDied = false
	}
	g.prefetched = nil
	g.prefetching = false
	g.mu.Unlock()
	g.broadcastState()

	numQ := g.Settings.NumQuestions
	diff := g.Settings.Difficulty
	if g.Settings.Mode == ModeElimination {
		numQ = 20
	}
	if g.Settings.Mode == ModeBattleRoyale || g.Settings.Mode == ModeSingleplayer {
		numQ = 20
		diff = g.Settings.StartDifficulty
	}
	logInfo("system", g.HostUser, "GAME_START", fmt.Sprintf("id=%s mode=%s topic=%s diff=%s players=%d websearch=%v", g.ID, g.Settings.Mode, g.Settings.Topic, diff, len(g.Players), g.Settings.WebSearch))
	qs, err := generateQuestions(g.Settings.Topic, diff, numQ, g.Settings.NumOptions, nil, g.Settings.WebSearch)
	if err != nil {
		g.setError(err.Error())
		return err
	}
	g.mu.Lock()
	g.Questions = qs
	g.CurrentQ = 0
	if g.Settings.ShowTutorial {
		g.Phase = PhaseTutorial
		g.mu.Unlock()
		g.broadcastState()
		return nil
	}
	if g.Settings.PlayIntro {
		g.Phase = PhaseIntro
		g.mu.Unlock()
		g.broadcastState()
		go func() { time.Sleep(time.Duration(getIntroDelay()) * time.Second); g.nextQuestion() }()
	} else {
		g.mu.Unlock()
		go g.nextQuestion()
	}
	return nil
}

func (g *Game) skipTutorial() {
	g.mu.Lock()
	if g.Phase != PhaseTutorial {
		g.mu.Unlock()
		return
	}
	if g.Settings.PlayIntro {
		g.Phase = PhaseIntro
		g.mu.Unlock()
		g.broadcastState()
		go func() { time.Sleep(time.Duration(getIntroDelay()) * time.Second); g.nextQuestion() }()
	} else {
		g.mu.Unlock()
		go g.nextQuestion()
	}
}

func (g *Game) triggerPrefetch() {
	if !g.isEndless() {
		return
	}
	rem := len(g.Questions) - g.CurrentQ
	if rem > 10 {
		return
	}
	g.prefetchMu.Lock()
	if g.prefetching || len(g.prefetched) > 0 || g.prefetchFailed {
		g.prefetchMu.Unlock()
		return
	}
	g.prefetching = true
	g.prefetchMu.Unlock()
	go func() {
		g.mu.Lock()
		prev := g.prevQTexts()
		nb := len(g.Questions) / 20
		diff := g.Settings.Difficulty
		if g.Settings.Mode == ModeBattleRoyale || g.Settings.Mode == ModeSingleplayer {
			diff = escalatedDifficulty(g.Settings.StartDifficulty, nb)
		}
		nO := g.Settings.NumOptions
		topic := g.Settings.Topic
		ws := g.Settings.WebSearch
		g.mu.Unlock()
		qs, err := generateQuestions(topic, diff, 20, nO, prev, ws)
		g.prefetchMu.Lock()
		if err != nil {
			g.prefetchFailed = true
			g.prefetchErr = err.Error()
			logWarn("system", "system", "PREFETCH_FAIL", err.Error())
		} else {
			g.prefetched = qs
		}
		g.prefetching = false
		g.prefetchMu.Unlock()
	}()
}

func (g *Game) nextQuestion() {
	g.mu.Lock()
	if g.Phase == PhaseEnd || g.Phase == PhaseError {
		g.mu.Unlock()
		return
	}
	active := g.countActive()
	if g.Settings.Mode == ModeSingleplayer && active <= 0 {
		g.Phase = PhaseEnd
		g.mu.Unlock()
		g.broadcastState()
		return
	}
	if g.isEndless() && g.Settings.Mode != ModeSingleplayer && active <= 1 {
		g.Phase = PhaseEnd
		g.mu.Unlock()
		g.broadcastState()
		return
	}
	if g.Settings.Mode == ModeClassic && (g.CurrentQ >= len(g.Questions) || active <= 1) {
		g.Phase = PhaseEnd
		g.mu.Unlock()
		g.broadcastState()
		return
	}
	if g.isEndless() && g.CurrentQ >= len(g.Questions) {
		g.prefetchMu.Lock()
		pf := g.prefetched
		g.prefetched = nil
		pfFailed := g.prefetchFailed
		pfErr := g.prefetchErr
		g.prefetchMu.Unlock()
		if len(pf) > 0 {
			g.Questions = append(g.Questions, pf...)
		} else if pfFailed {
			// Prefetch already failed — bail with that error
			msg := pfErr
			if msg == "" {
				msg = "Keine neuen Fragen verfügbar"
			}
			g.Phase = PhaseError
			g.ErrorMsg = msg
			g.mu.Unlock()
			g.broadcastState()
			return
		} else {
			g.Phase = PhaseRefill
			prev := g.prevQTexts()
			diff := g.currentDifficulty()
			ws := g.Settings.WebSearch
			g.mu.Unlock()
			g.broadcastState()
			nq, err := generateQuestions(g.Settings.Topic, diff, 20, g.Settings.NumOptions, prev, ws)
			g.mu.Lock()
			if err != nil {
				g.Phase = PhaseError
				g.ErrorMsg = err.Error()
				g.mu.Unlock()
				g.broadcastState()
				return
			}
			g.Questions = append(g.Questions, nq...)
			time.Sleep(1 * time.Second)
		}
	}
	for _, p := range g.Players {
		p.Answer = -1
		p.Answered = false
		p.JustLost = false
		p.JustDied = false
		if !p.Connected && p.Alive {
			p.Answered = true
		}
	}
	g.SomeoneLost = false
	g.AllWrong = false
	g.AllCorrect = false
	g.Phase = PhaseQuestion
	g.TimerEnd = time.Now().Add(time.Duration(g.Settings.TimePerQ) * time.Second)
	if g.timerCancel != nil {
		close(g.timerCancel)
	}
	g.timerCancel = make(chan struct{})
	cancel := g.timerCancel
	g.mu.Unlock()
	g.broadcastState()
	g.triggerPrefetch()
	go func() {
		t := time.NewTimer(time.Duration(g.Settings.TimePerQ)*time.Second + 500*time.Millisecond)
		defer t.Stop()
		select {
		case <-t.C:
			g.evaluateRound()
		case <-cancel:
		}
	}()
}

func (g *Game) submitAnswer(pid string, ans int) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.Phase != PhaseQuestion {
		return
	}
	p, ok := g.Players[pid]
	if !ok || !p.Alive || !p.Connected || p.Answered {
		return
	}
	p.Answer = ans
	p.Answered = true
	all := true
	for _, pl := range g.Players {
		if pl.Alive && pl.Connected && !pl.Answered {
			all = false
			break
		}
	}
	if all {
		if g.timerCancel != nil {
			close(g.timerCancel)
			g.timerCancel = nil
		}
		go g.evaluateRound()
	}
}

func (g *Game) evaluateRound() {
	g.mu.Lock()
	if g.Phase != PhaseQuestion {
		g.mu.Unlock()
		return
	}
	q := g.Questions[g.CurrentQ]
	sl := false
	correctCount := 0
	activeCount := 0
	for _, p := range g.Players {
		if !p.Alive {
			continue
		}
		activeCount++
		if !p.Connected || !p.Answered || p.Answer != q.Correct {
			p.Teeth--
			p.JustLost = true
			p.LostThisGame = true
			sl = true
		} else {
			correctCount++
		}
		if p.Teeth <= 0 {
			p.Teeth = 0
			p.Alive = false
			p.JustDied = true
			p.Eliminated = true
		}
	}
	g.SomeoneLost = sl
	g.AllWrong = activeCount > 0 && correctCount == 0
	g.AllCorrect = activeCount > 0 && correctCount == activeCount
	g.Phase = PhaseResults
	g.CurrentQ++
	g.mu.Unlock()
	g.broadcastState()
	go func() { time.Sleep(5 * time.Second); g.nextQuestion() }()
}

func (g *Game) forceEnd() {
	g.mu.Lock()
	if g.timerCancel != nil {
		close(g.timerCancel)
		g.timerCancel = nil
	}
	g.Phase = PhaseEnd
	g.mu.Unlock()
	g.broadcastState()
}

func (g *Game) getWinners() []map[string]string {
	type c struct {
		id, name string
		teeth    int
	}
	var all []c
	for _, p := range g.Players {
		all = append(all, c{p.ID, p.Name, p.Teeth})
	}
	sort.Slice(all, func(i, j int) bool { return all[i].teeth > all[j].teeth })
	if len(all) == 0 || all[0].teeth <= 0 {
		return nil
	}
	mx := all[0].teeth
	var w []map[string]string
	for _, x := range all {
		if x.teeth == mx {
			w = append(w, map[string]string{"id": x.id, "name": x.name})
		}
	}
	return w
}

func playerJSON(p *Player) map[string]interface{} {
	return map[string]interface{}{
		"id": p.ID, "name": p.Name, "teeth": p.Teeth, "maxTeeth": p.MaxTeeth,
		"alive": p.Alive, "answer": p.Answer, "answered": p.Answered,
		"justLost": p.JustLost, "justDied": p.JustDied, "lostThisGame": p.LostThisGame,
		"eliminated": p.Eliminated, "connected": p.Connected,
	}
}

func (g *Game) buildState() map[string]interface{} {
	ps := make([]map[string]interface{}, 0)
	for _, id := range g.PlayerOrder {
		if p, ok := g.Players[id]; ok {
			ps = append(ps, playerJSON(p))
		}
	}
	tQ := g.Settings.NumQuestions
	if g.isEndless() {
		tQ = len(g.Questions)
	}
	st := map[string]interface{}{
		"phase": g.Phase, "players": ps, "settings": g.Settings,
		"hostId": g.HostID, "gameId": g.ID, "inviteCode": g.InviteCode,
		"currentQuestion": g.CurrentQ, "totalQuestions": tQ, "timeLeft": 0, "someoneLost": g.SomeoneLost, "allWrong": g.AllWrong, "allCorrect": g.AllCorrect,
		"delegatedTo": g.DelegatedTo,
	}
	if g.Settings.Mode == ModeBattleRoyale || g.Settings.Mode == ModeSingleplayer {
		st["currentDifficulty"] = g.currentDifficulty()
	}
	switch g.Phase {
	case PhaseQuestion:
		if g.CurrentQ < len(g.Questions) {
			q := g.Questions[g.CurrentQ]
			st["question"] = map[string]interface{}{"text": q.Text, "options": q.Options, "index": g.CurrentQ}
		}
		r := time.Until(g.TimerEnd).Seconds()
		if r < 0 {
			r = 0
		}
		st["timeLeft"] = int(r)
	case PhaseResults:
		if g.CurrentQ > 0 && g.CurrentQ-1 < len(g.Questions) {
			q := g.Questions[g.CurrentQ-1]
			pr := make(map[string]string)
			for _, p := range g.Players {
				if p.JustLost {
					if !p.Answered || !p.Connected {
						pr[p.ID] = "timeout"
					} else {
						pr[p.ID] = "wrong"
					}
				} else if p.Alive {
					pr[p.ID] = "correct"
				}
			}
			st["results"] = map[string]interface{}{"correctAnswer": q.Correct, "playerResults": pr}
			st["question"] = map[string]interface{}{"text": q.Text, "options": q.Options, "index": g.CurrentQ - 1}
		}
	case PhaseEnd:
		st["winners"] = g.getWinners()
		if g.Settings.Mode == ModeSingleplayer {
			st["finalScore"] = g.CurrentQ
		}
	case PhaseError:
		st["errorMsg"] = g.ErrorMsg
	}
	return st
}

func (g *Game) broadcastState() {
	g.mu.Lock()
	st := g.buildState()
	g.mu.Unlock()
	d, _ := json.Marshal(map[string]interface{}{"type": "state", "payload": st})
	g.connMu.Lock()
	defer g.connMu.Unlock()
	for _, c := range g.connections {
		c.WriteMessage(websocket.TextMessage, d)
	}
}

func (g *Game) broadcastMsg(mt string, pl interface{}) {
	d, _ := json.Marshal(map[string]interface{}{"type": mt, "payload": pl})
	g.connMu.Lock()
	defer g.connMu.Unlock()
	for _, c := range g.connections {
		c.WriteMessage(websocket.TextMessage, d)
	}
}

func (g *Game) sendTo(pid, mt string, pl interface{}) {
	g.connMu.Lock()
	defer g.connMu.Unlock()
	if c, ok := g.connections[pid]; ok {
		d, _ := json.Marshal(map[string]interface{}{"type": mt, "payload": pl})
		c.WriteMessage(websocket.TextMessage, d)
	}
}

// setError puts the game into PhaseError with a user-facing message and broadcasts
func (g *Game) setError(msg string) {
	g.mu.Lock()
	if g.timerCancel != nil {
		close(g.timerCancel)
		g.timerCancel = nil
	}
	g.Phase = PhaseError
	g.ErrorMsg = msg
	g.mu.Unlock()
	g.broadcastState()
}
