#!/usr/bin/env bash

# ref: https://code.claude.com/docs/en/statusline
#
# Mirrors the native /context readout: current context-window usage toward the
# active model's window limit.
#
# Source of truth is the harness-provided `.context_window` object. Its figures
# match /context exactly: input-only (input + cache_creation + cache_read,
# EXCLUDING output_tokens) and current-not-cumulative as of Claude Code
# v2.1.132. We read `total_input_tokens` (the same input-only sum) for "used"
# and `context_window_size` for the limit, so 200k vs 1M models are handled
# automatically. A transcript fallback (identical input-only formula, sidechain
# entries excluded) covers pre-2.1.132 clients that lack `.context_window`.

# Read JSON input from stdin
input=$(cat)

MODEL_DISPLAY=$(echo "$input" | jq -r '.model.display_name')

# Nerd-font glyph prefixing the rate-limit group (Material Design Icons range,
# which the model glyph above confirms renders). Swap to taste; "" to disable.
RL_ICON="󰔟"

# Prefix for the Codex rate-limit group (a separate provider, read from
# ~/.codex session rollouts below). Reuses the same glyph and adds a "Codex"
# label so the two providers' limits stay distinguishable at a glance; when
# RL_ICON is disabled the label stands alone.
CODEX_RL_PREFIX="${RL_ICON:+$RL_ICON }Codex"

# 1234 -> 1.2K, 1000000 -> 1M, 1250000 -> 1.25M. Round to the magnitude's scale
# with bc, then let %g trim trailing zeros and any bare decimal point to save
# width (bc first because %g alone would keep 6 significant digits).
humanize() {
  local n=${1:-0}
  if [ "$n" -ge 1000000 ]; then
    printf "%gM" "$(echo "scale=2; $n/1000000" | bc)"
  elif [ "$n" -ge 1000 ]; then
    printf "%gK" "$(echo "scale=1; $n/1000" | bc)"
  else
    printf "%d" "$n"
  fi
}

# ANSI color for a 0-100 progress value: green below $2, yellow below $3, red at
# or above $3. Emits a real ESC byte so callers concatenate it directly.
tcolor() {
  local v=${1:-0} lo=$2 hi=$3
  if   [ "$v" -ge "$hi" ]; then printf "\033[31m"
  elif [ "$v" -ge "$lo" ]; then printf "\033[33m"
  else                          printf "\033[32m"
  fi
}

# Seconds -> "3d23h28m" / "3h28m" / "47m": full breakdown, leading zero units
# dropped. Used for elapsed time within a rate-limit window, so minutes are kept
# even at day scale (unlike a coarse "resets in ~3d" countdown).
fmt_dur() {
  local s=${1:-0}
  [ "$s" -lt 0 ] && s=0
  local d=$(( s / 86400 )) h=$(( (s % 86400) / 3600 )) m=$(( (s % 3600) / 60 ))
  if   [ "$d" -gt 0 ]; then printf "%dd%dh%dm" "$d" "$h" "$m"
  elif [ "$h" -gt 0 ]; then printf "%dh%dm" "$h" "$m"
  else                      printf "%dm" "$m"
  fi
}

# Minutes -> compact window label: 300 -> "5h", 10080 -> "7d", 90 -> "90m".
# Codex reports each window's length as `window_minutes`, so its labels are
# derived rather than hardcoded (Anthropic's payload omits it, hence the fixed
# 5h/7d passed to rl_segment for the Claude block).
win_label() {
  local m=${1:-0}
  if   [ $(( m % 1440 )) -eq 0 ]; then printf "%dd" "$(( m / 1440 ))"
  elif [ $(( m % 60 )) -eq 0 ]; then printf "%dh" "$(( m / 60 ))"
  else                                printf "%dm" "$m"
  fi
}

# Render one rate-limit window as "<pct>% <elapsed>/<label>", e.g. "28% 3h28m/5h".
# Elapsed is colored by time-through-window and pct by usage, so the two colors
# side by side reveal burn pace: usage redder than time == burning fast.
# Window length ($2, seconds) drives elapsed = window - (resets - now).
# Emits nothing when pct is absent; falls back to "<pct>% <label>" when resets_at
# is absent. When resets_at has already passed (NOW >= it), the window has rolled
# over and the carried pct describes an expired window — a stale SNAPSHOT source
# (Codex, written only at each API call), so we dim it and mark "~<label>" rather
# than report a wrong number. Live sources (Claude, injected fresh each render)
# keep resets_at in the future, so this branch never fires for them.
NOW=$(date +%s)
rl_segment() {
  local label=$1 window=$2 pct=$3 resets=$4
  [ -z "$pct" ] && return
  if [ -n "$resets" ] && [ "$resets" != "null" ] && [ "$NOW" -ge "$resets" ]; then
    printf "%b" "\033[90m~${label}\033[0m"; return
  fi
  local pct_int=${pct%%.*}; pct_int=${pct_int:-0}
  local pc; pc=$(tcolor "$pct_int" 50 80)
  if [ -n "$resets" ] && [ "$resets" != "null" ]; then
    local elapsed=$(( window - (resets - NOW) ))
    [ "$elapsed" -lt 0 ] && elapsed=0
    local tc; tc=$(tcolor "$(( elapsed * 100 / window ))" 50 80)
    printf "%b" "${pc}${pct_int}%\033[0m ${tc}$(fmt_dur "$elapsed")\033[0m\033[90m/${label}\033[0m"
  else
    printf "%b" "${pc}${pct_int}%\033[0m \033[90m${label}\033[0m"
  fi
}

if [ "$(echo "$input" | jq -r '(.context_window // null) != null')" = "true" ]; then
  # Current clients: read the pre-computed, /context-faithful figures directly.
  read -r used limit pct < <(echo "$input" | jq -r '
    .context_window as $c
    | ($c.context_window_size // 0) as $limit
    | ( $c.total_input_tokens
        // ( ($c.current_usage.input_tokens // 0)
             + ($c.current_usage.cache_creation_input_tokens // 0)
             + ($c.current_usage.cache_read_input_tokens // 0) ) ) as $used
    | ($c.used_percentage // (if $limit > 0 then $used * 100 / $limit else 0 end)) as $pct
    | "\($used) \($limit) \($pct)"')
else
  # Legacy fallback: derive the same input-only sum from the last main-thread
  # assistant turn in the transcript. Sidechain (sub-agent) turns are excluded
  # so we report the main conversation's context, not a sub-agent's.
  TRANSCRIPT_PATH=$(echo "$input" | jq -r '.transcript_path // empty')
  if [ -n "$TRANSCRIPT_PATH" ] && [ -f "$TRANSCRIPT_PATH" ]; then
    used=$(tail -n 400 "$TRANSCRIPT_PATH" 2>/dev/null | jq -s '
      map(select(.type == "assistant" and .message.usage and (.isSidechain != true))) |
      last | .message.usage |
      (.input_tokens // 0) +
      (.cache_creation_input_tokens // 0) +
      (.cache_read_input_tokens // 0)' 2>/dev/null)
  fi
  used=${used:-0}
  limit=0  # window size is absent from legacy payloads; show tokens only
  pct=0
fi

used=${used:-0}; limit=${limit:-0}; pct=${pct:-0}
used_h=$(humanize "$used")

if [ "$limit" -gt 0 ]; then
  limit_h=$(humanize "$limit")
  pct_int=${pct%%.*}; pct_int=${pct_int:-0}

  # Color the used figure by proximity to the limit (the /context concern); 30/60
  # preserve the prior absolute marks (300K/600K) on a 1M window while scaling to
  # whatever `context_window_size` reports. The percentage is elided to save
  # width — the absolute pair conveys it, and the color carries the warning.
  color=$(tcolor "$pct_int" 30 60)
  CONTEXT=$(echo -e "${color}${used_h}\033[0m/${limit_h}")
else
  CONTEXT=$(echo -e "\033[32m${used_h}\033[0m tkns")
fi

# Subscription rate limits (Claude.ai Pro/Max). `.rate_limits` is absent for
# non-subscribers and until the first API response of a session, and each window
# can be independently absent, so every field uses `// empty` and the block is
# appended only when at least one window rendered.
five=$(rl_segment "5h" 18000 \
  "$(echo "$input" | jq -r '.rate_limits.five_hour.used_percentage // empty')" \
  "$(echo "$input" | jq -r '.rate_limits.five_hour.resets_at // empty')")
week=$(rl_segment "7d" 604800 \
  "$(echo "$input" | jq -r '.rate_limits.seven_day.used_percentage // empty')" \
  "$(echo "$input" | jq -r '.rate_limits.seven_day.resets_at // empty')")
RATE=""
[ -n "$five" ] && RATE="$five"
# Comma, not a pipe: the two windows are one group of limit figures, so they
# read together rather than as separate sections.
[ -n "$week" ] && RATE="${RATE:+$RATE, }$week"
if [ -n "$RATE" ]; then
  pre=""; [ -n "$RL_ICON" ] && pre="${RL_ICON} "
  CONTEXT="${CONTEXT} | ${pre}${RATE}"
fi

# Codex subscription limits (OpenAI Codex CLI). Unlike Claude's rate_limits, which
# the harness injects live on stdin, Codex records them only in its session
# rollout files (~/.codex/sessions/YYYY/MM/DD/rollout-*.jsonl) as a snapshot at
# each API call. Read the newest rollout that actually carries one — newest by
# mtime, skipping a just-started session that hasn't hit the API yet — and take
# its last rate_limits entry. Silently absent when Codex isn't installed or has
# no session data. The snapshot's pct is exact until its window rolls over, after
# which rl_segment marks it stale (see there); resets_at stays absolute so the
# elapsed countdown keeps advancing correctly regardless of snapshot age.
codex_line=""
codex_files=$(find "$HOME/.codex/sessions" -type f -name 'rollout-*.jsonl' 2>/dev/null)
if [ -n "$codex_files" ]; then
  while IFS= read -r f; do
    [ -n "$f" ] || continue
    codex_line=$(grep 'rate_limits' "$f" 2>/dev/null | tail -1)
    [ -n "$codex_line" ] && break
  done < <(printf '%s\n' "$codex_files" | xargs ls -t 2>/dev/null)
fi
if [ -n "$codex_line" ]; then
  # primary = 5h window, secondary = 7d. Extract each field on its own rather than
  # packing them into one delimited line: a tab-joined `read` collapses empty
  # columns (tab is IFS whitespace), silently shifting the fields whenever one is
  # null. `// empty` yields "" for a missing field, which the guards below skip.
  codex_win() {
    local base=$1 pct wmin resets
    pct=$(echo "$codex_line" | jq -r "${base}.used_percent // empty")
    wmin=$(echo "$codex_line" | jq -r "${base}.window_minutes // empty")
    resets=$(echo "$codex_line" | jq -r "${base}.resets_at // empty")
    { [ -z "$pct" ] || [ -z "$wmin" ]; } && return
    case "$wmin" in ''|*[!0-9]*) return ;; esac  # guard the shell arithmetic below
    rl_segment "$(win_label "$wmin")" "$(( wmin * 60 ))" "$pct" "$resets"
  }
  c5=$(codex_win '.payload.rate_limits.primary')
  c7=$(codex_win '.payload.rate_limits.secondary')
  CODEX_RATE=""
  [ -n "$c5" ] && CODEX_RATE="$c5"
  [ -n "$c7" ] && CODEX_RATE="${CODEX_RATE:+$CODEX_RATE, }$c7"
  [ -n "$CODEX_RATE" ] && CONTEXT="${CONTEXT} | ${CODEX_RL_PREFIX} ${CODEX_RATE}"
fi

echo "󰚩 ${MODEL_DISPLAY} |  ${CONTEXT}"
