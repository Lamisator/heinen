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
		w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self' https://cdnjs.cloudflare.com; style-src 'self' 'unsafe-inline'; img-src 'self' data:; font-src 'self'; connect-src 'self' wss: ws:")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		next.ServeHTTP(w, r)
	})
}

func main() {
	initLog()
	initDB()

	// Require secure admin bootstrap: fail if no admin exists and no bootstrap credentials provided
	if countAdmins() == 0 {
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

	logInfo("system", "system", "STARTUP", fmt.Sprintf("Port=%d", Port))

	// Migrate legacy passwords from plaintext/salt:hash to bcrypt (one-time)
	migrateLegacyPasswords()

	// Periodic session cleanup
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			cleanExpiredSessions()
			cleanExpiredPlayerTokens()
		}
	}()


	log.Printf("🦷 Heinen auf http://localhost:%d", Port)
	mux := http.NewServeMux()

	// API routes
	mux.HandleFunc("/api/login", handleLogin)
	mux.HandleFunc("/api/logout", handleLogout)
	mux.HandleFunc("/api/me", handleMe)
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

	// Static and dynamic
	mux.HandleFunc("/sounds/", handleSoundFile)
	mux.Handle("/fonts/", http.StripPrefix("/fonts/", http.FileServer(http.Dir("fonts"))))
	mux.HandleFunc("/ws", handleWS)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, indexHTML)
	})

	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", Port),
		Handler:           securityHeadersMiddleware(mux),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}
