# GREEN Phase — Make the Test Pass

Make the one failing test from RED pass with the minimal code that does so. Quality is deliberately deferred: GREEN proves the behavior exists; REFACTOR improves its shape without re-asking whether it works. Mixing the two is how a green bar turns into an unreviewable diff.

- Add nothing the failing test does not demand — unbuilt functionality is unproven, and it is the next RED, not this GREEN.
- Tolerate ugly code for now: duplication, poor names, and hard-coded values are what REFACTOR exists for.

## Strategy — pick by your confidence in the implementation

- **Obvious Implementation**: type the real solution when it's clear and small enough to get right in one go. If it unexpectedly fails, your "obvious" was a guess — shrink the step and fall back to Fake It.
- **Fake It**: return the constant the test expects when the path to a real implementation is unclear; a passing test, even faked, turns an open design question into a concrete starting point.
- **Triangulation**: after faking, add a second test the fake cannot satisfy — two concrete examples pin down a generalization one example left ambiguous.

## Commit

A passing suite is the safe state — the one point in the cycle you can always return to. Run the whole suite, then commit (`feat:`/`fix:`) before touching structure in REFACTOR.
