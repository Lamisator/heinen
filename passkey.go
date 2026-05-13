package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
)

func initPasskeyTables() {
	db.Exec(`CREATE TABLE IF NOT EXISTS webauthn_credentials (
		id TEXT PRIMARY KEY,
		user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		name TEXT NOT NULL DEFAULT '',
		data TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`)
	db.Exec(`CREATE TABLE IF NOT EXISTS webauthn_sessions (
		id TEXT PRIMARY KEY,
		data TEXT NOT NULL,
		expires_at INTEGER NOT NULL
	)`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_wa_sessions_exp ON webauthn_sessions(expires_at)`)
}

func cleanExpiredPasskeySessions() {
	db.Exec("DELETE FROM webauthn_sessions WHERE expires_at < ?", time.Now().Unix())
}

// waForRequest creates a per-request WebAuthn instance using the Host header as RPID.
func waForRequest(r *http.Request) (*webauthn.WebAuthn, error) {
	host := r.Host
	if h, _, err := net.SplitHostPort(host); err == nil {
		host = h
	}
	return webauthn.New(&webauthn.Config{
		RPDisplayName: "Heinen – Das Zahnquiz",
		RPID:          host,
		RPOrigins: []string{
			"http://" + r.Host,
			"https://" + r.Host,
		},
	})
}

// wauthnUser implements webauthn.User for a DB user.
type wauthnUser struct {
	id          int64
	username    string
	credentials []webauthn.Credential
}

func (u *wauthnUser) WebAuthnID() []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(u.id))
	return b
}
func (u *wauthnUser) WebAuthnName() string        { return u.username }
func (u *wauthnUser) WebAuthnDisplayName() string { return u.username }
func (u *wauthnUser) WebAuthnCredentials() []webauthn.Credential {
	return u.credentials
}

func loadWAuthnUser(username string) (*wauthnUser, error) {
	var id int64
	if err := db.QueryRow("SELECT id FROM users WHERE username = ?", username).Scan(&id); err != nil {
		return nil, fmt.Errorf("user not found")
	}
	wu := &wauthnUser{id: id, username: username}
	rows, err := db.Query("SELECT data FROM webauthn_credentials WHERE user_id = ?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var raw string
		if err := rows.Scan(&raw); err != nil {
			continue
		}
		var cred webauthn.Credential
		if err := json.Unmarshal([]byte(raw), &cred); err != nil {
			continue
		}
		wu.credentials = append(wu.credentials, cred)
	}
	return wu, nil
}

func loadWAuthnUserByID(userID int64) (*wauthnUser, error) {
	var username string
	if err := db.QueryRow("SELECT username FROM users WHERE id = ?", userID).Scan(&username); err != nil {
		return nil, fmt.Errorf("user not found")
	}
	return loadWAuthnUser(username)
}

func storeWASession(id string, data *webauthn.SessionData) error {
	raw, err := json.Marshal(data)
	if err != nil {
		return err
	}
	exp := time.Now().Add(5 * time.Minute).Unix()
	_, err = db.Exec("INSERT OR REPLACE INTO webauthn_sessions (id, data, expires_at) VALUES (?, ?, ?)", id, string(raw), exp)
	return err
}

func loadWASession(id string) (*webauthn.SessionData, error) {
	var raw string
	var exp int64
	if err := db.QueryRow("SELECT data, expires_at FROM webauthn_sessions WHERE id = ?", id).Scan(&raw, &exp); err != nil {
		return nil, fmt.Errorf("session not found")
	}
	db.Exec("DELETE FROM webauthn_sessions WHERE id = ?", id)
	if time.Now().Unix() > exp {
		return nil, fmt.Errorf("session expired")
	}
	var sd webauthn.SessionData
	if err := json.Unmarshal([]byte(raw), &sd); err != nil {
		return nil, err
	}
	return &sd, nil
}

// --- Registration ---

func handlePasskeyRegisterBegin(w http.ResponseWriter, r *http.Request) {
	jr(w)
	u := getSessionUser(r)
	if u == "" {
		w.WriteHeader(401)
		json.NewEncoder(w).Encode(map[string]string{"error": "not logged in"})
		return
	}
	wn, err := waForRequest(r)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	wu, err := loadWAuthnUser(u)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	opts, sessionData, err := wn.BeginRegistration(wu,
		webauthn.WithResidentKeyRequirement(protocol.ResidentKeyRequirementRequired),
	)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	sid := generateToken()
	if err := storeWASession("reg:"+u+":"+sid, sessionData); err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	w.Header().Set("X-WA-Session", sid)
	json.NewEncoder(w).Encode(opts)
}

func handlePasskeyRegisterFinish(w http.ResponseWriter, r *http.Request) {
	jr(w)
	u := getSessionUser(r)
	if u == "" {
		w.WriteHeader(401)
		json.NewEncoder(w).Encode(map[string]string{"error": "not logged in"})
		return
	}
	sid := r.URL.Query().Get("sid")
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "Passkey"
	}
	if sid == "" {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]string{"error": "missing sid"})
		return
	}
	sessionData, err := loadWASession("reg:" + u + ":" + sid)
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]string{"error": "session invalid or expired"})
		return
	}
	wn, err := waForRequest(r)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	wu, err := loadWAuthnUser(u)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	cred, err := wn.FinishRegistration(wu, *sessionData, r)
	if err != nil {
		logWarn(getIP(r), u, "PASSKEY_REG_FAIL", err.Error())
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]string{"error": "Registrierung fehlgeschlagen: " + err.Error()})
		return
	}
	raw, err := json.Marshal(cred)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	credID := encodeBase64URL(cred.ID)
	var userID int64
	db.QueryRow("SELECT id FROM users WHERE username = ?", u).Scan(&userID)
	if _, err := db.Exec("INSERT OR REPLACE INTO webauthn_credentials (id, user_id, name, data) VALUES (?, ?, ?, ?)",
		credID, userID, name, string(raw)); err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	logInfo(getIP(r), u, "PASSKEY_REGISTER", "name="+name)
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}

// --- Login ---

func handlePasskeyLoginBegin(w http.ResponseWriter, r *http.Request) {
	jr(w)
	ip := getIP(r)
	if !limiter.CheckLoginRate(ip) {
		w.WriteHeader(429)
		json.NewEncoder(w).Encode(map[string]string{"error": "Zu viele Versuche"})
		return
	}
	wn, err := waForRequest(r)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	opts, sessionData, err := wn.BeginDiscoverableLogin()
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	sid := generateToken()
	if err := storeWASession("login:"+sid, sessionData); err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	w.Header().Set("X-WA-Session", sid)
	json.NewEncoder(w).Encode(opts)
}

func handlePasskeyLoginFinish(w http.ResponseWriter, r *http.Request) {
	jr(w)
	ip := getIP(r)
	sid := r.URL.Query().Get("sid")
	if sid == "" {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]string{"error": "missing sid"})
		return
	}
	sessionData, err := loadWASession("login:" + sid)
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]string{"error": "session invalid or expired"})
		return
	}
	wn, err := waForRequest(r)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	var loggedInUser string
	var foundCred *webauthn.Credential

	handler := func(rawID, userHandle []byte) (webauthn.User, error) {
		if len(userHandle) != 8 {
			return nil, fmt.Errorf("invalid user handle")
		}
		userID := int64(binary.BigEndian.Uint64(userHandle))
		wu, err := loadWAuthnUserByID(userID)
		if err != nil {
			return nil, err
		}
		loggedInUser = wu.username
		return wu, nil
	}

	cred, err := wn.FinishDiscoverableLogin(handler, *sessionData, r)
	if err != nil {
		logWarn(ip, "", "PASSKEY_LOGIN_FAIL", err.Error())
		limiter.RecordAuthFailure("passkey")
		w.WriteHeader(401)
		json.NewEncoder(w).Encode(map[string]string{"error": "Passkey-Authentifizierung fehlgeschlagen"})
		return
	}
	foundCred = cred

	// Update sign count
	raw, err := json.Marshal(foundCred)
	if err == nil {
		credID := encodeBase64URL(foundCred.ID)
		db.Exec("UPDATE webauthn_credentials SET data = ? WHERE id = ?", string(raw), credID)
	}

	limiter.RecordSuccess(ip, loggedInUser)
	logInfo(ip, loggedInUser, "PASSKEY_LOGIN_OK", "")
	t := createSession(loggedInUser)
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookie,
		Value:    t,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   86400 * 7,
	})
	json.NewEncoder(w).Encode(map[string]interface{}{"ok": true, "username": loggedInUser, "isAdmin": isAdmin(loggedInUser)})
}

// --- Credential management ---

func handlePasskeyCredentials(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		handlePasskeyList(w, r)
	case "DELETE":
		handlePasskeyDelete(w, r)
	default:
		w.WriteHeader(405)
	}
}

func handlePasskeyList(w http.ResponseWriter, r *http.Request) {
	jr(w)
	u := getSessionUser(r)
	if u == "" {
		w.WriteHeader(401)
		json.NewEncoder(w).Encode(map[string]string{"error": "not logged in"})
		return
	}
	rows, err := db.Query("SELECT id, name, created_at FROM webauthn_credentials WHERE user_id = (SELECT id FROM users WHERE username = ?) ORDER BY created_at DESC", u)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	defer rows.Close()
	type credRow struct {
		ID        string `json:"id"`
		Name      string `json:"name"`
		CreatedAt string `json:"createdAt"`
	}
	result := []credRow{}
	for rows.Next() {
		var cr credRow
		rows.Scan(&cr.ID, &cr.Name, &cr.CreatedAt)
		result = append(result, cr)
	}
	json.NewEncoder(w).Encode(result)
}

func handlePasskeyDelete(w http.ResponseWriter, r *http.Request) {
	jr(w)
	if r.Method != "DELETE" {
		w.WriteHeader(405)
		return
	}
	u := getSessionUser(r)
	if u == "" {
		w.WriteHeader(401)
		json.NewEncoder(w).Encode(map[string]string{"error": "not logged in"})
		return
	}
	var req struct{ ID string }
	json.NewDecoder(r.Body).Decode(&req)
	if req.ID == "" {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]string{"error": "id required"})
		return
	}
	res, err := db.Exec("DELETE FROM webauthn_credentials WHERE id = ? AND user_id = (SELECT id FROM users WHERE username = ?)", req.ID, u)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		w.WriteHeader(404)
		json.NewEncoder(w).Encode(map[string]string{"error": "not found"})
		return
	}
	logInfo(getIP(r), u, "PASSKEY_DELETE", "id="+req.ID)
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}

func encodeBase64URL(b []byte) string {
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"
	out := make([]byte, 0, (len(b)*4+2)/3)
	for i := 0; i < len(b); i += 3 {
		b0 := b[i]
		var b1, b2 byte
		if i+1 < len(b) {
			b1 = b[i+1]
		}
		if i+2 < len(b) {
			b2 = b[i+2]
		}
		out = append(out, chars[b0>>2])
		out = append(out, chars[(b0&0x3)<<4|b1>>4])
		if i+1 < len(b) {
			out = append(out, chars[(b1&0xF)<<2|b2>>6])
		}
		if i+2 < len(b) {
			out = append(out, chars[b2&0x3F])
		}
	}
	return string(out)
}
