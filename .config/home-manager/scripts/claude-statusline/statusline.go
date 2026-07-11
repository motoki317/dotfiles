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
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
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

// Nerd-font glyphs: the model marker, the token glyph prefixing the context
// count, and the timer-sand that prefixes each rate-limit group. glyphToken lives
// in the Private Use Area, where it renders blank in fonts without the nerd-font
// patch — hence the \U escape, so it can never be mistaken for whitespace and
// silently dropped (as it was when this was ported from bash).
const (
	glyphModel = "\U000F06A9" // 󰚩 md-robot
	glyphToken = "\U0000EDE8" // token/context marker
	glyphTimer = "\U000F051F" // 󰔟 md-timer-sand
	glyphCost  = "\U000F0114" // 󰄔 md-cash
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
	now := time.Now().Unix()
	cost := sessionCost(sessionID(p), now)
	fmt.Println(render(p, codex, cost, now))
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
func render(p Payload, codex *CodexRL, cost *costInfo, now int64) string {
	context := renderContext(p)

	if s := renderCost(cost); s != "" {
		context += " | " + s
	}
	if segs := claudeSegments(p, now); len(segs) > 0 {
		context += " | " + glyphTimer + " " + strings.Join(segs, ", ")
	}
	if segs := codexSegments(codex, now); len(segs) > 0 {
		context += " | " + glyphTimer + " Codex " + strings.Join(segs, ", ")
	}
	return glyphModel + " " + modelName(p.Model.DisplayName) + " | " + glyphToken + " " + context
}

// modelName strips a trailing context-window annotation from the harness display
// name: on context-variant models the harness appends " (1M context)" (or
// "(200K context)"), which duplicates the window size already shown in the
// context segment. Only a trailing parenthetical mentioning "context" is
// removed, so a meaningful suffix like "(New)" is preserved.
func modelName(display string) string {
	if i := strings.LastIndex(display, " ("); i >= 0 && strings.HasSuffix(display, ")") {
		inner := display[i+2 : len(display)-1]
		if strings.Contains(strings.ToLower(inner), "context") {
			return display[:i]
		}
	}
	return display
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
	var pctInt int64
	if p.ContextWindow != nil && p.ContextWindow.UsedPercentage != nil {
		pctInt = int64(*p.ContextWindow.UsedPercentage)
	} else {
		pctInt = used * 100 / limit
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
// degrades to no cost segment, never a fatal error. --offline uses ccusage's
// bundled pricing (no network); --since bounds the scan (see costScanDays).
func sessionCost(sid string, now int64) *costInfo {
	if sid == "" {
		return nil
	}
	bin, err := exec.LookPath("ccusage")
	if err != nil {
		return nil
	}
	since := time.Unix(now, 0).AddDate(0, 0, -costScanDays).Format("20060102")
	ctx, cancel := context.WithTimeout(context.Background(), costTimeout)
	defer cancel()
	out, err := exec.CommandContext(ctx, bin, "session", "--json", "--offline", "--since", since).Output()
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
	cost := paint(tcolor(int64(ci.USD), 10, 40), "$"+sigFig(ci.USD, 3))
	read := paint(tcolor(ci.CacheRead, 10_000_000, 50_000_000), humanize(ci.CacheRead))
	create := paint(tcolor(ci.CacheCreation, 400_000, 2_000_000), humanize(ci.CacheCreation))
	input := paint(tcolor(ci.Input, 20_000, 200_000), humanize(ci.Input))
	output := paint(tcolor(ci.Output, 100_000, 350_000), humanize(ci.Output))
	return glyphCost + " " + cost + " " +
		paint(cGray, "↑") + read + paint(cGray, "/") + create + paint(cGray, "/") + input + " " +
		paint(cGray, "↓") + output
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
	pctInt := int64(*pct)
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
		tc = tcolor(elapsed*100/windowSecs, 50, 80)
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

// humanize renders a token count with two significant figures, keeping a third
// for hundreds: 1234->"1.2K", 12340->"12K", 123400->"123K", 1.26e6->"1.3M".
// Counts under 1000 print verbatim.
func humanize(n int64) string {
	if n < 1000 {
		return strconv.FormatInt(n, 10)
	}
	f, exp := float64(n), 0
	for f >= 1000 {
		f /= 1000
		exp++
	}
	s := sigFig(f, 2)
	if s == "1000" { // mantissa rounded up into the next unit (e.g. 999600 -> 1M)
		s, exp = "1", exp+1
	}
	const units = "KMBT"
	if exp > len(units) { // >= 1e15; unreachable for token counts, but never index out of range
		exp = len(units)
	}
	return s + units[exp-1:exp]
}

// sigFig renders v with n significant figures, trailing zeros trimmed. Integer
// digits count toward n; a value below 1 gets n fractional digits — leading zeros
// aren't significant, and trimZeros drops the rest ($0.05, not $0.050). Tokens use
// n=2 for compactness; cost uses n=3, where the extra digit's precision earns its
// width. Values with more than n integer digits print whole (127, not 130).
func sigFig(v float64, n int) string {
	intDigits := 0
	for x := v; x >= 1; x /= 10 {
		intDigits++
	}
	dec := n - intDigits
	if dec < 0 {
		dec = 0
	}
	return trimZeros(v, dec)
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

// trueColor reports whether the terminal advertises 24-bit color, resolved once at
// startup. When it does, tcolor emits a smooth green→yellow→red gradient; when it
// doesn't, the 3-step ANSI palette — so color loses fidelity, never vanishes.
var trueColor = supportsTrueColor(os.Getenv)

// supportsTrueColor detects 24-bit terminals. COLORTERM=truecolor/24bit is the
// conventional signal, but Warp on WSL sets no COLORTERM, so the two terminals in
// use are also allowlisted by TERM_PROGRAM. Unknown terminals degrade to the 3-step
// palette rather than risk garbled escapes.
func supportsTrueColor(env func(string) string) bool {
	switch env("COLORTERM") {
	case "truecolor", "24bit":
		return true
	}
	switch env("TERM_PROGRAM") {
	case "WarpTerminal", "iTerm.app":
		return true
	}
	return false
}

// rgb is a 24-bit color. escape emits its ANSI truecolor foreground sequence, which
// cReset ("\x1b[0m") clears like any other SGR, so gradient spans concatenate safely.
type rgb struct{ r, g, b uint8 }

func (c rgb) escape() string { return fmt.Sprintf("\x1b[38;2;%d;%d;%dm", c.r, c.g, c.b) }

// mix linearly interpolates from c toward o by t in [0,1], rounded to nearest.
func (c rgb) mix(o rgb, t float64) rgb {
	return rgb{lerp8(c.r, o.r, t), lerp8(c.g, o.g, t), lerp8(c.b, o.b, t)}
}

func lerp8(a, b uint8, t float64) uint8 { return uint8(float64(a) + t*(float64(b)-float64(a)) + 0.5) }

// Gradient stops mirroring the 3-step palette's meaning: green well under the low
// mark, yellow at the midpoint, red at or above the high mark.
var (
	gradLow  = rgb{80, 200, 120}
	gradMid  = rgb{235, 200, 70}
	gradHigh = rgb{235, 90, 80}
)

// gradient maps v across [lo,hi] to a green→yellow→red 24-bit color. The two linear
// segments (green→yellow, then yellow→red) keep the midpoint a clean yellow instead
// of the muddy brown a direct green→red interpolation would pass through.
func gradient(v, lo, hi int64) string {
	t := norm(v, lo, hi)
	if t <= 0.5 {
		return gradLow.mix(gradMid, t*2).escape()
	}
	return gradMid.mix(gradHigh, (t-0.5)*2).escape()
}

// norm clamps v into [lo,hi] and returns its position there in [0,1]; a degenerate
// lo>=hi maps to 1 (the high color), matching tcolor's v>=hi branch.
func norm(v, lo, hi int64) float64 {
	switch {
	case hi <= lo || v >= hi:
		return 1
	case v <= lo:
		return 0
	default:
		return float64(v-lo) / float64(hi-lo)
	}
}

// tcolor maps a value to a color by proximity to the [lo,hi] band: a smooth
// green→yellow→red gradient on 24-bit terminals, else green below lo / yellow below
// hi / red at or above hi. Used both for 0-100 progress and token/cost magnitudes.
func tcolor(v, lo, hi int64) string {
	if trueColor {
		return gradient(v, lo, hi)
	}
	switch {
	case v >= hi:
		return cRed
	case v >= lo:
		return cYellow
	default:
		return cGreen
	}
}

// paint wraps s in an ANSI color and a reset, so colored spans concatenate safely.
func paint(color, s string) string { return color + s + cReset }

// fmtDur renders a duration as "3d23h" / "3h28m" / "47m", dropping leading zero
// units and, at day scale, the minutes — a multi-day window (the 7d rate limit)
// doesn't need minute precision, and dropping them keeps the line compact.
func fmtDur(s int64) string {
	if s < 0 {
		s = 0
	}
	d, h, m := s/86400, (s%86400)/3600, (s%3600)/60
	switch {
	case d > 0:
		return fmt.Sprintf("%dd%dh", d, h)
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
