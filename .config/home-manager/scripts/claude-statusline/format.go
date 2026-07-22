package main

import (
	"fmt"
	"strconv"
	"strings"
)

// humanize renders a token count with two significant figures, keeping a third
// for hundreds: 1234->"1.2K", 12340->"12K", 123400->"123K", 1.26e6->"1.3M".
// Counts under 1000 print verbatim.
func humanize(n int64) string {
	if n < 1000 {
		return strconv.FormatInt(n, 10)
	}
	f, exp := float64(n), 0
	for f >= 1000 {
		f /= 1000
		exp++
	}
	s := sigFig(f, 2)
	if s == "1000" { // mantissa rounded up into the next unit (e.g. 999600 -> 1M)
		s, exp = "1", exp+1
	}
	const units = "KMBT"
	if exp > len(units) { // >= 1e15; unreachable for token counts, but never index out of range
		exp = len(units)
	}
	return s + units[exp-1:exp]
}

// sigFig renders v with n significant figures, trailing zeros trimmed. Integer
// digits count toward n; a value below 1 gets n fractional digits — leading zeros
// aren't significant, and trimZeros drops the rest ($0.05, not $0.050). Tokens use
// n=2 for compactness; cost uses n=3, where the extra digit's precision earns its
// width. Values with more than n integer digits print whole (127, not 130).
func sigFig(v float64, n int) string {
	intDigits := 0
	for x := v; x >= 1; x /= 10 {
		intDigits++
	}
	dec := n - intDigits
	if dec < 0 {
		dec = 0
	}
	return trimZeros(v, dec)
}

// trimZeros formats v with at most scale decimals, dropping trailing zeros and a
// trailing point (2.50 -> "2.5", 1.00 -> "1").
func trimZeros(v float64, scale int) string {
	s := strconv.FormatFloat(v, 'f', scale, 64)
	if strings.Contains(s, ".") {
		s = strings.TrimRight(s, "0")
		s = strings.TrimRight(s, ".")
	}
	return s
}

// fmtDur renders a duration as "3d23h" / "3h28m" / "47m", dropping leading zero
// units and, at day scale, the minutes — a multi-day window (the 7d rate limit)
// doesn't need minute precision, and dropping them keeps the line compact.
func fmtDur(s int64) string {
	if s < 0 {
		s = 0
	}
	d, h, m := s/86400, (s%86400)/3600, (s%3600)/60
	switch {
	case d > 0:
		return fmt.Sprintf("%dd%dh", d, h)
	case h > 0:
		return fmt.Sprintf("%dh%dm", h, m)
	default:
		return fmt.Sprintf("%dm", m)
	}
}

// winLabel renders a window length in minutes as a compact label: 300 -> "5h",
// 10080 -> "7d", 90 -> "90m".
func winLabel(min int) string {
	switch {
	case min > 0 && min%1440 == 0:
		return fmt.Sprintf("%dd", min/1440)
	case min > 0 && min%60 == 0:
		return fmt.Sprintf("%dh", min/60)
	default:
		return fmt.Sprintf("%dm", min)
	}
}
