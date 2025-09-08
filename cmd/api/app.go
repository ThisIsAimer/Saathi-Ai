package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load(`cmd\api\.env`)
	if err != nil {
		fmt.Println("Failed to load env")
		return
	}

	apiKey := os.Getenv("GROQ_API_KEY")
	if apiKey == "" {
		fmt.Println("set GROQ_API_KEY")
		return
	}

	systemPrompt := os.Getenv("AI_SYSTEM_PROMPT")
	if systemPrompt == "" {
		fmt.Println("failed to get system prompt")
		return
	}

	url := "https://api.groq.com/openai/v1/chat/completions"

	reqBody := map[string]any{
		"model": "llama-3.1-8b-instant", // choose a model available to you
		"messages": []map[string]string{
			{"role": "system", "content": systemPrompt},
			{"role": "user", "content": "Explain goroutines."},
		},
		"temperature": 0.2,
	}

	b, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", url, bytes.NewReader(b))
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		panic(fmt.Sprintf("error: %s %s", resp.Status, string(body)))
	}

	// Groq's response JSON structure mirrors OpenAI-style responses.
	var out map[string]any
	err = json.NewDecoder(resp.Body).Decode(&out)
	if err != nil {
		panic(err)
	}

	choices, ok := out["choices"].([]any)

	if ok && len(choices) > 0 {
		ch0, ok := choices[0].(map[string]any)

		if ok {
			msg, ok := ch0["message"].(map[string]any)
			
			if ok {
				fmt.Println("assistant:", msg["content"])
			}
		}
	}
}
