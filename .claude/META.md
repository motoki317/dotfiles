# Maintaining the house assets

Claude Code auto-loads `~/.claude/CLAUDE.md` (plus `CLAUDE.local.md`) and every `.md` under `~/.claude/rules/`, recursively; skill descriptions also enter every session's system prompt. This file is neither, so it never enters agent context — it is documentation for whoever edits the shared guidance or skills.

## File roles

- `CLAUDE.md` — values only, the soul. No operating rules.
- `rules/code.md` — the shape of the artifact: what code to write and not write.
- `rules/process.md` — how work flows: the Orchestrator/Implementer roles and their loops.
- `rules/conventions.md` — mechanics: tool choices.
- `skills/` — occasion workflows, loaded only when invoked or matched; the description is always loaded.
- `refs/` — elaborations, recovery paths, and mechanics with zero standing cost; reachable only via a pointer from a rule or skill.
- `probes.md` — fixed steering probes that gate cuts and demotions.
- `~/.agents/skills/` — the cross-agent registry: real directories own external skills; relative symlinks expose Claude-owned skills.

One role per file. A new rule goes where its role says; if no role fits, question whether it belongs at all.

## Editing principles

Two grounding episodes: the verbose-agents episode (2026-07, commit 0fd7ce6 — aspirational values did not steer; operational devices did) and the backfill-dismissal incident (2026-08, commit 9fb5dbd — reviews detected the defect three times; undisciplined disposition of the findings shipped a 40-minute migration). Their synthesis: the effective axis is checkable vs aspirational, not general vs detailed. Compress by subsumption into checkable generalizations ("dismissal is a claim too"), never by abstraction into vibes ("be rigorous"). A standing calibration behind the writing devices: the user's edit passes cut 50–90% of agent-written prose (measured 2026-07-30) — the producing context cannot self-grade.

A line survives in any asset only if it is one of:
1. an environment fact the model cannot derive (`ax` exists, the `reviewer` agent is read-only, who pushes);
2. an arbitrary choice among defensibles (perl not sed, one run per plan) — state the bare choice, no justification;
3. a checkable override of a wrong model default (the ladder, dismissal-is-a-claim, cold-read) — phrased as an operational device: a stop condition, a named anti-pattern, an output contract, or a guardrail that makes the rule safe to follow hard; state tensions resolved, not one-sided.

Recurring fat, deleted on sight: workflow narration the model performs anyway; enumerated instances of a stated principle (keep the principle plus at most one anchor example); cross-file restatements; rationale inside rules — the why lives here and in commit messages.

## Placement — the tier model

- Tier 0, always loaded (`CLAUDE.md`, `rules/**` — rules/ loads recursively, so keep refs out of it): values, the seat contract and loop skeletons, and only those overrides whose occasion is not self-announcing — the wrong default silently succeeds, so nothing would ever prompt a lookup (sed works until it corrupts).
- Tier 1, skills (description always loaded, body on invocation): occasion workflows. The description is routing only — one sentence of what, plus the trigger condition; method belongs in the body.
- Tier 2, `refs/` and skills' `references/` (zero standing cost): elaborations, edge cases, recovery paths, rubrics.
- Demotion test: move a detail down only if its trigger fires before the mistake, from something unmissable — an error message, a file type, a loop step, a skill invocation.
- Each occasion routes from exactly one place; prefer the skill description over a rules line.

## Change discipline

- Subsumption budget: a lesson lands by replacing text — name the lines the new rule subsumes; net Tier-0 growth needs explicit justification.
- Rules are not evidence that behavior moved: gate every cut and demotion with `probes.md` — a change survives only if no probe flips.
- House rules and skills are procedures for a reader that cannot ask questions — write new English ones in STE form (imperative, condition before command, must/can; `skills/tech-writing`). Vendored and Japanese text keeps its source style; an existing English line keeps its style until rewritten wholesale.
- Dedup across files: a sentence living in two files is a future contradiction.
- In skill files, reference bundled files by absolute `$HOME/.claude/...` path, never relative — a subagent may not run from the skill's directory.

## Shared skill registry

- A real directory owns a skill. A symlink only exposes that skill to another agent; never mirror a symlink or replace another owner's real directory.
- After adding, deleting, or renaming a top-level skill under `~/.claude/skills`, run `sync-claude-skills` and include its link changes in the same commit.
- Before finishing an edit under `~/.claude/skills` or `~/.agents/skills`, run `sync-claude-skills --check`.

## Context-loading mechanics (verified 2026-07-11)

- Block-level HTML comments in `CLAUDE.md` are stripped before injection — documented, intended for maintainer notes. Comments inside code blocks are preserved; `Read` shows comments as-is.
- Whether comments in `rules/*.md` are also stripped is undocumented — do not rely on it; put notes here instead.
- Reference: https://code.claude.com/docs/en/memory.md

## Credits

- `rules/code.md` adapts the "ponytail" skill by Dietrich Gebert (MIT): https://github.com/DietrichGebert/ponytail — the ladder, no-unrequested-abstractions, output contract, and never-simplify-away guardrails come from there. Local changes: a proven-library rung (we prefer mature libraries over hand-rolling), root-cause fixes phrased as "the layer that owns the violated invariant", and idiom-not-verbosity repo matching.
- `skills/tech-writing` merges two sources into one language-neutral core plus language references (2026-08-03):
  - AminBlg/SimpleEnglish (MIT, commit 379728b51981b6d2ee1de0f201164483a9648972): https://github.com/AminBlg/SimpleEnglish — the ASD-STE100-derived sentence mechanics, `references/english.md`, and `scripts/ste_lint.py`. Local changes: dropped strict mode and the dictionary apparatus; modal ladder scoped so epistemic uncertainty and counterfactuals survive (the global ban would erase them); condition-first scoped to procedural text; repo spelling convention wins over American spelling. The linter is for before/after measurement only — aggregate counts, exits 0 regardless, not a proofreading tool.
  - k16shikano's Japanese writing norms (gist, formerly `skills/japanese-tech-writing`): argument rigor and redundancy devices went language-neutral into the core; 整形 and the LLM 口調 lists kept, excerpted and regrouped, in `references/japanese.md`; `references/manuscript.md` unchanged.
