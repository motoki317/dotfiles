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

# 1234 -> 1.2K, 1200000 -> 1.20M; matches the magnitude style /context uses.
humanize() {
  local n=${1:-0}
  if [ "$n" -ge 1000000 ]; then
    printf "%.2fM" "$(echo "scale=2; $n/1000000" | bc)"
  elif [ "$n" -ge 1000 ]; then
    printf "%.1fK" "$(echo "scale=1; $n/1000" | bc)"
  else
    printf "%d" "$n"
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

  # Color tracks proximity to the limit (the /context concern). 30/60 preserve
  # the prior absolute marks (300K/600K) on a 1M window while now scaling to
  # whatever `context_window_size` reports.
  if [ "$pct_int" -ge 60 ]; then
    color="\033[31m"
  elif [ "$pct_int" -ge 30 ]; then
    color="\033[33m"
  else
    color="\033[32m"
  fi

  CONTEXT=$(echo -e "${color}${used_h}/${limit_h} (${pct_int}%)\033[0m")
else
  CONTEXT=$(echo -e "\033[32m${used_h}\033[0m tkns")
fi

echo "󰚩 ${MODEL_DISPLAY} |  ${CONTEXT}"
