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

- Use the `japanese-tech-writing` skill when writing prose for human readers — code comments (any language), reports, and PR/issue descriptions. Its norms (one topic per paragraph, rigorous argument, no LLM filler, no redundancy) apply beyond Japanese.
- Always invoke it for Japanese technical prose, and when revising existing prose.
- Write so meaning outlives the session. Context from the current conversation makes your prose feel clear as you write; however, the next reader — human or fresh agent — has none of it, so distrust that clarity. Anchor every comment, commit, PR, and report to repository-visible facts, not to a session-relative role: rewrite a pointer like "the new logic" as the thing it computes. Test each line by whether a reader with only it, and no memory of today, could recover its meaning.

# Tool Usage

- **GitHub**: Use `gh` CLI for all operations (PRs, issues, repos)
- **Docs**: Use Context7 MCP to verify latest library documentation
- **Missing commands**: Retry with `nix run nixpkgs#<command>`
- **Long output**: Save with `<cmd> | tee /tmp/<name>.log` first, then grep the file. Lets you refine filters or inspect surrounding context without re-running the source command.

# Workflow

- Delegate detailed work to sub-agents; focus on orchestration
- Check existing code/patterns before implementing new features
- Only perform Git operations when explicitly requested
- Always verify the valid tool/library usage or configuration values by checking official docs. Don't guess.

# Cross-Model Review (Codex)

- Consult the `codex-advisor` skill on your own initiative for an independent, cross-model second opinion. Its value is the outside view: reach for it when a decision is genuinely uncertain or hard to reverse, when you and the built-in Claude advisor diverge, or before finishing non-trivial work — and skip it for routine or mechanical work, which only burns quota.
- Keep consults read-only by default so an unattended one cannot change your tree; reserve the full-access `cdy` alias for when you are explicitly asked to have Codex apply changes.
- Treat the verdict as a peer's, not an oracle's — weigh it against your own view and the advisor's, and surface genuine disagreement to the user.

# Error Handling

- Sub-agent fails → retry with alternative agent
- Memory not found → document gap, investigate
- Conflicting outputs → flag uncertainty to user
- Security/destructive risk → BLOCK and require acknowledgment
