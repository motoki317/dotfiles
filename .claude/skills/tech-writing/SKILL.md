---
name: tech-writing
description: Writing norms for human-facing prose in any language — sentence mechanics, argument rigor, redundancy, filler, plus a check mode. Load before you write or revise reports, PR/issue text, docs, or comments. Also load it for any request to make text read less machine-written. 日本語の技術文書、記事、書籍原稿の執筆と推敲、「AIっぽさを消して」「読みやすくして」という依頼にも使用する。
---

# Technical writing

Before you write, read the surface rules for the text's language:

- English: `$HOME/.claude/skills/tech-writing/references/english.md` — grammar, word caps, filler table, self-check
- Japanese: `$HOME/.claude/skills/tech-writing/references/japanese.md` — 整形、文の長さと読点、LLM口調リスト、比喩動詞の表、装置の日本語例、点検
- Book chapters and long articles in Japanese: `$HOME/.claude/skills/tech-writing/references/manuscript.md`. It sets where such text can use the staging that this file restricts. English long-form has no extra layer — this file and english.md apply.

## Classify the text

Classify each passage before you write it. Do not mix the two in one passage:

- **Procedural** — tells the reader what to do. Imperative mood. One instruction per sentence. Show the exact command, code, or setting that the reader uses.
- **Descriptive** — explains what a thing is or does. No imperative. One new point per sentence, one topic per paragraph.

## Sentence devices

- If a sentence carries two instructions or two separate points, split it.
- Put the condition before the command: "If the build fails, read the log." Readers and models drop a trailing condition.
- Write the condition and its result in place of "as needed" or 「適宜」: "If the server returns 404, retry up to three times."
- One name per concept through the whole document. Do not rotate synonyms (check/verify/confirm, 「検証」と「確認」): pick one and keep it.
- Describe an action with a verb, not a nominalization: "compress the file", not "perform compression". 「圧縮する」であって「圧縮を実施する」ではない。
- Name the actor of each action. For a decision, name who decided: "It was decided" → "The team decided". Replace a vague reference ("the other one", 「片方」「両者」) with its noun.
- Give a verb of intent or feeling only to a person. A tool, concept, or event is the subject only of an action that it performs: "The server rejects the request" stays. Rewrite "The architecture demands discipline" with the people who act and the rule that they follow. "The results show" is a convention and stays.
- Replace a figurative verb with the operation or state change that a reader can observe: "The cache fosters faster page loads" → "The cache makes pages load faster". 「デバッグに丸一日溶かした」→「デバッグに丸一日かかった」. A verb with its literal object is not figurative.
- Warnings: command or condition first, then the risk. Never put the instruction after the explanation.
- Error messages: what happened (past tense), the cause if known, then the fix as an imperative.

## Sentence length

- Keep each sentence within the caps of its language: english.md "Word caps", japanese.md 「文の長さと読点」. To meet a cap, split the sentence at a clause boundary.
- Do not add a short punchline sentence for emphasis: "Order is everything." 「順番が命だ。」
- Exception: in Japanese long-form reading material, a short sentence can be a deliberate beat. `/cognitive-rhythm-writing` sets where. The caps still apply.

## Requirement vs uncertainty

- A requirement is "must". A capability is "can". Do not write "should" in instructions — readers and models treat it as optional. A recommendation keeps its reason. State it as fact: "X avoids Y". If the choice stays with the reader, you can mark it "recommended" instead.
- Never convert genuine uncertainty into assertion. A hedge that carries epistemic content keeps its uncertainty. It marks an unverified fact, an inference from logs, a reader's likely doubt, a counterfactual, or the writer's own impression in a personal account. If evidence in the text grounds the claim, replace the hedge with a strong claim, without timid predicates. Otherwise, keep the hedge.
  - Keep: "The cache may still serve stale entries" (unverified inference). Once evidence in the text verifies it, assert it: "The cache serves stale entries until the TTL expires."

## Argument rigor

After you write a draft, check:

- Do not lump distinct things under one label. Keep distinct decisions, distinct causes, and different kinds of problems separate. Name each one. Then state how they relate.
- Before you group several concepts under one umbrella term, state in one sentence why they reduce to the same thing.
- Do not reduce a multi-cause event to one cause. Map each explanation to the part that it explains.
- When you claim causality, give the mechanism in one sentence: not "splitting by step makes changes ripple", but "each step shares the hand-off format, so a format change affects every step".
- Do not promise detection, guarantee, or resolution absolutely. State the condition under which it holds.
- Narrow every claim to what its example supports.
- Define a term before its first use. Keep one definition and one classification across the document.
- If you defer a point ("covered in the next section"), check that the target section resolves it.
- After a concession ("however", 「ただし」), advance the argument. Do not end on the concession.

## Redundancy

- One claim, once. Do not restate it in other words. Do not restate a word in parentheses: "raw output (the unmodified output)", 「不自然な日本語（いわゆるAI臭さ）」. A parenthesis that identifies a referent stays.
- Do not summarize what you just showed (an example, a log, a scene). Add only the one sentence that gives it meaning.
- Merge parallel facts with the same logical role into one sentence.
- Skip the intermediate steps that a reader can infer. If a multi-sentence argument compresses to one sentence, keep only that sentence.
- Do not stage a dialogue with an imagined reader, frame ideas meta-textually ("a natural continuation of this is…"), or add author's disclaimers. State the idea directly. A real reader question can stay a question.
- Do not negate a claim that no reader holds: "It is not just a cache, it is a contract" → "It is a contract". 「単なるXではなく、Yだ」→「Yだ」. Negate only a misreading that the reader is likely to make. Give the reason in one sentence.

## Filler

Delete a sentence or word that adds stance but no content — do not rephrase it. The categories (per-language word tables live in the language references):

- Announcements and self-labels: "In this chapter we explore…", "The key takeaway is…", 「重要なのは」.
- Wrap-ups: "In summary" when it only restates, a summary section in a short text, closing formulas ("I hope this helps").
- Unearned endings: a promise of payoff ("Once you master this, …"), or a lesson or call to action that the text does not support.
- Generic openers: "In today's fast-paced world", "Of course, this is not always the case".
- Empty intensifiers: "robust", "comprehensive", "crucial". Give the measurable property or delete.
- Grandiose abstractions: "a testament to", "a rich tapestry", 「真理」「境地」. State the plain fact.
- Words that sound like texture or insight: "deep dive", "valuable insights", 「解像度を上げる」「肌感」, and coined names without a definition. Write what was done, seen, or measured.
- Empty verbs: "delve into", "streamline".
- Padding connectives: "furthermore" chains. A single connective that marks a real turn or step stays.

## Format

- Write reasoning as paragraphs. Use a list only for parallel items: parameters, options, steps, a mapping. If the items of a "**Label**: text" list carry an argument, write a paragraph.
- Bold a term where you define it. This includes the term label of a glossary list (each item says what the term means). Bold at most one or two other spans per section. Put each at a logical pivot: a negation that prevents a misreading, or the section's conclusion. Do not bold a list label that is not a defined term.
- Do not use emoji.
- Use the punctuation and spacing of the text's language. japanese.md 「整形」 lists the habits from other languages to remove.

## Cut content, not grammar

To shorten a text, cut points and sentences, never grammar or needed context:

- Keep articles and "that" in English. Keep 助詞 in Japanese. No telegraph fragments: "Ensure file exists before running" → "Make sure that the file exists before you run the command."
- Do not cut the context that a reader needs to follow: scope, comparison axis, open questions.

## Untouchables

Leave exact: code, identifiers, commands, flags, file paths, quoted errors and logs, and product and proper names. Boilerplate that you must reproduce (for example, an attribution line) also stays exact. Exception: a quoted error or message is editable when that text itself is the artifact that you were asked to revise.

When you revise another writer's text, keep every constraint, number, and condition. Add no actor, number, setting, cause, or effect that the source does not state. If a rewrite needs a missing fact, ask the writer, or keep the statement neutral. The rewrites in this skill show the shape of a fix. In a revision, each fact comes from the source.

## Check mode (proofreading)

If the task is to check text rather than write it:

1. Run the language's mechanical pass — english.md "Self-check", japanese.md 「点検」. Classify each hit before you rewrite it: the patterns also match rule-following text, for example epistemic hedges.
2. Then check the judgment rules — classification, sentence devices, argument rigor, redundancy, filler, format.
3. Write each rewrite from the sentence's actor, action, and object. Do not swap the flagged word for a synonym: "robust" → "solid", or 「手触り」→「実態」, keeps the defect.
4. Name the device that each violation breaks. Use the heading of the section that states the rule, plus the rule's label or a few words of it: "Filler: Empty intensifiers", "Sentence devices: condition first". Report each violation as: device name — offending span (file:line) — compliant rewrite.
5. Report to an editor: no praise, no hedge without epistemic content, no summary.

If the task is to revise rather than report, apply the rewrites. Then run the mechanical pass again on the result. Stop after two reruns. Do not change a sentence only to clear a hit. Output the revised text. Then list each device that you applied, one line each.

The outcome test is separate: a reader with no session context reads the artifact as its real audience (`/cold-read` in Claude Code).
