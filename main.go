package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type Credentials struct {
	ClaudeAiOauth struct {
		AccessToken string `json:"accessToken"`
		ExpiresAt   int64  `json:"expiresAt"`
	} `json:"claudeAiOauth"`
}

type UsageBucket struct {
	Utilization float64 `json:"utilization"`
	ResetsAt    string  `json:"resets_at"`
}

type UsageResponse struct {
	FiveHour UsageBucket `json:"five_hour"`
	SevenDay UsageBucket `json:"seven_day"`
}

type WaybarOutput struct {
	Text       string `json:"text"`
	Tooltip    string `json:"tooltip"`
	Percentage int    `json:"percentage"`
}

func main() {
	output, err := run()
	if err != nil {
		output = &WaybarOutput{
			Text:    "󰚩  –",
			Tooltip: fmt.Sprintf("Error: %s", err),
		}
	}

	json.NewEncoder(os.Stdout).Encode(output)
}

func run() (*WaybarOutput, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("home dir: %w", err)
	}

	data, err := os.ReadFile(filepath.Join(home, ".claude", ".credentials.json"))
	if err != nil {
		return nil, fmt.Errorf("read credentials: %w", err)
	}

	var creds Credentials
	if err := json.Unmarshal(data, &creds); err != nil {
		return nil, fmt.Errorf("parse credentials: %w", err)
	}

	token := creds.ClaudeAiOauth.AccessToken
	if token == "" {
		return nil, fmt.Errorf("no OAuth token found")
	}

	if creds.ClaudeAiOauth.ExpiresAt > 0 && time.Now().UnixMilli() > creds.ClaudeAiOauth.ExpiresAt {
		return nil, fmt.Errorf("OAuth token expired")
	}

	req, err := http.NewRequest("GET", "https://api.anthropic.com/api/oauth/usage", nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("anthropic-beta", "oauth-2025-04-20")
	req.Header.Set("User-Agent", "claude-code/2.1.69")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API returned %d", resp.StatusCode)
	}

	var usage UsageResponse
	if err := json.NewDecoder(resp.Body).Decode(&usage); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	pct := int(math.Round(usage.SevenDay.Utilization))
	sessionPct := int(math.Round(usage.FiveHour.Utilization))

	tooltip := fmt.Sprintf("Session (5h): %d%%\nWeekly (7d): %d%%", sessionPct, pct)

	if usage.FiveHour.ResetsAt != "" {
		if t, err := time.Parse(time.RFC3339, usage.FiveHour.ResetsAt); err == nil {
			tooltip += fmt.Sprintf("\n\nSession resets: %s", t.Local().Format("15:04"))
		}
	}
	if usage.SevenDay.ResetsAt != "" {
		if t, err := time.Parse(time.RFC3339, usage.SevenDay.ResetsAt); err == nil {
			tooltip += fmt.Sprintf("\nWeekly resets: %s", t.Local().Format("Mon 15:04"))
		}
	}

	return &WaybarOutput{
		Text:       fmt.Sprintf("󰚩  %d%%", pct),
		Tooltip:    tooltip,
		Percentage: pct,
	}, nil
}
