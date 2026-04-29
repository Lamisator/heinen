package main

import (
	"encoding/json"
	"time"
)

type playerRoundResult struct {
	PlayerID    string
	PlayerName  string
	Answer      int
	Result      string // "correct", "wrong", "timeout"
	TeethBefore int
	TeethAfter  int
	Eliminated  bool
}

func initGameHistoryTables() {
	db.Exec(`CREATE TABLE IF NOT EXISTS game_history (
		id TEXT PRIMARY KEY,
		host_user TEXT NOT NULL,
		lobby_name TEXT NOT NULL,
		mode TEXT NOT NULL,
		topic TEXT NOT NULL,
		difficulty TEXT NOT NULL,
		start_difficulty TEXT NOT NULL DEFAULT '',
		num_teeth INTEGER NOT NULL,
		time_per_q INTEGER NOT NULL,
		num_options INTEGER NOT NULL,
		num_players INTEGER NOT NULL,
		started_at INTEGER NOT NULL,
		ended_at INTEGER,
		end_phase TEXT,
		winners TEXT
	)`)
	db.Exec(`CREATE TABLE IF NOT EXISTS game_rounds (
		game_id TEXT NOT NULL,
		round_number INTEGER NOT NULL,
		question_text TEXT NOT NULL,
		options TEXT NOT NULL,
		correct_answer INTEGER NOT NULL,
		difficulty TEXT NOT NULL DEFAULT '',
		PRIMARY KEY (game_id, round_number)
	)`)
	db.Exec(`CREATE TABLE IF NOT EXISTS game_player_rounds (
		game_id TEXT NOT NULL,
		round_number INTEGER NOT NULL,
		player_id TEXT NOT NULL,
		player_name TEXT NOT NULL,
		answer INTEGER NOT NULL DEFAULT -1,
		result TEXT NOT NULL,
		teeth_before INTEGER NOT NULL,
		teeth_after INTEGER NOT NULL,
		eliminated INTEGER NOT NULL DEFAULT 0,
		PRIMARY KEY (game_id, round_number, player_id)
	)`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_game_history_started ON game_history(started_at DESC)`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_game_rounds_game ON game_rounds(game_id)`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_gpr_game ON game_player_rounds(game_id)`)
}

func recordGameStart(histID, hostUser, lobbyName, mode, topic, difficulty, startDiff string, numTeeth, timePerQ, numOptions, numPlayers int) {
	db.Exec(`INSERT OR IGNORE INTO game_history
		(id, host_user, lobby_name, mode, topic, difficulty, start_difficulty, num_teeth, time_per_q, num_options, num_players, started_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		histID, hostUser, lobbyName, mode, topic, difficulty, startDiff,
		numTeeth, timePerQ, numOptions, numPlayers, time.Now().Unix())
}

func recordGameRound(histID string, roundNum int, q Question, difficulty string, players []playerRoundResult) {
	optJSON, _ := json.Marshal(q.Options)
	db.Exec(`INSERT OR IGNORE INTO game_rounds (game_id, round_number, question_text, options, correct_answer, difficulty) VALUES (?, ?, ?, ?, ?, ?)`,
		histID, roundNum, q.Text, string(optJSON), q.Correct, difficulty)
	for _, pr := range players {
		elim := 0
		if pr.Eliminated {
			elim = 1
		}
		db.Exec(`INSERT OR IGNORE INTO game_player_rounds
			(game_id, round_number, player_id, player_name, answer, result, teeth_before, teeth_after, eliminated)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			histID, roundNum, pr.PlayerID, pr.PlayerName, pr.Answer, pr.Result,
			pr.TeethBefore, pr.TeethAfter, elim)
	}
}

func recordGameEnd(histID, endPhase string, winners []map[string]string) {
	if histID == "" {
		return
	}
	var winnersJSON string
	if len(winners) > 0 {
		names := make([]string, 0, len(winners))
		for _, w := range winners {
			names = append(names, w["name"])
		}
		b, _ := json.Marshal(names)
		winnersJSON = string(b)
	}
	db.Exec(`UPDATE game_history SET ended_at = ?, end_phase = ?, winners = ? WHERE id = ?`,
		time.Now().Unix(), endPhase, winnersJSON, histID)
}
