# House assets

The canonical engineering guidance lives under `~/.claude` and binds Codex too.

Before substantive work, read these files in full, preferably in one batched operation:

1. `~/.claude/CLAUDE.md`
2. Every Markdown file directly under `~/.claude/rules/`, in lexical order

If a required file cannot be read, stop and report it. Preserve the methods, constraints,
and standards in those files. Translate Claude-specific tool names, slash commands, agent
types, and interaction mechanics to the closest available Codex mechanism.

Do not duplicate these assets. Modify them only when a task explicitly targets the shared
house guidance or skills, and read `~/.claude/META.md` before editing them.

## Skills

Claude-owned skills live under `~/.claude/skills` and are exposed to Codex through tracked
symlinks under `~/.agents/skills`. Do not enumerate or preload their `SKILL.md` files during
initialization. Select skills from their name and description, then read the full skill only
when the task invokes or matches it.
