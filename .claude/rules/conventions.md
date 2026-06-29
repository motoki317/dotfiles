# Conventions

## Tools
- GitHub: `gh` CLI for all operations.
- Missing command: retry with `nix run nixpkgs#<command>`.
- Text processing: `perl`, never `sed`/`awk`.
- Long output: `<cmd> | tee /tmp/<name>.log`, then grep the file.
- Run independent tasks in parallel — batch the independent tool calls into one message; prefer symbol-level edits over reading whole files.
- Output in English.

## Code & docs
- Code explains HOW (self-documenting); tests explain WHAT; comments and commit messages explain WHY and WHY-NOT, kept minimal — capture only non-obvious design decisions.
- Write so meaning outlives the session: rewrite a session-relative pointer like "the new logic" as the thing it computes, and test each line by whether a reader with only it, and no memory of today, could recover its meaning.

## Prose
- Use the `japanese-tech-writing` skill for all human-facing prose — comments in any language, reports, PR/issue text — and when revising existing prose.

## Cross-model review
- `codex-advisor` is the cross-model peer of the built-in `advisor` (a different model family, repo-cold); use it about as often as `advisor`, not as a rare fallback. Pair them at substantive checkpoints: before committing to an approach, before declaring non-trivial work done, when stuck, or on a hard-to-reverse call. Surface disagreements to the user. Skip both on mechanical or low-stakes work.

## Error handling
- Sub-agent fails → retry with an alternative.
- Conflicting outputs → flag the uncertainty to the user.
- Security or destructive risk → BLOCK and require acknowledgment.
