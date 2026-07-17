# REFACTOR Phase — Improve the Structure

Improve the code's structure while every test stays green; the green bar is the safety net that makes shape changes free. This phase loops — as many small refactors as the code warrants.

- Only refactor from a green bar; run the whole suite after every change.
- Structure only, never behavior: renames, extractions, moves, formatting, simplification. New functionality or a bug fix is the next RED, not a refactor — a reader must be able to trust that a `refactor:` commit touched no behavior (Tidy First).
- One small improvement per `refactor:` commit, separate from the GREEN commit: one rename, one extraction, or one move, roughly under 20 lines. Small steps mean a failed test points at a single change.
- On failure, don't debug the broken refactor — revert (`git checkout -- .`) and take a smaller step; the GREEN commit is always one command away.

Duplication is the primary signal that structure is missing: logic repeated between test and production code, magic values, similar patterns wanting a shared shape, names that hide intent.

Satisfied → the next behavior starts a new cycle at RED.
