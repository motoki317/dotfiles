# P2 — ベテラン現場担当 (Veteran operator)

> Fast, high-volume keyboard input — Tab nav, shortcuts, Enter mid-IME.

The power user who outruns the UI; generalizes to high-throughput / concurrent / scripted use (batch calls, automation that doesn't wait).

## What to probe
- **Keyboard-only** — whole task via Tab/Enter/shortcuts; focus order right, all controls reachable, no shortcut collision?
- **Enter mid-IME** — Enter with a CJK candidate open should commit text, not submit.
- **Ahead of the UI** — act before async load finishes; submit while a spinner is up.
- **Volume / batch** — large paste, many rows, big input, back-to-back requests; throughput holds, no truncation or lock contention?
- **Concurrency** — two sessions/requests on one record at once. (cross-note P4 CRUD)

## Expected result
The fast/bulk/keyboard path yields the same correct result as the slow/mouse path — no dropped input, premature submit, focus trap, or race.
