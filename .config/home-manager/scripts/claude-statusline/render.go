package main

import "strings"

type ContextWindow struct {
	TotalInputTokens  int64 `json:"total_input_tokens"`
	ContextWindowSize int64 `json:"context_window_size"`
}

func render(p Payload, claude, codex []Window, cost *costInfo, now int64, width int) string {
	segments := []string{glyphModel + " " + modelName(p.Model.DisplayName)}
	if s := renderContext(p.ContextWindow); s != "" {
		segments = append(segments, glyphToken+" "+s)
	}
	if s := renderCost(cost); s != "" {
		segments = append(segments, s)
	}
	if s := rateLimitSegment(claude, "", now); s != "" {
		segments = append(segments, s)
	}
	if s := rateLimitSegment(codex, "Codex ", now); s != "" {
		segments = append(segments, s)
	}
	return layout(segments, width)
}

// The display name's "(1M context)" suffix repeats the size that the context segment shows.
func modelName(display string) string {
	if i := strings.LastIndex(display, " ("); i >= 0 && strings.HasSuffix(display, ")") {
		inner := display[i+2 : len(display)-1]
		if strings.Contains(strings.ToLower(inner), "context") {
			return display[:i]
		}
	}
	return display
}

// The absolute usage and its color convey the percentage without extra text.
func renderContext(c *ContextWindow) string {
	if c == nil || c.ContextWindowSize <= 0 {
		return ""
	}
	used, limit := c.TotalInputTokens, c.ContextWindowSize
	pctInt := used * 100 / limit
	return tcolor(pctInt, 30, 60) + humanize(used) + cReset + "/" + humanize(limit)
}
