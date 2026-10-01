# P5: Migration specialist (移行担当者)

> Inject the old system's data — missing fields, malformed formats, char encoding, record-count match.

P5 feeds real legacy data through the new path, never clean synthetic data. P5 leads migration mode with P7. For a greenfield feature, P5 still asks: does this corrupt the existing data it touches?

## What to probe
- Real legacy shapes — rows with missing/null fields, deprecated formats, pre-validation values. Does the new path accept them, reject them, or silently mangle them?
- Format & encoding — different date formats, decimal separators, line endings, encodings (Shift_JIS↔UTF-8), trailing whitespace. Is there any mojibake or misparse?
- Count reconciliation — count in vs out. Does the new path drop, duplicate, or silently merge records?
- Outliers — longest value, oldest record, every-optional-field row, the known-bad production row.
- Idempotency — run it twice. Does the second run apply any change twice? Is it safe to resume after a partial failure?

## Expected result
The new path migrates, rejects with a reason, or reports every legacy record, and never loses or corrupts one. Counts reconcile, and a re-run is safe. Capture in/out counts and the reject log.
