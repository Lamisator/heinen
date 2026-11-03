package main

// Anmeldung über den Antwire-Login per OpenID Connect: Auth-Code + PKCE
// (öffentlicher Client, Client-ID = "heinen"), Redirect auf
// <Basis-URL>/auth/oidc/callback. Flows: "login" (Knopf), "silent"
// (prompt=none beim Seitenaufruf; ohne Hub-Sitzung geht es einfach als Gast
// weiter) und "link" (ein bestehendes Konto mit dem Antwire-Login verknüpfen).
// Registrierung bleibt zu: angelegt werden nur Konten, für die der Hub bürgt
// (er meldet nur freigeschaltete Konten an).

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

const (
	oidcCookie = "heinen_oidc"
	oidcTTL    = 10 * time.Minute
)

var oidcProv struct {
	mu       sync.Mutex
	provider *oidc.Provider
	verifier *oidc.IDTokenVerifier
	at       time.Time
}

// oidcProvider holt die Discovery des Hubs erst bei Bedarf und behält sie eine
// Stunde, damit ein gerade nicht erreichbarer Hub den Start nicht verhindert.
func oidcProvider(ctx context.Context) (*oidc.Provider, *oidc.IDTokenVerifier, error) {
	st := &oidcProv
	st.mu.Lock()
	defer st.mu.Unlock()
	if st.provider != nil && time.Since(st.at) < time.Hour {
		return st.provider, st.verifier, nil
	}
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	p, err := oidc.NewProvider(ctx, hubCfg.URL)
	if err != nil {
		if st.provider != nil {
			return st.provider, st.verifier, nil
		}
		return nil, nil, fmt.Errorf("identity provider unreachable: %w", err)
	}
	st.provider, st.verifier, st.at = p, p.Verifier(&oidc.Config{ClientID: hubCfg.ClientID}), time.Now()
	return st.provider, st.verifier, nil
}

func oauthConfig(r *http.Request, p *oidc.Provider) *oauth2.Config {
	return &oauth2.Config{ClientID: hubCfg.ClientID, Endpoint: p.Endpoint(), RedirectURL: baseURL(r) + "/auth/oidc/callback",
		Scopes: []string{oidc.ScopeOpenID, "profile", "email"}}
}

// oidcFlow liegt unter dem state-Wert in der Tabelle webauthn_sessions.
type oidcFlow struct {
	Flow     string `json:"flow"` // login, silent, link
	Bind     string `json:"bind"` // muss zum Cookie passen: derselbe Browser beginnt und beendet
	Verifier string `json:"verifier"`
	Nonce    string `json:"nonce"`
	User     string `json:"user,omitempty"` // link
	Next     string `json:"next,omitempty"`
}

func saveFlow(state string, f oidcFlow) error {
	raw, _ := json.Marshal(f)
	_, err := db.Exec("INSERT OR REPLACE INTO webauthn_sessions (id, data, expires_at) VALUES (?, ?, ?)", "oidc:"+state, string(raw), time.Now().Add(oidcTTL).Unix())
	return err
}

func takeFlow(state string) (oidcFlow, error) {
	var f oidcFlow
	var raw string
	var exp int64
	if state == "" || db.QueryRow("SELECT data, expires_at FROM webauthn_sessions WHERE id = ?", "oidc:"+state).Scan(&raw, &exp) != nil {
		return f, errors.New("flow not found")
	}
	db.Exec("DELETE FROM webauthn_sessions WHERE id = ?", "oidc:"+state)
	if time.Now().Unix() > exp || json.Unmarshal([]byte(raw), &f) != nil {
		return f, errors.New("flow expired")
	}
	return f, nil
}

var oidcStarts = struct {
	sync.Mutex
	m map[string][]time.Time
}{m: map[string][]time.Time{}}

// oidcRate erlaubt 30 Anmeldestarts je IP in 15 Minuten (eigener Zähler, damit
// stille Anmeldungen nicht das Kontingent des Passwort-Logins aufbrauchen).
func oidcRate(ip string) bool {
	oidcStarts.Lock()
	defer oidcStarts.Unlock()
	now := time.Now()
	var keep []time.Time
	for _, t := range oidcStarts.m[ip] {
		if now.Sub(t) < 15*time.Minute {
			keep = append(keep, t)
		}
	}
	if len(oidcStarts.m) > 10000 {
		oidcStarts.m = map[string][]time.Time{}
	}
	if len(keep) >= 30 {
		oidcStarts.m[ip] = keep
		return false
	}
	oidcStarts.m[ip] = append(keep, now)
	return true
}

// safeNext hält Weiterleitungen auf dieser Seite.
func safeNext(p string) string {
	if !strings.HasPrefix(p, "/") || strings.HasPrefix(p, "//") || strings.HasPrefix(p, "/\\") {
		return "/"
	}
	return p
}

// handleOIDCStart liefert die URL der Anmeldeseite des Hubs.
func handleOIDCStart(w http.ResponseWriter, r *http.Request) {
	jr(w)
	if r.Method != "POST" {
		w.WriteHeader(405)
		return
	}
	if !hubEnabled() {
		w.WriteHeader(404)
		json.NewEncoder(w).Encode(map[string]string{"error": "Antwire-Login ist nicht eingerichtet"})
		return
	}
	if !verifyCSRFToken(r) {
		w.WriteHeader(403)
		json.NewEncoder(w).Encode(map[string]string{"error": "CSRF token invalid"})
		return
	}
	if !oidcRate(getIP(r)) {
		w.WriteHeader(429)
		json.NewEncoder(w).Encode(map[string]string{"error": "Zu viele Versuche. Bitte später versuchen."})
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 8*1024)
	var in struct{ Flow, Next string }
	json.NewDecoder(r.Body).Decode(&in)
	f := oidcFlow{Flow: in.Flow, Bind: generateToken(), Verifier: generateToken(), Nonce: generateToken(), Next: safeNext(in.Next)}
	switch in.Flow {
	case "login":
	case "silent":
		// höchstens alle 30 Minuten, damit Gäste nicht im Kreis geschickt werden
		http.SetCookie(w, &http.Cookie{Name: hubTryCookie, Value: "1", Path: "/", HttpOnly: true, Secure: isHTTPS(r),
			SameSite: http.SameSiteLaxMode, MaxAge: int(hubTryEvery.Seconds())})
	case "link":
		if f.User = getSessionUser(r); f.User == "" {
			w.WriteHeader(401)
			json.NewEncoder(w).Encode(map[string]string{"error": "Bitte zuerst anmelden"})
			return
		}
	default:
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]string{"error": "unknown flow"})
		return
	}
	p, _, err := oidcProvider(r.Context())
	if err != nil {
		logWarn(getIP(r), "", "OIDC_UNREACHABLE", err.Error())
		w.WriteHeader(502)
		json.NewEncoder(w).Encode(map[string]string{"error": "Der Antwire-Login ist gerade nicht erreichbar."})
		return
	}
	state := generateToken()
	if err := saveFlow(state, f); err != nil {
		w.WriteHeader(500)
		return
	}
	http.SetCookie(w, &http.Cookie{Name: oidcCookie, Value: f.Bind, Path: "/auth/oidc/", HttpOnly: true, Secure: isHTTPS(r),
		SameSite: http.SameSiteLaxMode, MaxAge: int(oidcTTL.Seconds())})
	sum := sha256.Sum256([]byte(f.Verifier))
	opts := []oauth2.AuthCodeOption{
		oauth2.SetAuthURLParam("code_challenge", base64.RawURLEncoding.EncodeToString(sum[:])),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
		oidc.Nonce(f.Nonce),
	}
	if in.Flow == "silent" {
		opts = append(opts, oauth2.SetAuthURLParam("prompt", "none"))
	}
	json.NewEncoder(w).Encode(map[string]string{"url": oauthConfig(r, p).AuthCodeURL(state, opts...)})
}

func clearOIDCCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{Name: oidcCookie, Value: "", Path: "/auth/oidc/", HttpOnly: true, MaxAge: -1})
}

func oidcFail(w http.ResponseWriter, r *http.Request, code string) {
	clearOIDCCookie(w)
	http.Redirect(w, r, "/?sso_error="+url.QueryEscape(code), http.StatusFound)
}

func handleOIDCCallback(w http.ResponseWriter, r *http.Request) {
	if !hubEnabled() {
		http.Error(w, "Antwire-Login ist nicht eingerichtet", http.StatusNotFound)
		return
	}
	ip := getIP(r)
	q := r.URL.Query()
	f, ferr := takeFlow(q.Get("state"))
	if e := q.Get("error"); e != "" {
		if ferr == nil && f.Flow == "silent" {
			// beim Hub nicht angemeldet: als Gast weiter
			clearOIDCCookie(w)
			http.Redirect(w, r, safeNext(f.Next), http.StatusFound)
			return
		}
		logWarn(ip, "", "OIDC_PROVIDER_ERROR", e+": "+q.Get("error_description"))
		if e == "access_denied" || e == "login_required" {
			oidcFail(w, r, "cancelled")
		} else {
			oidcFail(w, r, "provider")
		}
		return
	}
	ck, err := r.Cookie(oidcCookie)
	if ferr != nil || err != nil || subtle.ConstantTimeCompare([]byte(ck.Value), []byte(f.Bind)) != 1 || q.Get("code") == "" {
		oidcFail(w, r, "expired")
		return
	}
	p, verifier, err := oidcProvider(r.Context())
	if err != nil {
		oidcFail(w, r, "provider")
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
	defer cancel()
	tok, err := oauthConfig(r, p).Exchange(ctx, q.Get("code"), oauth2.SetAuthURLParam("code_verifier", f.Verifier))
	if err != nil {
		logWarn(ip, "", "OIDC_TOKEN_FAIL", err.Error())
		oidcFail(w, r, "provider")
		return
	}
	rawID, _ := tok.Extra("id_token").(string)
	idt, err := verifier.Verify(ctx, rawID)
	if err != nil || idt.Nonce != f.Nonce {
		logWarn(ip, "", "OIDC_IDTOKEN_FAIL", fmt.Sprint(err))
		oidcFail(w, r, "provider")
		return
	}
	claims := map[string]interface{}{}
	idt.Claims(&claims)
	if ui, err := p.UserInfo(ctx, oauth2.StaticTokenSource(tok)); err == nil {
		extra := map[string]interface{}{}
		if ui.Claims(&extra) == nil {
			for k, v := range extra {
				if _, ok := claims[k]; !ok {
					claims[k] = v
				}
			}
		}
	}
	str := func(k string) string { v, _ := claims[k].(string); return strings.TrimSpace(v) }
	hubAdmin, _ := claims["antwire_admin"].(bool)
	iss, sub := idt.Issuer, idt.Subject

	if f.Flow == "link" {
		me := getSessionUser(r)
		if me == "" || me != f.User {
			oidcFail(w, r, "relogin")
			return
		}
		if _, err := db.Exec("UPDATE users SET oidc_issuer = ?, oidc_subject = ? WHERE username = ?", iss, sub, me); err != nil {
			oidcFail(w, r, "taken") // schon mit einem anderen Konto verknüpft
			return
		}
		logInfo(ip, me, "OIDC_LINK", "")
		clearOIDCCookie(w)
		http.Redirect(w, r, "/?sso=linked", http.StatusFound)
		return
	}

	var user string
	if db.QueryRow("SELECT username FROM users WHERE oidc_issuer = ? AND oidc_subject = ?", iss, sub).Scan(&user) != nil {
		user, err = provisionOIDC(iss, sub, str, hubAdmin)
		if err != nil {
			logError(ip, "", "OIDC_PROVISION_FAIL", err.Error())
			if f.Flow == "silent" {
				clearOIDCCookie(w)
				http.Redirect(w, r, safeNext(f.Next), http.StatusFound)
				return
			}
			oidcFail(w, r, "internal")
			return
		}
		logInfo(ip, user, "OIDC_USER_CREATE", fmt.Sprintf("admin=%v", hubAdmin))
	} else {
		user = syncProfile(ip, user, str("name"), hubAdmin) // Name und Admin-Recht des Hubs
	}
	limiter.RecordSuccess(ip, user)
	logInfo(ip, user, "OIDC_LOGIN_OK", f.Flow)
	setSessionCookie(w, createSession(user))
	clearOIDCCookie(w)
	http.Redirect(w, r, safeNext(f.Next), http.StatusFound)
}

// cleanUserName lässt 2–32 Zeichen ohne Steuerzeichen durch.
func cleanUserName(n string) (string, bool) {
	n = strings.Join(strings.Fields(n), " ")
	if c := utf8.RuneCountInString(n); c < 2 || c > 32 {
		return "", false
	}
	for _, r := range n {
		if unicode.IsControl(r) {
			return "", false
		}
	}
	return n, true
}

// provisionOIDC legt ein Konto für eine Identität an, für die der Hub bürgt.
// Ist der Name vergeben, wird durchnummeriert; ein Passwort gibt es nicht
// (Anmeldung über den Antwire-Login oder einen später angelegten Passkey).
func provisionOIDC(iss, sub string, str func(string) string, admin bool) (string, error) {
	pref := str("preferred_username")
	if i := strings.IndexByte(pref, '@'); i > 0 {
		pref = pref[:i] // Zitadel-Loginnamen sehen aus wie user@org.domain
	}
	base := ""
	for _, c := range []string{pref, str("given_name"), str("name"), strings.SplitN(str("email"), "@", 2)[0]} {
		if n, ok := cleanUserName(c); ok {
			base = n
			break
		}
	}
	if base == "" {
		sum := sha256.Sum256([]byte(sub))
		base = fmt.Sprintf("sso-%x", sum[:3])
	}
	av := 0
	if admin {
		av = 1
	}
	name := base
	for i := 2; i < 50; i++ {
		_, err := db.Exec("INSERT INTO users (username, password, is_admin, oidc_issuer, oidc_subject) VALUES (?, '', ?, ?, ?)", name, av, iss, sub)
		if err == nil {
			return name, nil
		}
		if !strings.Contains(err.Error(), "UNIQUE") {
			return "", err
		}
		name = fmt.Sprintf("%s %d", string([]rune(base)[:min(29, len([]rune(base)))]), i)
	}
	return "", errors.New("no free name")
}

// handleHubProfile übernimmt eine vom Hub geschickte Änderung von Name oder
// Admin-Recht: ein ES256-JWT wie ein ID-Token (Audience = "heinen") mit
// "antwire_event":"profile", geprüft mit dem Verifier der Anmeldung. Ohne
// CSRF-Token: das JWT ist der Nachweis.
func handleHubProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(405)
		return
	}
	if !hubEnabled() {
		http.NotFound(w, r)
		return
	}
	body, err := io.ReadAll(io.LimitReader(r.Body, 16<<10))
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	_, verifier, err := oidcProvider(r.Context())
	if err != nil {
		http.Error(w, "identity provider unreachable", http.StatusBadGateway)
		return
	}
	idt, err := verifier.Verify(r.Context(), strings.TrimSpace(string(body)))
	var cl struct {
		Event string `json:"antwire_event"`
		Name  string `json:"name"`
		Admin bool   `json:"antwire_admin"`
	}
	if err != nil || idt.Claims(&cl) != nil || cl.Event != "profile" {
		logWarn(getIP(r), "", "HUB_PROFILE_REJECTED", fmt.Sprint(err))
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}
	var user string
	if db.QueryRow("SELECT username FROM users WHERE oidc_issuer = ? AND oidc_subject = ?", idt.Issuer, idt.Subject).Scan(&user) == nil {
		syncProfile(getIP(r), user, cl.Name, cl.Admin)
	}
	w.WriteHeader(http.StatusNoContent) // auch ohne Konto hier
}

// syncProfile übernimmt den Namen des Hubs (durchnummeriert, falls hier
// vergeben) und macht Hub-Admins hier zu Admins (nie umgekehrt). HEINEN kennt
// nur den Benutzernamen, und der ist beim Passwort-Login die Kennung: umbenannt
// werden daher nur Konten ohne Passwort (vom Antwire-Login angelegt).
// Liefert den (neuen) Benutzernamen.
func syncProfile(ip, user, name string, admin bool) string {
	var pw string
	db.QueryRow("SELECT password FROM users WHERE username = ?", user).Scan(&pw)
	if n, ok := cleanUserName(name); ok && n != user && pw == "" {
		for i := 1; i < 10; i++ {
			cand := n
			if i > 1 {
				cand = fmt.Sprintf("%s %d", string([]rune(n)[:min(len([]rune(n)), 29)]), i)
			}
			if strings.EqualFold(cand, user) {
				break
			}
			tx, err := db.Begin()
			if err != nil {
				break
			}
			if _, err = tx.Exec("UPDATE users SET username = ? WHERE username = ?", cand, user); err == nil {
				_, err = tx.Exec("UPDATE sessions SET username = ? WHERE username = ?", cand, user)
			}
			if err != nil {
				tx.Rollback()
				continue // Name hier schon vergeben
			}
			tx.Commit()
			logInfo(ip, cand, "HUB_RENAME", "from="+user)
			user = cand
			break
		}
	}
	if admin && !isAdmin(user) {
		db.Exec("UPDATE users SET is_admin = 1 WHERE username = ?", user)
		logInfo(ip, user, "OIDC_ADMIN", "promoted by antwire_admin")
	}
	return user
}
