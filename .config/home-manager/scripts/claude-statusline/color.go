package main

import (
	"fmt"
	"os"
)

const (
	cRed    = "\x1b[31m"
	cYellow = "\x1b[33m"
	cGreen  = "\x1b[32m"
	cGray   = "\x1b[90m"
	cReset  = "\x1b[0m"
)

// Nerd Font glyphs render blank in other fonts, so escapes keep them visible in source.
const (
	glyphModel = "\U000F06A9"
	glyphToken = "\U0000EDE8"
	glyphTimer = "\U000F051F"
	glyphCost  = "\U000F0114"
)

var trueColor = supportsTrueColor(os.Getenv)

// Warp on WSL omits COLORTERM, so TERM_PROGRAM also matches the 24-bit terminals in use.
func supportsTrueColor(env func(string) string) bool {
	switch env("COLORTERM") {
	case "truecolor", "24bit":
		return true
	}
	switch env("TERM_PROGRAM") {
	case "WarpTerminal", "iTerm.app":
		return true
	}
	return false
}

type rgb struct{ r, g, b uint8 }

func (c rgb) escape() string { return fmt.Sprintf("\x1b[38;2;%d;%d;%dm", c.r, c.g, c.b) }

func (c rgb) mix(o rgb, t float64) rgb {
	return rgb{lerp8(c.r, o.r, t), lerp8(c.g, o.g, t), lerp8(c.b, o.b, t)}
}

func lerp8(a, b uint8, t float64) uint8 { return uint8(float64(a) + t*(float64(b)-float64(a)) + 0.5) }

var (
	gradLow  = rgb{80, 200, 120}
	gradMid  = rgb{235, 200, 70}
	gradHigh = rgb{235, 90, 80}
)

// The yellow midpoint avoids the brown of a direct green-to-red blend.
func gradient(v, lo, hi int64) string {
	t := norm(v, lo, hi)
	if t <= 0.5 {
		return gradLow.mix(gradMid, t*2).escape()
	}
	return gradMid.mix(gradHigh, (t-0.5)*2).escape()
}

func norm(v, lo, hi int64) float64 {
	switch {
	case hi <= lo || v >= hi:
		return 1
	case v <= lo:
		return 0
	default:
		return float64(v-lo) / float64(hi-lo)
	}
}

func tcolor(v, lo, hi int64) string {
	if trueColor {
		return gradient(v, lo, hi)
	}
	switch {
	case v >= hi:
		return cRed
	case v >= lo:
		return cYellow
	default:
		return cGreen
	}
}

func paint(color, s string) string { return color + s + cReset }
