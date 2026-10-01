---
name: investigate
description: Enumerate competing hypotheses, verify each with evidence, and suggest a fix but never apply it. Use to diagnose an error, failing test, or anomaly before any fix is written.
argument-hint: "<error-or-symptom>"
allowed-tools: [Bash, Read, Grep, Glob, Agent, AskUserQuestion]
---

# Rules
- Read-only: never modify files or apply a fix. Suggest the change, and hand it to the Implementer. If you orchestrate, use `/codex-work`.
- Judge from evidence (logs and code first), never speculation. Build the chain from symptom to cause. If you cannot confirm the chain, say so in the report.
- Enumerate at least 3 hypotheses before you investigate any. This rule guards against a failure mode with two forms. In one, you fixate on the most recent change. In the other, you iterate on a failed approach and never re-examine the assumptions.
- Ask the user only for what you cannot infer: scope boundaries, reproduction conditions, when the issue began.

# Workflow
1. Gather:
   - error messages, stack traces, and logs
   - recent changes (`git log`, `git diff`)
   - the related code paths, config, and runtime environment

   Fan out read-only agents only if the surfaces are independent.
2. Hypothesize — ≥3 candidates ranked by likelihood, each with the evidence that suggests it.
3. Verify — take the hypotheses most likely first. For each one, define the evidence that confirms or refutes it. Gather that evidence. Record Confirmed / Refuted / Inconclusive. After each new piece of evidence, re-rank the rest.
4. Report:
   - the root cause, with its evidence chain and your confidence in it
   - impact, and similar at-risk code
   - the suggested fix (for the Implementer), with its risk
   - the rejected hypotheses and why
   - if the cause is unconfirmed, the open questions
