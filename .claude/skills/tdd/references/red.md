# RED Phase — Write a Failing Test

Goal: write ONE small, focused test that fails — a precise question about behavior that does
not yet exist. The next phase's only job is to answer that question, so the sharper the test,
the smaller the GREEN step.

## Rules

- Occurs exactly once per cycle (**RED** → GREEN → REFACTOR).
- Add exactly one failing test, or one new failing assertion to an existing test.
- **No commit during RED** — a red bar is not a safe state. The checkpoint comes at GREEN.

## Writing the test

**Name it for behavior, not implementation.** The name should survive a rewrite of the code
under test, because it describes *what* is guaranteed, not *how*.

- Good: `shouldAuthenticateValidUser`, `test_empty_input_returns_error`
- Bad: `testAuth`, `test1`

**Keep it to one behavior.** If you catch yourself writing "and also…", that's a second test —
write it on a later cycle. One behavior per test keeps the failure unambiguous and the GREEN
step minimal.

**Make it fail for the right reason.** The failure must be a *missing implementation* — the
assertion not yet satisfied — not a typo, a syntax error, or an unresolved import. Run the test
and read the failure message: confirm the assertion is what failed, not the setup. A test that
fails for the wrong reason proves nothing in GREEN.

## When to split into separate tests

- It exercises different behavior, not just another input to the same behavior.
- The new assertion needs materially different setup.
- It's unclear which of two assertions would fail first — split so each failure is one signal.

## Next phase

You do **not** pick a GREEN strategy here — that choice depends on confidence you'll only have
once you're staring at the implementation. The three strategies (Fake It / Obvious / Triangulation)
and how to choose live in `$HOME/.claude/skills/tdd/references/green.md`.

## Exit checklist

- [ ] Name describes the expected behavior
- [ ] Scoped to exactly one thing
- [ ] Test fails — verified by running it
- [ ] Fails for the right reason (assertion, not setup)

→ Proceed to GREEN (`$HOME/.claude/skills/tdd/references/green.md`).
