package session

import (
	"encoding/json"
	"os"
	"time"
)

const stateFile = ".dojo_state.json"

// State represents an active training session
type State struct {
	ActiveIncidentID     string    `json:"active_incident_id"`
	StartTime            time.Time `json:"start_time"`
	VerificationAttempts int       `json:"verification_attempts"`
	HintLevel            int       `json:"hint_level"` // 0=none, 1=vague, 2=direction, 3=command, 4=root cause
	Mode                 string    `json:"mode"`       // "normal", "timed", "blind", "interview"
}

// LoadState reads the local state file, returning nil if it doesn't exist
func LoadState() (*State, error) {
	data, err := os.ReadFile(stateFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // No active session
		}
		return nil, err
	}
	var s State
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, err
	}
	return &s, nil
}

// SaveState persists the state to disk
func SaveState(s *State) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(stateFile, data, 0644)
}

// ClearState deletes the active session
func ClearState() error {
	return os.Remove(stateFile)
}
