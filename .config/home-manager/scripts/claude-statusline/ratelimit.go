package main

import (
	"fmt"
	"strings"
)

type Window struct {
	Minutes  int     `json:"minutes"`
	Used     float64 `json:"used"`
	ResetsAt *int64  `json:"resets_at"`
}

// Neither Claude's usage API nor the statusline JSON states window lengths.
var claudeWindowLengths = []struct {
	key     string
	minutes int
}{
	{"five_hour", 300},
	{"seven_day", 10080},
}

func (p Payload) claudeWindows() []Window {
	var windows []Window
	for _, length := range claudeWindowLengths {
		if w := p.RateLimits[length.key]; w.Used != nil {
			windows = append(windows, Window{Minutes: length.minutes, Used: *w.Used, ResetsAt: w.ResetsAt})
		}
	}
	return windows
}

func rateLimitSegment(windows []Window, prefix string, now int64) string {
	if len(windows) == 0 {
		return ""
	}
	parts := make([]string, 0, len(windows))
	for _, w := range windows {
		parts = append(parts, rlSegment(w, now))
	}
	return glyphTimer + " " + prefix + strings.Join(parts, ", ")
}

// Usage and elapsed time share the 50/80 color bands, so usage redder than time
// means the allowance is burning fast.
func rlSegment(w Window, now int64) string {
	label, windowSecs := winLabel(w.Minutes), int64(w.Minutes)*60
	if w.ResetsAt != nil && now >= *w.ResetsAt {
		return cGray + "~" + label + cReset
	}
	pctInt := int64(w.Used)
	pc := tcolor(pctInt, 50, 80)
	if w.ResetsAt == nil {
		return fmt.Sprintf("%s%d%%%s %s%s%s", pc, pctInt, cReset, cGray, label, cReset)
	}
	elapsed := windowSecs - (*w.ResetsAt - now)
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
