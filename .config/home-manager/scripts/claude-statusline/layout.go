package main

import (
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
)

// termWidth reads the per-pane width Claude Code injects as COLUMNS. Zero means
// unknown and retains single-line output; terminal queries cannot be used because
// statusline output is captured rather than attached to a TTY. This assumes
// statusLine.padding is zero; subtract configured padding here if that changes.
func termWidth() int {
	width, err := strconv.Atoi(os.Getenv("COLUMNS"))
	if err != nil || width <= 0 {
		return 0
	}
	return width
}

// visWidth returns the rendered cell width of s, skipping ANSI CSI escapes and
// summing each remaining rune's runeCells. It has to match what the terminal
// actually draws: when it undercounts, layout packs a row past the pane width and
// Claude Code right-truncates the overflow with an ellipsis (the bug this fixes —
// the earlier one-cell-per-rune count lost ~7 cells on a glyph-heavy line).
func visWidth(s string) int {
	width := 0
	for i := 0; i < len(s); {
		if s[i] == '\x1b' && i+1 < len(s) && s[i+1] == '[' {
			j := i + 2
			for j < len(s) && (s[j] < 0x40 || s[j] > 0x7e) {
				j++
			}
			if j < len(s) {
				i = j + 1
				continue
			}
		}
		r, size := utf8.DecodeRuneInString(s[i:])
		width += runeCells(r)
		i += size
	}
	return width
}

// runeCells returns the terminal cell width of r. Every rune this program prints
// is either ASCII or one of two double-width classes: nerd-font icons, which live
// in the Private Use Area (glyphModel/glyphToken/glyphTimer/glyphCost), and
// renderCost's ↑/↓, which are East-Asian Ambiguous and drawn wide here. That
// closed alphabet is why the full width table (go-runewidth) an arbitrary string
// would need — an external dep the stdlib-only build bars — is unnecessary: these
// ranges catch every wide rune actually emitted. Erring wide is the safe
// direction; a terminal that drew these narrow would only wrap a few cells early,
// where the reverse overflows and truncates.
func runeCells(r rune) int {
	switch {
	case r >= 0xE000 && r <= 0xF8FF: // BMP Private Use Area
		return 2
	case r >= 0xF0000 && r <= 0xFFFFD: // Supplementary Private Use Area-A (md icons)
		return 2
	case r == 0x2191 || r == 0x2193: // ↑ ↓
		return 2
	}
	return 1
}

// layout greedily packs segments into width-bounded rows while preserving their
// order. A non-positive width retains the legacy single-line output. A segment
// wider than width still gets its own row and overflows; per-segment splitting is
// intentionally out of scope.
func layout(segments []string, width int) string {
	const sep = " | " // one source for the join, its width term, and the concat
	if width <= 0 {
		return strings.Join(segments, sep)
	}

	var rows []string
	row := ""
	for i, segment := range segments {
		if i == 0 {
			row = segment
			continue
		}
		if visWidth(row)+visWidth(sep)+visWidth(segment) <= width {
			row += sep + segment
			continue
		}
		rows = append(rows, row)
		row = segment
	}
	if len(segments) > 0 {
		rows = append(rows, row)
	}
	return strings.Join(rows, "\n")
}
