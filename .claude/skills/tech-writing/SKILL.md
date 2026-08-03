---
name: tech-writing
description: Writing norms for human-facing prose in any language — reports, PR/issue text, docs, comments, articles — sentence mechanics (ASD-STE100-derived), argument rigor, redundancy elimination, filler bans, plus a check mode for proofreading. Use when writing or revising such prose. 日本語の技術文書・書籍原稿・章・草稿・記事・解説文を書くとき、または推敲・リライトするときも使用する。
---

# Technical writing

Before you write, read the surface rules for the text's language:

- English: `$HOME/.claude/skills/tech-writing/references/english.md` — grammar, word caps, filler table, self-check
- Japanese: `$HOME/.claude/skills/tech-writing/references/japanese.md` — 整形、LLM 口調リスト、装置の日本語例、点検
- Book chapters and long articles in Japanese: `$HOME/.claude/skills/tech-writing/references/manuscript.md`. English long-form has no extra layer — the core and english.md apply.

## Classify the text

Classify each passage before you write it. Do not mix the two in one passage:

- **Procedural** — tells the reader what to do. Imperative mood. One instruction per sentence.
- **Descriptive** — explains what a thing is or does. No imperative. One new point per sentence, one topic per paragraph.

## Sentence devices

- If a sentence carries two instructions or two separate points, split it.
- Put the condition before the command: "If the build fails, read the log." Readers and models drop a trailing condition.
- One name per concept through the whole document. Do not rotate synonyms (check/verify/confirm, 「検証」と「確認」): pick one and keep it.
- Describe an action with a verb, not a nominalization: "compress the file", not "perform compression". 「圧縮する」であって「圧縮を実施する」ではない。
- Warnings: command or condition first, then the risk. Never bury the instruction after the explanation.
- Error messages: what happened (past tense), the cause if known, then the fix as an imperative.

## Requirement vs uncertainty

- A requirement is "must". A capability is "can". Do not write "should" in instructions — readers and models treat it as optional. A recommendation keeps its reason: state it as fact ("X avoids Y"), or mark it "recommended" when the choice stays with the reader.
- Never convert genuine uncertainty into assertion. A hedge that carries epistemic content — an unverified fact, an inference from logs, a reader's likely doubt, a counterfactual — keeps its uncertainty. Delete a hedge only when evidence in the text grounds the claim, then state the claim strongly, without timid predicates.
  - Keep: "The cache may still serve stale entries" (unverified inference). Assert only once verified: "The cache serves stale entries until the TTL expires."
  - 悪い例：「提示し続けているかもしれない」を「提示し続けている」に変える。
    良い例：「提示し続けている可能性がある」— 不確実性を残して文を整える。

## Argument rigor

After drafting, check:

- Do not lump distinct things under one label. Keep distinct decisions, distinct causes, and different kinds of problems separate — name each, then state how they relate.
- Before you group several concepts under one umbrella term, state in one sentence why they reduce to the same thing.
- Do not reduce a multi-cause event to one cause. Map each explanation to the part it explains.
- When you claim causality, give the mechanism in one sentence: not "splitting by step makes changes ripple", but "each step shares the hand-off format, so a format change ripples".
- Do not promise detection, guarantee, or resolution absolutely. State the condition under which it holds.
- Narrow every claim to what its example supports.
- Define a term before its first use. Keep one definition and one classification across the document.
- If you defer a point ("covered in the next section"), confirm that the target section resolves it.
- After a concession ("however", 「ただし」), advance the argument. Do not end on the concession.

## Redundancy

- One claim, once. Do not restate it in other words.
- Do not summarize what you just showed (an example, a log, a scene). Add only the one sentence that gives it meaning.
- Merge parallel facts with the same logical role into one sentence.
- Skip the intermediate steps a reader can infer. If a multi-sentence argument compresses to one sentence, keep only that sentence.
- Do not stage a dialogue with an imagined reader, frame ideas meta-textually ("a natural continuation of this is…"), or add author's disclaimers. State the idea directly. A real reader question can stay a question.

## Filler

Delete a sentence or word that adds stance but no content — do not rephrase it. The categories (per-language word tables live in the language references):

- Announcements and wrap-ups: "In this chapter we explore…", "In summary" when it only restates.
- Empty intensifiers: "robust", "comprehensive", "crucial". Give the measurable property or delete.
- Empty verbs: "delve into", "streamline".
- Padding connectives: "furthermore" chains. A single connective that marks a real turn or step stays.

## Cut content, not grammar

Shorten by cutting points and sentences, never grammar or needed context:

- Keep articles and "that" in English. Keep 助詞 in Japanese. No telegraph fragments: "Ensure file exists before running" → "Make sure that the file exists before you run the command."
- Do not cut the context a reader needs to follow: scope, comparison axis, open questions.

## Untouchables

Leave exact: code, identifiers, commands, flags, file paths, quoted errors and logs, product and proper names. Exception: a quoted error or message is editable when that text itself is the artifact you were asked to revise.

## Check mode (proofreading)

When asked to check text instead of writing it:

1. Run the language's mechanical pass — english.md "Self-check", japanese.md 「点検」. Classify each hit before you rewrite it: the patterns also match rule-following text, for example epistemic hedges.
2. Then check the judgment rules — classification, argument rigor, redundancy, filler.
3. Report each violation as: device name — offending span (file:line) — compliant rewrite.
4. Report to an editor: no praise, no hedging, no summary.

The outcome test is separate: a reader with no session context reads the artifact as its real audience (`/cold-read` in Claude Code).
