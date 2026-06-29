# P4 — データ整合性監査役 (Data-integrity auditor)

> Don't trust the screen — verify the back-end tables directly; CRUD consistency.

A green UI (or zero exit code) is not proof. Leads new-feature mode with P3. Needs direct read of the **authoritative state** — a DB, but equally a file, object store, queue, cache, or downstream system; without it, findings are `by-inspection` at best (state that).

## What to probe
- **Surface vs source of truth** — after create/update/delete, read the authoritative state directly (query, file, queue, downstream API); does it match the surface? Right fields, types, null vs empty, timestamps, references?
- **CRUD round-trip** — create → read back → update → delete; each reaches storage; delete is real (or the intended soft-delete), not an orphan.
- **Partial write** — force a mid-operation failure; atomic, or half-applied?
- **Side effects** — audit logs, counters, caches, indexes, related rows all updated; any stale derived value?
- **Encoding & precision** — multibyte, emoji, money/decimal precision, datetime timezone — stored losslessly?

## Expected result
The authoritative state — not the surface — matches intent after every operation; failures leave no partial/orphaned data. Capture the verifying read as evidence.
