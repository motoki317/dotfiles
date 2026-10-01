# P3: Malicious operator (悪意ある操作者)

> Boundary/invalid/out-of-permission values, double-submit — do validation and exclusion control hold?

P3 assumes an adversary, or a user who does the forbidden thing. P3 leads new-feature mode with P4. P3 is the exercised counterpart of the `/feedback` Security viewpoint.

## What to probe
- Boundary — min−1, min, max, max+1, zero, empty, one over any limit.
- Invalid / malformed — wrong type/encoding, control chars, oversized payload. Where input reaches a sink, add injection probes (SQL/command/path/HTML).
- Out-of-permission — act as a denied role. Reach another tenant's record by id. Call the endpoint past the UI gate.
- Double-submit / replay — the same mutation twice (fast, back-then-resubmit, replayed request). Does idempotency or exclusion control hold, or is there a duplicate/partial write?
- Bypass client checks — disable JS, edit the request, skip the prior step. Does the server still enforce the checks?

## Expected result
The server rejects every invalid, over-privileged, or duplicated action without a data change. Capture the rejection (status/message) as evidence.
