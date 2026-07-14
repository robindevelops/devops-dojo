package session

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Leaderboard struct {
	TotalPoints      int `json:"total_points"`
	IncidentsResolved int `json:"incidents_resolved"`
	HintsUsed        int `json:"hints_used"`
}

func getLeaderboardPath() string {
	home, _ := os.UserHomeDir()
	dir := filepath.Join(home, ".dojo")
	os.MkdirAll(dir, 0755)
	return filepath.Join(dir, "leaderboard.json")
}

func LoadLeaderboard() (*Leaderboard, error) {
	data, err := os.ReadFile(getLeaderboardPath())
	if err != nil {
		if os.IsNotExist(err) {
			return &Leaderboard{}, nil
		}
		return nil, err
	}
	var l Leaderboard
	json.Unmarshal(data, &l)
	return &l, nil
}

func SaveLeaderboard(l *Leaderboard) error {
	data, _ := json.MarshalIndent(l, "", "  ")
	return os.WriteFile(getLeaderboardPath(), data, 0644)
}

func AddScore(points int, hints int) (*Leaderboard, error) {
	l, err := LoadLeaderboard()
	if err != nil {
		return nil, err
	}
	l.TotalPoints += points
	l.IncidentsResolved++
	l.HintsUsed += hints
	err = SaveLeaderboard(l)
	return l, err
}

func GetRank(points int) string {
	if points < 100 {
		return "Novice"
	}
	if points < 500 {
		return "Practitioner"
	}
	if points < 1000 {
		return "Sensei"
	}
	return "Chaos Master"
}
