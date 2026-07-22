// Command statusline renders the Claude Code status line: context-window usage,
// then Anthropic and OpenAI Codex rate limits. It reads the harness JSON payload
// on stdin and prints one line. home-manager builds and installs it on PATH.
//
// It replaces a bash script whose stringly-typed JSON handling had grown
// error-prone: nested extraction, filesystem discovery, timestamps, ANSI
// assembly, and stale-snapshot rules pushed past what shell does safely (a
// tab-delimited `read`, for one, silently misaligned fields when any was null).
// Decoding into typed structs removes that whole class of failure and makes the
// render logic unit-testable — see statusline_test.go.
//
// Design rules that hold throughout: every source is optional, so missing or
// malformed data degrades to absence, never a fatal error or a blank line.
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

// Payload is the harness JSON on stdin. Pointers mark optional objects so absent
// and zero are distinguishable.
type Payload struct {
	Model struct {
		DisplayName string `json:"display_name"`
	} `json:"model"`
	ContextWindow  *ContextWindow `json:"context_window"`
	RateLimits     *RateLimits    `json:"rate_limits"`
	TranscriptPath string         `json:"transcript_path"`
	SessionID      string         `json:"session_id"`
}

func main() {
	// --refresh-claude is the detached self-invocation spawned by
	// claudeUsageRateLimits. Like the Codex refresh below, it never shares the
	// render's stdin/stdout and only publishes a complete cache snapshot.
	//
	// --refresh-codex is the detached self-invocation spawned by codexRateLimits: it
	// fetches the live snapshot from the backend and writes the cache, then exits
	// without touching stdin/stdout, so it never interferes with a render.
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
	if data, err := readAll(os.Stdin); err == nil {
		_ = json.Unmarshal(data, &p) // best effort: a decode error leaves p zero
	}
	home, _ := os.UserHomeDir()
	now := time.Now().Unix()
	if live := claudeUsageRateLimits(now); live != nil {
		p.RateLimits = live
	}
	codex := codexRateLimits(home, now)
	cost := sessionCost(sessionID(p), now)
	fmt.Println(render(p, codex, cost, now, termWidth()))
}

func readAll(f *os.File) ([]byte, error) {
	var b strings.Builder
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 64*1024), 16*1024*1024)
	for sc.Scan() {
		b.WriteString(sc.Text())
	}
	return []byte(b.String()), sc.Err()
}
