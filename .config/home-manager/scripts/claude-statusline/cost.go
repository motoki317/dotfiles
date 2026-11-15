package main

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"
)

// The render waits this long for a fresh cost. A slower refresh finishes in the
// background so later renders can use its result.
const costTimeout = 2 * time.Second

const costRefreshTimeout = 60 * time.Second

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
	base := costCacheBase(sid)
	if base == "" {
		return nil
	}
	if cmd, err := spawnDetached("--refresh-cost", sid); err == nil {
		done := make(chan struct{})
		go func() {
			_ = cmd.Wait()
			close(done)
		}()
		timer := time.NewTimer(costTimeout)
		defer timer.Stop()
		select {
		case <-done:
		case <-timer.C:
		}
	}
	data, _ := os.ReadFile(base + ".json")
	return parseCost(data)
}

func costCacheBase(sid string) string {
	if sid == "" {
		return ""
	}
	for _, c := range sid {
		if !(c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z' || c == '-') {
			return ""
		}
	}
	return cachePath(filepath.Join("cost", sid))
}

func refreshCost(sid string) {
	base := costCacheBase(sid)
	if base == "" {
		return
	}
	bin, err := exec.LookPath("agtlog")
	if err != nil {
		return
	}
	if os.MkdirAll(filepath.Dir(base), 0o755) != nil {
		return
	}
	lock, err := os.OpenFile(base+".lock", os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		return
	}
	defer lock.Close()
	if syscall.Flock(int(lock.Fd()), syscall.LOCK_EX|syscall.LOCK_NB) != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), costRefreshTimeout)
	defer cancel()
	if data := agtlogCost(ctx, bin, sid); data != nil {
		_ = writeCache(base+".json", data)
	}
}

func agtlogCost(ctx context.Context, bin, sid string) []byte {
	out, err := exec.CommandContext(ctx, bin, "show", "claude:"+sid, "--no-events", "--offline").Output()
	if err != nil || parseCost(out) == nil {
		return nil
	}
	return out
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
