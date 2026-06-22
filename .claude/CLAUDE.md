# Critical Rules

- Use `perl` for text processing (never `sed` or `awk`)
- Execute independent tasks in parallel (multiple Task calls in one message)
- Use symbol-level operations over reading entire files
- Output in English

# Coding

- Codes should explain itself, HOW things are done
- Test codes should explain WHAT the code does
- Code comments and commit messages should explain WHY a code exists and WHY NOT a feature / an alternative was not selected — keep them minimal: capture only essential, non-obvious design decisions and avoid unnecessary length

# Writing

- Use the `japanese-tech-writing` skill for all human-facing prose — code comments (any language), reports, PR/issue text — not just Japanese; its norms apply across languages.
- Always invoke it for Japanese technical prose, and when revising existing prose.
- Write so meaning outlives the session. Context from the current conversation makes your prose feel clear as you write; however, the next reader — human or fresh agent — has none of it, so distrust that clarity. Anchor every comment, commit, PR, and report to repository-visible facts, not to a session-relative role: rewrite a pointer like "the new logic" as the thing it computes. Test each line by whether a reader with only it, and no memory of today, could recover its meaning.

# Tool Usage

- **GitHub**: Use `gh` CLI for all operations (PRs, issues, repos)
- **Missing commands**: Retry with `nix run nixpkgs#<command>`
- **Long output**: Save with `<cmd> | tee /tmp/<name>.log` first, then grep the file. Lets you refine filters or inspect surrounding context without re-running the source command.

# Workflow

- Delegate detailed work to sub-agents; focus on orchestration
- Check existing code/patterns before implementing new features
- Only perform Git operations when explicitly requested
- Always verify the valid tool/library usage or configuration values by checking official docs. Don't guess.

# Cross-Model Review (Codex)

- Reach for the `codex-advisor` skill proactively on genuinely uncertain or hard-to-reverse calls (the skill covers when, read-only vs. apply, and how to weigh the result).

# Error Handling

- Sub-agent fails → retry with alternative agent
- Memory not found → document gap, investigate
- Conflicting outputs → flag uncertainty to user
- Security/destructive risk → BLOCK and require acknowledgment
