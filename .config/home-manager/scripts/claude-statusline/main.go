package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Payload struct {
	Model struct {
		DisplayName string `json:"display_name"`
	} `json:"model"`
	ContextWindow *ContextWindow `json:"context_window"`
	RateLimits    map[string]struct {
		Used     *float64 `json:"used_percentage"`
		ResetsAt *int64   `json:"resets_at"`
	} `json:"rate_limits"`
	SessionID string `json:"session_id"`
}

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--refresh-claude":
			claudeSource.refresh()
			return
		case "--refresh-codex":
			codexSource.refresh()
			return
		}
	}

	var p Payload
	// encoding/json keeps filling the other fields after a type error, so a malformed
	// field drops only its own segment.
	_ = json.NewDecoder(os.Stdin).Decode(&p)
	now := time.Now().Unix()
	claude := claudeSource.current(now)
	if len(claude) == 0 {
		claude = p.claudeWindows()
	}
	codex := codexSource.current(now)
	cost := sessionCost(p.SessionID)
	fmt.Println(render(p, claude, codex, cost, now, termWidth()))
}
