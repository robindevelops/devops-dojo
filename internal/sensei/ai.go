package sensei

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/devops-dojo/cli/internal/session"
)

type geminiRequest struct {
	Contents []struct {
		Parts []struct {
			Text string `json:"text"`
		} `json:"parts"`
	} `json:"contents"`
}

// AskAI calls the Gemini API to get a dynamic hint.
func AskAI(state *session.State) (string, error) {
	apiKey := os.Getenv("DOJO_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("no DOJO_API_KEY found")
	}

	prompt := fmt.Sprintf("You are a Senior DevOps Mentor helping a junior engineer. They are currently facing this incident: '%s'. They have asked for a hint (this is hint #%d they have requested). Please provide a brief, helpful, and educational hint. Do not give them the exact code solution, but point them in the right direction based on their hint level.", state.ActiveIncidentID, state.HintLevel)

	reqBody := geminiRequest{}
	reqBody.Contents = []struct {
		Parts []struct {
			Text string `json:"text"`
		} `json:"parts"`
	}{
		{
			Parts: []struct {
				Text string `json:"text"`
			}{
				{Text: prompt},
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent?key=" + apiKey
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("AI API returned status %d: %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	// Parse candidates[0].content.parts[0].text
	candidates, ok := result["candidates"].([]interface{})
	if !ok || len(candidates) == 0 {
		return "", fmt.Errorf("no response candidates from AI")
	}

	candidate := candidates[0].(map[string]interface{})
	content, ok := candidate["content"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid AI response structure")
	}

	parts, ok := content["parts"].([]interface{})
	if !ok || len(parts) == 0 {
		return "", fmt.Errorf("empty response parts from AI")
	}

	textPart := parts[0].(map[string]interface{})
	text, ok := textPart["text"].(string)
	if !ok {
		return "", fmt.Errorf("failed to extract text from AI response")
	}

	return text, nil
}
