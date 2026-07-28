# House assets

This file bootstraps Codex; the canonical cross-agent engineering guidance remains under
`~/.claude`.

Before the first substantive task in a session, read these files in full, preferably in one
batch:

1. `~/.claude/CLAUDE.md`
2. Every Markdown file directly under `~/.claude/rules/`, in lexical order

If any required file cannot be read, stop and report it. Apply the guidance as Codex policy
for the session. Translate Claude-specific tools, slash commands, agent types, and interaction
mechanics by intent to the closest available Codex capability; do not emulate absent machinery.

Keep the source of truth under `~/.claude`; do not duplicate it in Codex files. Modify the
shared guidance or skills only when a task explicitly targets them, and read
`~/.claude/META.md` first.

## Skills

`sync-claude-skills` exposes compatible Claude-owned skills as tracked symlinks under Codex's
native discovery directory, `~/.agents/skills`; their canonical sources remain under
`~/.claude/skills`. The sync tool owns the compatibility boundary; do not scan the Claude source
tree for additional skills.

## Tool Calls

In Code Mode, within each bounded stage, run independent, functions.exec-available tool calls concurrently in one functions.exec call. Use await Promise.allSettled([...]) when partial results are useful, and inspect every result; use await Promise.all([...]) only when any failure should abort the batch. Keep dependencies, waits/resumes, approvals, conflicting or interdependent mutations, and adaptive investigations where each result may change the next step sequential. Do not split otherwise batchable inspections across outer tool calls.
