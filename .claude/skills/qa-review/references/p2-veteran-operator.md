# P2: Veteran operator (ベテラン現場担当)

> Fast, high-volume keyboard input — Tab nav, shortcuts, Enter mid-IME.

P2 is the power user who acts faster than the UI responds. P2 generalizes to high-throughput / concurrent / scripted use (batch calls, automation that does not wait).

## What to probe
- Keyboard-only — do the whole task via Tab/Enter/shortcuts. Is the focus order right? Is every control reachable? Does any shortcut collide?
- Enter mid-IME — Enter with a CJK candidate open must commit text, not submit.
- Ahead of the UI — act before an async load finishes. Submit while a spinner is up.
- Volume / batch — large paste, many rows, big input, back-to-back requests. Does throughput hold? Is there any truncation or lock contention?
- Concurrency — two sessions/requests on one record at once. (cross-note P4 CRUD)

## Expected result
The fast/bulk/keyboard path yields the same correct result as the slow/mouse path — no dropped input, premature submit, focus trap, or race.
