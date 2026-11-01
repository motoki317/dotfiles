package main

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

func pi(i int64) *int64 { return &i }

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
		{12340, "12K"},
		{123400, "123K"},
		{178289, "178K"},
		{1000000, "1M"},
		{1230000, "1.2M"},
		{1260000, "1.3M"},
		{1500000, "1.5M"},
		{999600, "1M"},
		{1_200_000_000, "1.2B"},
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
		{2.73, 3, "2.73"},
		{12.34, 3, "12.3"},
		{127.4, 3, "127"},
		{0.05, 3, "0.05"},
		{2.7, 2, "2.7"},
		{12.34, 2, "12"},
		{178.9, 2, "179"},
		{1.0, 2, "1"},
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
		{"", "WarpTerminal", true},
		{"", "iTerm.app", true},
		{"", "", false},
		{"256color", "Apple_Terminal", false},
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
		{"glyph plus ASCII", glyphTimer + " 5h", 5},
		{"truecolor escape", "\x1b[38;2;12;34;56mcolor" + cReset, 5},
		{"unterminated escape", "\x1b[31", 4},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := visWidth(c.in); got != c.want {
				t.Errorf("visWidth(%q) = %d, want %d", c.in, got, c.want)
			}
		})
	}
}

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
	if got := gradient(lo, lo, hi); got != gradLow.escape() {
		t.Errorf("low end = %q, want green stop", got)
	}
	if got := gradient(50, lo, hi); got != gradMid.escape() {
		t.Errorf("midpoint = %q, want yellow stop", got)
	}
	if got := gradient(hi, lo, hi); got != gradHigh.escape() {
		t.Errorf("high end = %q, want red stop", got)
	}
	if gradient(-10, lo, hi) != gradLow.escape() {
		t.Error("below lo should clamp to the green stop")
	}
	if gradient(999, lo, hi) != gradHigh.escape() {
		t.Error("above hi should clamp to the red stop")
	}
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
		{3*86400 + 23*3600 + 28*60, "3d23h"},
		{1*86400 + 21*3600 + 24*60, "1d21h"},
		{1*86400 + 0*3600 + 24*60, "1d0h"},
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

	fresh := rlSegment(Window{Minutes: 300, Used: 28, ResetsAt: pi(now + 5000)}, now)
	want := cGreen + "28%" + cReset + " " + cYellow + "3h36m" + cReset + cGray + "/5h" + cReset
	if fresh != want {
		t.Errorf("fresh:\n got %q\nwant %q", fresh, want)
	}

	if got := rlSegment(Window{Minutes: 300, Used: 93, ResetsAt: pi(now - 1)}, now); got != cGray+"~5h"+cReset {
		t.Errorf("stale = %q, want %q", got, cGray+"~5h"+cReset)
	}

	if got := rlSegment(Window{Minutes: 300, Used: 28}, now); got != cGreen+"28%"+cReset+" "+cGray+"5h"+cReset {
		t.Errorf("no-resets = %q", got)
	}

	if got := rlSegment(Window{Minutes: 10080, Used: 15.9}, now); !strings.Contains(got, "15%") {
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
		{"Opus 4.8", "Opus 4.8"},
		{"Claude 3.5 Sonnet (New)", "Claude 3.5 Sonnet (New)"},
		{"", ""},
	}
	for _, c := range cases {
		if got := modelName(c.in); got != c.want {
			t.Errorf("modelName(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestRenderContextElidesPercent(t *testing.T) {
	c := &ContextWindow{
		TotalInputTokens:  178200,
		ContextWindowSize: 1_000_000,
	}
	got := renderContext(c)
	if strings.Contains(got, "%") {
		t.Errorf("context should elide the percentage, got %q", got)
	}
	if !strings.Contains(got, "178K") || !strings.Contains(got, "/1M") {
		t.Errorf("context = %q, want 178K/1M", got)
	}
}

func TestRenderContextPayload(t *testing.T) {
	cases := []struct {
		name, data, context string
	}{
		{"absent context", `{}`, ""},
		{"null context", `{"context_window":null}`, ""},
		{"before first response", `{"context_window":{"total_input_tokens":0,"context_window_size":200000}}`, cGreen + "0" + cReset + "/200K"},
		{"ignores current usage", `{"context_window":{"context_window_size":200000,"current_usage":{"input_tokens":1000,"cache_creation_input_tokens":2000,"cache_read_input_tokens":3000}}}`, cGreen + "0" + cReset + "/200K"},
		{"missing limit", `{"context_window":{"total_input_tokens":1000}}`, ""},
		{"zero limit", `{"context_window":{"total_input_tokens":1000,"context_window_size":0}}`, ""},
		{"negative limit", `{"context_window":{"total_input_tokens":1000,"context_window_size":-1}}`, ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var p Payload
			if err := json.Unmarshal([]byte(c.data), &p); err != nil {
				t.Fatal(err)
			}
			p.Model.DisplayName = "Opus 4.8"
			want := glyphModel + " Opus 4.8"
			if c.context != "" {
				want += " | " + glyphToken + " " + c.context
			}
			if got := render(p, nil, nil, nil, 1_000_000, 0); got != want {
				t.Errorf("render = %q, want %q", got, want)
			}
		})
	}
}

func TestRenderContextColor(t *testing.T) {
	cases := []struct {
		used, limit int64
		color       string
	}{
		{59999, 200000, cGreen},
		{60000, 200000, cYellow},
		{119999, 200000, cYellow},
		{120000, 200000, cRed},
		{299999, 1000000, cGreen},
		{300000, 1000000, cYellow},
		{600000, 1000000, cRed},
	}
	for _, c := range cases {
		context := &ContextWindow{TotalInputTokens: c.used, ContextWindowSize: c.limit}
		if got := renderContext(context); !strings.HasPrefix(got, c.color) {
			t.Errorf("context %d/%d = %q, want color %q", c.used, c.limit, got, c.color)
		}
	}
}

func TestRenderFull(t *testing.T) {
	const now int64 = 1_000_000
	var p Payload
	p.Model.DisplayName = "Opus 4.8"
	p.ContextWindow = &ContextWindow{
		TotalInputTokens:  178200,
		ContextWindowSize: 1_000_000,
	}

	codex := []Window{
		{Minutes: 300, Used: 94, ResetsAt: pi(now + 3000)},
		{Minutes: 10080, Used: 16, ResetsAt: pi(now + 500000)},
	}

	got := render(p, nil, codex, nil, now, 0)
	for _, want := range []string{glyphModel + " Opus 4.8", " | " + glyphToken + " ", "178K", "/1M", glyphTimer + " Codex", "94%", "16%"} {
		if !strings.Contains(got, want) {
			t.Errorf("render missing %q in:\n%s", want, got)
		}
	}
}

func TestRenderCodexAbsent(t *testing.T) {
	var p Payload
	p.Model.DisplayName = "Opus 4.8"
	got := render(p, nil, nil, nil, 1_000_000, 0)
	if strings.Contains(got, "Codex") {
		t.Errorf("no Codex data should render no Codex group, got %q", got)
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
	if ci := sessionCost(""); ci != nil {
		t.Fatalf("cost without session ID = %+v, want nil", ci)
	}
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
	if ci := agtlogCost(context.Background(), bin, "missing"); ci != nil {
		t.Errorf("exit 3 should return nil, got %+v", ci)
	}

	blocking := "#!/bin/sh\nexec /bin/sleep 10\n"
	if err := os.WriteFile(bin, []byte(blocking), 0o755); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	if ci := agtlogCost(ctx, bin, "slow"); ci != nil {
		t.Errorf("timeout should return nil, got %+v", ci)
	}
}

func TestRenderCost(t *testing.T) {
	if renderCost(nil) != "" {
		t.Error("nil cost should render empty")
	}
	got := renderCost(&costInfo{USD: 2.73, Complete: true, Input: 7817, CacheCreation: 78245, CacheRead: 1700000, Output: 40734})
	for _, want := range []string{glyphCost, "$2.73", "1.7M", "78K", "7.8K", "41K"} {
		if !strings.Contains(got, want) {
			t.Errorf("renderCost missing %q in %q", want, got)
		}
	}
	if strings.Index(got, "1.7M") > strings.Index(got, "7.8K") {
		t.Errorf("cache-read should precede input, got %q", got)
	}
	if !strings.Contains(got, cGreen) && !strings.Contains(got, cYellow) && !strings.Contains(got, cRed) {
		t.Errorf("expected colored numbers, got %q", got)
	}
}

func TestRenderCostUnpriced(t *testing.T) {
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
	got := render(p, nil, nil, &costInfo{USD: 1.17, Complete: true, Input: 1000, CacheCreation: 2000, CacheRead: 3000, Output: 500}, 1_000_000, 0)
	if !strings.Contains(got, "$1.17") {
		t.Errorf("render should include the cost segment, got %q", got)
	}
}

func TestRenderWraps(t *testing.T) {
	const now int64 = 1_000_000
	var p Payload
	p.Model.DisplayName = "Opus 4.8"
	p.ContextWindow = &ContextWindow{
		TotalInputTokens:  178200,
		ContextWindowSize: 1_000_000,
	}
	claude := []Window{
		{Minutes: 300, Used: 28, ResetsAt: pi(now + 5000)},
		{Minutes: 10080, Used: 15, ResetsAt: pi(now + 500000)},
	}
	codex := []Window{
		{Minutes: 300, Used: 94, ResetsAt: pi(now + 3000)},
		{Minutes: 10080, Used: 16, ResetsAt: pi(now + 500000)},
	}
	cost := &costInfo{USD: 1.17, Complete: true, Input: 1000, Output: 500}

	narrow := render(p, claude, codex, cost, now, 40)
	if !strings.Contains(narrow, "\n") {
		t.Fatalf("width 40 should wrap, got %q", narrow)
	}
	for _, row := range strings.Split(narrow, "\n") {
		if visWidth(row) > 40 {
			t.Errorf("row %q has width %d, want <= 40", row, visWidth(row))
		}
	}

	unwrapped := render(p, claude, codex, cost, now, 0)
	wide := render(p, claude, codex, cost, now, 200)
	if strings.Contains(unwrapped, "\n") || strings.Contains(wide, "\n") {
		t.Errorf("non-wrapping widths produced a newline: width 0=%q, width 200=%q", unwrapped, wide)
	}
	if wide != unwrapped {
		t.Errorf("wide render differs from unwrapped render:\nwide %q\nunwrapped %q", wide, unwrapped)
	}
}

func TestUsageSourceCurrent(t *testing.T) {
	const now int64 = 1_000_000
	snapshot := []Window{{Minutes: 300, Used: 42, ResetsAt: pi(now + 9000)}}
	for _, source := range []usageSource{claudeSource, codexSource} {
		t.Run(source.file, func(t *testing.T) {
			cases := []struct {
				name    string
				age     *int64
				data    string
				want    []Window
				refresh bool
			}{
				{"missing", nil, "", nil, true},
				{"fresh", pi(0), "", snapshot, false},
				{"at interval", pi(source.interval), "", snapshot, false},
				{"stale", pi(source.interval + 1), "", snapshot, true},
				{"fresh Claude rate_limits cache", nil, `{"fetched_at":1000000,"rate_limits":{"seven_day":{"used_percentage":70}}}`, nil, false},
				{"fresh Codex rate_limits cache", nil, `{"fetched_at":1000000,"rate_limits":{"primary":{"used_percent":21,"window_minutes":10080}}}`, nil, false},
				{"stale Claude rate_limits cache", nil, `{"fetched_at":0,"rate_limits":{"seven_day":{"used_percentage":70}}}`, nil, true},
				{"stale Codex rate_limits cache", nil, `{"fetched_at":0,"rate_limits":{"primary":{"used_percent":21,"window_minutes":10080}}}`, nil, true},
				{"malformed", nil, `{`, nil, true},
				{"invalid cached type", nil, `{"fetched_at":1000000,"windows":"invalid"}`, nil, true},
				{"empty windows", nil, `{"fetched_at":1000000,"windows":[]}`, []Window{}, false},
			}
			for _, c := range cases {
				t.Run(c.name, func(t *testing.T) {
					t.Setenv("HOME", t.TempDir())
					t.Setenv("XDG_CACHE_HOME", t.TempDir())
					t.Setenv("PATH", t.TempDir())
					source := source
					source.arg = "-test.run=^$"
					if err := os.WriteFile(filepath.Join(os.Getenv("PATH"), source.bin), []byte("#!/bin/sh\nexit 0\n"), 0o755); err != nil {
						t.Fatal(err)
					}
					if c.age != nil {
						if err := source.write(&cacheFile{FetchedAt: now - *c.age, Windows: snapshot}); err != nil {
							t.Fatal(err)
						}
					} else if c.data != "" {
						if err := os.MkdirAll(filepath.Dir(source.cachePath()), 0o755); err != nil {
							t.Fatal(err)
						}
						if err := os.WriteFile(source.cachePath(), []byte(c.data), 0o644); err != nil {
							t.Fatal(err)
						}
					}
					if got := source.current(now); !reflect.DeepEqual(got, c.want) {
						t.Errorf("current = %+v, want %+v", got, c.want)
					}
					_, err := os.Stat(source.lockPath())
					if refreshed := err == nil; refreshed != c.refresh {
						t.Errorf("refresh attempted = %v, want %v (lock error: %v)", refreshed, c.refresh, err)
					}
				})
			}
		})
	}
}

func TestParseClaudeWindows(t *testing.T) {
	cases := []struct {
		name, api, harness string
		want               []Window
	}{
		{
			"both windows",
			`{"seven_day":{"utilization":54,"resets_at":"2026-07-23T06:00:00.272960+00:00"},"five_hour":{"utilization":15,"resets_at":"2026-07-22T03:10:00.272933+00:00"}}`,
			`{"seven_day":{"used_percentage":54,"resets_at":1784786400},"five_hour":{"used_percentage":15,"resets_at":1784689800}}`,
			[]Window{{Minutes: 300, Used: 15, ResetsAt: pi(1784689800)}, {Minutes: 10080, Used: 54, ResetsAt: pi(1784786400)}},
		},
		{
			"five hour only",
			`{"five_hour":{"utilization":23.5,"resets_at":"2026-07-22T03:10:00.272933+00:00"}}`,
			`{"five_hour":{"used_percentage":23.5,"resets_at":1784689800}}`,
			[]Window{{Minutes: 300, Used: 23.5, ResetsAt: pi(1784689800)}},
		},
		{
			"seven day only",
			`{"five_hour":null,"seven_day":{"utilization":70,"resets_at":"2026-07-23T06:00:00Z"},"limits":[]}`,
			`{"five_hour":null,"seven_day":{"used_percentage":70,"resets_at":1784786400}}`,
			[]Window{{Minutes: 10080, Used: 70, ResetsAt: pi(1784786400)}},
		},
		{"missing percent", `{"five_hour":{"resets_at":"2026-07-22T03:10:00Z"},"seven_day":{"utilization":0}}`, `{"five_hour":{"resets_at":1784689800},"seven_day":{"used_percentage":0}}`, []Window{{Minutes: 10080, Used: 0}}},
		{"null percent", `{"five_hour":{"utilization":null}}`, `{"five_hour":{"used_percentage":null}}`, nil},
		{"optional resets", `{"five_hour":{"utilization":15,"resets_at":""},"seven_day":{"utilization":54,"resets_at":"garbage"}}`, `{"five_hour":{"used_percentage":15},"seven_day":{"used_percentage":54,"resets_at":null}}`, []Window{{Minutes: 300, Used: 15}, {Minutes: 10080, Used: 54}}},
		{"absent windows", `{}`, `{}`, nil},
		{"null windows", `{"five_hour":null,"seven_day":null}`, `{"five_hour":null,"seven_day":null}`, nil},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := parseClaudeUsage([]byte(c.api)); !reflect.DeepEqual(got, c.want) {
				t.Errorf("API windows = %+v, want %+v", got, c.want)
			}
			var p Payload
			if err := json.Unmarshal([]byte(`{"rate_limits":`+c.harness+`}`), &p); err != nil {
				t.Fatal(err)
			}
			if got := p.claudeWindows(); !reflect.DeepEqual(got, c.want) {
				t.Errorf("harness windows = %+v, want %+v", got, c.want)
			}
		})
	}
	for _, data := range []string{"", "not json", `[]`, `{"five_hour":{"utilization":"invalid"}}`} {
		if got := parseClaudeUsage([]byte(data)); got != nil {
			t.Errorf("parseClaudeUsage(%q) = %+v, want nil", data, got)
		}
	}
}

func TestParseAppServerRateLimits(t *testing.T) {
	cases := []struct {
		name, data string
		want       []Window
	}{
		{"both windows", `{"rateLimits":{"primary":{"usedPercent":42,"windowDurationMins":300,"resetsAt":9999999999},"secondary":{"usedPercent":8,"windowDurationMins":10080,"resetsAt":9999999999}}}`, []Window{{Minutes: 300, Used: 42, ResetsAt: pi(9999999999)}, {Minutes: 10080, Used: 8, ResetsAt: pi(9999999999)}}},
		{"seven day primary", `{"rateLimits":{"primary":{"usedPercent":21,"windowDurationMins":10080,"resetsAt":1784786400},"secondary":null}}`, []Window{{Minutes: 10080, Used: 21, ResetsAt: pi(1784786400)}}},
		{"length determines label", `{"rateLimits":{"primary":{"usedPercent":21,"windowDurationMins":90}}}`, []Window{{Minutes: 90, Used: 21}}},
		{"missing length", `{"rateLimits":{"primary":{"usedPercent":42},"secondary":{"usedPercent":0,"windowDurationMins":10080}}}`, []Window{{Minutes: 10080, Used: 0}}},
		{"missing percent", `{"rateLimits":{"primary":{"windowDurationMins":300}}}`, nil},
		{"null length", `{"rateLimits":{"primary":{"usedPercent":42,"windowDurationMins":null}}}`, nil},
		{"null windows", `{"rateLimits":{"primary":null,"secondary":null}}`, nil},
		{"empty", `{}`, nil},
		{"malformed", `not json`, nil},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := parseAppServerRateLimits([]byte(c.data)); !reflect.DeepEqual(got, c.want) {
				t.Errorf("parseAppServerRateLimits = %+v, want %+v", got, c.want)
			}
		})
	}
}

func TestRenderSevenDayOnly(t *testing.T) {
	const now int64 = 1784484000
	var harness Payload
	if err := json.Unmarshal([]byte(`{"rate_limits":{"five_hour":null,"seven_day":{"used_percentage":70,"resets_at":1784786400}}}`), &harness); err != nil {
		t.Fatal(err)
	}
	api := parseClaudeUsage([]byte(`{"five_hour":null,"seven_day":{"utilization":70,"resets_at":"2026-07-23T06:00:00Z"}}`))
	codex := parseAppServerRateLimits([]byte(`{"rateLimits":{"primary":{"usedPercent":21,"windowDurationMins":10080,"resetsAt":1784786400},"secondary":null}}`))
	claudeSegment := glyphTimer + " " + cYellow + "70%" + cReset + " " + cYellow + "3d12h" + cReset + cGray + "/7d" + cReset
	codexSegment := glyphTimer + " Codex " + cGreen + "21%" + cReset + " " + cYellow + "3d12h" + cReset + cGray + "/7d" + cReset
	cases := []struct {
		name          string
		claude, codex []Window
		groups        string
	}{
		{"Claude API", api, nil, claudeSegment},
		{"Claude harness", harness.claudeWindows(), nil, claudeSegment},
		{"Codex", nil, codex, codexSegment},
		{"both sources", harness.claudeWindows(), codex, claudeSegment + " | " + codexSegment},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := render(Payload{}, c.claude, c.codex, nil, now, 0)
			want := render(Payload{}, nil, nil, nil, now, 0) + " | " + c.groups
			if got != want {
				t.Errorf("render = %q, want %q", got, want)
			}
		})
	}
}
