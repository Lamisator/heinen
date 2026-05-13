package main

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"net/http"
	"strings"
)

const CSRFCookie = "heinen_csrf"

// generateCSRFToken creates a cryptographically random CSRF token
func generateCSRFToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// setCSRFToken sets a CSRF token cookie and returns the token for client use
func setCSRFToken(w http.ResponseWriter) string {
	token := generateCSRFToken()
	http.SetCookie(w, &http.Cookie{
		Name:     CSRFCookie,
		Value:    token,
		Path:     "/",
		HttpOnly: false,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   86400 * 7,
	})
	return token
}

// verifyCSRFToken checks the double-submit cookie pattern.
// JSON requests are exempt because cross-origin form submissions cannot set Content-Type: application/json.
func verifyCSRFToken(r *http.Request) bool {
	if strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		return true
	}
	cookie, err := r.Cookie(CSRFCookie)
	if err != nil {
		return false
	}
	headerToken := r.Header.Get("X-CSRF-Token")
	if headerToken == "" {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(cookie.Value), []byte(headerToken)) == 1
}
