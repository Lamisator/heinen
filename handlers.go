package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func jr(w http.ResponseWriter) { w.Header().Set("Content-Type", "application/json") }

func handleLogin(w http.ResponseWriter, r *http.Request) {
	jr(w)
	if r.Method != "POST" {
		w.WriteHeader(405)
		return
	}
	if !verifyCSRFToken(r) {
		logWarn(getIP(r), "", "CSRF_FAIL", "login")
		w.WriteHeader(403)
		json.NewEncoder(w).Encode(map[string]string{"error": "CSRF token invalid"})
		return
	}
	ip := getIP(r)
	if !limiter.CheckLoginRate(ip) {
		logWarn(ip, "", "LOGIN_RATELIMIT", "IP blocked after too many attempts")
		w.WriteHeader(429)
		json.NewEncoder(w).Encode(map[string]string{"error": "Zu viele Versuche. Bitte später versuchen."})
		return
	}
	var req struct{ Username, Password string }
	json.NewDecoder(r.Body).Decode(&req)
	if limiter.CheckAccountLockout(req.Username) {
		logWarn(ip, req.Username, "LOGIN_LOCKED", "account locked after failures")
		w.WriteHeader(429)
		json.NewEncoder(w).Encode(map[string]string{"error": "Konto gesperrt. Bitte später versuchen."})
		return
	}
	if !authenticateUser(req.Username, req.Password) {
		limiter.RecordAuthFailure(req.Username)
		logWarn(ip, req.Username, "LOGIN_FAIL", "")
		w.WriteHeader(401)
		json.NewEncoder(w).Encode(map[string]string{"error": "Ungültige Zugangsdaten"})
		return
	}
	limiter.RecordSuccess(ip, req.Username)
	logInfo(ip, req.Username, "LOGIN_OK", "")
	t := createSession(req.Username)
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookie,
		Value:    t,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   86400 * 7,
	})
	json.NewEncoder(w).Encode(map[string]interface{}{"ok": true, "username": req.Username, "isAdmin": isAdmin(req.Username)})
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	jr(w)
	if c, err := r.Cookie(SessionCookie); err == nil {
		deleteSession(c.Value)
	}
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookie,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}

func handleMe(w http.ResponseWriter, r *http.Request) {
	jr(w)
	u := getSessionUser(r)
	if u == "" {
		w.WriteHeader(401)
		json.NewEncoder(w).Encode(map[string]string{"error": "not logged in"})
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"username": u, "isAdmin": isAdmin(u)})
}

func handleChangePassword(w http.ResponseWriter, r *http.Request) {
	jr(w)
	if r.Method != "POST" {
		w.WriteHeader(405)
		return
	}
	if !verifyCSRFToken(r) {
		logWarn(getIP(r), "", "CSRF_FAIL", "change-password")
		w.WriteHeader(403)
		json.NewEncoder(w).Encode(map[string]string{"error": "CSRF token invalid"})
		return
	}
	u := getSessionUser(r)
	if u == "" {
		w.WriteHeader(401)
		return
	}
	var req struct{ OldPassword, NewPassword string }
	json.NewDecoder(r.Body).Decode(&req)
	if !authenticateUser(u, req.OldPassword) {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]string{"error": "Altes Passwort falsch"})
		return
	}
	if req.NewPassword == "" {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]string{"error": "Leer"})
		return
	}
	db.Exec("UPDATE users SET password = ? WHERE username = ?", hashPassword(req.NewPassword), u)
	deleteUserSessions(u)
	logInfo(getIP(r), u, "PASSWORD_CHANGE", "")
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}

func handleUsers(w http.ResponseWriter, r *http.Request) {
	jr(w)
	u := getSessionUser(r)
	ip := getIP(r)
	if u == "" || !isAdmin(u) {
		w.WriteHeader(403)
		json.NewEncoder(w).Encode(map[string]string{"error": "Zugriff verweigert"})
		return
	}
	if r.Method != "GET" && !verifyCSRFToken(r) {
		logWarn(ip, u, "CSRF_FAIL", "users-"+r.Method)
		w.WriteHeader(403)
		json.NewEncoder(w).Encode(map[string]string{"error": "CSRF token invalid"})
		return
	}
	ac := countAdmins()
	switch r.Method {
	case "GET":
		rows, _ := db.Query("SELECT id, username, is_admin, password FROM users ORDER BY id")
		if rows == nil {
			w.WriteHeader(500)
			return
		}
		defer rows.Close()
		users := make([]map[string]interface{}, 0)
		for rows.Next() {
			var id, admin int
			var un, pw string
			rows.Scan(&id, &un, &admin, &pw)
			users = append(users, map[string]interface{}{
				"id": id, "username": un, "isAdmin": admin == 1, "isOnlyAdmin": admin == 1 && ac <= 1,
				"passwordHashType": getPasswordHashType(pw),
			})
		}
		json.NewEncoder(w).Encode(users)
	case "POST":
		var req struct {
			Username, Password string
			IsAdmin            bool
		}
		json.NewDecoder(r.Body).Decode(&req)
		if req.Username == "" || req.Password == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(map[string]string{"error": "Felder ausfüllen"})
			return
		}
		av := 0
		if req.IsAdmin {
			av = 1
		}
		if _, err := db.Exec("INSERT INTO users (username,password,is_admin) VALUES (?,?,?)", req.Username, hashPassword(req.Password), av); err != nil {
			w.WriteHeader(409)
			json.NewEncoder(w).Encode(map[string]string{"error": "Existiert bereits"})
			return
		}
		logInfo(ip, u, "USER_CREATE", fmt.Sprintf("username=%s admin=%v", req.Username, req.IsAdmin))
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	case "DELETE":
		var req struct{ ID int }
		json.NewDecoder(r.Body).Decode(&req)
		var tu string
		var ta int
		db.QueryRow("SELECT username,is_admin FROM users WHERE id = ?", req.ID).Scan(&tu, &ta)
		if tu == u {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(map[string]string{"error": "Selbstlöschung"})
			return
		}
		if ta == 1 && ac <= 1 {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(map[string]string{"error": "Letzte*r Admin"})
			return
		}
		db.Exec("DELETE FROM users WHERE id = ?", req.ID)
		logInfo(ip, u, "USER_DELETE", "deleted="+tu)
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	case "PUT":
		var req struct {
			ID       int
			Password string
			IsAdmin  *bool
		}
		json.NewDecoder(r.Body).Decode(&req)
		if req.Password != "" {
			var targetUser string
			db.QueryRow("SELECT username FROM users WHERE id = ?", req.ID).Scan(&targetUser)
			db.Exec("UPDATE users SET password = ? WHERE id = ?", hashPassword(req.Password), req.ID)
			deleteUserSessions(targetUser)
			logInfo(ip, u, "USER_PW_RESET", fmt.Sprintf("id=%d user=%s", req.ID, targetUser))
		}
		if req.IsAdmin != nil {
			if !*req.IsAdmin {
				var ca int
				db.QueryRow("SELECT is_admin FROM users WHERE id = ?", req.ID).Scan(&ca)
				if ca == 1 && ac <= 1 {
					w.WriteHeader(400)
					json.NewEncoder(w).Encode(map[string]string{"error": "Letzte*r Admin"})
					return
				}
			}
			av := 0
			if *req.IsAdmin {
				av = 1
			}
			db.Exec("UPDATE users SET is_admin = ? WHERE id = ?", av, req.ID)
			logInfo(ip, u, "USER_ROLE", fmt.Sprintf("id=%d admin=%v", req.ID, *req.IsAdmin))
		}
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	}
}

func handleAIConfig(w http.ResponseWriter, r *http.Request) {
	jr(w)
	u := getSessionUser(r)
	if u == "" || !isAdmin(u) {
		w.WriteHeader(403)
		return
	}
	if r.Method != "GET" && !verifyCSRFToken(r) {
		logWarn(getIP(r), u, "CSRF_FAIL", "ai-config")
		w.WriteHeader(403)
		json.NewEncoder(w).Encode(map[string]string{"error": "CSRF token invalid"})
		return
	}
	switch r.Method {
	case "GET":
		oK := os.Getenv("OPENAI_API_KEY")
		if oK == "" {
			oK = getSetting("openai_api_key")
		}
		aK := os.Getenv("ANTHROPIC_API_KEY")
		if aK == "" {
			aK = getSetting("anthropic_api_key")
		}
		oM := ""
		if oK != "" {
			if len(oK) > 8 {
				oM = oK[:4] + "..." + oK[len(oK)-4:]
			} else {
				oM = "****"
			}
		}
		aM := ""
		if aK != "" {
			if len(aK) > 8 {
				aM = aK[:4] + "..." + aK[len(aK)-4:]
			} else {
				aM = "****"
			}
		}
		json.NewEncoder(w).Encode(map[string]string{
			"provider": getSetting("ai_provider"), "model": getSetting("ai_model"),
			"openaiKey": oM, "anthropicKey": aM, "introDelay": getSetting("intro_delay"),
			"vol_intro": getSetting("vol_intro"), "vol_background": getSetting("vol_background"),
			"vol_wrong": getSetting("vol_wrong"), "vol_answer": getSetting("vol_answer"),
			"vol_hurry": getSetting("vol_hurry"), "vol_timeout": getSetting("vol_timeout"),
			"vol_question": getSetting("vol_question"),
			"vol_allwrong": getSetting("vol_allwrong"), "vol_allcorrect": getSetting("vol_allcorrect"),
		})
	case "POST":
		var req struct {
			Provider, Model, OpenaiKey, AnthropicKey, IntroDelay                                    string
			VolIntro, VolBackground, VolWrong, VolAnswer, VolHurry, VolTimeout, VolQuestion, VolAllwrong, VolAllcorrect string
		}
		json.NewDecoder(r.Body).Decode(&req)
		if req.Provider != "" {
			setSetting("ai_provider", req.Provider)
		}
		if req.Model != "" {
			setSetting("ai_model", req.Model)
		}
		if req.OpenaiKey != "" {
			setSetting("openai_api_key", req.OpenaiKey)
		}
		if req.AnthropicKey != "" {
			setSetting("anthropic_api_key", req.AnthropicKey)
		}
		if req.IntroDelay != "" {
			setSetting("intro_delay", req.IntroDelay)
		}
		if req.VolIntro != "" {
			setSetting("vol_intro", req.VolIntro)
		}
		if req.VolBackground != "" {
			setSetting("vol_background", req.VolBackground)
		}
		if req.VolWrong != "" {
			setSetting("vol_wrong", req.VolWrong)
		}
		if req.VolAnswer != "" {
			setSetting("vol_answer", req.VolAnswer)
		}
		if req.VolHurry != "" {
			setSetting("vol_hurry", req.VolHurry)
		}
		if req.VolTimeout != "" {
			setSetting("vol_timeout", req.VolTimeout)
		}
		if req.VolQuestion != "" {
			setSetting("vol_question", req.VolQuestion)
		}
		if req.VolAllwrong != "" {
			setSetting("vol_allwrong", req.VolAllwrong)
		}
		if req.VolAllcorrect != "" {
			setSetting("vol_allcorrect", req.VolAllcorrect)
		}
		logInfo(getIP(r), u, "SETTINGS_CHANGE", "")
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	}
}

func handleTestAI(w http.ResponseWriter, r *http.Request) {
	jr(w)
	u := getSessionUser(r)
	if u == "" || !isAdmin(u) {
		w.WriteHeader(403)
		return
	}
	var req struct{ Provider, Key, Model string }
	json.NewDecoder(r.Body).Decode(&req)
	if req.Key == "" {
		if req.Provider == "anthropic" {
			req.Key = os.Getenv("ANTHROPIC_API_KEY")
			if req.Key == "" {
				req.Key = getSetting("anthropic_api_key")
			}
		} else {
			req.Key = os.Getenv("OPENAI_API_KEY")
			if req.Key == "" {
				req.Key = getSetting("openai_api_key")
			}
		}
	}
	if req.Model == "" {
		req.Model = getSetting("ai_model")
	}
	if req.Key == "" {
		json.NewEncoder(w).Encode(map[string]interface{}{"ok": false, "message": "Kein Key"})
		return
	}
	logInfo(getIP(r), u, "AI_TEST", fmt.Sprintf("provider=%s", req.Provider))
	ok, msg := testAIKey(req.Provider, req.Key, req.Model)
	json.NewEncoder(w).Encode(map[string]interface{}{"ok": ok, "message": msg})
}

func handleSounds(w http.ResponseWriter, r *http.Request) {
	jr(w)
	u := getSessionUser(r)
	if u == "" || !isAdmin(u) {
		w.WriteHeader(403)
		return
	}
	if r.Method != "GET" && !verifyCSRFToken(r) {
		logWarn(getIP(r), u, "CSRF_FAIL", "sounds-"+r.Method)
		w.WriteHeader(403)
		json.NewEncoder(w).Encode(map[string]string{"error": "CSRF token invalid"})
		return
	}
	switch r.Method {
	case "GET":
		res := map[string]interface{}{}
		for _, k := range soundTypes {
			res[k] = soundExists(k)
		}
		json.NewEncoder(w).Encode(res)
	case "POST":
		r.Body = http.MaxBytesReader(w, r.Body, 20*1024*1024)
		if err := r.ParseMultipartForm(20 << 20); err != nil {
			w.WriteHeader(413)
			json.NewEncoder(w).Encode(map[string]string{"error": "Datei zu groß"})
			return
		}
		st := r.FormValue("type")
		valid := false
		for _, t := range soundTypes {
			if t == st {
				valid = true
				break
			}
		}
		if !valid {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(map[string]string{"error": "Typ ungültig"})
			return
		}
		file, header, err := r.FormFile("file")
		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(map[string]string{"error": "Datei fehlt"})
			return
		}
		defer file.Close()
		ext := strings.ToLower(filepath.Ext(header.Filename))
		if ext != ".mp3" && ext != ".wav" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(map[string]string{"error": "Nur .mp3/.wav"})
			return
		}
		oe := getSetting(st + "_ext")
		if oe != "" && oe != ext {
			os.Remove(filepath.Join("sounds", st+oe))
		}
		dst := filepath.Join("sounds", st+ext)
		out, err := os.Create(dst)
		if err != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(map[string]string{"error": "Speicherfehler"})
			return
		}
		defer out.Close()
		written, err := io.CopyN(out, file, 20*1024*1024+1)
		if err != nil && err != io.EOF {
			os.Remove(dst)
			w.WriteHeader(413)
			json.NewEncoder(w).Encode(map[string]string{"error": "Datei zu groß"})
			return
		}
		if written > 20*1024*1024 {
			os.Remove(dst)
			w.WriteHeader(413)
			json.NewEncoder(w).Encode(map[string]string{"error": "Datei zu groß"})
			return
		}
		setSetting(st+"_ext", ext)
		logInfo(getIP(r), u, "SOUND_UPLOAD", fmt.Sprintf("type=%s size=%d", st, written))
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	case "DELETE":
		var req struct{ Type string }
		json.NewDecoder(r.Body).Decode(&req)
		if req.Type != "" {
			e := getSetting(req.Type + "_ext")
			if e != "" {
				os.Remove(filepath.Join("sounds", req.Type+e))
				setSetting(req.Type+"_ext", "")
			}
		}
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	}
}

func handleGlobalSounds(w http.ResponseWriter, r *http.Request) {
	jr(w)
	res := map[string]string{}
	for _, k := range soundTypes {
		if soundExists(k) {
			res[k] = "/sounds/" + k + getSetting(k+"_ext") + "?t=" + fmt.Sprint(time.Now().Unix())
		}
	}
	res["vol_intro"] = getSetting("vol_intro")
	res["vol_background"] = getSetting("vol_background")
	res["vol_wrong"] = getSetting("vol_wrong")
	res["vol_answer"] = getSetting("vol_answer")
	res["vol_hurry"] = getSetting("vol_hurry")
	res["vol_timeout"] = getSetting("vol_timeout")
	res["vol_question"] = getSetting("vol_question")
	res["vol_allwrong"] = getSetting("vol_allwrong")
	res["vol_allcorrect"] = getSetting("vol_allcorrect")
	json.NewEncoder(w).Encode(res)
}

func handleSoundFile(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/sounds/")
	if name == "" {
		http.NotFound(w, r)
		return
	}
	ok := false
	for _, p := range soundTypes {
		if strings.HasPrefix(name, p) {
			ok = true
			break
		}
	}
	if !ok {
		http.NotFound(w, r)
		return
	}
	fp := filepath.Join("sounds", name)
	if _, err := os.Stat(fp); err != nil {
		http.NotFound(w, r)
		return
	}
	ext := strings.ToLower(filepath.Ext(name))
	if ext == ".mp3" {
		w.Header().Set("Content-Type", "audio/mpeg")
	} else if ext == ".wav" {
		w.Header().Set("Content-Type", "audio/wav")
	}
	w.Header().Set("Cache-Control", "public, max-age=3600")
	http.ServeFile(w, r, fp)
}

func handleTutorial(w http.ResponseWriter, r *http.Request) {
	jr(w)
	data, err := os.ReadFile("tutorial.md")
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{"content": "# Tutorial\nKeine tutorial.md."})
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"content": string(data)})
}

// handleLogs returns structured log entries with optional search and level filter
// Query params: search, level, limit
func handleLogs(w http.ResponseWriter, r *http.Request) {
	jr(w)
	u := getSessionUser(r)
	if u == "" || !isAdmin(u) {
		w.WriteHeader(403)
		return
	}
	if r.Method != "GET" && !verifyCSRFToken(r) {
		logWarn(getIP(r), u, "CSRF_FAIL", "logs-"+r.Method)
		w.WriteHeader(403)
		json.NewEncoder(w).Encode(map[string]string{"error": "CSRF token invalid"})
		return
	}
	switch r.Method {
	case "GET":
		entries, _ := parseLogFile()
		search := r.URL.Query().Get("search")
		level := r.URL.Query().Get("level")
		filtered := filterLogs(entries, search, level)
		// Limit to last 1000 entries for performance
		limit := 1000
		if len(filtered) > limit {
			filtered = filtered[:limit]
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"entries": filtered,
			"total":   len(filtered),
		})
	case "DELETE":
		logMu.Lock()
		if logFile != nil {
			logFile.Close()
		}
		os.WriteFile(LogFile, []byte{}, 0644)
		logFile, _ = os.OpenFile(LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		logMu.Unlock()
		logInfo(getIP(r), u, "LOGS_CLEAR", "")
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	}
}

func handleLogsExport(w http.ResponseWriter, r *http.Request) {
	u := getSessionUser(r)
	if u == "" || !isAdmin(u) {
		w.WriteHeader(403)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename=heinen.log")
	data, _ := os.ReadFile(LogFile)
	w.Write(data)
}

func handleLobbies(w http.ResponseWriter, r *http.Request) {
	jr(w)
	gamesMu.Lock()
	defer gamesMu.Unlock()
	lobbies := make([]map[string]interface{}, 0)
	for _, g := range games {
		g.mu.Lock()
		if g.Phase == PhaseLobby && g.Settings.LobbyMode == LobbyOpen {
			connected := 0
			for _, p := range g.Players {
				if p.Connected {
					connected++
				}
			}
			lobbies = append(lobbies, map[string]interface{}{
				"inviteCode": g.InviteCode, "lobbyName": g.Settings.LobbyName,
				"lobbyMode": g.Settings.LobbyMode, "topic": g.Settings.Topic,
				"mode": g.Settings.Mode, "players": connected, "host": g.HostUser,
			})
		}
		g.mu.Unlock()
	}
	json.NewEncoder(w).Encode(lobbies)
}
