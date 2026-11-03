package main

// Antwire-Login (games.antwire.net) und Lobby-Umzüge zwischen den Antwire-Spielen.
//
// Login: Der Hub ist ein OpenID-Connect-Provider vor Zitadel mit eigener
// Kontoliste (HEINEN_HUB_URL). Konten, für die der Hub bürgt, werden hier beim
// ersten Anmelden angelegt; Hub-Admins sind hier auch Admins. Solange niemand
// angemeldet ist, versucht der Browser höchstens alle 30 Minuten eine stille
// Anmeldung (prompt=none). Passwort- und Passkey-Login bleiben unverändert.
//
// Umzüge: Der*die Host wählt ein anderes Spiel; jeder Platz bekommt einen
// zufälligen Schlüssel, und jeder Browser wird mit seinem Schlüssel im
// URL-Fragment zum anderen Spiel geschickt. Der Browser des*der Host bringt die
// Liste der Spieler*innen mit (Namen, SHA-256 jedes Schlüssels) und eröffnet
// dort per POST /api/rooms/import eine neue Lobby; die anderen holen ihren
// Platz mit ihrem Schlüssel ab (POST /api/move/claim, wiederholt, bis die
// Lobby steht). Kein Spiel spricht mit einem anderen, niemand kann einen
// fremden Platz übernehmen.
//
// Hier angekommene Plätze sind reservierte, noch nicht verbundene Spieler*innen
// der neuen Lobby; abgeholt wird ein Platz mit einem Reconnect-Token über den
// WebSocket ("reconnect"), wie beim Wiederverbinden.

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	hubTryCookie = "heinen_sso_try"
	hubTryEvery  = 30 * time.Minute
	seatTTL      = 15 * time.Minute
	hubGameID    = "heinen"
)

var hubCfg struct {
	URL      string // https://games.antwire.net; leer = kein Antwire-Login, keine Umzüge
	ClientID string
	BaseURL  string // https://heinen.antwire.net; leer = aus dem Request
	Label    string
}

func hubEnabled() bool { return hubCfg.URL != "" }

func initHub() {
	hubCfg.URL = strings.TrimRight(os.Getenv("HEINEN_HUB_URL"), "/")
	hubCfg.BaseURL = strings.TrimRight(os.Getenv("HEINEN_BASE_URL"), "/")
	hubCfg.ClientID = os.Getenv("HEINEN_HUB_CLIENT_ID")
	if hubCfg.ClientID == "" {
		hubCfg.ClientID = hubGameID
	}
	hubCfg.Label = os.Getenv("HEINEN_OIDC_BUTTON_LABEL")
	if hubCfg.Label == "" {
		hubCfg.Label = "Antwire-Login"
	}
	// Verknüpfung der Konten mit dem Antwire-Login (Zitadel-Subject)
	db.Exec(`ALTER TABLE users ADD COLUMN oidc_issuer TEXT NOT NULL DEFAULT ''`)
	db.Exec(`ALTER TABLE users ADD COLUMN oidc_subject TEXT NOT NULL DEFAULT ''`)
	db.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS idx_users_oidc ON users(oidc_issuer, oidc_subject) WHERE oidc_subject != ''`)
}

// baseURL ist die öffentliche Adresse dieses Spiels.
func baseURL(r *http.Request) string {
	if hubCfg.BaseURL != "" {
		return hubCfg.BaseURL
	}
	scheme := "http"
	if r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}
	return scheme + "://" + r.Host
}

func isHTTPS(r *http.Request) bool { return strings.HasPrefix(baseURL(r), "https://") }

// handleConfig verrät dem Frontend, ob und wie der Antwire-Login läuft.
func handleConfig(w http.ResponseWriter, r *http.Request) {
	jr(w)
	out := map[string]interface{}{}
	if hubEnabled() {
		_, err := r.Cookie(hubTryCookie)
		u := getSessionUser(r)
		out["sso"] = map[string]string{"label": hubCfg.Label}
		out["hub"] = map[string]interface{}{"url": hubCfg.URL, "silent": u == "" && err != nil, "open": u == "" && openHosting(r.Context())}
	}
	json.NewEncoder(w).Encode(out)
}

func setSessionCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookie,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   86400 * 7,
	})
}

// ─── Spieleliste des Hubs ───────────────────────────────────────────────────────

type hubGame struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
	Move bool   `json:"move"`
	Min  int    `json:"min"`
	Max  int    `json:"max"`
	Bots bool   `json:"bots"`
	Up   bool   `json:"up"`
}

var hubList struct {
	mu    sync.Mutex
	games []hubGame
	at    time.Time
}

// hubGames ist die Spieleliste des Hubs, fünf Minuten zwischengespeichert.
func hubGames(ctx context.Context) ([]hubGame, error) {
	hubList.mu.Lock()
	if hubList.games != nil && time.Since(hubList.at) < 5*time.Minute {
		g := hubList.games
		hubList.mu.Unlock()
		return g, nil
	}
	hubList.mu.Unlock()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, hubCfg.URL+"/api/games", nil)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var list []hubGame
	if res.StatusCode != http.StatusOK || json.NewDecoder(res.Body).Decode(&list) != nil {
		return nil, fmt.Errorf("hub game list: %s", res.Status)
	}
	hubList.mu.Lock()
	hubList.games, hubList.at = list, time.Now()
	hubList.mu.Unlock()
	return list, nil
}

// handleMoveTargets listet die Spiele, in die eine Lobby umziehen kann.
func handleMoveTargets(w http.ResponseWriter, r *http.Request) {
	jr(w)
	out := []hubGame{}
	if hubEnabled() {
		list, err := hubGames(r.Context())
		if err != nil {
			w.WriteHeader(502)
			json.NewEncoder(w).Encode(map[string]string{"error": "Die Antwire-Spieleliste ist gerade nicht erreichbar."})
			return
		}
		for _, g := range list {
			if g.Move && g.ID != hubGameID {
				out = append(out, g)
			}
		}
	}
	json.NewEncoder(w).Encode(out)
}

// ─── Umzug hinaus ───────────────────────────────────────────────────────────────

type movePlayer struct {
	Name string `json:"name"`
	Key  string `json:"kh,omitempty"` // hex SHA-256 des Platzschlüssels
	Bot  bool   `json:"bot,omitempty"`
}

type movePlan struct {
	From    string       `json:"from"`
	Host    string       `json:"host,omitempty"` // Name des*der Host, für Host*s ohne Konto
	Players []movePlayer `json:"players"`
}

func keyHash(k string) string {
	s := sha256.Sum256([]byte(k))
	return hex.EncodeToString(s[:])
}

// moveGame schickt alle verbundenen Spieler*innen der Lobby in ein anderes
// Spiel; die Lobby hier wird danach geschlossen. Nur in der Lobby und nach
// Spielende, nur durch den*die Host.
func moveGame(ctx context.Context, g *Game, playerID, target string) error {
	if !hubEnabled() {
		return fmt.Errorf("Umziehen ist hier nicht eingerichtet.")
	}
	list, err := hubGames(ctx)
	if err != nil {
		return fmt.Errorf("Die Antwire-Spieleliste ist gerade nicht erreichbar.")
	}
	var to *hubGame
	for i := range list {
		if list[i].ID == target && list[i].Move && target != hubGameID {
			to = &list[i]
		}
	}
	if to == nil {
		return fmt.Errorf("Dieses Spiel nimmt gerade keine Lobbys an.")
	}
	g.mu.Lock()
	if playerID != g.HostID {
		g.mu.Unlock()
		return fmt.Errorf("Nur der*die Host kann umziehen.")
	}
	if g.Phase != PhaseLobby && g.Phase != PhaseEnd {
		g.mu.Unlock()
		return fmt.Errorf("Umziehen geht nur in der Lobby oder nach Spielende.")
	}
	plan := movePlan{From: hubGameID}
	keys := map[string]string{}
	humans := 0
	for _, pid := range g.PlayerOrder {
		p := g.Players[pid]
		if p == nil || !p.Connected {
			continue
		}
		humans++
		if pid == g.HostID {
			plan.Host = p.Name
			continue
		}
		k := generateToken()
		keys[pid] = k
		plan.Players = append(plan.Players, movePlayer{Name: p.Name, Key: keyHash(k)})
	}
	if to.Max > 0 && humans > to.Max {
		g.mu.Unlock()
		return fmt.Errorf("Für %s sind zu viele Leute in der Lobby (max. %d).", to.Name, to.Max)
	}
	data, _ := json.Marshal(plan)
	hostURL := to.URL + "/m#host=" + base64.RawURLEncoding.EncodeToString(data)
	host, code := g.HostID, g.InviteCode
	g.mu.Unlock()

	// Erst niemanden mehr hereinlassen, dann alle losschicken
	gamesMu.Lock()
	if games[code] == g {
		delete(games, code)
	}
	gamesMu.Unlock()
	g.connMu.Lock()
	for pid, c := range g.connections {
		url := to.URL
		if pid == host {
			url = hostURL
		} else if k, ok := keys[pid]; ok {
			url = to.URL + "/m#k=" + k
		}
		d, _ := json.Marshal(map[string]interface{}{"type": "moved", "payload": map[string]string{"url": url, "game": to.Name}})
		c.WriteMessage(websocket.TextMessage, d)
	}
	g.connMu.Unlock()
	logInfo("system", g.HostUser, "GAME_MOVED", fmt.Sprintf("invite=%s to=%s players=%d", code, to.ID, humans))
	return nil
}

// ─── Umzug herein ───────────────────────────────────────────────────────────────

type moveSeat struct {
	Code, Token string
	Exp         time.Time
}

var moveSeats = struct {
	sync.Mutex
	m map[string]moveSeat // Hash des Schlüssels -> reservierter Platz
}{m: map[string]moveSeat{}}

// moveName kürzt Namen wie im Beitrittsformular (max. 20 Zeichen).
func moveName(n string) string {
	n = strings.Join(strings.Fields(n), " ")
	if utf8.RuneCountInString(n) > 20 {
		n = string([]rune(n)[:20])
	}
	if n == "" {
		n = "Gast"
	}
	return n
}

// handleImport eröffnet eine Lobby für Spieler*innen, die aus einem anderen
// Spiel hierher umziehen. Der*die angemeldete Host bekommt den ersten Platz;
// die Plätze der anderen bleiben reserviert, bis sie abgeholt werden.
func handleImport(w http.ResponseWriter, r *http.Request) {
	jr(w)
	if r.Method != "POST" {
		w.WriteHeader(405)
		return
	}
	if !verifyCSRFToken(r) {
		w.WriteHeader(403)
		json.NewEncoder(w).Encode(map[string]string{"error": "CSRF token invalid"})
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 64*1024)
	var in movePlan
	if json.NewDecoder(r.Body).Decode(&in) != nil || len(in.Players) > 64 {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]string{"error": "Ungültiger Umzug"})
		return
	}
	user, hostName, anon := getSessionUser(r), "", false
	if user != "" {
		hostName = moveName(user)
	} else {
		// ohne Konto nur, solange der Hub das Hosten ohne Anmeldung erlaubt
		if strings.TrimSpace(in.Host) == "" {
			in.Host = "Host"
		}
		hu, err := anonHost(r.Context(), anonIP(r), in.Host)
		if err != nil {
			w.WriteHeader(401)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		user, hostName, anon = hu, moveName(in.Host), true
	}
	hostID := uuid.New().String()[:12]
	g := newGame(hostID, hostName, user)
	g.Anon, g.HostIP = anon, anonIP(r)
	g.Players[hostID].Connected = false // verbindet sich gleich per Reconnect-Token
	type claim struct {
		hash, token string
	}
	var claims []claim
	for _, p := range in.Players {
		if p.Bot || len(p.Key) != 64 {
			continue // HEINEN hat keine Computergegner*innen
		}
		if g.countParticipants() >= 64 {
			break
		}
		pid := uuid.New().String()[:12]
		g.Players[pid] = &Player{ID: pid, Name: moveName(p.Name), Teeth: g.Settings.NumTeeth, MaxTeeth: g.Settings.NumTeeth, Alive: true, Answer: -1}
		g.PlayerOrder = append(g.PlayerOrder, pid)
		claims = append(claims, claim{strings.ToLower(p.Key), issueReconnectToken(pid, g.ID)})
	}
	hostToken := issueReconnectToken(hostID, g.ID)
	gamesMu.Lock()
	games[g.InviteCode] = g
	gamesMu.Unlock()
	exp := time.Now().Add(seatTTL)
	moveSeats.Lock()
	for k, st := range moveSeats.m {
		if time.Now().After(st.Exp) {
			delete(moveSeats.m, k)
		}
	}
	for _, c := range claims {
		moveSeats.m[c.hash] = moveSeat{Code: g.InviteCode, Token: c.token, Exp: exp}
	}
	moveSeats.Unlock()
	// Kommt niemand an, verschwindet die Lobby wieder
	code := g.InviteCode
	time.AfterFunc(seatTTL, func() { cleanupIfEmpty(code) })
	logInfo(getIP(r), user, "GAME_IMPORT", fmt.Sprintf("invite=%s from=%s players=%d", code, moveName(in.From), len(claims)))
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(map[string]string{"code": code, "token": hostToken})
}

// handleClaim gibt einer umgezogenen Person ihren Platz, sobald der*die Host
// die Lobby eröffnet hat (bis dahin 202 {"pending": true}).
func handleClaim(w http.ResponseWriter, r *http.Request) {
	jr(w)
	if r.Method != "POST" {
		w.WriteHeader(405)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 4*1024)
	var in struct {
		Key string `json:"key"`
	}
	json.NewDecoder(r.Body).Decode(&in)
	// kein Rate-Limit: die Browser fragen wiederholt, ein Schlüssel hat 256 Zufallsbits
	h := keyHash(strings.TrimSpace(in.Key))
	moveSeats.Lock()
	st, ok := moveSeats.m[h]
	if ok {
		delete(moveSeats.m, h)
	}
	moveSeats.Unlock()
	if !ok || time.Now().After(st.Exp) {
		w.WriteHeader(202)
		json.NewEncoder(w).Encode(map[string]bool{"pending": true})
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"code": st.Code, "token": st.Token})
}

// ─── Hosten ohne Anmeldung ──────────────────────────────────────────────────────
//
// Solange der Hub es erlaubt (GET /api/policy, 30 s zwischengespeichert, nicht
// erreichbar = aus), kann jede*r mit einem Namen eine Lobby eröffnen. Host ist,
// wer die Lobby eröffnet hat: erkannt an der WebSocket-Verbindung bzw. am
// Reconnect-Token des Platzes, nie über die Tabelle users.
//
// Die KI-Fragen kosten Geld, daher gelten für Lobbys ohne Konto Grenzen:
// höchstens 60 solche Lobbys gleichzeitig, 10 neue je IP in 10 Minuten,
// 6 Spielstarts je IP in der Stunde, keine Internet-Recherche und in den
// endlosen Modi höchstens 100 Fragen je Spiel (danach ist das Spiel aus).

const (
	maxAnonGames      = 60
	anonMaxQuestions  = 100
	anonCreatePer10m  = 10
	anonStartsPerHour = 6
)

var openPolicy struct {
	sync.Mutex
	open bool
	at   time.Time
}

// openHosting: darf gerade jede*r ohne Konto eine Lobby eröffnen?
func openHosting(ctx context.Context) bool {
	if !hubEnabled() {
		return false
	}
	openPolicy.Lock()
	if time.Since(openPolicy.at) < 30*time.Second {
		o := openPolicy.open
		openPolicy.Unlock()
		return o
	}
	openPolicy.Unlock()
	ctx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()
	open := false
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, hubCfg.URL+"/api/policy", nil)
	if res, err := http.DefaultClient.Do(req); err == nil {
		var p struct {
			Open bool `json:"open_hosting"`
		}
		if res.StatusCode == http.StatusOK && json.NewDecoder(res.Body).Decode(&p) == nil {
			open = p.Open
		}
		res.Body.Close()
	}
	openPolicy.Lock()
	openPolicy.open, openPolicy.at = open, time.Now()
	openPolicy.Unlock()
	return open
}

// anonHits zählt Ereignisse je IP in einem Zeitfenster (Lobbys, Spielstarts).
var anonHits = struct {
	sync.Mutex
	m map[string][]time.Time
}{m: map[string][]time.Time{}}

func anonHit(key string, max int, window time.Duration) bool {
	anonHits.Lock()
	defer anonHits.Unlock()
	now := time.Now()
	var keep []time.Time
	for _, t := range anonHits.m[key] {
		if now.Sub(t) < window {
			keep = append(keep, t)
		}
	}
	if len(anonHits.m) > 20000 {
		anonHits.m = map[string][]time.Time{}
	}
	if len(keep) >= max {
		anonHits.m[key] = keep
		return false
	}
	anonHits.m[key] = append(keep, now)
	return true
}

// anonHost prüft, ob jemand ohne Konto mit diesem Namen eine Lobby eröffnen
// darf; liefert den Host-Eintrag für Logs/Verlauf oder eine Fehlermeldung.
func anonHost(ctx context.Context, ip, name string) (string, error) {
	if !openHosting(ctx) {
		return "", fmt.Errorf("Nicht eingeloggt")
	}
	if strings.TrimSpace(name) == "" {
		return "", fmt.Errorf("Namen eingeben")
	}
	gamesMu.Lock()
	n := 0
	for _, g := range games {
		if g.Anon {
			n++
		}
	}
	gamesMu.Unlock()
	if n >= maxAnonGames {
		return "", fmt.Errorf("Gerade sind zu viele Lobbys offen. Bitte versuche es später noch einmal.")
	}
	if !anonHit("create:"+ip, anonCreatePer10m, 10*time.Minute) {
		return "", fmt.Errorf("Zu viele neue Lobbys. Bitte später versuchen.")
	}
	return "ohne Konto: " + moveName(name), nil
}

// anonIP ist die IP für die Grenzen ohne Konto. Anders als getIP traut sie
// X-Forwarded-For auch von privaten Adressen (Reverse-Proxy im Docker-Netz),
// sonst teilten sich alle Gäste die IP des Proxys.
func anonIP(r *http.Request) string {
	remote, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		remote = r.RemoteAddr
	}
	if ip := net.ParseIP(remote); ip != nil && (ip.IsLoopback() || ip.IsPrivate()) {
		if f := r.Header.Get("X-Forwarded-For"); f != "" {
			return strings.TrimSpace(strings.Split(f, ",")[0])
		}
		if f := r.Header.Get("X-Real-IP"); f != "" {
			return strings.TrimSpace(f)
		}
	}
	return remote
}
