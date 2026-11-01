package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	claudeUsageFetchTimeout  = 15 * time.Second
	claudeUsageResponseLimit = 1 << 20
)

var claudeSource = usageSource{
	file: "claude-usage.json", bin: "claude", arg: "--refresh-claude", interval: 45, fetch: fetchClaudeUsage,
}

// Reading usage does not consume the plan's allowance, but the endpoint answers
// frequent polling with HTTP 429.
// Claude Code owns token refresh, so this function only reads the stored token.
func fetchClaudeUsage() []Window {
	ctx, cancel := context.WithTimeout(context.Background(), claudeUsageFetchTimeout)
	defer cancel()
	token := claudeAccessToken(ctx)
	if token == "" {
		return nil
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.anthropic.com/api/oauth/usage", nil)
	if err != nil {
		return nil
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("User-Agent", "claude-statusline")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil
	}
	data, err := io.ReadAll(io.LimitReader(resp.Body, claudeUsageResponseLimit+1))
	if err != nil || len(data) > claudeUsageResponseLimit {
		return nil
	}
	return parseClaudeUsage(data)
}

func claudeAccessToken(ctx context.Context) string {
	var data []byte
	var err error
	if runtime.GOOS == "darwin" {
		data, err = exec.CommandContext(ctx, "security", "find-generic-password", "-s", "Claude Code-credentials", "-w").Output()
	} else {
		home, homeErr := os.UserHomeDir()
		if homeErr != nil {
			return ""
		}
		data, err = os.ReadFile(filepath.Join(home, ".claude", ".credentials.json"))
	}
	if err != nil {
		return ""
	}
	var credentials struct {
		ClaudeAIOAuth struct {
			AccessToken string `json:"accessToken"`
		} `json:"claudeAiOauth"`
	}
	if json.Unmarshal(data, &credentials) != nil {
		return ""
	}
	return strings.TrimSpace(credentials.ClaudeAIOAuth.AccessToken)
}

func parseClaudeUsage(data []byte) []Window {
	var usage map[string]json.RawMessage
	if json.Unmarshal(data, &usage) != nil {
		return nil
	}
	var windows []Window
	for _, length := range claudeWindowLengths {
		data, ok := usage[length.key]
		if !ok {
			continue
		}
		var w struct {
			Utilization *float64 `json:"utilization"`
			ResetsAt    string   `json:"resets_at"`
		}
		if json.Unmarshal(data, &w) != nil {
			return nil
		}
		if w.Utilization == nil {
			continue
		}
		window := Window{Minutes: length.minutes, Used: *w.Utilization}
		if reset, err := time.Parse(time.RFC3339, w.ResetsAt); err == nil {
			unix := reset.Unix()
			window.ResetsAt = &unix
		}
		windows = append(windows, window)
	}
	return windows
}
