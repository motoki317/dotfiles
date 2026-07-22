package main

import (
	"fmt"
	"os"
)

// ANSI colors. Real ESC bytes, so rendered segments concatenate directly.
const (
	cRed    = "\x1b[31m"
	cYellow = "\x1b[33m"
	cGreen  = "\x1b[32m"
	cGray   = "\x1b[90m"
	cReset  = "\x1b[0m"
)

// Nerd-font glyphs: the model marker, the token glyph prefixing the context
// count, and the timer-sand that prefixes each rate-limit group. glyphToken lives
// in the Private Use Area, where it renders blank in fonts without the nerd-font
// patch — hence the \U escape, so it can never be mistaken for whitespace and
// silently dropped (as it was when this was ported from bash).
const (
	glyphModel = "\U000F06A9" // 󰚩 md-robot
	glyphToken = "\U0000EDE8" // token/context marker
	glyphTimer = "\U000F051F" // 󰔟 md-timer-sand
	glyphCost  = "\U000F0114" // 󰄔 md-cash
)

// trueColor reports whether the terminal advertises 24-bit color, resolved once at
// startup. When it does, tcolor emits a smooth green→yellow→red gradient; when it
// doesn't, the 3-step ANSI palette — so color loses fidelity, never vanishes.
var trueColor = supportsTrueColor(os.Getenv)

// supportsTrueColor detects 24-bit terminals. COLORTERM=truecolor/24bit is the
// conventional signal, but Warp on WSL sets no COLORTERM, so the two terminals in
// use are also allowlisted by TERM_PROGRAM. Unknown terminals degrade to the 3-step
// palette rather than risk garbled escapes.
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

// rgb is a 24-bit color. escape emits its ANSI truecolor foreground sequence, which
// cReset ("\x1b[0m") clears like any other SGR, so gradient spans concatenate safely.
type rgb struct{ r, g, b uint8 }

func (c rgb) escape() string { return fmt.Sprintf("\x1b[38;2;%d;%d;%dm", c.r, c.g, c.b) }

// mix linearly interpolates from c toward o by t in [0,1], rounded to nearest.
func (c rgb) mix(o rgb, t float64) rgb {
	return rgb{lerp8(c.r, o.r, t), lerp8(c.g, o.g, t), lerp8(c.b, o.b, t)}
}

func lerp8(a, b uint8, t float64) uint8 { return uint8(float64(a) + t*(float64(b)-float64(a)) + 0.5) }

// Gradient stops mirroring the 3-step palette's meaning: green well under the low
// mark, yellow at the midpoint, red at or above the high mark.
var (
	gradLow  = rgb{80, 200, 120}
	gradMid  = rgb{235, 200, 70}
	gradHigh = rgb{235, 90, 80}
)

// gradient maps v across [lo,hi] to a green→yellow→red 24-bit color. The two linear
// segments (green→yellow, then yellow→red) keep the midpoint a clean yellow instead
// of the muddy brown a direct green→red interpolation would pass through.
func gradient(v, lo, hi int64) string {
	t := norm(v, lo, hi)
	if t <= 0.5 {
		return gradLow.mix(gradMid, t*2).escape()
	}
	return gradMid.mix(gradHigh, (t-0.5)*2).escape()
}

// norm clamps v into [lo,hi] and returns its position there in [0,1]; a degenerate
// lo>=hi maps to 1 (the high color), matching tcolor's v>=hi branch.
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

// tcolor maps a value to a color by proximity to the [lo,hi] band: a smooth
// green→yellow→red gradient on 24-bit terminals, else green below lo / yellow below
// hi / red at or above hi. Used both for 0-100 progress and token/cost magnitudes.
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

// paint wraps s in an ANSI color and a reset, so colored spans concatenate safely.
func paint(color, s string) string { return color + s + cReset }
