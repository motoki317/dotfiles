# P1: New user (新人ユーザー)

> Reads nothing, operates on intuition — does it break under misclicks, empty submits, rapid repeat-clicks?

P1 is the careless, impatient first-timer. On a non-GUI surface, P1 is the naive client (wrong arg order, omitted input, retry after a hang).

## What to probe
- Empty / default submit — submit with nothing filled, or with only defaults. Does the target reject it cleanly?
- Wrong order — actions out of sequence (submit before select, cancel mid-flow, back then resubmit). Does the state stay coherent?
- Impatient retry — double-click submit, resend an apparently-hung request. Does any duplicate effect occur? (cross-note P3 double-submit)
- Ignored guidance — skip the step, dismiss the hint, paste where the target expects typed input. Does the target recover, or does it wedge?
- Obvious-wrong input — letters in a number field, far-future date, 1-char name. Can a non-reader understand the message?

## Expected result
The target handles each naive action, or refuses it with a plain message. The outcome is never a stack trace, a silent no-op, or corrupted state.
