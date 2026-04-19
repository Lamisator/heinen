package main

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strings"
	"time"
)

var difficultyPrompts = map[string]string{
	"leicht": "LEICHT – Allgemeinwissen, das die meisten kennen.",
	"mittel": "MITTELSCHWER – man muss nachdenken.",
	"schwer": "SCHWER – nur mit gutem Wissen.",
	"extrem": "EXTREM SCHWER – echtes Expertenwissen.",
}

var difficultyOrder = []string{"leicht", "mittel", "schwer", "extrem"}

func difficultyIndex(d string) int {
	for i, v := range difficultyOrder {
		if v == d {
			return i
		}
	}
	return 1
}

func escalatedDifficulty(start string, batch int) string {
	idx := difficultyIndex(start) + batch
	if idx >= len(difficultyOrder) {
		idx = len(difficultyOrder) - 1
	}
	return difficultyOrder[idx]
}

func buildPrompt(topic, diff string, count, nO int, prev []string) string {
	dt := difficultyPrompts[diff]
	if dt == "" {
		dt = difficultyPrompts["mittel"]
	}
	excl := ""
	if len(prev) > 0 {
		mx := len(prev)
		if mx > 40 {
			mx = 40
		}
		r := prev[len(prev)-mx:]
		excl = "\n\nBereits gestellte Fragen – KEINE ähnlichen generieren:\n"
		for i, q := range r {
			excl += fmt.Sprintf("%d. %s\n", i+1, q)
		}
	}
	return fmt.Sprintf(`Generiere genau %d Quiz-Fragen zum Thema "%s". %d Antwortmöglichkeiten, eine richtig. Deutsch, abwechslungsreich.
Schwierigkeit: %s
REGELN: Konkrete Fragen/Antworten, KEINE Platzhalter. Antwort NICHT trivial ableitbar. Alle falschen Antworten plausibel, ähnlich spezifisch. NICHT durch Ausschluss erratbar. Falsche dürfen sich NICHT durch Länge/Stil unterscheiden.%s
NUR JSON-Array: [{"text":"Frage?","options":["A","B","C","D"],"correct":0}]`, count, topic, nO, dt, excl)
}

func generateQuestions(topic, diff string, count, nO int, prev []string, webSearch bool) ([]Question, error) {
	prov := getSetting("ai_provider")
	model := getSetting("ai_model")
	key := ""
	if prov == "anthropic" {
		key = getSetting("anthropic_api_key")
	} else {
		key = getSetting("openai_api_key")
		prov = "openai"
	}
	if key == "" {
		logWarn("system", "system", "AI_NO_KEY", "no API key configured")
		return nil, fmt.Errorf("Kein API-Schlüssel im Admin-Panel konfiguriert. Bitte an den*die Administrator*in wenden.")
	}
	rc := count + 5
	if rc > 50 {
		rc = 50
	}
	prompt := buildPrompt(topic, diff, rc, nO, prev)
	effectiveModel := model
	if webSearch && prov == "openai" {
		effectiveModel = webSearchModelFor(model)
	}
	logInfo("system", "system", "AI_CALL", fmt.Sprintf("provider=%s model=%s diff=%s count=%d web_search=%v", prov, effectiveModel, diff, rc, webSearch))
	var text string
	var err error
	if prov == "anthropic" {
		text, err = callAnthropic(key, effectiveModel, prompt)
	} else {
		text, err = callOpenAI(key, effectiveModel, prompt, webSearch)
	}
	if err != nil {
		logError("system", "system", "AI_ERROR", err.Error())
		return nil, fmt.Errorf("KI-Anbieter antwortet nicht: %s", err.Error())
	}
	logInfo("system", "system", "AI_OK", fmt.Sprintf("len=%d", len(text)))
	rawResponse := text // keep original for error messages
	text = strings.TrimSpace(text)
	if strings.HasPrefix(text, "```") {
		ls := strings.Split(text, "\n")
		if len(ls) > 2 {
			text = strings.Join(ls[1:len(ls)-1], "\n")
		}
	}
	text = stripCitations(text)
	var qs []Question
	if err := json.Unmarshal([]byte(text), &qs); err != nil {
		logError("system", "system", "AI_PARSE_ERROR", fmt.Sprintf("err=%s response=%s", err.Error(), rawResponse))
		return nil, fmt.Errorf("KI-Antwort konnte nicht verarbeitet werden (ungültiges Format).\n\nParse-Fehler: %s\n\nAntwort der KI:\n%s", err.Error(), aiPreview(rawResponse))
	}
	fl := make([]Question, 0, len(qs))
	for _, q := range qs {
		if q.Correct < 0 || q.Correct >= len(q.Options) {
			q.Correct = 0
		}
		if !isPlaceholder(q) {
			fl = append(fl, q)
		}
	}
	qs = fl
	if len(qs) == 0 {
		logError("system", "system", "AI_NO_VALID_Q", fmt.Sprintf("response=%s", rawResponse))
		return nil, fmt.Errorf("KI-Antwort enthielt keine verwertbaren Fragen.\n\nAntwort der KI:\n%s", aiPreview(rawResponse))
	}
	// If fewer than requested, that's still OK — use what we got
	if len(qs) > count {
		qs = qs[:count]
	}
	for i := range qs {
		qs[i] = shuffleOpts(qs[i])
	}
	return qs, nil
}

// webSearchModelFor maps a configured OpenAI model to the corresponding search-capable model.
// When web search is enabled via Chat Completions, only specific model variants support it.
func webSearchModelFor(model string) string {
	m := strings.ToLower(model)
	switch {
	case strings.Contains(m, "mini") || strings.Contains(m, "nano"):
		return "gpt-4o-mini-search-preview"
	case strings.HasPrefix(m, "gpt-5"):
		return "gpt-5-search-api"
	default:
		return "gpt-4o-search-preview"
	}
}

// stripCitations removes inline citation markers that web-search models inject,
// which would otherwise break JSON parsing.
func stripCitations(s string) string {
	// Unicode reference brackets 【...】
	for {
		start := strings.Index(s, "\u3010")
		if start < 0 {
			break
		}
		end := strings.Index(s[start:], "\u3011")
		if end < 0 {
			break
		}
		s = s[:start] + s[start+end+len("\u3011"):]
	}
	return s
}

func shuffleOpts(q Question) Question {
	n := len(q.Options)
	if n <= 1 {
		return q
	}
	idx := make([]int, n)
	for i := range idx {
		idx[i] = i
	}
	for i := n - 1; i > 0; i-- {
		jB, _ := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		j := int(jB.Int64())
		idx[i], idx[j] = idx[j], idx[i]
	}
	no := make([]string, n)
	nc := 0
	for ni, oi := range idx {
		no[ni] = q.Options[oi]
		if oi == q.Correct {
			nc = ni
		}
	}
	q.Options = no
	q.Correct = nc
	return q
}

func callOpenAI(key, model, prompt string, webSearch bool) (string, error) {
	payload := map[string]interface{}{
		"model":                 model,
		"max_completion_tokens": 4096,
		"messages":              []map[string]interface{}{{"role": "user", "content": prompt}},
	}
	if webSearch {
		// web_search_options is required by search-capable models; empty object enables default behavior
		payload["web_search_options"] = map[string]interface{}{}
	}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+key)
	resp, err := (&http.Client{Timeout: 120 * time.Second}).Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	rb, _ := io.ReadAll(resp.Body)
	var r struct {
		Choices []struct{ Message struct{ Content string } }
		Error   *struct{ Message string }
	}
	json.Unmarshal(rb, &r)
	if r.Error != nil {
		return "", fmt.Errorf("%s", r.Error.Message)
	}
	if len(r.Choices) == 0 {
		return "", fmt.Errorf("no choices")
	}
	return r.Choices[0].Message.Content, nil
}

func callAnthropic(key, model, prompt string) (string, error) {
	body, _ := json.Marshal(map[string]interface{}{"model": model, "max_tokens": 4096, "messages": []map[string]interface{}{{"role": "user", "content": prompt}}})
	req, _ := http.NewRequest("POST", "https://api.anthropic.com/v1/messages", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", key)
	req.Header.Set("anthropic-version", "2023-06-01")
	resp, err := (&http.Client{Timeout: 90 * time.Second}).Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	rb, _ := io.ReadAll(resp.Body)
	var r struct {
		Content []struct{ Type, Text string }
		Error   *struct{ Message string }
	}
	json.Unmarshal(rb, &r)
	if r.Error != nil {
		return "", fmt.Errorf("%s", r.Error.Message)
	}
	for _, c := range r.Content {
		if c.Type == "text" {
			return c.Text, nil
		}
	}
	return "", fmt.Errorf("no text")
}

func testAIKey(prov, key, model string) (bool, string) {
	var err error
	if prov == "anthropic" {
		_, err = callAnthropic(key, model, "Sage nur: OK")
	} else {
		_, err = callOpenAI(key, model, "Sage nur: OK", false)
	}
	if err != nil {
		return false, err.Error()
	}
	return true, "Verbindung erfolgreich!"
}

func isPlaceholder(q Question) bool {
	phs := []string{"erste antwort", "zweite antwort", "dritte antwort", "vierte antwort", "antwort 1", "antwort 2", "antwort 3", "antwort 4", "option a", "option b", "option c", "option d"}
	for _, o := range q.Options {
		l := strings.ToLower(strings.TrimSpace(o))
		for _, ph := range phs {
			if l == ph {
				return true
			}
		}
	}
	return false
}

// aiPreview returns a truncated version of an AI response for user-facing error messages.
// The full response is always logged separately.
func aiPreview(s string) string {
	if len(s) > 800 {
		return s[:800] + "\n… (gekürzt, vollständige Antwort im Log)"
	}
	return s
}
