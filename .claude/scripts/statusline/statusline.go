// Command statusline renders the Claude Code status line: context-window usage,
// then Anthropic and OpenAI Codex rate limits. It reads the harness JSON payload
// on stdin and prints one line. The statusline.sh launcher builds and caches it.
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
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

// ANSI colors. Real ESC bytes, so rendered segments concatenate directly.
const (
	cRed    = "\x1b[31m"
	cYellow = "\x1b[33m"
	cGreen  = "\x1b[32m"
	cGray   = "\x1b[90m"
	cReset  = "\x1b[0m"
)

// Nerd-font glyphs (Material Design Icons range): the model marker, and the
// timer-sand that prefixes each rate-limit group.
const (
	glyphModel = "\U000F06A9" // 󰚩
	glyphTimer = "\U000F051F" // 󰔟
)

// Anthropic's payload omits window lengths, so they are fixed here. Codex's
// payload carries window_minutes, so its lengths and labels are derived instead.
const (
	fiveHourSecs = 5 * 60 * 60
	sevenDaySecs = 7 * 24 * 60 * 60
)

// codexScanLimit bounds how many of the most-recently-modified rollout files are
// inspected. Rate limits are account-global, so the freshest snapshot is in one
// of the last-written sessions; scanning a handful and picking the newest event
// stays correct when several sessions are active without reading all ~hundreds.
const codexScanLimit = 8

// Payload is the harness JSON on stdin. Pointers mark optional objects so absent
// and zero are distinguishable.
type Payload struct {
	Model struct {
		DisplayName string `json:"display_name"`
	} `json:"model"`
	ContextWindow  *ContextWindow `json:"context_window"`
	RateLimits     *RateLimits    `json:"rate_limits"`
	TranscriptPath string         `json:"transcript_path"`
}

// ContextWindow carries the harness-computed, /context-faithful figures. Newer
// clients supply TotalInputTokens directly; older ones only CurrentUsage.
type ContextWindow struct {
	TotalInputTokens  *int64       `json:"total_input_tokens"`
	ContextWindowSize *int64       `json:"context_window_size"`
	UsedPercentage    *float64     `json:"used_percentage"`
	CurrentUsage      CurrentUsage `json:"current_usage"`
}

// CurrentUsage is the input-only breakdown (output tokens excluded, matching
// /context): summed when TotalInputTokens is absent.
type CurrentUsage struct {
	InputTokens              int64 `json:"input_tokens"`
	CacheCreationInputTokens int64 `json:"cache_creation_input_tokens"`
	CacheReadInputTokens     int64 `json:"cache_read_input_tokens"`
}

// RateLimits holds Anthropic's two rate-limit windows, injected live each render.
type RateLimits struct {
	FiveHour *ClaudeWindow `json:"five_hour"`
	SevenDay *ClaudeWindow `json:"seven_day"`
}

// ClaudeWindow is one Anthropic rate-limit window, injected live each render.
type ClaudeWindow struct {
	UsedPercentage *float64 `json:"used_percentage"`
	ResetsAt       *int64   `json:"resets_at"`
}

// CodexRL and CodexWindow mirror the rate_limits object Codex writes into its
// session rollout files. window_minutes is explicit, unlike Anthropic's payload.
type CodexRL struct {
	Primary   *CodexWindow `json:"primary"`
	Secondary *CodexWindow `json:"secondary"`
}

type CodexWindow struct {
	UsedPercent   *float64 `json:"used_percent"`
	WindowMinutes *int     `json:"window_minutes"`
	ResetsAt      *int64   `json:"resets_at"`
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

func main() {
	var p Payload
	if data, err := readAll(os.Stdin); err == nil {
		_ = json.Unmarshal(data, &p) // best effort: a decode error leaves p zero
	}
	home, _ := os.UserHomeDir()
	codex := codexRateLimits(home)
	fmt.Println(render(p, codex, time.Now().Unix()))
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

// render builds the whole status line. Pure given its inputs (now is passed, not
// read from the clock) so every branch is testable.
func render(p Payload, codex *CodexRL, now int64) string {
	context := renderContext(p)

	if segs := claudeSegments(p, now); len(segs) > 0 {
		context += " | " + glyphTimer + " " + strings.Join(segs, ", ")
	}
	if segs := codexSegments(codex, now); len(segs) > 0 {
		context += " | " + glyphTimer + " Codex " + strings.Join(segs, ", ")
	}
	return glyphModel + " " + p.Model.DisplayName + " |  " + context
}

// renderContext shows current context-window usage: "<used>/<limit>" with used
// colored by proximity to the limit, or "<used> tkns" when the window size is
// unknown (legacy payloads). The percentage itself is elided — the absolute pair
// and the color convey it.
func renderContext(p Payload) string {
	used, limit := ctxUsedLimit(p)
	usedH := humanize(used)
	if limit <= 0 {
		return cGreen + usedH + cReset + " tkns"
	}
	pctInt := 0
	if p.ContextWindow != nil && p.ContextWindow.UsedPercentage != nil {
		pctInt = int(*p.ContextWindow.UsedPercentage)
	} else {
		pctInt = int(used * 100 / limit)
	}
	// 30/60 keep the prior absolute marks (300K/600K on a 1M window) while scaling
	// to whatever context_window_size reports, so 200k and 1M models both work.
	return tcolor(pctInt, 30, 60) + usedH + cReset + "/" + humanize(limit)
}

// ctxUsedLimit returns input-only used tokens and the window size, matching
// /context exactly (input + cache_creation + cache_read, excluding output). It
// prefers the harness-computed figures and falls back to summing the last
// main-thread assistant turn in the transcript for pre-2.1.132 clients that lack
// context_window (limit stays 0, so only the token count shows).
func ctxUsedLimit(p Payload) (used, limit int64) {
	if c := p.ContextWindow; c != nil {
		if c.ContextWindowSize != nil {
			limit = *c.ContextWindowSize
		}
		if c.TotalInputTokens != nil {
			used = *c.TotalInputTokens
		} else {
			used = c.CurrentUsage.InputTokens +
				c.CurrentUsage.CacheCreationInputTokens +
				c.CurrentUsage.CacheReadInputTokens
		}
		return used, limit
	}
	return transcriptUsed(p.TranscriptPath), 0
}

// transcriptUsed sums the input-only tokens of the last main-thread assistant
// turn in a transcript file. Sidechain (sub-agent) turns are skipped so the count
// reflects the main conversation. Returns 0 on any absence or error.
func transcriptUsed(path string) int64 {
	if path == "" {
		return 0
	}
	f, err := os.Open(path)
	if err != nil {
		return 0
	}
	defer f.Close()

	var last int64
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 64*1024), 16*1024*1024)
	for sc.Scan() {
		var e struct {
			Type        string `json:"type"`
			IsSidechain bool   `json:"isSidechain"`
			Message     struct {
				Usage *struct {
					InputTokens              int64 `json:"input_tokens"`
					CacheCreationInputTokens int64 `json:"cache_creation_input_tokens"`
					CacheReadInputTokens     int64 `json:"cache_read_input_tokens"`
				} `json:"usage"`
			} `json:"message"`
		}
		if json.Unmarshal(sc.Bytes(), &e) != nil {
			continue
		}
		if e.Type == "assistant" && !e.IsSidechain && e.Message.Usage != nil {
			u := e.Message.Usage
			last = u.InputTokens + u.CacheCreationInputTokens + u.CacheReadInputTokens
		}
	}
	return last
}

// claudeSegments renders Anthropic's 5h and 7d windows, skipping any that is
// absent (non-subscribers, or before the first API response of a session).
func claudeSegments(p Payload, now int64) []string {
	var out []string
	if p.RateLimits == nil {
		return out
	}
	if w := p.RateLimits.FiveHour; w != nil {
		if s := rlSegment("5h", fiveHourSecs, w.UsedPercentage, w.ResetsAt, now); s != "" {
			out = append(out, s)
		}
	}
	if w := p.RateLimits.SevenDay; w != nil {
		if s := rlSegment("7d", sevenDaySecs, w.UsedPercentage, w.ResetsAt, now); s != "" {
			out = append(out, s)
		}
	}
	return out
}

// codexSegments renders Codex's primary (5h) and secondary (7d) windows.
func codexSegments(rl *CodexRL, now int64) []string {
	var out []string
	if rl == nil {
		return out
	}
	for _, w := range []*CodexWindow{rl.Primary, rl.Secondary} {
		if w == nil || w.UsedPercent == nil || w.WindowMinutes == nil {
			continue
		}
		if s := rlSegment(winLabel(*w.WindowMinutes), int64(*w.WindowMinutes)*60, w.UsedPercent, w.ResetsAt, now); s != "" {
			out = append(out, s)
		}
	}
	return out
}

// rlSegment renders one rate-limit window as "<pct>% <elapsed>/<label>", e.g.
// "28% 3h28m/5h". Percent is colored by usage and elapsed by time-through-window,
// so the two colors side by side reveal burn pace: usage redder than time means
// burning fast. Returns "" when pct is absent.
//
// resetsAt drives elapsed = window - (resetsAt - now). Once it has passed
// (now >= resetsAt) the window has rolled over and any carried percent describes
// an expired window — a stale snapshot (Codex records them only at each API
// call), shown dimmed as "~<label>" rather than as a wrong number. Live sources
// (Anthropic, injected fresh each render) keep resetsAt in the future, so that
// branch never fires for them. When resetsAt is absent, only "<pct>% <label>".
func rlSegment(label string, windowSecs int64, pct *float64, resetsAt *int64, now int64) string {
	if pct == nil {
		return ""
	}
	if resetsAt != nil && now >= *resetsAt {
		return cGray + "~" + label + cReset
	}
	pctInt := int(*pct)
	pc := tcolor(pctInt, 50, 80)
	if resetsAt == nil {
		return fmt.Sprintf("%s%d%%%s %s%s%s", pc, pctInt, cReset, cGray, label, cReset)
	}
	elapsed := windowSecs - (*resetsAt - now)
	if elapsed < 0 {
		elapsed = 0
	}
	tc := cGreen
	if windowSecs > 0 {
		tc = tcolor(int(elapsed*100/windowSecs), 50, 80)
	}
	return fmt.Sprintf("%s%d%%%s %s%s%s%s/%s%s",
		pc, pctInt, cReset, tc, fmtDur(elapsed), cReset, cGray, label, cReset)
}

// codexRateLimits returns the freshest Codex rate-limit snapshot, or nil when
// Codex is not installed or has no session data. It reads the newest rollout
// files by mtime and picks the entry with the newest event timestamp, so a
// just-started session that has not hit the API yet does not hide older-but-real
// data, and concurrent sessions resolve to the genuinely latest event.
func codexRateLimits(home string) *CodexRL {
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

// humanize renders a token count compactly: 1234 -> "1.2K", 1000000 -> "1M",
// 1250000 -> "1.25M". Trailing zeros and a bare decimal point are trimmed to save
// width.
func humanize(n int64) string {
	switch {
	case n >= 1_000_000:
		return trimZeros(float64(n)/1_000_000, 2) + "M"
	case n >= 1000:
		return trimZeros(float64(n)/1000, 1) + "K"
	default:
		return strconv.FormatInt(n, 10)
	}
}

// trimZeros formats v with at most scale decimals, dropping trailing zeros and a
// trailing point (2.50 -> "2.5", 1.00 -> "1").
func trimZeros(v float64, scale int) string {
	s := strconv.FormatFloat(v, 'f', scale, 64)
	if strings.Contains(s, ".") {
		s = strings.TrimRight(s, "0")
		s = strings.TrimRight(s, ".")
	}
	return s
}

// tcolor maps a 0-100 progress value to green below lo, yellow below hi, red at
// or above hi.
func tcolor(v, lo, hi int) string {
	switch {
	case v >= hi:
		return cRed
	case v >= lo:
		return cYellow
	default:
		return cGreen
	}
}

// fmtDur renders a duration as "3d23h28m" / "3h28m" / "47m", dropping leading
// zero units. Minutes are kept even at day scale (this is elapsed within a
// rate-limit window, not a coarse countdown).
func fmtDur(s int64) string {
	if s < 0 {
		s = 0
	}
	d, h, m := s/86400, (s%86400)/3600, (s%3600)/60
	switch {
	case d > 0:
		return fmt.Sprintf("%dd%dh%dm", d, h, m)
	case h > 0:
		return fmt.Sprintf("%dh%dm", h, m)
	default:
		return fmt.Sprintf("%dm", m)
	}
}

// winLabel renders a window length in minutes as a compact label: 300 -> "5h",
// 10080 -> "7d", 90 -> "90m".
func winLabel(min int) string {
	switch {
	case min > 0 && min%1440 == 0:
		return fmt.Sprintf("%dd", min/1440)
	case min > 0 && min%60 == 0:
		return fmt.Sprintf("%dh", min/60)
	default:
		return fmt.Sprintf("%dm", min)
	}
}
