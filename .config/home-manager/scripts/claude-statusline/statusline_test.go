package main

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func pf(f float64) *float64 { return &f }
func pi(i int64) *int64     { return &i }
func pn(i int) *int         { return &i }

// TestMain pins the 3-step palette so the many exact color-escape assertions below
// are independent of the developer's terminal (which may set TERM_PROGRAM). Tests
// for the truecolor path toggle trueColor themselves or call gradient directly.
func TestMain(m *testing.M) {
	trueColor = false
	os.Exit(m.Run())
}

func TestHumanize(t *testing.T) {
	cases := []struct {
		in   int64
		want string
	}{
		{0, "0"},
		{999, "999"},
		{1000, "1K"},
		{1234, "1.2K"},
		{12340, "12K"},   // 2 sig figs drops the trailing .3
		{123400, "123K"}, // hundreds digit kept whole (the "3 as exception" case)
		{178289, "178K"},
		{1000000, "1M"},
		{1230000, "1.2M"},
		{1260000, "1.3M"},
		{1500000, "1.5M"},
		{999600, "1M"},          // mantissa rounds up into the next unit
		{1_200_000_000, "1.2B"}, // billions covered
	}
	for _, c := range cases {
		if got := humanize(c.in); got != c.want {
			t.Errorf("humanize(%d) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestSigFig(t *testing.T) {
	cases := []struct {
		v    float64
		n    int
		want string
	}{
		{2.73, 3, "2.73"},  // cost: 3 sig figs
		{12.34, 3, "12.3"}, // one integer digit consumed, two decimals -> one
		{127.4, 3, "127"},  // three integer digits -> whole
		{0.05, 3, "0.05"},  // sub-dollar: trailing zero trimmed, not "0.050"
		{2.7, 2, "2.7"},    // token mantissa: 2 sig figs
		{12.34, 2, "12"},   //
		{178.9, 2, "179"},  // hundreds keep the third digit (integer, rounded)
		{1.0, 2, "1"},      // trailing zeros trimmed
	}
	for _, c := range cases {
		if got := sigFig(c.v, c.n); got != c.want {
			t.Errorf("sigFig(%v,%d) = %q, want %q", c.v, c.n, got, c.want)
		}
	}
}

func TestTrimZeros(t *testing.T) {
	cases := []struct {
		v     float64
		scale int
		want  string
	}{
		{1.0, 2, "1"},
		{2.50, 2, "2.5"},
		{178.2, 1, "178.2"},
		{1.25, 2, "1.25"},
	}
	for _, c := range cases {
		if got := trimZeros(c.v, c.scale); got != c.want {
			t.Errorf("trimZeros(%v,%d) = %q, want %q", c.v, c.scale, got, c.want)
		}
	}
}

func TestTcolor(t *testing.T) {
	cases := []struct {
		v      int64
		lo, hi int64
		want   string
	}{
		{0, 50, 80, cGreen},
		{49, 50, 80, cGreen},
		{50, 50, 80, cYellow},
		{79, 50, 80, cYellow},
		{80, 50, 80, cRed},
		{100, 50, 80, cRed},
	}
	for _, c := range cases {
		if got := tcolor(c.v, c.lo, c.hi); got != c.want {
			t.Errorf("tcolor(%d,%d,%d) = %q, want %q", c.v, c.lo, c.hi, got, c.want)
		}
	}
}

func TestSupportsTrueColor(t *testing.T) {
	cases := []struct {
		colorterm, termProgram string
		want                   bool
	}{
		{"truecolor", "", true},
		{"24bit", "", true},
		{"", "WarpTerminal", true}, // Warp on WSL sets no COLORTERM
		{"", "iTerm.app", true},
		{"", "", false},
		{"256color", "Apple_Terminal", false}, // 256-color terminal, not 24-bit
	}
	for _, c := range cases {
		env := func(k string) string {
			switch k {
			case "COLORTERM":
				return c.colorterm
			case "TERM_PROGRAM":
				return c.termProgram
			}
			return ""
		}
		if got := supportsTrueColor(env); got != c.want {
			t.Errorf("supportsTrueColor(CT=%q,TP=%q) = %v, want %v", c.colorterm, c.termProgram, got, c.want)
		}
	}
}

func TestVisWidth(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want int
	}{
		{"plain ASCII", "hello", 5},
		{"colored text", cGreen + "green" + cReset, 5},
		{"supplementary PUA icon is 2 cells", glyphModel, 2},
		{"BMP PUA icon is 2 cells", glyphToken, 2},
		{"ambiguous arrows are 2 cells each", "↑↓", 4},
		{"glyph plus ASCII", glyphTimer + " 5h", 5}, // 2 + 1 + 1 + 1
		{"truecolor escape", "\x1b[38;2;12;34;56mcolor" + cReset, 5},
		{"unterminated escape", "\x1b[31", 4}, // no final byte: the bytes count as visible
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := visWidth(c.in); got != c.want {
				t.Errorf("visWidth(%q) = %d, want %d", c.in, got, c.want)
			}
		})
	}
}

// TestLayoutWrapsGlyphHeavyLine reproduces the reported overflow: a glyph-heavy
// line whose real width (100 cells) exceeds a 97-column pane. Before runeCells
// taught visWidth that the PUA icons and ↑/↓ are two cells, this line measured 93,
// so layout left it on one row and Claude Code truncated it with an ellipsis. It
// must now measure over the pane width and wrap, with every row inside the pane.
func TestLayoutWrapsGlyphHeavyLine(t *testing.T) {
	segments := []string{
		glyphModel + " Opus 4.8",
		glyphToken + " 222K/1M",
		glyphCost + " $9.9 ↑6.9M/429K/78 ↓86K",
		glyphTimer + " ~5h, 52% 5d1h/7d",
		glyphTimer + " Codex 18% 3h7m/5h",
	}
	const pane = 97
	if w := visWidth(strings.Join(segments, " | ")); w <= pane {
		t.Fatalf("joined width = %d, want > %d so the line must wrap", w, pane)
	}
	rows := strings.Split(layout(segments, pane), "\n")
	if len(rows) < 2 {
		t.Errorf("layout did not wrap: got %d row(s)", len(rows))
	}
	for _, row := range rows {
		if visWidth(row) > pane {
			t.Errorf("row %q has width %d, want <= %d", row, visWidth(row), pane)
		}
	}
}

func TestLayout(t *testing.T) {
	cases := []struct {
		name     string
		segments []string
		width    int
		want     string
		rowsFit  bool
	}{
		{"unknown width", []string{"one", "two"}, 0, "one | two", false},
		{"wide enough", []string{"one", "two"}, 9, "one | two", true},
		{"one over", []string{"one", "two"}, 8, "one\ntwo", true},
		{"wraps greedily", []string{"aa", "bb", "ccc"}, 8, "aa | bb\nccc", true},
		{"packs after wrap", []string{"aa", "bb", "ccc", "d"}, 8, "aa | bb\nccc | d", true},
		{"measures visible width", []string{cGreen + "aaa" + cReset, "bbb"}, 9, cGreen + "aaa" + cReset + " | bbb", true},
		{"over-wide segment", []string{"aa", "toolong", "b"}, 4, "aa\ntoolong\nb", false},
		{"empty", nil, 40, "", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := layout(c.segments, c.width)
			if got != c.want {
				t.Errorf("layout(%q, %d) = %q, want %q", c.segments, c.width, got, c.want)
			}
			if c.rowsFit {
				for _, row := range strings.Split(got, "\n") {
					if visWidth(row) > c.width {
						t.Errorf("row %q has width %d, want <= %d", row, visWidth(row), c.width)
					}
				}
			}
		})
	}
}

func TestTermWidth(t *testing.T) {
	cases := []struct {
		name  string
		value string
		unset bool
		want  int
	}{
		{"unset", "", true, 0},
		{"positive", "94", false, 94},
		{"zero", "0", false, 0},
		{"non-numeric", "abc", false, 0},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Setenv("COLUMNS", c.value)
			if c.unset {
				if err := os.Unsetenv("COLUMNS"); err != nil {
					t.Fatal(err)
				}
			}
			if got := termWidth(); got != c.want {
				t.Errorf("termWidth() with COLUMNS=%q = %d, want %d", c.value, got, c.want)
			}
		})
	}
}

func TestGradient(t *testing.T) {
	const lo, hi = int64(0), int64(100)
	// Endpoints and midpoint land exactly on the palette stops.
	if got := gradient(lo, lo, hi); got != gradLow.escape() {
		t.Errorf("low end = %q, want green stop", got)
	}
	if got := gradient(50, lo, hi); got != gradMid.escape() {
		t.Errorf("midpoint = %q, want yellow stop", got)
	}
	if got := gradient(hi, lo, hi); got != gradHigh.escape() {
		t.Errorf("high end = %q, want red stop", got)
	}
	// Out-of-band values clamp to the endpoint colors.
	if gradient(-10, lo, hi) != gradLow.escape() {
		t.Error("below lo should clamp to the green stop")
	}
	if gradient(999, lo, hi) != gradHigh.escape() {
		t.Error("above hi should clamp to the red stop")
	}
	// Every output is a 24-bit foreground escape.
	if !strings.HasPrefix(gradient(25, lo, hi), "\x1b[38;2;") {
		t.Errorf("expected a 24-bit escape, got %q", gradient(25, lo, hi))
	}
}

func TestTcolorTrueColor(t *testing.T) {
	trueColor = true
	defer func() { trueColor = false }()
	if got := tcolor(75, 50, 80); !strings.HasPrefix(got, "\x1b[38;2;") {
		t.Errorf("truecolor tcolor = %q, want a 24-bit escape", got)
	}
}

func TestFmtDur(t *testing.T) {
	cases := []struct {
		s    int64
		want string
	}{
		{-5, "0m"},
		{0, "0m"},
		{47 * 60, "47m"},
		{3*3600 + 28*60, "3h28m"},
		{3*86400 + 23*3600 + 28*60, "3d23h"}, // minutes dropped at day scale
		{1*86400 + 21*3600 + 24*60, "1d21h"}, // 1d21h24m -> 1d21h
		{1*86400 + 0*3600 + 24*60, "1d0h"},   // hours kept even at zero
	}
	for _, c := range cases {
		if got := fmtDur(c.s); got != c.want {
			t.Errorf("fmtDur(%d) = %q, want %q", c.s, got, c.want)
		}
	}
}

func TestWinLabel(t *testing.T) {
	cases := []struct {
		min  int
		want string
	}{
		{300, "5h"},
		{10080, "7d"},
		{60, "1h"},
		{1440, "1d"},
		{90, "90m"},
	}
	for _, c := range cases {
		if got := winLabel(c.min); got != c.want {
			t.Errorf("winLabel(%d) = %q, want %q", c.min, got, c.want)
		}
	}
}

func TestRlSegment(t *testing.T) {
	const now int64 = 1_000_000

	// Fresh: 28% used (green), 13000s/18000s = 72% elapsed (yellow), "3h36m/5h".
	fresh := rlSegment("5h", 18000, pf(28), pi(now+5000), now)
	want := cGreen + "28%" + cReset + " " + cYellow + "3h36m" + cReset + cGray + "/5h" + cReset
	if fresh != want {
		t.Errorf("fresh:\n got %q\nwant %q", fresh, want)
	}

	// Rolled-over snapshot (resets in the past) -> dimmed "~5h", no number.
	if got := rlSegment("5h", 18000, pf(93), pi(now-1), now); got != cGray+"~5h"+cReset {
		t.Errorf("stale = %q, want %q", got, cGray+"~5h"+cReset)
	}

	// No resetsAt -> "<pct>% <label>" fallback.
	if got := rlSegment("5h", 18000, pf(28), nil, now); got != cGreen+"28%"+cReset+" "+cGray+"5h"+cReset {
		t.Errorf("no-resets = %q", got)
	}

	// Absent percent -> empty (skipped by caller).
	if got := rlSegment("5h", 18000, nil, pi(now+5000), now); got != "" {
		t.Errorf("nil pct = %q, want empty", got)
	}

	// Fractional percent truncates like the old ${pct%%.*}.
	if got := rlSegment("7d", 604800, pf(15.9), nil, now); !strings.Contains(got, "15%") {
		t.Errorf("fractional pct: got %q, want it to contain 15%%", got)
	}
}

func TestModelName(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"Opus 4.8 (1M context)", "Opus 4.8"},
		{"Sonnet 4.5 (200K context)", "Sonnet 4.5"},
		{"Opus 4.8", "Opus 4.8"},                               // no annotation, untouched
		{"Claude 3.5 Sonnet (New)", "Claude 3.5 Sonnet (New)"}, // non-context parenthetical preserved
		{"", ""},
	}
	for _, c := range cases {
		if got := modelName(c.in); got != c.want {
			t.Errorf("modelName(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestRenderContextElidesPercent(t *testing.T) {
	var p Payload
	p.ContextWindow = &ContextWindow{
		TotalInputTokens:  pi(178200),
		ContextWindowSize: pi(1_000_000),
		UsedPercentage:    pf(17.8),
	}
	got := renderContext(p)
	if strings.Contains(got, "%") {
		t.Errorf("context should elide the percentage, got %q", got)
	}
	if !strings.Contains(got, "178K") || !strings.Contains(got, "/1M") {
		t.Errorf("context = %q, want 178K/1M", got)
	}
}

func TestRenderContextLegacyNoLimit(t *testing.T) {
	var p Payload // no ContextWindow, no transcript
	got := renderContext(p)
	if !strings.Contains(got, "tkns") {
		t.Errorf("legacy context = %q, want a 'tkns' suffix", got)
	}
}

func TestRenderFull(t *testing.T) {
	const now int64 = 1_000_000
	var p Payload
	p.Model.DisplayName = "Opus 4.8"
	p.ContextWindow = &ContextWindow{
		TotalInputTokens:  pi(178200),
		ContextWindowSize: pi(1_000_000),
		UsedPercentage:    pf(17.8),
	}

	codex := &CodexRL{
		Primary:   &CodexWindow{UsedPercent: pf(94), WindowMinutes: pn(300), ResetsAt: pi(now + 3000)},
		Secondary: &CodexWindow{UsedPercent: pf(16), WindowMinutes: pn(10080), ResetsAt: pi(now + 500000)},
	}

	got := render(p, codex, nil, now, 0)
	// A color reset sits between "178.2K" and "/1M", so check the parts separately.
	// The token glyph must prefix the context (regression: it was dropped once).
	for _, want := range []string{glyphModel + " Opus 4.8", " | " + glyphToken + " ", "178K", "/1M", glyphTimer + " Codex", "94%", "16%"} {
		if !strings.Contains(got, want) {
			t.Errorf("render missing %q in:\n%s", want, got)
		}
	}
}

func TestRenderCodexAbsent(t *testing.T) {
	var p Payload
	p.Model.DisplayName = "Opus 4.8"
	got := render(p, nil, nil, 1_000_000, 0)
	if strings.Contains(got, "Codex") {
		t.Errorf("no Codex data should render no Codex group, got %q", got)
	}
}

func TestSessionID(t *testing.T) {
	// Explicit field wins.
	if got := sessionID(Payload{SessionID: "abc"}); got != "abc" {
		t.Errorf("explicit = %q, want abc", got)
	}
	// Falls back to the transcript filename stem.
	p := Payload{TranscriptPath: "/home/u/.claude/projects/x/b21a0283-29b2.jsonl"}
	if got := sessionID(p); got != "b21a0283-29b2" {
		t.Errorf("fallback = %q, want b21a0283-29b2", got)
	}
}

func TestParseCost(t *testing.T) {
	data := []byte(`{"session":{
		"tokens":{"uncached_input":7817,"output":40734,"cache_write":78295,"cache_read":1770917},
		"cost":{"usd":2.7258,"complete":true,"estimated":false,"missing_pricing":[]}
	}}`)
	ci := parseCost(data)
	if ci == nil {
		t.Fatal("expected a session summary")
	}
	if ci.USD != 2.7258 || ci.Input != 7817 || ci.Output != 40734 ||
		ci.CacheCreation != 78295 || ci.CacheRead != 1770917 || !ci.Complete {
		t.Errorf("parsed = %+v", ci)
	}
	// A structured not-found payload must decode as absence, never zero cost.
	if parseCost([]byte(`{"error":{"code":"not_found","message":"session not found"}}`)) != nil {
		t.Error("not-found response should be nil")
	}
	if parseCost([]byte("not json")) != nil {
		t.Error("malformed input should be nil")
	}
}

func TestAgtlogCostCommand(t *testing.T) {
	dir := t.TempDir()
	bin := filepath.Join(dir, "agtlog")
	argsPath := filepath.Join(dir, "args")
	fake := `#!/bin/sh
printf 'call\n' >> "$AGTLOG_TEST_ARGS"
printf '%s\n' "$@" >> "$AGTLOG_TEST_ARGS"
printf '%s\n' '{"session":{"tokens":{},"cost":{"usd":1.25,"complete":true}}}'
`
	if err := os.WriteFile(bin, []byte(fake), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("AGTLOG_TEST_ARGS", argsPath)
	t.Setenv("PATH", dir)
	ci := sessionCost("00000000-0000-4000-8000-000000000001")
	if ci == nil || ci.USD != 1.25 || !ci.Complete {
		t.Fatalf("agtlogCost = %+v", ci)
	}
	args, err := os.ReadFile(argsPath)
	if err != nil {
		t.Fatal(err)
	}
	want := "call\nshow\nclaude:00000000-0000-4000-8000-000000000001\n--no-events\n--offline\n"
	if string(args) != want {
		t.Errorf("agtlog arguments = %q, want %q", args, want)
	}
}

func TestAgtlogCostFailures(t *testing.T) {
	dir := t.TempDir()
	bin := filepath.Join(dir, "agtlog")
	notFound := `#!/bin/sh
printf '%s\n' '{"error":{"code":"not_found"}}' >&2
exit 3
`
	if err := os.WriteFile(bin, []byte(notFound), 0o755); err != nil {
		t.Fatal(err)
	}
	if ci := agtlogCost(bin, "missing"); ci != nil {
		t.Errorf("exit 3 should return nil, got %+v", ci)
	}

	blocking := "#!/bin/sh\nexec /bin/sleep 10\n"
	if err := os.WriteFile(bin, []byte(blocking), 0o755); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	if ci := agtlogCostContext(ctx, bin, "slow"); ci != nil {
		t.Errorf("timeout should return nil, got %+v", ci)
	}
	if costTimeout != 2*time.Second {
		t.Errorf("costTimeout = %v, want 2s", costTimeout)
	}
}

func TestRenderCost(t *testing.T) {
	if renderCost(nil) != "" {
		t.Error("nil cost should render empty")
	}
	got := renderCost(&costInfo{USD: 2.73, Complete: true, Input: 7817, CacheCreation: 78245, CacheRead: 1700000, Output: 40734})
	// Cost is 3-sig-fig ($2.73, precision matters); tokens are humanized at 2-sig-fig.
	// Order is cache-read, cache-creation, input, then output. Color codes break up
	// the ↑ group, so assert the pieces, not one contiguous substring.
	for _, want := range []string{glyphCost, "$2.73", "1.7M", "78K", "7.8K", "41K"} {
		if !strings.Contains(got, want) {
			t.Errorf("renderCost missing %q in %q", want, got)
		}
	}
	// ↑ group is ordered largest-first: cache-read precedes input.
	if strings.Index(got, "1.7M") > strings.Index(got, "7.8K") {
		t.Errorf("cache-read should precede input, got %q", got)
	}
	// Numbers are colored, not a uniform gray block.
	if !strings.Contains(got, cGreen) && !strings.Contains(got, cYellow) && !strings.Contains(got, cRed) {
		t.Errorf("expected colored numbers, got %q", got)
	}
}

func TestRenderCostUnpriced(t *testing.T) {
	// An incomplete agtlog total with no priced tokens is unknown, not a real $0.
	got := renderCost(&costInfo{USD: 0, Input: 7, CacheCreation: 39800, CacheRead: 148136, Output: 22715})
	if !strings.Contains(got, "$?") {
		t.Errorf("unpriced session should render $?, got %q", got)
	}
	if strings.Contains(got, "$0") {
		t.Errorf("should not render a misleading $0, got %q", got)
	}
	if !strings.Contains(got, "148K") {
		t.Errorf("tokens should still render, got %q", got)
	}
}

func TestRenderCostCompleteZero(t *testing.T) {
	got := renderCost(&costInfo{Complete: true})
	if !strings.Contains(got, "$0") || strings.Contains(got, "$?") || strings.Contains(got, "!") {
		t.Errorf("complete zero cost should render $0, got %q", got)
	}
}

func TestRenderCostPartial(t *testing.T) {
	got := renderCost(&costInfo{USD: 2.73})
	want := "$2.73" + cReset + cGray + "!" + cReset
	if !strings.Contains(got, want) {
		t.Errorf("partial cost should render a gray ! after the total, got %q", got)
	}
	if strings.Contains(got, "$?") {
		t.Errorf("partial nonzero cost should render its known total, got %q", got)
	}
}

func TestRenderWithCost(t *testing.T) {
	var p Payload
	p.Model.DisplayName = "Opus 4.8"
	got := render(p, nil, &costInfo{USD: 1.17, Complete: true, Input: 1000, CacheCreation: 2000, CacheRead: 3000, Output: 500}, 1_000_000, 0)
	if !strings.Contains(got, "$1.17") { // cost at 3 sig figs keeps the cents
		t.Errorf("render should include the cost segment, got %q", got)
	}
}

func TestRenderWraps(t *testing.T) {
	const now int64 = 1_000_000
	var p Payload
	p.Model.DisplayName = "Opus 4.8"
	p.ContextWindow = &ContextWindow{
		TotalInputTokens:  pi(178200),
		ContextWindowSize: pi(1_000_000),
		UsedPercentage:    pf(17.8),
	}
	p.RateLimits = &RateLimits{
		FiveHour: &ClaudeWindow{UsedPercentage: pf(28), ResetsAt: pi(now + 5000)},
		SevenDay: &ClaudeWindow{UsedPercentage: pf(15), ResetsAt: pi(now + 500000)},
	}
	codex := &CodexRL{
		Primary:   &CodexWindow{UsedPercent: pf(94), WindowMinutes: pn(300), ResetsAt: pi(now + 3000)},
		Secondary: &CodexWindow{UsedPercent: pf(16), WindowMinutes: pn(10080), ResetsAt: pi(now + 500000)},
	}
	cost := &costInfo{USD: 1.17, Complete: true, Input: 1000, Output: 500}

	narrow := render(p, codex, cost, now, 40)
	if !strings.Contains(narrow, "\n") {
		t.Fatalf("width 40 should wrap, got %q", narrow)
	}
	for _, row := range strings.Split(narrow, "\n") {
		if visWidth(row) > 40 {
			t.Errorf("row %q has width %d, want <= 40", row, visWidth(row))
		}
	}

	unwrapped := render(p, codex, cost, now, 0)
	wide := render(p, codex, cost, now, 200)
	if strings.Contains(unwrapped, "\n") || strings.Contains(wide, "\n") {
		t.Errorf("non-wrapping widths produced a newline: width 0=%q, width 200=%q", unwrapped, wide)
	}
	if wide != unwrapped {
		t.Errorf("wide render differs from legacy render:\nwide %q\nlegacy %q", wide, unwrapped)
	}
}

// writeLines writes JSONL to a temp rollout file and returns its path.
func writeLines(t *testing.T, lines ...string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "rollout-test.jsonl")
	if err := os.WriteFile(p, []byte(strings.Join(lines, "\n")+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestLastRateLimitEventStructural(t *testing.T) {
	// A prose line mentioning "rate_limits" both before AND after the real event,
	// plus a malformed line — none should be selected, and the scan must not abort.
	path := writeLines(t,
		`{"type":"event_msg","timestamp":"2026-07-10T07:00:00Z","payload":{"type":"agent_message","text":"tell me about rate_limits"}}`,
		`{"type":"event_msg","timestamp":"2026-07-10T07:40:00Z","payload":{"type":"token_count","rate_limits":{"primary":{"used_percent":93.0,"window_minutes":300,"resets_at":1783677171},"secondary":{"used_percent":15.0,"window_minutes":10080,"resets_at":1784242529}}}}`,
		`{"type":"event_msg","timestamp":"2026-07-10T07:41:00Z"`, // truncated / malformed
		`{"type":"response_item","payload":{"content":"more discussion of rate_limits here"}}`,
	)
	rl, ts, ok := lastRateLimitEvent(path)
	if !ok {
		t.Fatal("expected a rate_limits event to be found")
	}
	if rl.Primary == nil || rl.Primary.UsedPercent == nil || *rl.Primary.UsedPercent != 93.0 {
		t.Errorf("selected the wrong object: %+v", rl.Primary)
	}
	if ts.IsZero() {
		t.Error("timestamp should have parsed")
	}
}

func TestLastRateLimitEventPicksLast(t *testing.T) {
	path := writeLines(t,
		`{"type":"event_msg","timestamp":"2026-07-10T07:00:00Z","payload":{"type":"token_count","rate_limits":{"primary":{"used_percent":10.0,"window_minutes":300,"resets_at":1}}}}`,
		`{"type":"event_msg","timestamp":"2026-07-10T08:00:00Z","payload":{"type":"token_count","rate_limits":{"primary":{"used_percent":42.0,"window_minutes":300,"resets_at":2}}}}`,
	)
	rl, _, ok := lastRateLimitEvent(path)
	if !ok || rl.Primary == nil || *rl.Primary.UsedPercent != 42.0 {
		t.Errorf("expected the last event (42%%), got %+v", rl)
	}
}

func TestLastRateLimitEventNone(t *testing.T) {
	path := writeLines(t,
		`{"type":"event_msg","payload":{"type":"agent_message","text":"no limits here"}}`,
	)
	if _, _, ok := lastRateLimitEvent(path); ok {
		t.Error("expected no event")
	}
}

func TestCodexSegmentsMissingWindow(t *testing.T) {
	// Only the 5h window present; 7d absent -> one segment.
	rl := &CodexRL{Primary: &CodexWindow{UsedPercent: pf(42), WindowMinutes: pn(300), ResetsAt: pi(9_999_999_999)}}
	segs := codexSegments(rl, 1_000_000)
	if len(segs) != 1 {
		t.Fatalf("want 1 segment, got %d: %v", len(segs), segs)
	}
	if !strings.Contains(segs[0], "42%") || !strings.Contains(segs[0], "/5h") {
		t.Errorf("segment = %q", segs[0])
	}
}

func TestCodexRateLimitsNewestByEventTimestamp(t *testing.T) {
	// Two rollout files. The one written LAST (newer mtime) carries an OLDER event
	// timestamp; selection must prefer the newer event, not the newer file.
	home := t.TempDir()
	dir := filepath.Join(home, ".codex", "sessions", "2026", "07", "10")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	older := `{"type":"event_msg","timestamp":"2026-07-10T06:00:00Z","payload":{"type":"token_count","rate_limits":{"primary":{"used_percent":10.0,"window_minutes":300,"resets_at":9999999999}}}}`
	newer := `{"type":"event_msg","timestamp":"2026-07-10T09:00:00Z","payload":{"type":"token_count","rate_limits":{"primary":{"used_percent":80.0,"window_minutes":300,"resets_at":9999999999}}}}`
	// File A holds the newer EVENT; write it first.
	if err := os.WriteFile(filepath.Join(dir, "rollout-a.jsonl"), []byte(newer+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	// File B holds the older event; write it second so it has the newer mtime.
	if err := os.WriteFile(filepath.Join(dir, "rollout-b.jsonl"), []byte(older+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	rl := codexRateLimitsFromLogs(home)
	if rl == nil || rl.Primary == nil || *rl.Primary.UsedPercent != 80.0 {
		t.Errorf("expected the newer event (80%%), got %+v", rl)
	}
}

func TestCodexRateLimitsAbsent(t *testing.T) {
	if rl := codexRateLimitsFromLogs(t.TempDir()); rl != nil {
		t.Errorf("no ~/.codex should yield nil, got %+v", rl)
	}
}

func TestParseClaudeUsage(t *testing.T) {
	cases := []struct {
		name                  string
		data                  string
		fiveUsed, sevenUsed   *float64
		fiveReset, sevenReset *int64
	}{
		{
			"both windows",
			`{"five_hour":{"utilization":15,"resets_at":"2026-07-22T03:10:00.272933+00:00"},"seven_day":{"utilization":54,"resets_at":"2026-07-23T06:00:00.272960+00:00"}}`,
			pf(15), pf(54), pi(1784689800), pi(1784786400),
		},
		{
			"five hour only",
			`{"five_hour":{"utilization":23.5,"resets_at":"2026-07-22T03:10:00.272933+00:00"}}`,
			pf(23.5), nil, pi(1784689800), nil,
		},
		{
			"invalid resets remain optional",
			`{"five_hour":{"utilization":15,"resets_at":""},"seven_day":{"utilization":54,"resets_at":"garbage"}}`,
			pf(15), pf(54), nil, nil,
		},
		{"empty body", "", nil, nil, nil, nil},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			rl := parseClaudeUsage([]byte(c.data))
			if c.fiveUsed == nil && c.sevenUsed == nil {
				if rl != nil {
					t.Fatalf("parseClaudeUsage(%q) = %+v, want nil", c.data, rl)
				}
				return
			}
			if rl == nil {
				t.Fatal("parseClaudeUsage returned nil")
			}
			assertClaudeWindow(t, "five-hour", rl.FiveHour, c.fiveUsed, c.fiveReset)
			assertClaudeWindow(t, "seven-day", rl.SevenDay, c.sevenUsed, c.sevenReset)
		})
	}
}

func TestClaudeUsageRequiresUsableWindow(t *testing.T) {
	cases := []struct {
		name string
		rl   *RateLimits
		want bool
	}{
		{"nil", nil, false},
		{"empty", &RateLimits{}, false},
		{"reset only", &RateLimits{FiveHour: &ClaudeWindow{ResetsAt: pi(123)}}, false},
		{"five hour", &RateLimits{FiveHour: &ClaudeWindow{UsedPercentage: pf(12)}}, true},
		{"seven day", &RateLimits{SevenDay: &ClaudeWindow{UsedPercentage: pf(34)}}, true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := usableClaudeRateLimits(c.rl)
			if c.want && got != c.rl {
				t.Fatalf("usable rate limits = %+v, want original snapshot %+v", got, c.rl)
			}
			if !c.want && got != nil {
				t.Fatalf("unusable rate limits = %+v, want nil payload fallback", got)
			}
		})
	}
}

func assertClaudeWindow(t *testing.T, name string, got *ClaudeWindow, wantUsed *float64, wantReset *int64) {
	t.Helper()
	if wantUsed == nil {
		if got != nil {
			t.Errorf("%s = %+v, want nil", name, got)
		}
		return
	}
	if got == nil || got.UsedPercentage == nil || *got.UsedPercentage != *wantUsed {
		t.Fatalf("%s utilization mismatch: %+v", name, got)
	}
	if (got.ResetsAt == nil) != (wantReset == nil) ||
		got.ResetsAt != nil && *got.ResetsAt != *wantReset {
		t.Errorf("%s reset mismatch: %+v", name, got)
	}
}

func TestParseAppServerRateLimits(t *testing.T) {
	// Both windows present: camelCase field names map to CodexWindow.
	both := `{"rateLimits":{"primary":{"usedPercent":42,"windowDurationMins":300,"resetsAt":9999999999},
		"secondary":{"usedPercent":8,"windowDurationMins":10080,"resetsAt":9999999999}}}`
	rl := parseAppServerRateLimits([]byte(both))
	if rl == nil || rl.Primary == nil || *rl.Primary.UsedPercent != 42 || *rl.Primary.WindowMinutes != 300 {
		t.Fatalf("primary mismatch: %+v", rl)
	}
	if rl.Secondary == nil || *rl.Secondary.WindowMinutes != 10080 {
		t.Fatalf("secondary mismatch: %+v", rl)
	}

	// This account's shape: a single primary window, secondary null. Must not be dropped.
	primaryOnly := `{"rateLimits":{"primary":{"usedPercent":8,"windowDurationMins":10080,"resetsAt":1784608504},"secondary":null}}`
	rl = parseAppServerRateLimits([]byte(primaryOnly))
	if rl == nil || rl.Primary == nil || rl.Secondary != nil {
		t.Fatalf("primary-only mismatch: %+v", rl)
	}

	// No usable window -> nil, so the segment is omitted rather than rendered empty.
	for _, empty := range []string{`{"rateLimits":{"primary":null,"secondary":null}}`, `{}`, `not json`} {
		if rl := parseAppServerRateLimits([]byte(empty)); rl != nil {
			t.Errorf("parseAppServerRateLimits(%q) = %+v, want nil", empty, rl)
		}
	}
}
