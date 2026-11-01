package main

import (
	"fmt"
	"strconv"
	"strings"
)

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
	if s == "1000" {
		s, exp = "1", exp+1
	}
	const units = "KMBT"
	if exp > len(units) {
		exp = len(units)
	}
	return s + units[exp-1:exp]
}

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

func trimZeros(v float64, scale int) string {
	s := strconv.FormatFloat(v, 'f', scale, 64)
	if strings.Contains(s, ".") {
		s = strings.TrimRight(s, "0")
		s = strings.TrimRight(s, ".")
	}
	return s
}

// Minute precision on multi-day windows costs space without helping gauge usage.
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
