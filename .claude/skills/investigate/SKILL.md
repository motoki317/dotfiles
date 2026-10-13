---
name: investigate
description: Enumerate competing hypotheses and confirm the cause with evidence. Use to diagnose an error, failing test, or anomaly before any fix is written.
argument-hint: "<error-or-symptom>"
---

# Rules
- Edit nothing until one hypothesis is Confirmed. If none is, report that and the open questions.
- Before you test any hypothesis, write down at least 3, ranked, each with the evidence that suggests it. One must not be the most recent change. One must place the cause outside this repo: search the dependency's upstream issues.
- Test the most likely first. Mark each Confirmed, Refuted, or Inconclusive with its evidence, and re-rank after each result.
- A Confirmed cause must explain every reported symptom under the user's reported conditions.

# Report
- the root cause, with its evidence chain and your confidence
- the reported symptoms that it leaves unexplained
- the suggested fix and its risk, plus similar at-risk code
- the rejected hypotheses and why
