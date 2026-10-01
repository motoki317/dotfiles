# House assets

This file bootstraps Codex. The canonical cross-agent engineering guidance remains under
`~/.claude`.

Before the first substantive task in a session, read these files in full, preferably in one
batch:

1. `~/.claude/CLAUDE.md`
2. Every Markdown file directly under `~/.claude/rules/`, in lexical order

If you cannot read a required file, stop and report it. Apply the guidance as Codex policy for
the session. Translate Claude-specific tools, slash commands, agent types, and interaction
mechanics by intent to the closest available Codex capability. Do not emulate absent machinery.

Keep the source of truth under `~/.claude`. Do not duplicate it in Codex files. Unless a task
explicitly targets the shared guidance or skills, do not modify them. Before you modify them, read
`~/.claude/META.md`.

## Skills

`sync-claude-skills` exposes compatible Claude-owned skills as tracked symlinks under Codex's
native discovery directory, `~/.agents/skills`. Their canonical sources remain under
`~/.claude/skills`. The sync tool owns the compatibility boundary. Do not scan the Claude source
tree for additional skills.
