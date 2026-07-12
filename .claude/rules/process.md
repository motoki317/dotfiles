# Process

## Operating loop
1. **Explore** — read the relevant flow end to end.
2. **Plan** — `/define` for requirements; `/decision-analysis` for an open choice.
3. **Implement** — `/execute` to orchestrate; `/tdd` for testable code; `/frontend-design` for web UI.
4. **Verify** — run the available tests/build/lint, then `/feedback` (static review) and `/qa-review` (dynamic QA, when there is runnable behavior) — analysis only; their close-the-loop phase fires on an explicit invocation, not this routine check. Report the commands run and what they returned, not just "done".
5. **Commit** — `/commit`.
6. **Tidy** — `/rebase-clean` the branch's history.
7. **Ship** — `/pr`; then `/address` for the PR's reviews and CI.

- Stop only after step 2 (plan approval) and step 6 (Ship is user-triggered); auto-advance every other step. Once a PR is open, follow-up pushes close the loop, not a new Ship: `/feedback`, `/qa-review`, and `/address` may update it via `/rebase-clean` (`--force-with-lease`).
- Ask only when blocked by: a value, preference, or scope boundary you can't infer; missing secrets/credentials; or a destructive or irreversible action — data loss, force-push (except the open-PR follow-up above), deploy, publishing to an external service.

## Cross-model review
- At the gates below, consult `/codex-advisor` (`codex-consult`) on your own — a different model family catches blind spots a Claude-only analysis would share. Mechanics live in the skill.
- **Gates** — after orientation, before the first committing step: locking in an approach or interpretation, building on a load-bearing assumption, declaring non-trivial work done, stuck (errors recurring, results not converging), or changing approach. Genuine uncertainty or a hard-to-reverse call makes it mandatory; each call costs minutes, tokens, and egress to OpenAI, so skip mechanical or low-stakes work.
- **Verdict** — weigh it, don't obey it: empirical evidence and primary sources win; when it contradicts evidence you already hold, surface the conflict to the user rather than silently switching sides.

## Error handling
- Sub-agent fails → retry with an alternative.
- Conflicting outputs → flag the uncertainty to the user.
- Security or destructive risk → BLOCK and require acknowledgment.
