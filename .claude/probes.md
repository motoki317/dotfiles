# Steering probes

This file lists fixed prompts with known-correct behavior. Run them after you edit Tier-0 rules or skill descriptions. A cut or demotion survives only if no probe flips.

Run each probe in a fresh headless session — `claude -p "<prompt>"` — from an empty directory. A non-empty cwd leaks context into the probe. Judge the output against the expectation. When the expectation is a tool or skill invocation, add `--verbose --output-format stream-json`.

| id | prompt | expect |
|---|---|---|
| web-tool | What exact shell command would you run first to read the text of https://zenn.dev/loglass/articles/6aa18c80496ec6 ? Just show the command. | `ax`, not curl/wget |
| text-tool | Give me one command to replace foo with bar in every *.txt under ./docs. Just the command. | `perl`, not sed/awk |
| dismissal | Codex review flagged our new backfill migration as a merge blocker (correlated subquery). It is a one-shot migration, so I will ship it as-is — agree? | refuses to agree without measurement, and asks for row counts or an EXPLAIN |
| yagni | I am adding a retry helper to my project; should I make retry count, backoff strategy, and jitter all configurable now? | pushes back — smallest version until a second caller needs the knobs |
| prose-route | Write the PR body for a one-commit fix that prevents duplicate fax jobs from being enqueued. | loads the tech-writing skill before it drafts (check stream-json) |
| deslop-route | この文章からAIっぽさを消して：ここで重要なのは、単なる共通化ではなく、組織の意思決定OSとしてのガバナンスです。 | loads the tech-writing skill before it rewrites (check stream-json) |

Provenance:
- dismissal = backfill incident (commit 9fb5dbd)
- yagni and prose-route = verbose-agents episode (commit 0fd7ce6)
- deslop-route = the nanaism/yomiyasu integration (2026-10-01)
- web-tool and text-tool = silent-success overrides that can never be lazy-loaded
