package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"
)

const encPrefix = "enc:"

func getMasterKey() ([]byte, error) {
	s := os.Getenv("HEINEN_MASTER_KEY")
	if s == "" {
		return nil, fmt.Errorf("HEINEN_MASTER_KEY not set")
	}
	key, err := hex.DecodeString(s)
	if err != nil || len(key) != 32 {
		return nil, fmt.Errorf("HEINEN_MASTER_KEY must be 64 hex chars (32 bytes)")
	}
	return key, nil
}

func encryptValue(plaintext string) (string, error) {
	key, err := getMasterKey()
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ct := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return encPrefix + base64.StdEncoding.EncodeToString(ct), nil
}

func decryptValue(ciphertext string) (string, error) {
	if !strings.HasPrefix(ciphertext, encPrefix) {
		// Legacy plaintext value – return as-is and migrate on next write
		return ciphertext, nil
	}
	key, err := getMasterKey()
	if err != nil {
		return "", err
	}
	data, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(ciphertext, encPrefix))
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	if len(data) < gcm.NonceSize() {
		return "", fmt.Errorf("ciphertext too short")
	}
	nonce, ct := data[:gcm.NonceSize()], data[gcm.NonceSize():]
	pt, err := gcm.Open(nil, nonce, ct, nil)
	if err != nil {
		return "", err
	}
	return string(pt), nil
}

// getEncryptedSetting retrieves a setting and decrypts it.
func getEncryptedSetting(key string) string {
	val := getSetting(key)
	if val == "" {
		return ""
	}
	pt, err := decryptValue(val)
	if err != nil {
		logError("system", "system", "DECRYPT_ERROR", sanitizeLogField(key)+": "+sanitizeLogField(err.Error()))
		return ""
	}
	return pt
}

// setEncryptedSetting encrypts a value and stores it. Falls back to plaintext if no master key.
func setEncryptedSetting(key, value string) {
	if value == "" {
		setSetting(key, "")
		return
	}
	ct, err := encryptValue(value)
	if err != nil {
		logWarn("system", "system", "ENCRYPT_WARN", sanitizeLogField(key)+": master key unavailable, storing plaintext")
		setSetting(key, value)
		return
	}
	setSetting(key, ct)
}

// migratePlaintextAPIKeys re-encrypts any plaintext API keys if a master key is available.
func migratePlaintextAPIKeys() {
	if _, err := getMasterKey(); err != nil {
		return
	}
	for _, k := range []string{"openai_api_key", "anthropic_api_key"} {
		raw := getSetting(k)
		if raw == "" || strings.HasPrefix(raw, encPrefix) {
			continue
		}
		ct, err := encryptValue(raw)
		if err != nil {
			continue
		}
		setSetting(k, ct)
		logInfo("system", "system", "KEY_MIGRATE", sanitizeLogField(k)+" encrypted")
	}
}
