package db

import (
	"database/sql"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// InitDB initializes the SQLite database connection and schema
func InitDB() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	dir := filepath.Join(home, ".dojo")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	dbPath := filepath.Join(dir, "dojo.db")

	var dbErr error
	db, dbErr = sql.Open("sqlite3", dbPath)
	if dbErr != nil {
		return dbErr
	}

	schema := `
	CREATE TABLE IF NOT EXISTS player_stats (
		id INTEGER PRIMARY KEY,
		total_xp INTEGER DEFAULT 0,
		incidents_resolved INTEGER DEFAULT 0,
		hints_used INTEGER DEFAULT 0
	);

	CREATE TABLE IF NOT EXISTS sessions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		incident_id TEXT,
		difficulty TEXT,
		mode TEXT,
		start_time DATETIME,
		end_time DATETIME,
		status TEXT,
		hints_used INTEGER,
		time_taken_seconds INTEGER,
		xp_gained INTEGER
	);
	`
	// Simple migration for new 'mode' column (ignore error if column already exists)
	db.Exec("ALTER TABLE sessions ADD COLUMN mode TEXT DEFAULT 'normal'")
	
	_, err = db.Exec(schema)
	if err != nil {
		return err
	}

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM player_stats").Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		_, err = db.Exec("INSERT INTO player_stats (id, total_xp, incidents_resolved, hints_used) VALUES (1, 0, 0, 0)")
		if err != nil {
			return err
		}
	}

	return nil
}

// SessionRecord represents a completed or failed session
type SessionRecord struct {
	IncidentID       string
	Difficulty       string
	Mode             string
	StartTime        time.Time
	EndTime          time.Time
	Status           string
	HintsUsed        int
	TimeTakenSeconds int
	XPGained         int
}

// SaveSession saves a session and updates the player stats if resolved
func SaveSession(s SessionRecord) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
		INSERT INTO sessions (incident_id, difficulty, mode, start_time, end_time, status, hints_used, time_taken_seconds, xp_gained)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, s.IncidentID, s.Difficulty, s.Mode, s.StartTime, s.EndTime, s.Status, s.HintsUsed, s.TimeTakenSeconds, s.XPGained)
	if err != nil {
		return err
	}

	if s.Status == "resolved" {
		_, err = tx.Exec(`
			UPDATE player_stats
			SET total_xp = total_xp + ?,
			    incidents_resolved = incidents_resolved + 1,
			    hints_used = hints_used + ?
			WHERE id = 1
		`, s.XPGained, s.HintsUsed)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// PlayerStats represents the overall stats
type PlayerStats struct {
	TotalXP           int
	IncidentsResolved int
	HintsUsed         int
}

// GetPlayerStats retrieves the current stats
func GetPlayerStats() (*PlayerStats, error) {
	var stats PlayerStats
	err := db.QueryRow("SELECT total_xp, incidents_resolved, hints_used FROM player_stats WHERE id = 1").
		Scan(&stats.TotalXP, &stats.IncidentsResolved, &stats.HintsUsed)
	if err != nil {
		return nil, err
	}
	return &stats, nil
}

// GetRank returns a rank based on XP
func GetRank(xp int) string {
	if xp < 100 {
		return "Novice"
	}
	if xp < 500 {
		return "Practitioner"
	}
	if xp < 1500 {
		return "Sensei"
	}
	return "Chaos Master"
}
