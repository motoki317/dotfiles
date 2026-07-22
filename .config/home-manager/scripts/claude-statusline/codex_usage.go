package main

import (
	"bufio"
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"time"
)

// codexScanLimit bounds how many of the most-recently-modified rollout files are
// inspected. Rate limits are account-global, so the freshest snapshot is in one
// of the last-written sessions; scanning a handful and picking the newest event
// stays correct when several sessions are active without reading all ~hundreds.
const codexScanLimit = 8

var codexSource = usageSource[CodexRL]{
	file: "codex-rate-limits.json", bin: "codex", arg: "--refresh-codex", interval: 45, fetch: fetchCodexRateLimits,
}

// codexRateLimits returns the Codex rate-limit snapshot for rendering. It reads the
// cache maintained by the background refresh; when the cache is missing or older than
// the source interval it kicks off a detached refresh (non-blocking) that updates the
// cache for the next render. The log scan is the fallback, so the segment still appears
// before the first successful fetch and when the backend is unreachable.
func codexRateLimits(home string, now int64) *CodexRL {
	cache := codexSource.read()
	if cache == nil || now-cache.FetchedAt > codexSource.interval {
		codexSource.spawn(now)
	}
	if cache != nil {
		return cache.RateLimits
	}
	if home == "" {
		return nil
	}
	return codexRateLimitsFromLogs(home)
}

// fetchCodexRateLimits drives `codex app-server` over stdio to call account/rateLimits/read
// — a live GET to the backend's accounts/check, which is authoritative and account-global.
// The app-server owns auth and token refresh, so this stays a thin JSON-RPC client. Returns
// nil on any failure (codex missing, logged out, timeout, malformed reply).
func fetchCodexRateLimits() *CodexRL {
	bin, err := exec.LookPath("codex")
	if err != nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), codexFetchTimeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin, "app-server")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil
	}
	if cmd.Start() != nil {
		return nil
	}
	defer func() { _ = cmd.Process.Kill(); _ = cmd.Wait() }()

	enc := json.NewEncoder(stdin)
	sc := bufio.NewScanner(stdout)
	sc.Buffer(make([]byte, 0, 64*1024), 16*1024*1024)

	// Handshake: initialize (id 0), await its result, then initialized + the read (id 1).
	_ = enc.Encode(rpcMsg{JSONRPC: "2.0", ID: pint(0), Method: "initialize",
		Params: map[string]any{"clientInfo": map[string]any{"name": "claude-statusline", "version": "0"}}})
	if awaitRPCResult(sc, 0) == nil {
		return nil
	}
	_ = enc.Encode(rpcMsg{JSONRPC: "2.0", Method: "initialized"})
	_ = enc.Encode(rpcMsg{JSONRPC: "2.0", ID: pint(1), Method: "account/rateLimits/read", Params: map[string]any{}})
	result := awaitRPCResult(sc, 1)
	if result == nil {
		return nil
	}
	return parseAppServerRateLimits(result)
}

// codexFetchTimeout caps the whole app-server round-trip; a hang degrades to no refresh
// (the previous cache stands) rather than a lingering process.
const codexFetchTimeout = 15 * time.Second

// rpcMsg is one outbound JSON-RPC message. A nil ID omits the field, marking a notification.
type rpcMsg struct {
	JSONRPC string `json:"jsonrpc"`
	ID      *int   `json:"id,omitempty"`
	Method  string `json:"method"`
	Params  any    `json:"params,omitempty"`
}

func pint(i int) *int { return &i }

// awaitRPCResult scans line-delimited JSON-RPC messages until one carries id with a
// result, returning that raw result. Nil on EOF/timeout or an error reply. Interleaved
// notifications (no id) and other ids are skipped.
func awaitRPCResult(sc *bufio.Scanner, id int) json.RawMessage {
	for sc.Scan() {
		var m struct {
			ID     *int            `json:"id"`
			Result json.RawMessage `json:"result"`
		}
		if json.Unmarshal(sc.Bytes(), &m) != nil {
			continue
		}
		if m.ID != nil && *m.ID == id && m.Result != nil {
			return m.Result
		}
	}
	return nil
}

// parseAppServerRateLimits converts account/rateLimits/read's result into CodexRL. The
// app-server uses camelCase (usedPercent, windowDurationMins, resetsAt), distinct from
// the snake_case the rollout logs carry, so this cannot reuse CodexRL's own unmarshal.
func parseAppServerRateLimits(result json.RawMessage) *CodexRL {
	var r struct {
		RateLimits struct {
			Primary   *appWindow `json:"primary"`
			Secondary *appWindow `json:"secondary"`
		} `json:"rateLimits"`
	}
	if json.Unmarshal(result, &r) != nil {
		return nil
	}
	primary, secondary := r.RateLimits.Primary.toCodex(), r.RateLimits.Secondary.toCodex()
	if primary == nil && secondary == nil {
		return nil
	}
	return &CodexRL{Primary: primary, Secondary: secondary}
}

// appWindow is one window as account/rateLimits/read reports it. toCodex maps it to the
// render-side CodexWindow, dropping a window whose used percent is absent.
type appWindow struct {
	UsedPercent        *float64 `json:"usedPercent"`
	WindowDurationMins *int     `json:"windowDurationMins"`
	ResetsAt           *int64   `json:"resetsAt"`
}

func (w *appWindow) toCodex() *CodexWindow {
	if w == nil || w.UsedPercent == nil {
		return nil
	}
	return &CodexWindow{UsedPercent: w.UsedPercent, WindowMinutes: w.WindowDurationMins, ResetsAt: w.ResetsAt}
}

// codexRateLimitsFromLogs returns the freshest Codex rate-limit snapshot found in the
// rollout logs, or nil when Codex is not installed or has no session data. It reads the
// newest rollout files by mtime and picks the entry with the newest event timestamp, so
// a just-started session that has not hit the API yet does not hide older-but-real data,
// and concurrent sessions resolve to the genuinely latest event.
//
// This is the cold-start/offline fallback for codexRateLimits: the logged snapshot is a
// passive echo of the last /responses call, so it lags a window rollover and misses
// consumption from other clients. codexRateLimits prefers the live backend snapshot.
func codexRateLimitsFromLogs(home string) *CodexRL {
	if home == "" {
		return nil
	}
	files, _ := filepath.Glob(filepath.Join(home, ".codex", "sessions", "*", "*", "*", "rollout-*.jsonl"))
	if len(files) == 0 {
		return nil
	}

	type stamped struct {
		path string
		mod  time.Time
	}
	cand := make([]stamped, 0, len(files))
	for _, f := range files {
		if st, err := os.Stat(f); err == nil {
			cand = append(cand, stamped{f, st.ModTime()})
		}
	}
	sort.Slice(cand, func(i, j int) bool { return cand[i].mod.After(cand[j].mod) })

	var best *CodexRL
	var bestTS time.Time
	for i := 0; i < len(cand) && i < codexScanLimit; i++ {
		rl, ts, ok := lastRateLimitEvent(cand[i].path)
		if !ok {
			continue
		}
		if best == nil || ts.After(bestTS) {
			best, bestTS = rl, ts
		}
	}
	return best
}

// lastRateLimitEvent returns the last structurally-valid token_count event that
// carries rate_limits in a rollout file, with its parsed timestamp. Decoding each
// line as JSON (rather than grepping the "rate_limits" substring) avoids matching
// the string inside user text or tool output; malformed or truncated lines are
// skipped, so a partially-written final line never aborts the scan.
func lastRateLimitEvent(path string) (*CodexRL, time.Time, bool) {
	f, err := os.Open(path)
	if err != nil {
		return nil, time.Time{}, false
	}
	defer f.Close()

	var rl *CodexRL
	var ts time.Time
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 64*1024), 16*1024*1024)
	for sc.Scan() {
		var ev codexEvent
		if json.Unmarshal(sc.Bytes(), &ev) != nil {
			continue
		}
		if ev.Type != "event_msg" || ev.Payload.Type != "token_count" || ev.Payload.RateLimits == nil {
			continue
		}
		rl = ev.Payload.RateLimits
		ts, _ = time.Parse(time.RFC3339, ev.Timestamp) // zero time on parse failure
	}
	if rl == nil {
		return nil, time.Time{}, false
	}
	return rl, ts, true
}

// codexEvent is one line of a rollout JSONL file. Only event_msg/token_count
// lines carrying rate_limits are of interest.
type codexEvent struct {
	Type      string `json:"type"`
	Timestamp string `json:"timestamp"`
	Payload   struct {
		Type       string   `json:"type"`
		RateLimits *CodexRL `json:"rate_limits"`
	} `json:"payload"`
}
