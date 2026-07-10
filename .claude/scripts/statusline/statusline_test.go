package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func pf(f float64) *float64 { return &f }
func pi(i int64) *int64     { return &i }
func pn(i int) *int         { return &i }

func TestHumanize(t *testing.T) {
	cases := []struct {
		in   int64
		want string
	}{
		{0, "0"},
		{999, "999"},
		{1000, "1K"},
		{1234, "1.2K"},
		{178289, "178.3K"}, // rounds (bash's bc truncated to 178.2K); rounding is intended
		{1000000, "1M"},
		{1250000, "1.25M"},
		{1500000, "1.5M"},
	}
	for _, c := range cases {
		if got := humanize(c.in); got != c.want {
			t.Errorf("humanize(%d) = %q, want %q", c.in, got, c.want)
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
		v      int
		lo, hi int
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

func TestFmtDur(t *testing.T) {
	cases := []struct {
		s    int64
		want string
	}{
		{-5, "0m"},
		{0, "0m"},
		{47 * 60, "47m"},
		{3*3600 + 28*60, "3h28m"},
		{3*86400 + 23*3600 + 28*60, "3d23h28m"},
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
	if !strings.Contains(got, "178.2K") || !strings.Contains(got, "/1M") {
		t.Errorf("context = %q, want 178.2K/1M", got)
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

	got := render(p, codex, now)
	// A color reset sits between "178.2K" and "/1M", so check the parts separately.
	for _, want := range []string{glyphModel + " Opus 4.8", " |  ", "178.2K", "/1M", glyphTimer + " Codex", "94%", "16%"} {
		if !strings.Contains(got, want) {
			t.Errorf("render missing %q in:\n%s", want, got)
		}
	}
}

func TestRenderCodexAbsent(t *testing.T) {
	var p Payload
	p.Model.DisplayName = "Opus 4.8"
	got := render(p, nil, 1_000_000)
	if strings.Contains(got, "Codex") {
		t.Errorf("no Codex data should render no Codex group, got %q", got)
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
	rl := codexRateLimits(home)
	if rl == nil || rl.Primary == nil || *rl.Primary.UsedPercent != 80.0 {
		t.Errorf("expected the newer event (80%%), got %+v", rl)
	}
}

func TestCodexRateLimitsAbsent(t *testing.T) {
	if rl := codexRateLimits(t.TempDir()); rl != nil {
		t.Errorf("no ~/.codex should yield nil, got %+v", rl)
	}
}
