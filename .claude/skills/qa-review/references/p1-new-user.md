# P1 — 新人ユーザー (New user)

> Reads nothing, operates on intuition — does it break under misclicks, empty submits, rapid repeat-clicks?

The careless, impatient first-timer; on a non-GUI surface, the naive client (wrong arg order, omitted input, retry after a hang).

## What to probe
- **Empty / default submit** — submit with nothing filled or only defaults; clean rejection?
- **Wrong order** — actions out of sequence (submit before select, cancel mid-flow, back then resubmit); state stays coherent?
- **Impatient retry** — double-click submit, resend an apparently-hung request; any duplicate effect? (cross-note P3 double-submit)
- **Ignored guidance** — skip the step, dismiss the hint, paste where typing was expected; recovers or wedges?
- **Obvious-wrong input** — letters in a number field, far-future date, 1-char name; message understandable to a non-reader?

## Expected result
Naive actions are handled or refused with a plain message — never a stack trace, silent no-op, or corrupted state.
