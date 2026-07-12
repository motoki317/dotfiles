# House assets

Standing guidance for every Codex run on this machine. The engineering assets are shared
with Claude Code and live under `~/.claude`; they bind Codex too.

Mandatory initialization before any edit or commit:

1. Read `~/.claude/rules/code.md` (coding principles) and `~/.claude/rules/conventions.md`
   (tool and prose conventions) in full.
2. List `~/.claude/skills/` and read `commit/SKILL.md` (commit-message format); read
   `tdd/SKILL.md` for testable code, `japanese-tech-writing/SKILL.md` before writing
   human-facing prose, `frontend-design/SKILL.md` for web UI, and any other skill relevant
   to the task.
3. If a mandatory asset cannot be read, stop and report it.

These are read-only references written for Claude Code — treat tool mechanics as advisory,
the methods and standards as binding; never modify them. Do not begin coding until this
initialization is complete.
