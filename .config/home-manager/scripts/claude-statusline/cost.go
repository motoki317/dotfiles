package main

import (
	"context"
	"encoding/json"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// agtlog resolves a session ref directly, so lookup needs no date window and
// cannot undercount a long-running session. costTimeout keeps a hung subprocess
// from stalling the statusline. A price-table change can force a three-second
// full parse, so one render can omit cost. Claude Code retries 60 seconds later.
const costTimeout = 2 * time.Second

// Input is uncached input. CacheCreation and CacheRead are cache writes and reads.
// Output is generated tokens.
type costInfo struct {
	USD           float64
	Complete      bool
	Input         int64
	CacheCreation int64
	CacheRead     int64
	Output        int64
}

// Session totals already include subagents.
// The omitted estimated field is Codex-only. claude: refs never use ChatGPT-plan estimates.
type agtlogShow struct {
	Session *struct {
		Tokens struct {
			UncachedInput int64 `json:"uncached_input"`
			Output        int64 `json:"output"`
			CacheWrite    int64 `json:"cache_write"`
			CacheRead     int64 `json:"cache_read"`
		} `json:"tokens"`
		Cost struct {
			USD      float64 `json:"usd"`
			Complete bool    `json:"complete"`
		} `json:"cost"`
	} `json:"session"`
}

// sessionID resolves the session's UUID: the payload field when present, else the
// transcript filename stem (<uuid>.jsonl), which equals it. The fallback keeps
// cost working on any client whose payload omits session_id.
func sessionID(p Payload) string {
	if p.SessionID != "" {
		return p.SessionID
	}
	return strings.TrimSuffix(filepath.Base(p.TranscriptPath), ".jsonl")
}

// sessionCost returns this session's cost via agtlog, or nil when agtlog is not
// installed, times out, errors, or has no entry for the session — every one
// degrades to no cost segment, never a fatal error.
func sessionCost(sid string) *costInfo {
	if sid == "" {
		return nil
	}
	bin, err := exec.LookPath("agtlog")
	if err != nil {
		return nil
	}
	return agtlogCost(bin, sid)
}

// agtlogCost runs one offline summary query and returns nil on any failure.
func agtlogCost(bin, sid string) *costInfo {
	ctx, cancel := context.WithTimeout(context.Background(), costTimeout)
	defer cancel()
	return agtlogCostContext(ctx, bin, sid)
}

func agtlogCostContext(ctx context.Context, bin, sid string) *costInfo {
	out, err := exec.CommandContext(ctx, bin, "show", "claude:"+sid, "--no-events", "--offline").Output()
	if err != nil {
		return nil
	}
	return parseCost(out)
}

// parseCost returns nil unless the payload contains a decodable session summary.
func parseCost(data []byte) *costInfo {
	var result agtlogShow
	if json.Unmarshal(data, &result) != nil || result.Session == nil {
		return nil
	}
	s := result.Session
	return &costInfo{
		USD:           s.Cost.USD,
		Complete:      s.Cost.Complete,
		Input:         s.Tokens.UncachedInput,
		CacheCreation: s.Tokens.CacheWrite,
		CacheRead:     s.Tokens.CacheRead,
		Output:        s.Tokens.Output,
	}
}

// renderCost renders "<glyph> $<cost> ↑<cache-read>/<cache-write>/<input> ↓<output>".
// The ↑ group is tokens sent, ordered largest-first: the cached prefix re-read,
// the delta written to cache, then fresh uncached input; ↓ is generated output.
// Cost and every token are colored by magnitude (per-field bands ≈ recent-session
// p50/p90), the arrows and slashes dimmed. Empty when cost is absent.
func renderCost(ci *costInfo) string {
	if ci == nil {
		return ""
	}
	costColor, costStr := tcolor(int64(ci.USD), 10, 40), "$"+sigFig(ci.USD, 3)
	if !ci.Complete && ci.USD == 0 {
		costColor, costStr = cGray, "$?"
	}
	cost := paint(costColor, costStr)
	if !ci.Complete && ci.USD > 0 {
		cost += paint(cGray, "!")
	}
	read := paint(tcolor(ci.CacheRead, 10_000_000, 50_000_000), humanize(ci.CacheRead))
	create := paint(tcolor(ci.CacheCreation, 400_000, 2_000_000), humanize(ci.CacheCreation))
	input := paint(tcolor(ci.Input, 20_000, 200_000), humanize(ci.Input))
	output := paint(tcolor(ci.Output, 100_000, 350_000), humanize(ci.Output))
	return glyphCost + " " + cost + " " +
		paint(cGray, "↑") + read + paint(cGray, "/") + create + paint(cGray, "/") + input + " " +
		paint(cGray, "↓") + output
}
