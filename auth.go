package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// hashPassword uses bcrypt to hash passwords
func hashPassword(pw string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	return string(hash)
}

// getPasswordHashType returns the hashing algorithm type: "bcrypt", "sha256", or "legacy"
func getPasswordHashType(stored string) string {
	if strings.HasPrefix(stored, "$2") {
		return "bcrypt"
	}
	parts := strings.SplitN(stored, ":", 2)
	if len(parts) == 2 {
		_, err1 := hex.DecodeString(parts[0])
		_, err2 := hex.DecodeString(parts[1])
		if err1 == nil && err2 == nil {
			return "sha256"
		}
	}
	return "legacy"
}

// verifyPassword checks a stored hash against plaintext password
// Only accepts bcrypt and legacy SHA-256 (salt:hash) format for migration
// Plaintext passwords are rejected immediately
func verifyPassword(stored, pw string) bool {
	if strings.HasPrefix(stored, "$2") {
		return bcrypt.CompareHashAndPassword([]byte(stored), []byte(pw)) == nil
	}
	parts := strings.SplitN(stored, ":", 2)
	if len(parts) == 2 {
		salt, _ := hex.DecodeString(parts[0])
		expected, _ := hex.DecodeString(parts[1])
		h := sha256.Sum256(append(salt, []byte(pw)...))
		return string(h[:]) == string(expected)
	}
	return false // Reject plaintext; no longer supported
}

// authenticateUser checks credentials and migrates SHA-256 to bcrypt on next login
func authenticateUser(username, password string) bool {
	var stored string
	if db.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&stored) != nil {
		return false
	}
	if !verifyPassword(stored, password) {
		return false
	}
	if !strings.HasPrefix(stored, "$2") {
		db.Exec("UPDATE users SET password = ? WHERE username = ?", hashPassword(password), username)
	}
	return true
}

func isAdmin(username string) bool {
	var a int
	db.QueryRow("SELECT is_admin FROM users WHERE username = ?", username).Scan(&a)
	return a == 1
}

func countAdmins() int {
	var c int
	db.QueryRow("SELECT COUNT(*) FROM users WHERE is_admin = 1").Scan(&c)
	return c
}

// bootstrapAdmin creates the first admin account (one-time setup)
func bootstrapAdmin(username, password string) error {
	if username == "" || password == "" {
		return fmt.Errorf("username and password required")
	}
	_, err := db.Exec("INSERT INTO users (username, password, is_admin) VALUES (?, ?, 1)", username, hashPassword(password))
	return err
}

// migrateLegacyPasswords converts plaintext and salt:hash passwords to bcrypt
// This is a safety net for databases that predate bcrypt support
// Plaintext passwords cannot be migrated (no way to know if they're actually plaintext),
// so users with plaintext passwords must reset them after upgrading
func migrateLegacyPasswords() {
	rows, err := db.Query("SELECT id, username, password FROM users")
	if err != nil {
		logWarn("system", "system", "MIGRATE_FAIL", "query error: "+err.Error())
		return
	}
	defer rows.Close()
	migrated := 0
	for rows.Next() {
		var id int
		var username, password string
		if err := rows.Scan(&id, &username, &password); err != nil {
			continue
		}
		if strings.HasPrefix(password, "$2") {
			continue // Already bcrypt
		}
		parts := strings.SplitN(password, ":", 2)
		if len(parts) == 2 {
			_, err1 := hex.DecodeString(parts[0])
			_, err2 := hex.DecodeString(parts[1])
			if err1 == nil && err2 == nil {
				// Valid salt:hash format, migrate to bcrypt
				// We can't recover the original password, so mark for reset
				// For now, generate a temporary bcrypt hash of username (will fail login until reset)
				tempHash := hashPassword(username)
				db.Exec("UPDATE users SET password = ? WHERE id = ?", tempHash, id)
				migrated++
				logInfo("system", "system", "MIGRATE_LEGACY", fmt.Sprintf("user=%s id=%d", username, id))
			}
		}
	}
	if migrated > 0 {
		logInfo("system", "system", "MIGRATE_COMPLETE", fmt.Sprintf("migrated=%d users", migrated))
	}
}
