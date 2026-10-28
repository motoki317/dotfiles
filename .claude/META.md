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
1. an environment fact the model cannot derive (`ax` exists, Codex has no `reviewer` agent type, who pushes).
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
- `~/.agents/skills/agent-browser` is a local rewrite (2026-10-13) of the upstream vercel-labs stub. `npx skills update` can overwrite it.

## Context-loading mechanics (verified 2026-07-11)

- Claude Code strips block-level HTML comments in `CLAUDE.md` before injection — documented, intended for maintainer notes. Claude Code preserves comments inside code blocks. `Read` shows comments as-is.
- Whether Claude Code also strips comments in `rules/*.md` is undocumented — do not rely on it. Put notes here instead.
- Reference: https://code.claude.com/docs/en/memory.md

## Credits

- `rules/code.md` adapts the "ponytail" skill by Dietrich Gebert (MIT): https://github.com/DietrichGebert/ponytail — the ladder, no-unrequested-abstractions, and never-simplify-away guardrails come from there. Local changes: a proven-library rung (we prefer mature libraries over hand-rolled code), root-cause fixes phrased as "the layer that owns the violated invariant", and idiom-not-verbosity repo matching. No-unrequested-abstractions and idiom-not-verbosity now live in `skills/quality-review/rubric.md` (Simplicity).
- The single-source-of-truth rule in `rules/code.md` and `refs/single-source-of-truth.md` adapts pstack by poteto (MIT, commit 6b78e9f4e46c75d1697d837e51f07c2692e60b24): https://github.com/cursor/plugins/tree/main/pstack
  - Sources: its Model the Domain principle (enum, table, dispatch point, state machine), "derive instead of sync" from Minimize Reader Load, and the "Hand-synced list" red flag of its architect skill.
  - Not adopted: Migrate Callers Then Delete Legacy APIs, the guard test of Fix Root Causes, and Attack the Premise. On 2026-10-19, one-off probes that `probes.md` does not keep showed that Claude and Codex already followed each one in 4 of 4 runs without a new rule.
- `skills/quality-review` (2026-10-13) merges the former `feedback`, which reviewed against the ISO/IEC 25010:2023 quality model, and `qa-review`, whose scenario angles adapt Nexta's "7 QA personas": https://zenn.dev/nexta_/articles/be13a2395a5d2a
- `skills/tech-writing` merges two sources into one language-neutral core plus language references (2026-08-03):
  - AminBlg/SimpleEnglish (MIT, commit 379728b51981b6d2ee1de0f201164483a9648972): https://github.com/AminBlg/SimpleEnglish — the ASD-STE100-derived sentence mechanics and `references/english.md`. Local changes:
    - We dropped strict mode and the dictionary apparatus.
    - We scoped the modal ladder so that epistemic uncertainty and counterfactuals survive (the global ban would erase them).
    - We scoped condition-first to procedural text.
    - A participle or gerund that acts as an adjective or a noun is correct, and a counterfactual keeps "would have".
    - We dropped the spelling rule and the paragraph cap (2026-10-07).
    - We dropped `scripts/ste_lint.py` (2026-10-13): no instruction referenced it.
  - k16shikano's Japanese writing norms (https://gist.github.com/k16shikano/fd287c3133457c4fd8f5601d34aa817d, formerly `skills/japanese-tech-writing`):
    - Its argument rigor and redundancy devices went language-neutral into the core.
    - Its 整形 and LLM口調 lists lived in `references/japanese.md` until 2026-10-03. That day the user chose to rebuild japanese.md from yomiyasu alone, which dropped the rules that only k16shikano held: one sentence per line, no 中黒 for lists, 「**用語**：説明」 glossaries, its LLM口調 words that yomiyasu lacks (正面から系, 空虚な形容, 空虚な動詞, 〜において), and the Japanese forms of the argument devices.
    - `references/manuscript.md` stayed unchanged until the 2026-10-01 yomiyasu integration, which removed its Japanese–Latin spaces and two lines that the core now holds.
    - We dropped its permission to soften a grounded claim for tone (「語調を整えるための意図的な緩和は許す」), because it contradicts the core hedge rule (2026-10-01).
    - On 2026-10-07 we cut manuscript.md to the rules that change a draft. We merged its unused-detail rules into one line, dropped its rules that the core holds, and dropped the person-name exception for historical figures.
- `skills/tech-writing` integrates nanaism/yomiyasu (MIT, and `LICENSE-yomiyasu-Unicode` for the Unicode-derived data in its linter): https://github.com/nanaism/yomiyasu — a Japanese rewrite skill against machine-sounding prose. Its author explains the design in https://zenn.dev/algoartis/articles/0b1c731881b25c (read 2026-10-07). The first integration pinned commit 30ee6041c328ce21d38a7963f667e079a93d7a12 (2026-10-01). The current pin is commit 8f77b7aa8f19c3f718b511e71a3eee5d06eb3013 (v1.0.8, 2026-10-07):
  - Its principle 1 (「誰が・何を・どうした」が1文で完結する, `references/gemini-syntax.md`) is the core's Complete check. Its 意味の保持, 情報の不増補, and セルフラベリングと否定対比の整理 (`SKILL.md`) are the core's rules for revising another writer's text. The English counterparts of its devices come from https://en.wikipedia.org/wiki/Wikipedia:Signs_of_AI_writing (read 2026-10-01).
  - `references/japanese.md` comes from its `SKILL.md` (本体), `references/gemini-syntax.md` (構文), `references/slop-catalog.md` (語彙), `references/domains/`, and the linter's word list alone. yomiyasu targets revision, so japanese.md marks its revision-only rules with 「書き直すときは」. Section sources, by the headings of commit 8d5abeeb:
    - 文書の立場と文末: 本体「文書の立場と文末」, 構文 原則2, domains/tech §3.
    - 誰が何をどうしたか: 本体「主語と目的語の明確化」「擬人化の解消」, 構文 原則1, 3, 13. v1.0.8 added コピー調.
    - 比喩動詞と英語直訳調: 本体「比喩動詞の具体化」, 構文 原則4, 語彙 §1–2.
    - 置き換える語: 語彙 §3–5, domains/essay §1–2, the linter's SLOP_WORDS. v1.0.8 added the role-inflating nouns.
    - 削る前置きと結び: 本体「セルフラベリングと否定対比の整理」, 構文 原則5–6, 語彙 §6.
    - 文の長さと読点: 本体「文長と読点の調整」, 構文 原則7–8.
    - 記号と空白: 本体「装飾記号と不要な空白の排除」, 構文 原則9, 11–12, 14.
    - 点検: 本体「2. 実行手順」Step 3–4.
  - The core's tense rule for a finished change, current behavior, and a plan comes from v1.0.8 (本体「文の働き」).
  - `scripts/yomiyasu_lint.py`, `scripts/yomiyasu_diff.py`, and `scripts/markdown_visibility.py` (the linter imports it) are vendored byte-identical. The 8d5abeeb linter got 26 of upstream's 34 counted bold-rendering fixtures wrong. The 8f77b7aa linter gets all 34 right.
  - Dropped: the tech/business/essay modes (their unique rules went into the shared sections), the 20% kanji-ratio target (unsourced, unmeasured), the 15% list cap and the bold-per-1,000-characters cap (the core's list and bold rules govern), the output format, the academic citations, and the linter's 0–100 score (Goodhart). On 2026-10-07 we also dropped the 勧め default for an undecidable stance, the replacement words for the flagged nouns, the 30–45-character average sentence length (no script measures it), and the essay words 結末 and 結晶.
  - Not adopted from v1.0.8: keep the source's Japanese–Latin spacing. We keep the no-space rule.
  - Not generalized: the em-dash ban (native English punctuation that the house style uses) and the trailing-colon ban (English lead-in colons are native).
  - Local changes:
    - We merged the rules that its files repeat and shortened the examples, so that SKILL.md plus japanese.md stay within 24,316 bytes, their size before the rebuild. openclaw loads both once per session.
    - We replaced suggested rewrites that our own rules flag, for example 掘り下げる for 踏み込む.
    - The idiom exemption covers only idioms with a body part (「骨が折れる」). A dictionary lists the figurative senses of 踏み込む and 深掘り, so a dictionary test keeps what the table rewrites.
    - The core's grounded-claim rule asserts a hedge that evidence grounds. A revision keeps the source's strength instead (意味の保持).
    - Enumeration 読点 do not count toward the limit of two. The Japanese–Latin spacing rule also covers code spans.
- `skills/tech-writing` adapts the method of "explain-before-shorten" by Kazuki Ota (https://zenn.dev/microsoft/articles/explain-before-shorten, published 2026-10-26, read 2026-10-28). The article states no license, so we took its method and wrote our own text and examples:
  - Adopted: its order (decide the point, explain it as if aloud, then shorten) as SKILL.md "Write" steps 1–4. Its check that the short text keeps every condition of the long explanation is step 4. Its table of over-shortened phrases is "Write the event, not a label".
  - The label principle replaced these former bullets: write the action as a verb, name the actor, replace a figurative verb, replace "as needed" with a condition, name each cause of an aggregate problem, and give the mechanism of a causal claim.
  - Not adopted: its slide render check (Own the outcome covers it), its output and missing-information tables (the revision rules cover them), and its trigger phrases (the description stays routing-only).
  - A/B test on 2026-10-28: four tasks (a Japanese slide, an English PR body, an English on-call message, a Japanese wiki section), two fresh `claude -p` runs per variant, and a blind judge per task. The first draft lost both Japanese tasks to over-compression, which step 4 now catches. The final version ranked above both old-skill outputs on all four tasks (8 of 8). Codex followed the final version on both Japanese tasks.
