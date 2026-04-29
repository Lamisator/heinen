package main

import (
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/sha3"
)

var questionsDB *sql.DB

func initQuestionsDB() {
	var err error
	questionsDB, err = sql.Open("sqlite3", "questions.db?_journal_mode=WAL")
	if err != nil {
		logError("system", "system", "DB_ERROR", "Failed to open questions.db: "+err.Error())
		return
	}
	questionsDB.Exec(`CREATE TABLE IF NOT EXISTS questions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		text TEXT NOT NULL,
		options TEXT NOT NULL,
		correct_answer INTEGER NOT NULL,
		difficulty TEXT NOT NULL DEFAULT '',
		topic TEXT NOT NULL DEFAULT '',
		ai_provider TEXT NOT NULL DEFAULT '',
		ai_model TEXT NOT NULL DEFAULT '',
		num_options INTEGER NOT NULL DEFAULT 0,
		asked_at INTEGER NOT NULL
	)`)
	// Migrations for existing installs
	questionsDB.Exec("ALTER TABLE questions ADD COLUMN num_options INTEGER NOT NULL DEFAULT 0")
	questionsDB.Exec("ALTER TABLE questions ADD COLUMN duplicate_ok INTEGER NOT NULL DEFAULT 0")
	questionsDB.Exec("ALTER TABLE questions ADD COLUMN report_count INTEGER NOT NULL DEFAULT 0")
	questionsDB.Exec(`CREATE INDEX IF NOT EXISTS idx_questions_asked ON questions(asked_at DESC)`)
	questionsDB.Exec(`CREATE INDEX IF NOT EXISTS idx_questions_topic ON questions(topic)`)
	questionsDB.Exec(`CREATE INDEX IF NOT EXISTS idx_questions_difficulty ON questions(difficulty)`)
	questionsDB.Exec(`CREATE INDEX IF NOT EXISTS idx_questions_numopts ON questions(num_options)`)
	questionsDB.Exec(`CREATE INDEX IF NOT EXISTS idx_questions_dupok ON questions(duplicate_ok)`)
	questionsDB.Exec(`CREATE INDEX IF NOT EXISTS idx_questions_reported ON questions(report_count)`)
	questionsDB.Exec("ALTER TABLE questions ADD COLUMN hash TEXT NOT NULL DEFAULT ''")
	questionsDB.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS idx_questions_hash ON questions(hash) WHERE hash != ''`)
	questionsDB.Exec(`CREATE TABLE IF NOT EXISTS duplicate_reviewed_pairs (
		id_a INTEGER NOT NULL,
		id_b INTEGER NOT NULL,
		PRIMARY KEY(id_a, id_b)
	)`)
}

// computeQuestionHash returns a SHA3-256 hash of the question's canonical content
// (normalized text + sorted options + correct option text). This hash is stable
// regardless of option ordering, making it suitable for deduplication and sync.
func computeQuestionHash(text string, options []string, correct int) string {
	correctText := ""
	if correct >= 0 && correct < len(options) {
		correctText = options[correct]
	}
	sorted := make([]string, len(options))
	copy(sorted, options)
	sort.Strings(sorted)
	var sb strings.Builder
	sb.WriteString(strings.TrimSpace(text))
	sb.WriteByte(0)
	sb.WriteString(correctText)
	sb.WriteByte(0)
	sb.WriteString(strings.Join(sorted, "\x00"))
	h := sha3.Sum256([]byte(sb.String()))
	return hex.EncodeToString(h[:])
}

func saveQuestionsBulk(qs []Question, difficulty, topic, aiProvider, aiModel string, numOptions int) {
	if questionsDB == nil || len(qs) == 0 {
		return
	}
	tx, err := questionsDB.Begin()
	if err != nil {
		logError("system", "system", "DB_ERROR", "saveQuestionsBulk begin: "+err.Error())
		return
	}
	stmt, err := tx.Prepare(`INSERT OR IGNORE INTO questions (text, options, correct_answer, difficulty, topic, ai_provider, ai_model, num_options, asked_at, hash) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		tx.Rollback()
		logError("system", "system", "DB_ERROR", "saveQuestionsBulk prepare: "+err.Error())
		return
	}
	defer stmt.Close()
	now := time.Now().Unix()
	var inserted int64
	for _, q := range qs {
		optJSON, _ := json.Marshal(q.Options)
		nO := len(q.Options)
		if nO == 0 {
			nO = numOptions
		}
		h := computeQuestionHash(q.Text, q.Options, q.Correct)
		res, err := stmt.Exec(q.Text, string(optJSON), q.Correct, difficulty, topic, aiProvider, aiModel, nO, now, h)
		if err != nil {
			logError("system", "system", "DB_ERROR", "saveQuestionsBulk exec: "+err.Error())
			continue
		}
		if n, _ := res.RowsAffected(); n > 0 {
			inserted++
		}
	}
	tx.Commit()
	if ignored := int64(len(qs)) - inserted; ignored > 0 {
		logInfo("system", "system", "DEDUP_HASH", fmt.Sprintf("topic=%s diff=%s exact_dupes_ignored=%d", topic, difficulty, ignored))
	}
}

func updateQuestion(id int64, text string, options []string, correctAnswer, numOptions int, difficulty, topic, aiProvider, aiModel string) error {
	if questionsDB == nil {
		return nil
	}
	optJSON, _ := json.Marshal(options)
	nO := numOptions
	if nO == 0 {
		nO = len(options)
	}
	h := computeQuestionHash(text, options, correctAnswer)
	_, err := questionsDB.Exec(
		`UPDATE questions SET text=?, options=?, correct_answer=?, difficulty=?, topic=?, ai_provider=?, ai_model=?, num_options=?, hash=? WHERE id=?`,
		text, string(optJSON), correctAnswer, difficulty, topic, aiProvider, aiModel, nO, h, id,
	)
	return err
}

func insertManualQuestion(text string, options []string, correctAnswer, numOptions int, difficulty, topic string) error {
	if questionsDB == nil {
		return nil
	}
	optJSON, _ := json.Marshal(options)
	nO := numOptions
	if nO == 0 {
		nO = len(options)
	}
	h := computeQuestionHash(text, options, correctAnswer)
	_, err := questionsDB.Exec(
		`INSERT OR IGNORE INTO questions (text, options, correct_answer, difficulty, topic, ai_provider, ai_model, num_options, asked_at, hash) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		text, string(optJSON), correctAnswer, difficulty, topic, "Manuell", "", nO, time.Now().Unix(), h,
	)
	return err
}

type duplicatePair struct {
	A     questionRecord `json:"a"`
	B     questionRecord `json:"b"`
	Score float64        `json:"score"`
}

type questionRecord struct {
	ID            int64    `json:"id"`
	Text          string   `json:"text"`
	Options       []string `json:"options"`
	CorrectAnswer int      `json:"correctAnswer"`
	Difficulty    string   `json:"difficulty"`
	Topic         string   `json:"topic"`
	AIProvider    string   `json:"aiProvider"`
	AIModel       string   `json:"aiModel"`
	NumOptions    int      `json:"numOptions"`
	AskedAt       int64    `json:"askedAt"`
	DuplicateOk   bool     `json:"duplicateOk"`
	ReportCount   int      `json:"reportCount"`
	Hash          string   `json:"hash"`
}

func queryQuestions(search, topic, difficulty, aiProvider string, numOptions, limit, offset int) ([]questionRecord, int, error) {
	if questionsDB == nil {
		return nil, 0, nil
	}

	where := "WHERE 1=1"
	var filterArgs []interface{}

	if search != "" {
		idStr := strings.TrimPrefix(strings.TrimSpace(search), "#")
		if sid, err := strconv.ParseInt(idStr, 10, 64); err == nil && sid > 0 {
			where += " AND id = ?"
			filterArgs = append(filterArgs, sid)
		} else {
			where += " AND text LIKE ?"
			filterArgs = append(filterArgs, "%"+search+"%")
		}
	}
	if topic != "" {
		where += " AND topic LIKE ?"
		filterArgs = append(filterArgs, "%"+topic+"%")
	}
	if difficulty != "" {
		where += " AND difficulty = ?"
		filterArgs = append(filterArgs, difficulty)
	}
	if aiProvider != "" {
		where += " AND ai_provider = ?"
		filterArgs = append(filterArgs, aiProvider)
	}
	if numOptions > 0 {
		where += " AND num_options = ?"
		filterArgs = append(filterArgs, numOptions)
	}

	var total int
	countRow := questionsDB.QueryRow("SELECT COUNT(*) FROM questions "+where, filterArgs...)
	if err := countRow.Scan(&total); err != nil {
		return nil, 0, err
	}

	pageArgs := append(filterArgs, limit, offset)
	rows, err := questionsDB.Query(
		"SELECT id, text, options, correct_answer, difficulty, topic, ai_provider, ai_model, num_options, asked_at, duplicate_ok, report_count, hash FROM questions "+where+" ORDER BY asked_at DESC LIMIT ? OFFSET ?",
		pageArgs...,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var records []questionRecord
	for rows.Next() {
		var rec questionRecord
		var optJSON string
		if err := rows.Scan(&rec.ID, &rec.Text, &optJSON, &rec.CorrectAnswer, &rec.Difficulty, &rec.Topic, &rec.AIProvider, &rec.AIModel, &rec.NumOptions, &rec.AskedAt, &rec.DuplicateOk, &rec.ReportCount, &rec.Hash); err != nil {
			continue
		}
		json.Unmarshal([]byte(optJSON), &rec.Options)
		records = append(records, rec)
	}
	return records, total, nil
}

func queryReportedQuestions() ([]questionRecord, error) {
	if questionsDB == nil {
		return nil, nil
	}
	rows, err := questionsDB.Query(
		`SELECT id, text, options, correct_answer, difficulty, topic, ai_provider, ai_model, num_options, asked_at, duplicate_ok, report_count, hash
		 FROM questions WHERE report_count > 0 ORDER BY report_count DESC, asked_at DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var records []questionRecord
	for rows.Next() {
		var rec questionRecord
		var optJSON string
		if err := rows.Scan(&rec.ID, &rec.Text, &optJSON, &rec.CorrectAnswer, &rec.Difficulty, &rec.Topic, &rec.AIProvider, &rec.AIModel, &rec.NumOptions, &rec.AskedAt, &rec.DuplicateOk, &rec.ReportCount, &rec.Hash); err != nil {
			continue
		}
		json.Unmarshal([]byte(optJSON), &rec.Options)
		records = append(records, rec)
	}
	return records, nil
}

func setDuplicateOk(id int64, ok bool) error {
	if questionsDB == nil {
		return nil
	}
	v := 0
	if ok {
		v = 1
	}
	_, err := questionsDB.Exec("UPDATE questions SET duplicate_ok = ? WHERE id = ?", v, id)
	return err
}

func reportQuestion(text string) error {
	if questionsDB == nil {
		return nil
	}
	_, err := questionsDB.Exec("UPDATE questions SET report_count = report_count + 1 WHERE text = ?", text)
	return err
}

func clearQuestionReport(id int64) error {
	if questionsDB == nil {
		return nil
	}
	_, err := questionsDB.Exec("UPDATE questions SET report_count = 0 WHERE id = ?", id)
	return err
}

func markPairReviewed(idA, idB int64) error {
	if questionsDB == nil {
		return nil
	}
	if idA > idB {
		idA, idB = idB, idA
	}
	_, err := questionsDB.Exec(
		"INSERT OR IGNORE INTO duplicate_reviewed_pairs(id_a, id_b) VALUES(?, ?)",
		idA, idB,
	)
	return err
}

type topicStat struct {
	Topic  string         `json:"topic"`
	Count  int            `json:"count"`
	ByDiff map[string]int `json:"byDiff,omitempty"`
}

func getTopicsWithStats() ([]topicStat, error) {
	if questionsDB == nil {
		return nil, nil
	}
	rows, err := questionsDB.Query(`
		SELECT topic, COUNT(*) as c,
		       SUM(CASE WHEN difficulty='leicht' THEN 1 ELSE 0 END),
		       SUM(CASE WHEN difficulty='mittel' THEN 1 ELSE 0 END),
		       SUM(CASE WHEN difficulty='schwer' THEN 1 ELSE 0 END),
		       SUM(CASE WHEN difficulty='extrem' THEN 1 ELSE 0 END)
		FROM questions WHERE topic != ''
		GROUP BY topic ORDER BY c DESC LIMIT 300`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []topicStat
	for rows.Next() {
		var t topicStat
		var l, m, s, e int
		if err := rows.Scan(&t.Topic, &t.Count, &l, &m, &s, &e); err != nil {
			continue
		}
		t.ByDiff = map[string]int{}
		if l > 0 {
			t.ByDiff["leicht"] = l
		}
		if m > 0 {
			t.ByDiff["mittel"] = m
		}
		if s > 0 {
			t.ByDiff["schwer"] = s
		}
		if e > 0 {
			t.ByDiff["extrem"] = e
		}
		result = append(result, t)
	}
	return result, nil
}

func getQuestionStats() (total int, byDifficulty map[string]int, err error) {
	byDifficulty = make(map[string]int)
	if questionsDB == nil {
		return
	}
	questionsDB.QueryRow("SELECT COUNT(*) FROM questions").Scan(&total)
	rows, e := questionsDB.Query("SELECT difficulty, COUNT(*) FROM questions GROUP BY difficulty")
	if e != nil {
		err = e
		return
	}
	defer rows.Close()
	for rows.Next() {
		var d string
		var c int
		if rows.Scan(&d, &c) == nil {
			byDifficulty[d] = c
		}
	}
	return
}

func deleteQuestionDB(id int64) error {
	if questionsDB == nil {
		return nil
	}
	_, err := questionsDB.Exec("DELETE FROM questions WHERE id = ?", id)
	return err
}

func deleteAllQuestions(filter struct{ topic, difficulty, aiProvider string }) (int64, error) {
	if questionsDB == nil {
		return 0, nil
	}
	where := "WHERE 1=1"
	var args []interface{}
	if filter.topic != "" {
		where += " AND topic LIKE ?"
		args = append(args, "%"+filter.topic+"%")
	}
	if filter.difficulty != "" {
		where += " AND difficulty = ?"
		args = append(args, filter.difficulty)
	}
	if filter.aiProvider != "" {
		where += " AND ai_provider = ?"
		args = append(args, filter.aiProvider)
	}
	res, err := questionsDB.Exec("DELETE FROM questions "+where, args...)
	if err != nil {
		return 0, err
	}
	n, _ := res.RowsAffected()
	return n, nil
}

func normalizeQText(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	var b strings.Builder
	prevSpace := false
	for _, r := range s {
		if unicode.IsSpace(r) {
			if !prevSpace {
				b.WriteRune(' ')
			}
			prevSpace = true
		} else {
			b.WriteRune(r)
			prevSpace = false
		}
	}
	return b.String()
}

func levenshtein(a, b string) int {
	ra, rb := []rune(a), []rune(b)
	la, lb := len(ra), len(rb)
	if la == 0 {
		return lb
	}
	if lb == 0 {
		return la
	}
	if la > lb {
		ra, rb = rb, ra
		la, lb = lb, la
	}
	prev := make([]int, la+1)
	curr := make([]int, la+1)
	for i := range prev {
		prev[i] = i
	}
	for j := 1; j <= lb; j++ {
		curr[0] = j
		for i := 1; i <= la; i++ {
			cost := 1
			if ra[i-1] == rb[j-1] {
				cost = 0
			}
			del := prev[i] + 1
			ins := curr[i-1] + 1
			sub := prev[i-1] + cost
			curr[i] = del
			if ins < curr[i] {
				curr[i] = ins
			}
			if sub < curr[i] {
				curr[i] = sub
			}
		}
		prev, curr = curr, prev
	}
	return prev[la]
}

// deStopwords are common German words that carry little semantic meaning.
var deStopwords = map[string]bool{
	"der": true, "die": true, "das": true, "des": true, "dem": true, "den": true,
	"ein": true, "eine": true, "einer": true, "einem": true, "einen": true, "eines": true,
	"und": true, "oder": true, "aber": true, "auch": true, "noch": true,
	"ist": true, "sind": true, "war": true, "waren": true, "wird": true, "werden": true,
	"hat": true, "haben": true, "hatte": true, "hatten": true, "sein": true,
	"von": true, "zu": true, "in": true, "im": true, "an": true, "am": true,
	"auf": true, "bei": true, "mit": true, "nach": true, "aus": true, "für": true,
	"als": true, "wie": true, "was": true, "wer": true, "wo": true, "wann": true,
	"welche": true, "welcher": true, "welches": true, "welchem": true, "welchen": true,
	"nicht": true, "kein": true, "keine": true, "keiner": true,
	"es": true, "er": true, "sie": true, "wir": true, "ihr": true,
	"sich": true, "man": true, "durch": true, "über": true, "unter": true,
	"zwischen": true, "beim": true, "zum": true, "zur": true, "vom": true,
	"diesem": true, "dieser": true, "dieses": true,
	"ihre": true, "ihrer": true,
}

// tokenizeText splits a normalized text into meaningful word tokens,
// removing stopwords and very short tokens.
func tokenizeText(s string) map[string]struct{} {
	tokens := make(map[string]struct{})
	for _, word := range strings.Fields(s) {
		word = strings.Trim(word, ".,!?;:\"'()[]")
		if len([]rune(word)) < 3 {
			continue
		}
		if deStopwords[word] {
			continue
		}
		tokens[word] = struct{}{}
	}
	return tokens
}

// jaccardSim returns the Jaccard similarity between two token sets: |A∩B| / |A∪B|.
func jaccardSim(a, b map[string]struct{}) float64 {
	if len(a) == 0 && len(b) == 0 {
		return 1.0
	}
	intersection := 0
	for w := range a {
		if _, ok := b[w]; ok {
			intersection++
		}
	}
	union := len(a) + len(b) - intersection
	if union == 0 {
		return 0
	}
	return float64(intersection) / float64(union)
}

// levSimilarity returns a [0,1] similarity score via Levenshtein edit distance.
func levSimilarity(a, b string) float64 {
	la, lb := len([]rune(a)), len([]rune(b))
	maxLen := la
	if lb > maxLen {
		maxLen = lb
	}
	if maxLen == 0 {
		return 1.0
	}
	return 1.0 - float64(levenshtein(a, b))/float64(maxLen)
}

// strSimilarity combines character-level Levenshtein and word-level Jaccard similarity.
// It returns the higher of the two scores so that both near-identical phrasing
// (caught by Levenshtein) and same-fact paraphrasing (caught by Jaccard) are detected.
func strSimilarity(normA, normB string, tokA, tokB map[string]struct{}) float64 {
	jaccard := jaccardSim(tokA, tokB)
	lev := levSimilarity(normA, normB)
	if jaccard > lev {
		return jaccard
	}
	return lev
}

type preparedQuestion struct {
	rec  questionRecord
	norm string
	toks map[string]struct{}
}

func findDuplicatePairs(threshold float64) ([]duplicatePair, error) {
	if questionsDB == nil {
		return []duplicatePair{}, nil
	}
	// Load already-reviewed pairs to skip them.
	reviewed := make(map[[2]int64]struct{})
	if pr, e := questionsDB.Query("SELECT id_a, id_b FROM duplicate_reviewed_pairs"); e == nil {
		defer pr.Close()
		for pr.Next() {
			var a, b int64
			if pr.Scan(&a, &b) == nil {
				reviewed[[2]int64{a, b}] = struct{}{}
			}
		}
	}

	rows, err := questionsDB.Query(
		`SELECT id, text, options, correct_answer, difficulty, topic,
		        ai_provider, ai_model, num_options, asked_at, duplicate_ok, report_count
		 FROM questions ORDER BY asked_at DESC LIMIT 3000`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var all []preparedQuestion
	for rows.Next() {
		var rec questionRecord
		var optJSON string
		if err := rows.Scan(&rec.ID, &rec.Text, &optJSON,
			&rec.CorrectAnswer, &rec.Difficulty, &rec.Topic,
			&rec.AIProvider, &rec.AIModel, &rec.NumOptions, &rec.AskedAt,
			&rec.DuplicateOk, &rec.ReportCount,
		); err != nil {
			continue
		}
		json.Unmarshal([]byte(optJSON), &rec.Options)
		norm := normalizeQText(rec.Text)
		all = append(all, preparedQuestion{rec: rec, norm: norm, toks: tokenizeText(norm)})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	seen := make(map[[2]int64]struct{})
	var pairs []duplicatePair

	// jaccardPreFilter is set below the Levenshtein threshold so we don't miss
	// pairs where Jaccard is moderate but Levenshtein would push score over threshold.
	jaccardPreFilter := threshold - 0.2
	if jaccardPreFilter < 0 {
		jaccardPreFilter = 0
	}

	for i := 0; i < len(all); i++ {
		for j := i + 1; j < len(all); j++ {
			if len(pairs) >= 200 {
				goto done
			}

			// Fast Jaccard pre-filter: skip pairs with clearly insufficient word overlap.
			jaccard := jaccardSim(all[i].toks, all[j].toks)
			if jaccard < jaccardPreFilter {
				// Also skip Levenshtein if string lengths differ too much.
				la, lb := len([]rune(all[i].norm)), len([]rune(all[j].norm))
				maxLen := la
				if lb > maxLen {
					maxLen = lb
				}
				absDiff := la - lb
				if absDiff < 0 {
					absDiff = -absDiff
				}
				if maxLen == 0 || float64(absDiff)/float64(maxLen) > (1-threshold) {
					continue
				}
			}

			score := strSimilarity(all[i].norm, all[j].norm, all[i].toks, all[j].toks)
			if score < threshold {
				continue
			}

			idA, idB := all[i].rec.ID, all[j].rec.ID
			if idA > idB {
				idA, idB = idB, idA
			}
			key := [2]int64{idA, idB}
			if _, exists := seen[key]; exists {
				continue
			}
			if _, isReviewed := reviewed[key]; isReviewed {
				continue
			}
			seen[key] = struct{}{}
			pairs = append(pairs, duplicatePair{
				A:     all[i].rec,
				B:     all[j].rec,
				Score: score,
			})
		}
	}
done:
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].Score > pairs[j].Score
	})
	if pairs == nil {
		pairs = []duplicatePair{}
	}
	return pairs, nil
}

// loadPreparedForTopicDiff loads up to maxLoad existing questions for the given
// topic+difficulty as preparedQuestion values ready for similarity comparison.
func loadPreparedForTopicDiff(topic, difficulty string, maxLoad int) ([]preparedQuestion, error) {
	if questionsDB == nil {
		return nil, nil
	}
	if maxLoad <= 0 || maxLoad > 500 {
		maxLoad = 500
	}
	rows, err := questionsDB.Query(
		`SELECT id, text, options, correct_answer, difficulty, topic, ai_provider, ai_model, num_options, asked_at, duplicate_ok, report_count
		 FROM questions WHERE topic = ? AND difficulty = ? ORDER BY asked_at DESC LIMIT ?`,
		topic, difficulty, maxLoad,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []preparedQuestion
	for rows.Next() {
		var rec questionRecord
		var optJSON string
		if err := rows.Scan(&rec.ID, &rec.Text, &optJSON,
			&rec.CorrectAnswer, &rec.Difficulty, &rec.Topic,
			&rec.AIProvider, &rec.AIModel, &rec.NumOptions, &rec.AskedAt,
			&rec.DuplicateOk, &rec.ReportCount,
		); err != nil {
			continue
		}
		json.Unmarshal([]byte(optJSON), &rec.Options)
		norm := normalizeQText(rec.Text)
		result = append(result, preparedQuestion{rec: rec, norm: norm, toks: tokenizeText(norm)})
	}
	return result, nil
}

// filterNewQuestions removes questions from newQs that are too similar to any
// question in existing. Returns the filtered slice and the number removed.
// Uses the same two-stage Jaccard pre-filter + strSimilarity logic as findDuplicatePairs.
func filterNewQuestions(newQs []Question, existing []preparedQuestion, threshold float64) ([]Question, int) {
	if len(existing) == 0 || threshold <= 0 {
		return newQs, 0
	}
	jaccardPre := threshold - 0.15
	if jaccardPre < 0 {
		jaccardPre = 0
	}
	var kept []Question
	removed := 0
	for _, q := range newQs {
		norm := normalizeQText(q.Text)
		toks := tokenizeText(norm)
		duplicate := false
		for _, ex := range existing {
			jaccard := jaccardSim(toks, ex.toks)
			if jaccard < jaccardPre {
				la, lb := len([]rune(norm)), len([]rune(ex.norm))
				maxLen := la
				if lb > maxLen {
					maxLen = lb
				}
				diff := la - lb
				if diff < 0 {
					diff = -diff
				}
				if maxLen == 0 || float64(diff)/float64(maxLen) > (1-threshold) {
					continue
				}
			}
			if strSimilarity(norm, ex.norm, toks, ex.toks) >= threshold {
				duplicate = true
				break
			}
		}
		if duplicate {
			removed++
		} else {
			kept = append(kept, q)
		}
	}
	if kept == nil {
		kept = []Question{}
	}
	return kept, removed
}

// recentQTextsForTopicDiff returns the texts of the most recent `max` questions
// for the given topic+difficulty, used to populate the AI prompt exclusion list.
func recentQTextsForTopicDiff(topic, difficulty string, max int) ([]string, error) {
	if questionsDB == nil {
		return nil, nil
	}
	rows, err := questionsDB.Query(
		`SELECT text FROM questions WHERE topic = ? AND difficulty = ? ORDER BY asked_at DESC LIMIT ?`,
		topic, difficulty, max,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var texts []string
	for rows.Next() {
		var t string
		if rows.Scan(&t) == nil {
			texts = append(texts, t)
		}
	}
	return texts, nil
}
