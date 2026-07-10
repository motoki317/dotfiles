# Process

## Operating loop
- Default flow: Explore → Plan → Implement → Verify → Commit → Ship. Auto-advance through Implement, Verify, and Commit on your own; Ship stays user-triggered — first publishing a branch and opening a PR. Once a PR is open, follow-up pushes to it are part of closing the loop, not a new Ship: the `feedback`/`qa-review`/`address` skills may update it via `/rebase-clean` (`--force-with-lease`).
- Verify before declaring done: run the available tests/build/lint, then `feedback` (static review) and `qa-review` (dynamic QA, when there is runnable behavior) — for analysis only; their close-the-loop phase fires on an explicit `/feedback` or `/qa-review`, not on this routine check. Report the commands run and what they returned, not just "done".
- Ask only when blocked by: a value or preference you can't infer, missing secrets/credentials, a scope boundary you can't infer, or a destructive or irreversible action — data loss, force-push (except the open-PR follow-up push above), deploy, or publishing to an external service.

## Cross-model review
- `codex-consult` (the codex-advisor skill) is your outside-view advisor — a different model family, so it catches blind spots a Claude-only analysis would share. Trigger it yourself at the gates below; don't wait to be asked.
- **Gates.** In the gap after orientation and before the first committing step: locking in an approach or interpretation, building on a load-bearing assumption, declaring non-trivial work done, stuck (errors recurring, results not converging), or changing approach. Genuine uncertainty or a hard-to-reverse call makes it mandatory; each call costs minutes, tokens, and egress to OpenAI, so skip mechanical or low-stakes work — consulting less often than a free tool is expected.
- **How.** `codex-consult --context` sends a redacted session snapshot. Codex cannot see your thinking or injected context (CLAUDE.md, memory) — when the reasoning is what needs checking, write it into the brief and point Codex at the files. Run gates foreground and wait for the verdict; background only a non-gating opinion.
- **Verdict.** Weigh it, don't obey it: empirical evidence and primary sources win, and when the verdict contradicts evidence you already hold, surface the conflict to the user rather than silently switching sides.

## Error handling
- Sub-agent fails → retry with an alternative.
- Conflicting outputs → flag the uncertainty to the user.
- Security or destructive risk → BLOCK and require acknowledgment.
