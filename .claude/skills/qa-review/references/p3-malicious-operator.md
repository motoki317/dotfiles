# P3 — 悪意ある操作者 (Malicious operator)

> Boundary/invalid/out-of-permission values, double-submit — do validation and exclusion control hold?

Assumes an adversary, or a user who does the forbidden thing. Leads new-feature mode with P4; pairs with the `/feedback` Security viewpoint, but exercised.

## What to probe
- **Boundary** — min−1, min, max, max+1, zero, empty, one over any limit.
- **Invalid / malformed** — wrong type/encoding, control chars, oversized payload; injection probes (SQL/command/path/HTML) where input reaches a sink.
- **Out-of-permission** — act as a denied role; reach another tenant's record by id; call the endpoint past the UI gate.
- **Double-submit / replay** — same mutation twice (fast, back-then-resubmit, replayed request); idempotency / exclusion control, or a duplicate/partial write?
- **Bypass client checks** — disable JS, edit the request, skip the prior step; does the server still enforce?

## Expected result
Every invalid, over-privileged, or duplicated action is rejected server-side with no data change. Capture the rejection (status/message) as evidence.
