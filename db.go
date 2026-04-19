package main

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

var soundTypes = []string{"intro_sound", "background_sound", "wrong_sound", "answer_sound", "hurry_sound", "timeout_sound", "question_sound", "allwrong_sound", "allcorrect_sound"}

func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "heinen.db?_journal_mode=WAL")
	if err != nil {
		log.Fatal(err)
	}
	db.Exec(`CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY AUTOINCREMENT, username TEXT UNIQUE NOT NULL, password TEXT NOT NULL, is_admin INTEGER DEFAULT 0)`)
	db.Exec(`CREATE TABLE IF NOT EXISTS settings (key TEXT PRIMARY KEY, value TEXT)`)

	initSessionTable()

	var count int
	db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if count == 0 {
		db.Exec("INSERT INTO users (username, password, is_admin) VALUES (?, ?, 1)", "admin", "admin")
		logInfo("system", "system", "INIT", "Standard-Admin erstellt")
	}
	os.MkdirAll("sounds", 0755)
	os.MkdirAll("fonts", 0755)

	defaults := map[string]string{
		"ai_provider": "openai", "ai_model": "gpt-5.4-mini", "intro_delay": "4",
		"vol_intro": "0.6", "vol_background": "0.2", "vol_wrong": "0.6", "vol_answer": "0.6",
		"vol_hurry": "0.5", "vol_timeout": "0.6", "vol_question": "0.5",
		"vol_allwrong": "0.6", "vol_allcorrect": "0.6",
	}
	for k, v := range defaults {
		if getSetting(k) == "" {
			setSetting(k, v)
		}
	}
}

func getSetting(key string) string {
	var v string
	db.QueryRow("SELECT value FROM settings WHERE key = ?", key).Scan(&v)
	return v
}

func setSetting(key, value string) {
	_, err := db.Exec("INSERT OR REPLACE INTO settings (key, value) VALUES (?, ?)", key, value)
	if err != nil {
		logError("system", "system", "DB_ERROR", "setSetting "+key+": "+err.Error())
	}
}

func soundPath(st string) string {
	ext := getSetting(st + "_ext")
	if ext == "" {
		return ""
	}
	return filepath.Join("sounds", st+ext)
}

func soundExists(st string) bool {
	p := soundPath(st)
	if p == "" {
		return false
	}
	_, err := os.Stat(p)
	return err == nil
}
