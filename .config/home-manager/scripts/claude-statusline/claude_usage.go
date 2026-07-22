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

// The timeout prevents a detached child from lingering indefinitely when credentials
// or the network stop responding. The body cap leaves ample room for ignored future
// fields without trusting an external response with unbounded memory.
const (
	claudeUsageFetchTimeout  = 15 * time.Second
	claudeUsageResponseLimit = 1 << 20
)

var claudeSource = usageSource[RateLimits]{
	file: "claude-usage.json", bin: "claude", arg: "--refresh-claude", interval: 45, fetch: fetchClaudeUsage,
}

// claudeUsageRateLimits returns the live Anthropic snapshot cached by the detached
// refresher. Unlike Codex, Claude has no account-global log fallback; nil deliberately
// leaves the harness payload untouched in main until a live fetch succeeds.
func claudeUsageRateLimits(now int64) *RateLimits {
	cache := claudeSource.read()
	var live *RateLimits
	if cache != nil {
		live = usableClaudeRateLimits(cache.RateLimits)
	}
	if cache == nil || live == nil || now-cache.FetchedAt > claudeSource.interval {
		claudeSource.spawn(now)
	}
	return live
}

// usableClaudeRateLimits protects the payload precedence seam from a syntactically
// valid but partial cache left by an interrupted upgrade or manual edit. Successful
// endpoint parsing already guarantees this invariant before any cache write.
func usableClaudeRateLimits(rl *RateLimits) *RateLimits {
	if rl == nil {
		return nil
	}
	fiveHour := rl.FiveHour != nil && rl.FiveHour.UsedPercentage != nil
	sevenDay := rl.SevenDay != nil && rl.SevenDay.UsedPercentage != nil
	if !fiveHour && !sevenDay {
		return nil
	}
	return rl
}

// fetchClaudeUsage reads Anthropic's unmetered account-global usage status with the
// access token already owned by Claude Code. It neither refreshes nor persists auth;
// every credential, request, status, and body failure leaves the payload fallback intact.
func fetchClaudeUsage() *RateLimits {
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

// claudeAccessToken reads Claude Code's current OAuth credential without taking
// ownership of its lifecycle. macOS stores the JSON blob in Keychain; other platforms
// use the home-directory credential file. Any read or shape error disables live usage.
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

// parseClaudeUsage converts Anthropic's usage response into the harness-shaped model
// the render path already consumes. The endpoint may omit either window, so an absent
// or unusable window is dropped independently and an entirely empty response is nil.
func parseClaudeUsage(data []byte) *RateLimits {
	var usage struct {
		FiveHour *claudeUsageWindow `json:"five_hour"`
		SevenDay *claudeUsageWindow `json:"seven_day"`
	}
	if json.Unmarshal(data, &usage) != nil {
		return nil
	}
	fiveHour, sevenDay := usage.FiveHour.toClaude(), usage.SevenDay.toClaude()
	if fiveHour == nil && sevenDay == nil {
		return nil
	}
	return &RateLimits{FiveHour: fiveHour, SevenDay: sevenDay}
}

type claudeUsageWindow struct {
	Utilization *float64 `json:"utilization"`
	ResetsAt    string   `json:"resets_at"`
}

// toClaude leaves an absent or malformed reset unset: utilization is still useful on
// its own, and every source in the statusline degrades field-by-field rather than
// discarding a whole segment because one optional value is invalid.
func (w *claudeUsageWindow) toClaude() *ClaudeWindow {
	if w == nil || w.Utilization == nil {
		return nil
	}
	result := &ClaudeWindow{UsedPercentage: w.Utilization}
	if reset, err := time.Parse(time.RFC3339, w.ResetsAt); err == nil {
		unix := reset.Unix()
		result.ResetsAt = &unix
	}
	return result
}
