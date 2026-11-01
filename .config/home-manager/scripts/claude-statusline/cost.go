package main

import (
	"context"
	"encoding/json"
	"os/exec"
	"time"
)

// After its price table changes, agtlog reparses for about 3 seconds. The timeout
// drops cost from that one render instead of delaying the status line.
const costTimeout = 2 * time.Second

type costInfo struct {
	USD           float64
	Complete      bool
	Input         int64
	CacheCreation int64
	CacheRead     int64
	Output        int64
}

// agtlog includes subagents in session totals, so no separate aggregation is needed.
// The cost's estimated flag applies only to Codex sessions, so it is not decoded.
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

func sessionCost(sid string) *costInfo {
	if sid == "" {
		return nil
	}
	bin, err := exec.LookPath("agtlog")
	if err != nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), costTimeout)
	defer cancel()
	return agtlogCost(ctx, bin, sid)
}

func agtlogCost(ctx context.Context, bin, sid string) *costInfo {
	out, err := exec.CommandContext(ctx, bin, "show", "claude:"+sid, "--no-events", "--offline").Output()
	if err != nil {
		return nil
	}
	return parseCost(out)
}

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

// Each field's two tcolor thresholds approximate its p50 and p90 over past sessions.
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
