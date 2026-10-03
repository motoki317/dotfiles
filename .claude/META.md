# Maintaining the house assets

Claude Code auto-loads `~/.claude/CLAUDE.md` (plus `CLAUDE.local.md`) and every `.md` under `~/.claude/rules/`, recursively. Skill descriptions also enter every session's system prompt. This file is neither, so it never enters agent context — it is documentation for whoever edits the shared guidance or skills.

## File roles

- `CLAUDE.md` — values only. No operating rules.
- `rules/code.md` — the shape of the artifact: what code to write and not write.
- `rules/process.md` — how work flows: the Orchestrator/Implementer roles and their loops.
- `rules/conventions.md` — mechanics: tool choices.
- `skills/` — occasion workflows. A skill's body loads only when the skill is invoked or matched. Its description loads in every session.
- `refs/` — elaborations, recovery paths, and mechanics with zero standing cost, reachable only via a pointer from a rule or skill.
- `probes.md` — fixed steering probes that gate cuts and demotions.
- `~/.agents/skills/` — the cross-agent registry: real directories own external skills, and relative symlinks expose Claude-owned skills.

One role per file. A new rule goes to the file with the matching role. If no role fits, question whether it belongs at all.

## Editing principles

Two episodes ground these principles. In the verbose-agents episode (2026-07, commit 0fd7ce6), aspirational values did not steer, but operational devices did. In the backfill-dismissal incident (2026-08, commit 9fb5dbd), reviews detected the defect three times, but undisciplined disposition of the findings shipped a 40-minute migration. Their synthesis: the effective axis is checkable vs aspirational, not general vs detailed. Compress by subsumption into checkable generalizations ("dismissal is a claim too"), never by abstraction into vibes ("be rigorous").

A standing calibration behind the writing devices: the user's edit passes cut 50–90% of agent-written prose (measured 2026-07-30). A writer inside the context that produced a text cannot grade it.

A line survives in any asset only if it is one of:
1. an environment fact the model cannot derive (`ax` exists, the `reviewer` agent is read-only, who pushes).
2. an arbitrary choice among defensibles (perl not sed, one plan task per run) — state the bare choice, no justification.
3. a checkable override of a wrong model default (the ladder, dismissal-is-a-claim, cold-read). Phrase it as one of these operational devices:
   - a stop condition
   - a named anti-pattern
   - an output contract
   - a guardrail that makes the rule safe to follow hard

   State tensions resolved, not one-sided.

Delete this recurring fat on sight:
- workflow narration that the model performs anyway
- enumerated instances of a stated principle (keep the principle plus at most one anchor example)
- cross-file restatements
- rationale inside rules — the why lives here and in commit messages

## Placement — the tier model

- Tier 0, always loaded (`CLAUDE.md`, `rules/**`): values, the seat contract and loop skeletons, and only those overrides whose occasion is not self-announcing. In such an occasion, the wrong default silently succeeds, so nothing would ever prompt a lookup (sed works until it corrupts). `rules/` loads recursively, so keep refs out of it.
- Tier 1, skills (description always loaded, body on invocation): occasion workflows. The description is routing only — one sentence of what, plus the trigger condition. Method belongs in the body.
- Tier 2, `refs/` and skills' `references/` (zero standing cost): elaborations, edge cases, recovery paths, rubrics.
- Demotion test: if a detail's trigger fires before the mistake and comes from something unmissable, you can move the detail down. Otherwise, keep it in place. Examples of an unmissable source: an error message, a file type, a loop step, a skill invocation.
- Each occasion routes from exactly one place. Prefer the skill description over a rules line.

## Change discipline

- Subsumption budget: a lesson replaces text — name the lines that the new rule subsumes. Net Tier-0 growth needs explicit justification.
- Rules are not evidence that behavior changed. Gate every cut and demotion with `probes.md`: a change survives only if no probe flips.
- House rules and skills are procedures for a reader that cannot ask questions. Write them in the form that `skills/tech-writing` prescribes for their language. English takes STE form (imperative, condition before command, must/can). Japanese follows `references/japanese.md`. Vendored text keeps its source style.
- Models imitate the style of the prompt, so a writing skill teaches by example as much as by rule. Before you finish an edit to a non-vendored asset, run the mechanical pass (english.md "Self-check", japanese.md 「点検」) on the changed lines. Fix each hit that is a violation.
- Dedup across files: a sentence that lives in two files is a future contradiction.
- In skill files, reference bundled files by absolute `$HOME/.claude/...` path, never relative. A subagent does not always run from the skill's directory.

## Shared skill registry

- A real directory owns a skill. A symlink only exposes that skill to another agent. Never mirror a symlink or replace another owner's real directory.
- After you add, delete, or rename a top-level skill under `~/.claude/skills`, run `sync-claude-skills`. Include its link changes in the same commit.
- Before you finish an edit under `~/.claude/skills` or `~/.agents/skills`, run `sync-claude-skills --check`.

## Context-loading mechanics (verified 2026-07-11)

- Claude Code strips block-level HTML comments in `CLAUDE.md` before injection — documented, intended for maintainer notes. Claude Code preserves comments inside code blocks. `Read` shows comments as-is.
- Whether Claude Code also strips comments in `rules/*.md` is undocumented — do not rely on it. Put notes here instead.
- Reference: https://code.claude.com/docs/en/memory.md

## Credits

- `rules/code.md` adapts the "ponytail" skill by Dietrich Gebert (MIT): https://github.com/DietrichGebert/ponytail — the ladder, no-unrequested-abstractions, output contract, and never-simplify-away guardrails come from there. Local changes: a proven-library rung (we prefer mature libraries over hand-rolled code), root-cause fixes phrased as "the layer that owns the violated invariant", and idiom-not-verbosity repo matching.
- `skills/tech-writing` merges two sources into one language-neutral core plus language references (2026-08-03):
  - AminBlg/SimpleEnglish (MIT, commit 379728b51981b6d2ee1de0f201164483a9648972): https://github.com/AminBlg/SimpleEnglish — the ASD-STE100-derived sentence mechanics, `references/english.md`, and `scripts/ste_lint.py`. The linter is for before/after measurement only — aggregate counts, exits 0 regardless, not a proofreading tool. Local changes:
    - We dropped strict mode and the dictionary apparatus.
    - We scoped the modal ladder so that epistemic uncertainty and counterfactuals survive (the global ban would erase them).
    - We scoped condition-first to procedural text.
    - A repo spelling convention wins over American spelling.
  - k16shikano's Japanese writing norms (https://gist.github.com/k16shikano/fd287c3133457c4fd8f5601d34aa817d, formerly `skills/japanese-tech-writing`):
    - Its argument rigor and redundancy devices went language-neutral into the core.
    - Its 整形 and LLM口調 lists lived in `references/japanese.md` until 2026-10-03. That day the user chose to rebuild japanese.md from yomiyasu alone, which dropped the rules that only k16shikano held: one sentence per line, no 中黒 for lists, 「**用語**：説明」 glossaries, its LLM口調 words that yomiyasu lacks (正面から系, 空虚な形容, 空虚な動詞, 〜において), and the Japanese forms of the argument devices.
    - `references/manuscript.md` stayed unchanged until the 2026-10-01 yomiyasu integration, which removed its Japanese–Latin spaces and two lines that the core now holds.
    - We dropped its permission to soften a grounded claim for tone (「語調を整えるための意図的な緩和は許す」), because it contradicts the core hedge rule (2026-10-01).
- `skills/tech-writing` integrates nanaism/yomiyasu (MIT): https://github.com/nanaism/yomiyasu — a Japanese rewrite skill against machine-sounding prose. The first integration pinned commit 30ee6041c328ce21d38a7963f667e079a93d7a12 (2026-10-01). The current pin is commit 8d5abeebe2dd20c2db005deaddcc50be43c59c0a (2026-10-03):
  - Its principle 1 (「誰が・何を・どうした」が1文で完結する, `references/gemini-syntax.md`) is the core's Complete check. Its 意味の保持, 情報の不増補, and セルフラベリングと否定対比の整理 (`SKILL.md`) are the core's rules for revising another writer's text. The English counterparts of its devices come from https://en.wikipedia.org/wiki/Wikipedia:Signs_of_AI_writing (read 2026-10-01).
  - `references/japanese.md` comes from its `SKILL.md`, `references/gemini-syntax.md`, `references/slop-catalog.md`, `references/domains/`, and the linter's word list alone. Each section cites its source headings. yomiyasu targets revision, so japanese.md marks its revision-only rules with 「書き直すときは」.
  - `scripts/yomiyasu_lint.py` and `scripts/yomiyasu_diff.py` are vendored byte-identical.
  - Dropped: the tech/business/essay modes (their unique rules went into the shared sections), the 20% kanji-ratio target (unsourced, unmeasured), the 15% list cap and the bold-per-1,000-characters cap (the core's list and bold rules govern), the output format, the academic citations, and the linter's 0–100 score (Goodhart).
  - Not generalized: the em-dash ban (native English punctuation that the house style uses) and the trailing-colon ban (English lead-in colons are native).
  - Local changes:
    - We merged the rules that its files repeat and shortened the examples, so that SKILL.md plus japanese.md stay within 24,316 bytes, their size before the rebuild. openclaw loads both once per session.
    - We replaced suggested rewrites that our own rules flag: 掘り下げる for 踏み込む, and 不可欠な and 根幹となる for load-bearing.
    - The core's grounded-claim rule asserts a hedge that evidence grounds. A revision keeps the source's strength instead (意味の保持).
    - We replaced its connective example 「設定を変更した後は」, which drops the fact that someone changed the setting, with 「そのため」.
    - We stated two conventions that its own text follows but leaves unstated: enumeration 読点 do not count toward the limit of two, and code spans are exempt from the Japanese–Latin spacing rule.
