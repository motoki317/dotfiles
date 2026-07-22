package main

import (
	"context"
	"encoding/json"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// costScanDays bounds how many days of history ccusage scans (--since), so render
// latency stays flat as the transcript archive grows rather than scaling with all
// of it. The active session is always recent, so this never hides it — but a
// single session running longer than costScanDays undercounts; widen it if that
// matters. costTimeout caps the ccusage subprocess so a hang degrades to no cost
// segment instead of stalling the whole line.
const (
	costScanDays = 3
	costTimeout  = 2 * time.Second
)

// costInfo is one session's cumulative spend and token usage, as computed by
// ccusage. Input is fresh (uncached) input; CacheCreation and CacheRead are the
// cache-write and cache-read halves of cached input; Output is generated tokens.
type costInfo struct {
	USD           float64
	Input         int64
	CacheCreation int64
	CacheRead     int64
	Output        int64
}

// tokens is the session's total token count. Zero cost with nonzero tokens is the
// signal that a pricing source couldn't price the model (a priced session is never
// exactly 0 with tokens), driving both the online escalation and the $? fallback.
func (ci *costInfo) tokens() int64 {
	return ci.Input + ci.CacheCreation + ci.CacheRead + ci.Output
}

// ccusageSession mirrors the fields this tool reads from `ccusage session --json`.
// ccusage keys each entry by period == the Claude Code session_id and folds every
// sub-agent (the flat <session>/subagents/*.jsonl files) and every compaction (all
// in the one main transcript) into that entry, so one lookup is the whole-session
// total — the reason cost is delegated here rather than summed in-process: ccusage
// owns the per-model pricing, which is not the classic $15/$75 for current models.
type ccusageSession struct {
	Session []struct {
		Period              string  `json:"period"`
		TotalCost           float64 `json:"totalCost"`
		InputTokens         int64   `json:"inputTokens"`
		OutputTokens        int64   `json:"outputTokens"`
		CacheCreationTokens int64   `json:"cacheCreationTokens"`
		CacheReadTokens     int64   `json:"cacheReadTokens"`
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

// sessionCost returns this session's cost via ccusage, or nil when ccusage is not
// installed, times out, errors, or has no entry for the session — every one
// degrades to no cost segment, never a fatal error.
//
// It reads offline first (ccusage's embedded price table: no network, ~60ms). That
// table is baked in at ccusage build time, so it lags newly released models — one it
// can't price yields cost 0 with real token counts. Only then does it pay for one
// authoritative online fetch (LiteLLM), which prices new models correctly. So every
// already-priced session keeps the fast offline path; only a model too new for the
// installed ccusage triggers the network, and only until ccusage embeds its price.
func sessionCost(sid string, now int64) *costInfo {
	if sid == "" {
		return nil
	}
	bin, err := exec.LookPath("ccusage")
	if err != nil {
		return nil
	}
	since := time.Unix(now, 0).AddDate(0, 0, -costScanDays).Format("20060102")

	ci := ccusageCost(bin, since, sid, true)
	if ci != nil && ci.USD == 0 && ci.tokens() > 0 {
		if online := ccusageCost(bin, since, sid, false); online != nil && online.USD > 0 {
			return online
		}
	}
	return ci
}

// ccusageCost runs one `ccusage session --json` query bounded by costTimeout and
// returns the session's totals, or nil on any failure. offline selects ccusage's
// embedded price table over a live LiteLLM fetch; --since bounds the scan (see
// costScanDays).
func ccusageCost(bin, since, sid string, offline bool) *costInfo {
	args := []string{"session", "--json", "--since", since}
	if offline {
		args = append(args, "--offline")
	}
	ctx, cancel := context.WithTimeout(context.Background(), costTimeout)
	defer cancel()
	out, err := exec.CommandContext(ctx, bin, args...).Output()
	if err != nil {
		return nil
	}
	return parseCost(out, sid)
}

// parseCost extracts the session's totals from ccusage `session --json` output.
// Returns nil when the payload can't be decoded or holds no matching session.
func parseCost(data []byte, sid string) *costInfo {
	var cu ccusageSession
	if json.Unmarshal(data, &cu) != nil {
		return nil
	}
	for _, s := range cu.Session {
		if s.Period == sid {
			return &costInfo{
				USD:           s.TotalCost,
				Input:         s.InputTokens,
				CacheCreation: s.CacheCreationTokens,
				CacheRead:     s.CacheReadTokens,
				Output:        s.OutputTokens,
			}
		}
	}
	return nil
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
	// Cost 0 with real tokens means no source could price the model — sessionCost
	// already tried the live fetch, so this is a genuine unknown (e.g. offline, or a
	// model no pricing source knows yet), not a false $0. Show it as such.
	if ci.USD == 0 && ci.tokens() > 0 {
		costColor, costStr = cGray, "$?"
	}
	cost := paint(costColor, costStr)
	read := paint(tcolor(ci.CacheRead, 10_000_000, 50_000_000), humanize(ci.CacheRead))
	create := paint(tcolor(ci.CacheCreation, 400_000, 2_000_000), humanize(ci.CacheCreation))
	input := paint(tcolor(ci.Input, 20_000, 200_000), humanize(ci.Input))
	output := paint(tcolor(ci.Output, 100_000, 350_000), humanize(ci.Output))
	return glyphCost + " " + cost + " " +
		paint(cGray, "↑") + read + paint(cGray, "/") + create + paint(cGray, "/") + input + " " +
		paint(cGray, "↓") + output
}
