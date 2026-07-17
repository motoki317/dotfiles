# RED Phase — Write a Failing Test

Write ONE small failing test — a precise question about behavior that does not yet exist. The sharper the question, the smaller the GREEN step that answers it.

- Add exactly one failing test, or one new failing assertion to an existing test.
- No commit during RED — a red bar is not a safe state; the checkpoint comes at GREEN.
- Name it for behavior, not implementation (`test_empty_input_returns_error`, not `test1`): the name should survive a rewrite of the code under test.
- One behavior per test. "And also…" is a second test on a later cycle. Split when a new assertion needs materially different setup, or when it's unclear which of two assertions would fail first.
- Make it fail for the right reason: the unsatisfied assertion, not a typo, syntax error, or missing import. Run it and read the failure message — a test failing on setup proves nothing in GREEN.

The GREEN strategy (Fake It / Obvious / Triangulation) is chosen in GREEN, once you're staring at the implementation — not here.
