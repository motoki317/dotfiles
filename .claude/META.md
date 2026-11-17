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
- Models imitate the style of the prompt, so a writing skill teaches by example as much as by rule. Before you finish an edit to a non-vendored asset, run the check of `skills/tech-writing` (english.md "Self-check", japanese.md 「点検」) on the changed lines. Fix each hit that is a violation.
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
    - The 2026-11-17 ablation (below) kept the active voice, no contractions or semicolons, and the modal rules. Without the skill, writers broke them in all four English outputs. It dropped simple tenses, participle clauses, Latin abbreviations, the word caps, the filler table, and the term sets. No output in any arm used a filler-table word, and the no-skill outputs that broke the other rules ranked as high as the skill outputs.
  - k16shikano's Japanese writing norms (https://gist.github.com/k16shikano/fd287c3133457c4fd8f5601d34aa817d, formerly `skills/japanese-tech-writing`):
    - Its argument rigor and redundancy devices went language-neutral into the core.
    - Its 整形 and LLM口調 lists lived in `references/japanese.md` until 2026-10-03. That day the user chose to rebuild japanese.md from yomiyasu alone, which dropped the rules that only k16shikano held: one sentence per line, no 中黒 for lists, 「**用語**：説明」 glossaries, its LLM口調 words that yomiyasu lacks (正面から系, 空虚な形容, 空虚な動詞, 〜において), and the Japanese forms of the argument devices.
    - `references/manuscript.md` stayed unchanged until the 2026-10-01 yomiyasu integration, which removed its Japanese–Latin spaces and two lines that the core now holds.
    - We dropped its permission to soften a grounded claim for tone (「語調を整えるための意図的な緩和は許す」), because it contradicts the core hedge rule (2026-10-01).
    - On 2026-10-07 we cut manuscript.md to the rules that change a draft. We merged its unused-detail rules into one line, dropped its rules that the core holds, and dropped the person-name exception for historical figures.
    - The 2026-11-17 ablation cut manuscript.md to 整形 and 演出と見出し. Two no-skill article sections broke none of the paragraph, argument, reader-load, and viewpoint rules, but broke the drama or heading rules in both audits (two and three hits). 読者への誠実さ and the appositive and new-example rules had no occasion in the test, and we cut them untested. The term rule moved to japanese.md「語を選ぶ」.
- `skills/tech-writing` integrates nanaism/yomiyasu (MIT, and `LICENSE-yomiyasu-Unicode` for the Unicode-derived data in its linter): https://github.com/nanaism/yomiyasu — a Japanese rewrite skill against machine-sounding prose. Its author explains the design in https://zenn.dev/algoartis/articles/0b1c731881b25c (read 2026-10-07). Earlier pins were commit 30ee6041c328ce21d38a7963f667e079a93d7a12 (2026-10-01) and commit 8f77b7aa8f19c3f718b511e71a3eee5d06eb3013 (v1.0.8, 2026-10-07). The current pin is commit db78b9bc5e868549b517a037df7a4d1c93edac23 (v3.1.0, 2026-11-12). Its `LICENSE` and `UNICODE-LICENSE.txt` did not change.
  - `scripts/yomiyasu_lint.py`, `scripts/yomiyasu_diff.py`, and `scripts/markdown_visibility.py` (the lint imports it) are vendored byte-identical from its `skills/yomiyasu/scripts/`. Not vendored: `yomiyasu_check.py` and `yomiyasu_hook.py`, which keep records in `.yomiyasu/` for its full revision workflow, and `references/lists/*.txt`, which the lint has built in.
  - v3 lint and diff refuse a file argument outside the working directory (exit 2). japanese.md「点検」 therefore passes the lint its input on stdin and runs the diff from the scratch folder. `--doc-type` turns on its pause_comma rule.
  - `words/` is a local word list for `--words`: English terms that the assistant translated word by word in Slack (前景, 背景作業, 背景処理, 背景で〜, 直し, 見張り, 偽), with the word to use. The owner chose to list only literal translations and to leave figures of speech such as 道連れ to the principles. On a sample of 10 such phrases, v3.1.0 alone flagged 3. With the list, it flagged all but one inside 「」, which the lint skips. The list has no hit in upstream's 181-file test corpus.
  - japanese.md sources: 「語を選ぶ」 is the owner's principle (2026-11-17), and it matches upstream `references/calques.md` §2. In 「表記」, the no-space rule is a local choice (upstream keeps the source spacing), and the personification line comes from 本体「比喩動詞の具体化と擬人化の解消」. The lint checks neither, so 「点検」 adds a perl check for the spaces.
  - The core's principle 3 comes from its 意味の保持 and 情報の不増補. Upstream measured its own value as meaning preservation (CHANGELOG 3.0.0「測ったこと」): AI judges preferred no skill for readability 23 to 1, and no skill changed meaning 14 times against 2. The tense clause of principle 3 (a finished change, current behavior, or a plan) comes from v1.0.8 (本体「文の働き」).
  - Dropped by the 2026-11-17 ablation: 文書の立場と文末, the 「も」 rule, the metaphor-verb table, the replaced nouns, the openers and closers, and every rule of 記号と空白 except the spacing. Writers without the skill followed them, or the lint flags them (dashes, trailing colons, non-rendering bold, three equal endings in a row, and part of the slop words). The 体言止め rule became a case of principle 1. The sentence-length and comma rules went too. Writers without the skill broke the 60-character limit 6 times with no lower readability score, and writers under the old comma rules wrote more pause commas than writers without them (11 against 4). The lint's pause_comma covers that case.
  - Earlier drops: the tech/business/essay modes, the 20% kanji-ratio target (unsourced, unmeasured), the list and bold caps, the output format, the academic citations, and the linter's 0–100 score (Goodhart). We do not use its `excess_list` and `excess_bold` findings.
  - Not generalized: the em-dash ban (native English punctuation that the house style uses) and the trailing-colon ban (English lead-in colons are native).
- The 2026-11-17 ablation rebuilt `skills/tech-writing` as three principles with one to three examples each, plus the lint, and cut `skills/cold-read` the same way:
  - Method: six fixed tasks (a Japanese Slack report, a Japanese PR body, a Japanese deslop revision, an English PR body, a cold read of a README, and a Japanese article section). Two `claude -p --safe-mode` writers per task ran for each arm: the old skill, the new skill, and no skill. Safe mode loads no user skill, CLAUDE.md, or rule, and the skill arms got the skill body in the system prompt. A blind `claude -p --safe-mode` judge scored each task, and a blind auditor counted the old rules that each output broke.
  - Result on the five writing tasks: the arms tied on readability (mean rank 3.5 of 6 each). The judges listed 6 meaning slips for the old skill, 0 for the new skill, and 13 for no skill.
  - A tested rule stayed only if writers without it lost readability, meaning, or term choice, or broke an owner choice (the spacing, the English mechanics). Writers without the skill wrote restatements, lists for reasoning, and bold labels, and the judges did not rank them lower for it.
  - The examples steer hard. Two writers copied a hedge example word for word when it came from their task, so the examples come from outside the tasks.
  - cold-read lost its three-reviewer step. Three reviewers doubled the cost of a README cold read ($1.04 against $0.57) and did not change the result in two reruns.
  - cold-read keeps "For style repairs, follow tech-writing". With "For the rewrite" instead, README writers left out facts that only the code shows: 2 of 10 READMEs mentioned the `--dry-run` flag, against 6 of 6 with the old wording.
  - The prose shrank from 24,638 to 9,546 bytes.
- `skills/tech-writing` adapts the method of "explain-before-shorten" by Kazuki Ota (https://zenn.dev/microsoft/articles/explain-before-shorten, published 2026-10-26, read 2026-10-28). The article states no license, so we took its method and wrote our own text and examples:
  - Adopted: its order (choose the points, then explain each as if aloud) and its table of over-shortened phrases, as SKILL.md principle 1. The 2026-11-17 rebuild dropped the shorten-then-restore steps, because principle 1 now says not to force a point short.
  - The label principle replaced these former bullets: write the action as a verb, name the actor, replace a figurative verb, replace "as needed" with a condition, name each cause of an aggregate problem, and give the mechanism of a causal claim.
  - Not adopted: its slide render check (Own the outcome covers it), its output and missing-information tables (the revision rules cover them), and its trigger phrases (the description stays routing-only).
  - A/B test on 2026-10-28: four tasks (a Japanese slide, an English PR body, an English on-call message, a Japanese wiki section), two fresh `claude -p` runs per variant, and a blind judge per task. The first draft lost both Japanese tasks to over-compression, which step 4 now catches. The final version ranked above both old-skill outputs on all four tasks (8 of 8). Codex followed the final version on both Japanese tasks.
