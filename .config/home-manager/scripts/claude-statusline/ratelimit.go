package main

import "fmt"

// Anthropic's payload omits window lengths, so they are fixed here. Codex's
// payload carries window_minutes, so its lengths and labels are derived instead.
const (
	fiveHourSecs = 5 * 60 * 60
	sevenDaySecs = 7 * 24 * 60 * 60
)

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
