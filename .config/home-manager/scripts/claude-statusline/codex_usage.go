package main

import (
	"bufio"
	"context"
	"encoding/json"
	"os/exec"
	"time"
)

var codexSource = usageSource{
	file: "codex-rate-limits.json", bin: "codex", arg: "--refresh-codex", interval: 45, fetch: fetchCodexRateLimits,
}

// fetchCodexRateLimits drives `codex app-server` over stdio to call account/rateLimits/read
// — a live GET to the backend's accounts/check, which is authoritative and account-global.
// The app-server owns auth and token refresh, so this stays a thin JSON-RPC client. Returns
// nil on any failure (codex missing, logged out, timeout, malformed reply).
func fetchCodexRateLimits() []Window {
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

func parseAppServerRateLimits(result json.RawMessage) []Window {
	var r struct {
		RateLimits struct {
			Primary   *appWindow `json:"primary"`
			Secondary *appWindow `json:"secondary"`
		} `json:"rateLimits"`
	}
	if json.Unmarshal(result, &r) != nil {
		return nil
	}
	var windows []Window
	for _, w := range []*appWindow{r.RateLimits.Primary, r.RateLimits.Secondary} {
		if w != nil && w.UsedPercent != nil && w.WindowDurationMins != nil {
			windows = append(windows, Window{Minutes: *w.WindowDurationMins, Used: *w.UsedPercent, ResetsAt: w.ResetsAt})
		}
	}
	return windows
}

type appWindow struct {
	UsedPercent        *float64 `json:"usedPercent"`
	WindowDurationMins *int     `json:"windowDurationMins"`
	ResetsAt           *int64   `json:"resetsAt"`
}
