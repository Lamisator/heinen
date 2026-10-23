package main

import (
	"net/http"
	"sync"
	"time"
)

// Session events push a "logout" event to every open browser tab whose session
// was revoked (password change, admin reset, role change, user deletion, logout
// in another tab), so the UI can drop the login state immediately.

type sessionSub struct {
	token, username string
	ch              chan struct{}
}

var (
	sessionSubs   = make(map[*sessionSub]struct{})
	sessionSubsMu sync.Mutex
)

// notifyRevokedSessions signals every subscriber of username whose session no
// longer exists; with username "" only the given token is checked.
func notifyRevokedSessions(username, token string) {
	sessionSubsMu.Lock()
	defer sessionSubsMu.Unlock()
	for s := range sessionSubs {
		if (username != "" && s.username != username) || (username == "" && s.token != token) {
			continue
		}
		var n int
		db.QueryRow("SELECT COUNT(*) FROM sessions WHERE token = ?", s.token).Scan(&n)
		if n == 0 {
			select {
			case s.ch <- struct{}{}:
			default:
			}
		}
	}
}

func handleSessionEvents(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie(SessionCookie)
	user := getSessionUser(r)
	if err != nil || user == "" {
		w.WriteHeader(401)
		return
	}
	rc := http.NewResponseController(w)
	// The server-wide read/write timeouts would cut the stream after 10 s
	rc.SetWriteDeadline(time.Time{})
	rc.SetReadDeadline(time.Time{})
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("X-Accel-Buffering", "no")

	sub := &sessionSub{token: c.Value, username: user, ch: make(chan struct{}, 1)}
	sessionSubsMu.Lock()
	sessionSubs[sub] = struct{}{}
	sessionSubsMu.Unlock()
	defer func() {
		sessionSubsMu.Lock()
		delete(sessionSubs, sub)
		sessionSubsMu.Unlock()
	}()

	logout := func() {
		w.Write([]byte("event: logout\ndata: {}\n\n"))
		rc.Flush()
	}
	w.Write([]byte(": ok\n\n"))
	if rc.Flush() != nil {
		return
	}
	tk := time.NewTicker(25 * time.Second)
	defer tk.Stop()
	for {
		select {
		case <-r.Context().Done():
			return
		case <-sub.ch:
			logout()
			return
		case <-tk.C:
			// Also catches expiry and anything revoked without a notification
			if getSessionUser(r) == "" {
				logout()
				return
			}
			w.Write([]byte(": ping\n\n"))
			if rc.Flush() != nil {
				return
			}
		}
	}
}
