package main

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"
)

const sessionDuration = 7 * 24 * time.Hour

func initSessionTable() {
	db.Exec(`CREATE TABLE IF NOT EXISTS sessions (
		token TEXT PRIMARY KEY,
		username TEXT NOT NULL,
		created_at INTEGER NOT NULL,
		expires_at INTEGER NOT NULL
	)`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_sessions_expires ON sessions(expires_at)`)
	// Clean expired sessions on startup
	cleanExpiredSessions()
}

func generateToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// createSession persists a new session in the DB
func createSession(username string) string {
	token := generateToken()
	now := time.Now().Unix()
	exp := time.Now().Add(sessionDuration).Unix()
	_, err := db.Exec("INSERT INTO sessions (token, username, created_at, expires_at) VALUES (?, ?, ?, ?)", token, username, now, exp)
	if err != nil {
		logError("system", "system", "SESSION_CREATE_FAIL", err.Error())
	}
	return token
}

// getSessionUser looks up the user from a session cookie, returns "" if invalid/expired
func getSessionUser(r *http.Request) string {
	c, err := r.Cookie(SessionCookie)
	if err != nil {
		return ""
	}
	var username string
	var expires int64
	err = db.QueryRow("SELECT username, expires_at FROM sessions WHERE token = ?", c.Value).Scan(&username, &expires)
	if err != nil {
		return ""
	}
	if time.Now().Unix() > expires {
		db.Exec("DELETE FROM sessions WHERE token = ?", c.Value)
		return ""
	}
	return username
}

// deleteSession removes a session from the DB
func deleteSession(token string) {
	db.Exec("DELETE FROM sessions WHERE token = ?", token)
}

// cleanExpiredSessions removes all expired sessions
func cleanExpiredSessions() {
	res, err := db.Exec("DELETE FROM sessions WHERE expires_at < ?", time.Now().Unix())
	if err == nil {
		if n, _ := res.RowsAffected(); n > 0 {
			logInfo("system", "system", "SESSION_CLEANUP", "removed_expired="+intToStr(int(n)))
		}
	}
}

func intToStr(i int) string {
	if i == 0 {
		return "0"
	}
	buf := make([]byte, 0, 12)
	neg := false
	if i < 0 {
		neg = true
		i = -i
	}
	for i > 0 {
		buf = append([]byte{byte('0' + i%10)}, buf...)
		i /= 10
	}
	if neg {
		buf = append([]byte{'-'}, buf...)
	}
	return string(buf)
}
