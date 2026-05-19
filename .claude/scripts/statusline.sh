#!/usr/bin/env bash

# ref: https://zenn.dev/makotan/articles/db34bb6860cda5
#      https://zenn.dev/pnd/articles/claude-code-statusline

# Read JSON input from stdin
input=$(cat)

MODEL_DISPLAY=$(echo "$input" | jq -r '.model.display_name')
TRANSCRIPT_PATH=$(echo "$input" | jq -r '.transcript_path')

# display_name is the only stdin field that exposes the 1M variant — the
# transcript's .model field strips the [1m] suffix and can't be used here.
if [[ "$MODEL_DISPLAY" == *"1M"* ]]; then
  CONTEXT_MAX=1000000
else
  CONTEXT_MAX=200000
fi
# 85% matches Claude Code's current auto-compact default (~83.5%, rounded).
COMPACTION_THRESHOLD=$((CONTEXT_MAX * 85 / 100))

if [ -z "$TRANSCRIPT_PATH" ] || [ ! -f "$TRANSCRIPT_PATH" ]; then
  TOKEN_COUNT="_ tkns"
else
  total_tokens=$(tail -n 100 "$TRANSCRIPT_PATH" 2>/dev/null |
    jq -s 'map(select(.type == "assistant" and .message.usage)) |
    last |
    .message.usage |
    (.input_tokens // 0) +
    (.output_tokens // 0) +
    (.cache_creation_input_tokens // 0) +
  (.cache_read_input_tokens // 0)' 2>/dev/null)

  total_tokens=${total_tokens:-0}

  if [ "$total_tokens" -ge 1000000 ]; then
    millions=$(echo "scale=2; $total_tokens/1000000" | bc)
    token_display=$(printf "%.2fM" "$millions")
  elif [ "$total_tokens" -ge 1000 ]; then
    thousands=$(echo "scale=1; $total_tokens/1000" | bc)
    token_display=$(printf "%.1fK" "$thousands")
  else
    token_display="$total_tokens"
  fi

  # 1M models degrade well before compaction (Chroma context-rot research:
  # quality drops past ~300K, retrieval unreliable past ~600K), so color
  # tracks quality risk on 1M and compaction proximity on legacy 200k.
  if [ "$CONTEXT_MAX" -eq 1000000 ]; then
    if [ "$total_tokens" -ge 600000 ]; then
      color="\033[31m"
    elif [ "$total_tokens" -ge 300000 ]; then
      color="\033[33m"
    else
      color="\033[32m"
    fi
  else
    percentage=$((total_tokens * 100 / COMPACTION_THRESHOLD))
    if [ "$percentage" -ge 90 ]; then
      color="\033[31m"
    elif [ "$percentage" -ge 70 ]; then
      color="\033[33m"
    else
      color="\033[32m"
    fi
  fi

  TOKEN_COUNT=$(echo -e "${color}${token_display}\033[0m tkns")
fi

echo "󰚩 ${MODEL_DISPLAY} |  ${TOKEN_COUNT}"
