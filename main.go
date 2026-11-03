package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	Port          = 8671
	SessionCookie = "heinen_session"
	LogFile       = "heinen.log"
)

func securityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data:; font-src 'self'; connect-src 'self' wss: ws:")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		next.ServeHTTP(w, r)
	})
}

func main() {
	initLog()
	initDB()
	initQuestionsDB()
	initHub()

	// Require secure admin bootstrap: fail if no admin exists and no bootstrap credentials provided
	// (with the Antwire login, the first hub admin to sign in becomes admin here)
	if countAdmins() == 0 && !(hubEnabled() && os.Getenv("HEINEN_ADMIN_PASSWORD") == "") {
		adminPass := os.Getenv("HEINEN_ADMIN_PASSWORD")
		if adminPass == "" {
			log.Fatal("ERROR: No admin account exists. Bootstrap via HEINEN_ADMIN_PASSWORD=<password> HEINEN_ADMIN_USER=<username> (default: admin)")
		}
		adminUser := os.Getenv("HEINEN_ADMIN_USER")
		if adminUser == "" {
			adminUser = "admin"
		}
		if err := bootstrapAdmin(adminUser, adminPass); err != nil {
			log.Fatal("ERROR: Failed to bootstrap admin:", err)
		}
		logInfo("system", "system", "INIT", "Admin account bootstrapped: "+adminUser)
	}

	addr := os.Getenv("HEINEN_ADDR")
	if addr == "" {
		addr = fmt.Sprintf(":%d", Port)
	}
	logInfo("system", "system", "STARTUP", "Addr="+addr+" Hub="+hubCfg.URL)

	// Migrate legacy passwords from plaintext/salt:hash to bcrypt (one-time)
	migrateLegacyPasswords()

	// Migrate any plaintext API keys to encrypted storage
	migratePlaintextAPIKeys()

	// Periodic session and rate-limiter cleanup
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			cleanExpiredSessions()
			cleanExpiredPlayerTokens()
			cleanExpiredPasskeySessions()
			limiter.CleanExpired()
		}
	}()


	log.Printf("🦷 Heinen auf %s", addr)
	mux := http.NewServeMux()

	// API routes
	mux.HandleFunc("/api/login", handleLogin)
	mux.HandleFunc("/api/logout", handleLogout)
	mux.HandleFunc("/api/me", handleMe)
	mux.HandleFunc("/api/session-events", handleSessionEvents)
	mux.HandleFunc("/api/change-password", handleChangePassword)
	mux.HandleFunc("/api/users", handleUsers)
	mux.HandleFunc("/api/ai-config", handleAIConfig)
	mux.HandleFunc("/api/test-ai", handleTestAI)
	mux.HandleFunc("/api/sounds", handleSounds)
	mux.HandleFunc("/api/global-sounds", handleGlobalSounds)
	mux.HandleFunc("/api/tutorial", handleTutorial)
	mux.HandleFunc("/api/logs", handleLogs)
	mux.HandleFunc("/api/logs/export", handleLogsExport)
	mux.HandleFunc("/api/lobbies", handleLobbies)
	mux.HandleFunc("/api/admin/games", handleGameHistory)
	mux.HandleFunc("/api/admin/games/", handleGameHistory)
	mux.HandleFunc("/api/admin/questions", handleAdminQuestions)
	mux.HandleFunc("/api/admin/questions/generate", handleAdminGenerateQuestions)
	mux.HandleFunc("/api/questions/report", handleReportQuestion)

	// Passkey / WebAuthn
	mux.HandleFunc("/api/passkey/register/begin", handlePasskeyRegisterBegin)
	mux.HandleFunc("/api/passkey/register/finish", handlePasskeyRegisterFinish)
	mux.HandleFunc("/api/passkey/login/begin", handlePasskeyLoginBegin)
	mux.HandleFunc("/api/passkey/login/finish", handlePasskeyLoginFinish)
	mux.HandleFunc("/api/passkey/credentials", handlePasskeyCredentials)

	// Antwire-Login + Lobby-Umzüge (hub.go, oidc.go)
	mux.HandleFunc("/api/config", handleConfig)
	mux.HandleFunc("/api/oidc/start", handleOIDCStart)
	mux.HandleFunc("/auth/oidc/callback", handleOIDCCallback)
	mux.HandleFunc("/api/move-targets", handleMoveTargets)
	mux.HandleFunc("/api/rooms/import", handleImport)
	mux.HandleFunc("/api/move/claim", handleClaim)

	// Static and dynamic
	mux.HandleFunc("/sounds/", handleSoundFile)
	mux.Handle("/fonts/", http.StripPrefix("/fonts/", http.FileServer(http.Dir("fonts"))))
	mux.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))
	mux.HandleFunc("/ws", handleWS)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		setCSRFToken(w)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, indexHTML)
	})

	server := &http.Server{
		Addr:              addr,
		Handler:           securityHeadersMiddleware(mux),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}
