package main

import (
	"bufio"
	"encoding/json"
	"os"
	"strings"
)

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

// render builds the whole status line. Pure given its inputs (now and width are
// passed rather than read from process state) so every branch is testable.
func render(p Payload, codex *CodexRL, cost *costInfo, now int64, width int) string {
	segments := []string{
		glyphModel + " " + modelName(p.Model.DisplayName),
		glyphToken + " " + renderContext(p),
	}
	if s := renderCost(cost); s != "" {
		segments = append(segments, s)
	}
	if segs := claudeSegments(p, now); len(segs) > 0 {
		segments = append(segments, glyphTimer+" "+strings.Join(segs, ", "))
	}
	if segs := codexSegments(codex, now); len(segs) > 0 {
		segments = append(segments, glyphTimer+" Codex "+strings.Join(segs, ", "))
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
