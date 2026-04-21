package main

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// hashPassword uses bcrypt to hash passwords
func hashPassword(pw string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	return string(hash)
}

// verifyPassword checks a stored hash against plaintext password
// Supports both bcrypt and legacy SHA-256 (salt:hash) format
func verifyPassword(stored, pw string) bool {
	if strings.HasPrefix(stored, "$2") {
		return bcrypt.CompareHashAndPassword([]byte(stored), []byte(pw)) == nil
	}
	parts := strings.SplitN(stored, ":", 2)
	if len(parts) != 2 {
		return stored == pw
	}
	salt, _ := hex.DecodeString(parts[0])
	expected, _ := hex.DecodeString(parts[1])
	h := sha256.Sum256(append(salt, []byte(pw)...))
	return string(h[:]) == string(expected)
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
