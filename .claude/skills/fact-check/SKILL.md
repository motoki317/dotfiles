---
name: fact-check
description: Verify claims against authoritative sources — Context7 for library/framework APIs, WebSearch/WebFetch for standards and cited URLs. Use when asked to fact-check, verify claims, validate documentation, or check source references.
---

# Fact Check

## Process
1. **Extract claims** — the verifiable assertions: library API behaviour, documentation references, standard compliance, version-specific behaviour.
2. **Select source** — library/framework API → Context7 (`resolve-library-id`, then `query-docs`; trust score 7+); web standard or spec → WebSearch (w3.org, MDN); cited URL → WebFetch; general technical fact → WebSearch.
3. **Assess** — score confidence 0–100 against the source and flag anything below 80. Cross-reference disputed claims with a second source; note version context for version-specific claims.

## Report per flagged claim
Claim / where it was made / evidence found (with its source) / confidence / recommended correction.

## Rules
- Never mark a claim verified without an actual source check; document the evidence source every time.
- Don't over-verify obvious facts.
