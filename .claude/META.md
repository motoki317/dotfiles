# Maintaining the house assets

Claude Code auto-loads only `~/.claude/CLAUDE.md` (plus `CLAUDE.local.md`) and every `.md` under `~/.claude/rules/`, recursively. This file is neither, so it never enters agent context — it is documentation for whoever edits the shared guidance or skills.

## File roles

- `CLAUDE.md` — values only, the soul. No operating rules.
- `rules/code.md` — the shape of the artifact: what code to write and not write.
- `rules/process.md` — how work flows: the Orchestrator/Implementer roles, their loops, cross-model review gates.
- `rules/conventions.md` — mechanics: tools, prose.
- `skills/` — task-specific workflows, loaded only when invoked or matched.
- `../.agents/skills/` — the cross-agent registry: real directories own external skills; relative symlinks expose Claude-owned skills.

One role per file. A new rule goes where its role says; if no role fits, question whether it belongs at all.

## Editing principles

Learned 2026-07 while diagnosing why agents stayed verbose despite the values ("the verbose-agents episode", commit 0fd7ce6):

- Aspirational values do not steer generation; operational devices do. Prefer an ordered ladder with a stop condition, named anti-patterns ("no interface with one implementation"), an output contract, and explicit guardrails that make a rule safe to follow hard.
- State tensions resolved, not one-sided. "Rigor over expedience" alone licensed over-building; "rigor in understanding, economy in the artifact" bounds both directions.
- Every line must pay behavioral rent. Before adding one, name the agent behavior it changes; after editing, compare diffs on a few fixed representative prompts — rules are not evidence that behavior moved.
- Dedup across files: a sentence living in two files is a future contradiction.
- House rules and skills are procedures for a reader that cannot ask questions — write new English ones in STE form (imperative, condition before command, must/can; `skills/tech-writing`). Vendored and Japanese text keeps its source style.
- Budget: always-loaded total was 90 lines as of 2026-07-17. Steady growth is the pitfall these principles exist to prevent.

## Shared skill registry

- A real directory owns a skill. A symlink only exposes that skill to another agent; never mirror a symlink or replace another owner's real directory.
- After adding, deleting, or renaming a top-level skill under `~/.claude/skills`, run `sync-claude-skills` and include its link changes in the same commit.
- Before finishing an edit under `~/.claude/skills` or `~/.agents/skills`, run `sync-claude-skills --check`.

## Context-loading mechanics (verified 2026-07-11)

- Block-level HTML comments in `CLAUDE.md` are stripped before injection — documented, intended for maintainer notes. Comments inside code blocks are preserved; `Read` shows comments as-is.
- Whether comments in `rules/*.md` are also stripped is undocumented — do not rely on it; put notes here instead.
- Reference: https://code.claude.com/docs/en/memory.md

## Credits

- `rules/code.md` adapts the "ponytail" skill by Dietrich Gebert (MIT): https://github.com/DietrichGebert/ponytail — the ladder, no-unrequested-abstractions, output contract, and never-simplify-away guardrails come from there. Local changes: a proven-library rung (we prefer mature libraries over hand-rolling), root-cause fixes phrased as "the layer that owns the violated invariant", and idiom-not-verbosity repo matching — the last three refined by an OpenAI Codex cross-review (2026-07-11).
- `skills/tech-writing` merges two sources into one language-neutral core plus language references (2026-08-03):
  - AminBlg/SimpleEnglish (MIT, commit 379728b51981b6d2ee1de0f201164483a9648972): https://github.com/AminBlg/SimpleEnglish — the ASD-STE100-derived sentence mechanics, `references/english.md`, and `scripts/ste_lint.py`. Local changes: dropped strict mode and the dictionary apparatus; modal ladder scoped so epistemic uncertainty and counterfactuals survive (the global ban would erase them — flagged independently by a Codex cross-review, 2026-08-03); condition-first scoped to procedural text; repo spelling convention wins over American spelling. The linter is for before/after measurement only — aggregate counts, exits 0 regardless, not a proofreading tool.
  - k16shikano's Japanese writing norms (gist, formerly `skills/japanese-tech-writing`): argument rigor and redundancy devices went language-neutral into the core; 整形 and the LLM 口調 lists kept, excerpted and regrouped, in `references/japanese.md`; `references/manuscript.md` unchanged.
