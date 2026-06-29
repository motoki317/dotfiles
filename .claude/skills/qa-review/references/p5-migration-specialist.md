# P5 — 移行担当者 (Migration specialist)

> Inject the old system's data — missing fields, malformed formats, char encoding, record-count match.

Feeds real legacy data through the new path, never clean synthetic data. Leads migration mode with P7. For a greenfield feature, still asks: does this corrupt the existing data it touches?

## What to probe
- **Real legacy shapes** — rows with missing/null fields, deprecated formats, pre-validation values; accepted, rejected, or silently mangled?
- **Format & encoding** — differing date formats, decimal separators, line endings, encodings (Shift_JIS↔UTF-8), trailing whitespace; any mojibake or misparse?
- **Count reconciliation** — count in vs out; records dropped, duplicated, or silently merged?
- **Outliers** — longest value, oldest record, every-optional-field row, the known-bad production row.
- **Idempotency** — run twice; double-applied? safe to resume after a partial failure?

## Expected result
Every legacy record is migrated, rejected-with-reason, or reported — never lost or corrupted; counts reconcile; re-running is safe. Capture in/out counts and the reject log.
