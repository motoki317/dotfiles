package main

import (
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
)

// Claude Code captures stdout, so terminal queries cannot provide the pane width.
// Using the full COLUMNS width assumes statusLine.padding is zero.
func termWidth() int {
	width, err := strconv.Atoi(os.Getenv("COLUMNS"))
	if err != nil || width <= 0 {
		return 0
	}
	return width
}

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

// The Private Use Area icons and the arrows draw two cells wide. Undercounting them
// overflows a row, which Claude Code truncates. These ranges cover every glyph this
// program prints, so no Unicode width table is needed.
func runeCells(r rune) int {
	switch {
	case r >= 0xE000 && r <= 0xF8FF:
		return 2
	case r >= 0xF0000 && r <= 0xFFFFD:
		return 2
	case r == 0x2191 || r == 0x2193:
		return 2
	}
	return 1
}

func layout(segments []string, width int) string {
	const sep = " | "
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
