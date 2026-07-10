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
- `codex-advisor` (the `codex-consult` command) is your outside-view advisor — a different model family, so it catches blind spots a Claude-only analysis would share. Consult it on the cadence below.
- **When to fire.** Orientation first — finding files, reading, reproducing is not substantive work and needs no consult. Consult in the gap *after* orientation and *before* the first committing step: before locking in an approach or interpretation, before building on a load-bearing assumption, before declaring non-trivial work done, when stuck (errors recurring, results not converging), or when changing approach. Genuine uncertainty or a hard-to-reverse call turns it from optional into do-it-now.
- **How often.** At those checkpoints only — each call costs minutes, tokens, and egress to OpenAI, so deliberately skip mechanical or low-stakes work where a quick self-check suffices. Consulting less often than a free inline tool is expected, not a lapse.
- **How to consult.** `codex-consult --context` sends a redacted session snapshot — enough for a briefless trajectory review. Codex cannot see your thinking or injected context (CLAUDE.md, memory), so when your *reasoning* is what needs checking, put it in the brief and point Codex at the files. At a real gate, run it foreground and wait for the verdict; background only a non-gating opinion.
- **Handling the verdict.** Weigh it seriously; don't apply it blindly. If empirical evidence or a primary source contradicts it, adapt. When it disagrees with evidence you already hold, surface the conflict to the user rather than silently switching sides.

## Error handling
- Sub-agent fails → retry with an alternative.
- Conflicting outputs → flag the uncertainty to the user.
- Security or destructive risk → BLOCK and require acknowledgment.
