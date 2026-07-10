#!/usr/bin/env bash

# Statusline launcher. The real logic lives in statusline/statusline.go; this
# stays the stable entry point named in ~/.claude/settings.json, so the Go program
# needs no settings change and no separate build step.
#
# It builds a cached per-arch binary only when the source is newer, then execs it
# (passing stdin through). The build is atomic — compile to a tempfile, then
# rename — and if a build fails the previous known-good binary keeps running, so a
# broken edit never blanks the statusline. Mirrors the codex-advisor helpers'
# self-building approach, hardened for a UI component invoked on every render.

set -u

here="$(cd "$(dirname "$0")" && pwd)"
src="$here/statusline"
cache="${XDG_CACHE_HOME:-$HOME/.cache}/claude-statusline"
bin="$cache/bin-$(uname -m)"

if [ ! -x "$bin" ] || [ "$src/statusline.go" -nt "$bin" ]; then
  if command -v go >/dev/null 2>&1; then
    mkdir -p "$cache"
    tmp="$bin.tmp.$$"  # unique per process: concurrent renders never clash on it
    if ( cd "$src" && go build -o "$tmp" . ) 2>/dev/null; then
      mv -f "$tmp" "$bin"
    else
      rm -f "$tmp"  # build failed: fall through to the existing $bin if present
    fi
  fi
fi

if [ -x "$bin" ]; then
  exec "$bin"
fi

# No usable binary (first run with no Go toolchain, or a first build that failed):
# emit a short diagnostic instead of a blank line so the cause is visible.
exec printf '%s\n' "⚠ statusline: build unavailable — run: (cd $src && go build .)"
