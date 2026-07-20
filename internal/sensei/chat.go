package sensei

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// AskAgent calls the Gemini API as a conversational agent.
func AskAgent(prompt string) (string, error) {
	apiKey := os.Getenv("DOJO_API_KEY")
	if apiKey == "" {
		return "I need a DOJO_API_KEY environment variable to connect to my brain!", nil
	}

	systemPrompt := "You are Dojo AI, an expert DevOps engineering assistant. You are chatting with a user in a terminal interface. Be helpful, concise, and provide clear code snippets or commands when applicable."
	
	fullPrompt := systemPrompt + "\n\nUser: " + prompt + "\n\nAgent:"

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
				{Text: fullPrompt},
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
		return "", fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	candidates, ok := result["candidates"].([]interface{})
	if !ok || len(candidates) == 0 {
		return "Hmm, I couldn't formulate a response.", nil
	}

	candidate, ok := candidates[0].(map[string]interface{})
	if !ok {
		return "Invalid response format.", nil
	}
	content, ok := candidate["content"].(map[string]interface{})
	if !ok {
		return "Invalid response structure from my brain.", nil
	}

	parts, ok := content["parts"].([]interface{})
	if !ok || len(parts) == 0 {
		return "Empty response.", nil
	}

	textPart, ok := parts[0].(map[string]interface{})
	if !ok {
		return "Failed to parse text part.", nil
	}
	text, ok := textPart["text"].(string)
	if !ok {
		return "Failed to extract text.", nil
	}

	return text, nil
}
