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

	// Periodic session cleanup
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			cleanExpiredSessions()
		}
	}()

	// API routes
	http.HandleFunc("/api/login", handleLogin)
	http.HandleFunc("/api/logout", handleLogout)
	http.HandleFunc("/api/me", handleMe)
	http.HandleFunc("/api/change-password", handleChangePassword)
	http.HandleFunc("/api/users", handleUsers)
	http.HandleFunc("/api/ai-config", handleAIConfig)
	http.HandleFunc("/api/test-ai", handleTestAI)
	http.HandleFunc("/api/sounds", handleSounds)
	http.HandleFunc("/api/global-sounds", handleGlobalSounds)
	http.HandleFunc("/api/tutorial", handleTutorial)
	http.HandleFunc("/api/logs", handleLogs)
	http.HandleFunc("/api/logs/export", handleLogsExport)
	http.HandleFunc("/api/lobbies", handleLobbies)

	// Static and dynamic
	http.HandleFunc("/sounds/", handleSoundFile)
	http.Handle("/fonts/", http.StripPrefix("/fonts/", http.FileServer(http.Dir("fonts"))))
	http.HandleFunc("/ws", handleWS)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, indexHTML)
	})

	log.Printf("🦷 Heinen auf http://localhost:%d", Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", Port), nil))
}
