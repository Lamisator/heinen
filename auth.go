package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

// hashPassword: SHA-256 with 16-byte salt, format "salt_hex:hash_hex"
func hashPassword(pw string) string {
	salt := make([]byte, 16)
	rand.Read(salt)
	h := sha256.Sum256(append(salt, []byte(pw)...))
	return hex.EncodeToString(salt) + ":" + hex.EncodeToString(h[:])
}

// verifyPassword checks a stored hash against a plaintext password
// Falls back to plaintext comparison for legacy passwords
func verifyPassword(stored, pw string) bool {
	parts := strings.SplitN(stored, ":", 2)
	if len(parts) != 2 {
		return stored == pw
	}
	salt, _ := hex.DecodeString(parts[0])
	expected, _ := hex.DecodeString(parts[1])
	h := sha256.Sum256(append(salt, []byte(pw)...))
	return string(h[:]) == string(expected)
}

// authenticateUser checks credentials and migrates legacy plaintext passwords
func authenticateUser(username, password string) bool {
	var stored string
	if db.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&stored) != nil {
		return false
	}
	if !verifyPassword(stored, password) {
		return false
	}
	if !strings.Contains(stored, ":") {
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
