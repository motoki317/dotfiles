---
name: investigate
description: Enumerate competing hypotheses and verify each with evidence; suggests a fix but never applies it. Use to diagnose an error, failing test, or anomaly before any fix is written.
argument-hint: "<error-or-symptom>"
allowed-tools: [Bash, Read, Grep, Glob, Agent, AskUserQuestion]
---

# Rules
- Read-only: never modify files or apply a fix; suggest the change and hand it to the Implementer (`/codex-work` when orchestrating).
- Judge from evidence (logs and code first), never speculation; build the chain from symptom to cause, and report honestly when it can't be confirmed.
- Enumerate at least 3 hypotheses before investigating any — tunnel vision on the most recent change, or iterating on a failed approach instead of re-examining assumptions, is the failure mode this guards against.
- Ask the user only for what you can't infer: scope boundaries, reproduction conditions, when the issue began.

# Workflow
1. **Gather** — error messages, stack traces, and logs; recent changes (`git log`, `git diff`); the related code paths, config, and runtime environment. Fan out read-only agents only when the surfaces are genuinely independent.
2. **Hypothesize** — ≥3 candidates ranked by likelihood, each with the evidence that suggests it.
3. **Verify** — for each, most likely first: define what would confirm or refute it, gather that evidence, record Confirmed / Refuted / Inconclusive. Re-rank the rest whenever evidence lands.
4. **Report** — the root cause with its evidence chain and your confidence in it; impact and similar at-risk code; the suggested fix (for the Implementer) with its risk; the rejected hypotheses and why; open questions when the cause isn't confirmed.
