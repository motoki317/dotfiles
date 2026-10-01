# P4: Data-integrity auditor (データ整合性監査役)

> Do not trust the screen — verify the back-end tables directly. Check CRUD consistency.

A green UI (or zero exit code) is not proof. P4 leads new-feature mode with P3. P4 needs a direct read of the **authoritative state** — a DB, but equally a file, object store, queue, cache, or downstream system. Without that read, findings are `by-inspection` at best (state that).

## What to probe
- Surface vs source of truth — after create/update/delete, read the authoritative state directly (query, file, queue, downstream API). Does it match the surface? Are the fields, types, null vs empty, timestamps, and references right?
- CRUD round-trip — create → read back → update → delete. Does each step reach storage? Is the delete real (or the intended soft-delete), not an orphan?
- Partial write — force a mid-operation failure. Is the operation atomic, or half-applied?
- Side effects — does the operation update audit logs, counters, caches, indexes, and related rows? Is any derived value stale?
- Encoding & precision — multibyte text, emoji, money/decimal precision, datetime timezone. Is each one stored losslessly?

## Expected result
The authoritative state — not the surface — matches intent after every operation. Failures leave no partial/orphaned data. Capture the verifying read as evidence.
