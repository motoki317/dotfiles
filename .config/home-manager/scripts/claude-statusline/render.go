package main

import "strings"

type ContextWindow struct {
	TotalInputTokens  int64 `json:"total_input_tokens"`
	ContextWindowSize int64 `json:"context_window_size"`
}

// render builds the whole status line. Pure given its inputs (now and width are
// passed rather than read from process state) so every branch is testable.
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

// The absolute usage and its color convey the percentage without extra text.
func renderContext(c *ContextWindow) string {
	if c == nil || c.ContextWindowSize <= 0 {
		return ""
	}
	used, limit := c.TotalInputTokens, c.ContextWindowSize
	pctInt := used * 100 / limit
	// 30/60 keep the prior absolute marks (300K/600K on a 1M window) while scaling
	// to whatever context_window_size reports, so 200k and 1M models both work.
	return tcolor(pctInt, 30, 60) + humanize(used) + cReset + "/" + humanize(limit)
}
