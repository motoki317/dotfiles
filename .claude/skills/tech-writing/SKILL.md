---
name: tech-writing
description: Writing norms for human-facing prose in any language — sentence mechanics, argument rigor, redundancy, filler, plus a check mode. Load before you write or revise reports, PR/issue text, docs, or comments. Also load it for any request to make text read less machine-written. 日本語の技術文書、記事、書籍原稿の執筆と推敲、「AIっぽさを消して」「読みやすくして」という依頼にも使用する。
---

# Technical writing

Every sentence must pass two checks:

1. **Complete**: on one read, a reader can tell who does what and when, under which condition, and how sure the writer is.
2. **Necessary**: deleting the sentence, or any word in it, costs the reader information.

Read the reference for the text's language before you write. Where a reference differs from this file, the reference wins. For another language, apply this file alone.

- English: `$HOME/.claude/skills/tech-writing/references/english.md`
- Japanese: `$HOME/.claude/skills/tech-writing/references/japanese.md`
- Japanese book chapters and published articles (a blog, Zenn): also `$HOME/.claude/skills/tech-writing/references/manuscript.md`. It wins over japanese.md.

The mechanical pass is english.md "Self-check" or japanese.md「点検」. Run it before you deliver a doc, a report file, a code comment, or a PR or issue body. A report sent as a chat message is exempt. Before the pass, write a draft that has no file to a scratch file. Each match is a hit. If the sentence of a hit breaks a rule, fix the sentence. Otherwise, leave it.

## Complete

- Write procedural text (what the reader does) in the imperative, one instruction per sentence. Give the exact command, code, or setting.
- Write descriptive text (what a thing is or does) without the imperative.
- Put the condition before the command. In a warning, put the condition before the risk. Replace "as needed" with the condition and its result: "If the server returns 404, retry up to three times."
- Name the actor: "It was decided" → "The team decided". Write the action as a verb: "perform compression" → "compress".
- Replace a figurative verb with the operation or state change that a reader can observe: "The cache fosters faster page loads" → "The cache makes pages load faster". A verb of intent or feeling given to a thing is figurative: "The function wants a string" → "The function takes a string". "The server rejects the request" and "The results show" are correct.
- Mark a finished change, current behavior, and a plan by tense: "removed", "retries", "will remove". In a change report, give the old and the new state: "Retries stop after three attempts (previously unlimited)", not "Updated the retry policy."
- Make every pronoun, demonstrative, and connective point to something that the reader can identify. "However" must match the actual relation.
- Define a term before its first use. Use one name per concept: do not rotate check/verify/confirm.
- State each claim at the strength of its evidence:
  - Assert a claim that you verified or that evidence in the text grounds.
  - Keep a hedge only on an unverified fact, an inference, a counterfactual, or a personal impression.
  - Write a recommendation as the fact that supports it: "X avoids Y", not "You should use X". When you offer options for the reader to choose from, recommend one with its reason: "I recommend X because Y."
  - State the condition under which a guarantee holds.
- Name each decision, cause, or problem separately: "A combination of problems caused the outage" → name each problem. Map each cause to the part of the event that it explains.
- For a causal claim, give the mechanism in one sentence: "each step shares the hand-off format, so a format change affects every step".

## Necessary

Delete each sentence or word that fails this check. Do not swap it for a synonym.

- To shorten a text, cut whole points. Keep the grammar words of the points that remain: "Check config before deploy" → "Check the config before you deploy."
- Delete a label that announces a point ("Key takeaway:") and an opener about the subject area ("In today's fast-paced world"). Also delete a lesson or call to action that no fact in the text supports.
- If an intensifier stands for a measurable property, write the property: "robust" → "retries three times, then stops". Otherwise, delete the intensifier.
- State one claim once. Do not retell an example or a log after you show it. Do not restate a claim in parentheses: "raw output (the unmodified output)". A parenthesis that identifies a referent stays.
- Do not stage a dialogue with an imagined reader: "You may wonder why. The reason is simple." Do not end on a punchline: "Order is everything."
- If the text before it invites a wrong reading, negate that reading. Then say why the reading is wrong. Otherwise, state the positive claim alone: "This is not just a refactor, it is a bug fix" → "This is a bug fix".

## Length and format

- Split a sentence that exceeds the length limit of its language reference. Split it at a clause boundary where a connective can start the second sentence: reason, means, order, or contrast. Do not split a purpose clause from the conclusion that it governs.
- Write reasoning as paragraphs, one topic each. Use a list only for parallel items: parameters, options, steps, or a mapping.
- Use bold only for a term where you define it and for at most two other spans under each heading. Without headings, the whole text counts as one heading's span.

## Revising another writer's text

This section and the revision rules of a reference win over every other rule. Your own earlier draft is not another writer's text.

Leave exact: code, identifiers, commands, paths, quoted errors and logs, and names. Exception: a quoted message that is itself the text to revise.

Keep four properties of each sentence that you keep:

1. Claim: the point and its logical relations.
2. Weight: what the sentence puts first or plays down. Keep a negation that the text before it invites. Do not turn a negated item into an addition ("A as well as B"). Do not rank items that the source does not rank.
3. Strength: an assertion stays an assertion, and a hedge stays a hedge.
4. Function and time: evaluation, explanation, request, or impression, and a finished change, current behavior, or a plan.

Add no fact that the source does not state, such as an actor, a number, or a cause. Do not narrow a vague word to a specific fact. You can restore an actor or referent that the nearby text names. If a rewrite needs a missing fact, keep the statement as broad as the source. Then ask the writer.

If a label or intensifier carries the sentence's evaluation, move it into the predicate: "What matters most is speed" → "Speed matters most". Otherwise, delete it. When you replace a figurative word, keep its attitude, such as regret or understatement.

## Check and revise

1. Run the mechanical pass. Then test each sentence against the rules of this file and its language reference.
2. To check: report each violation with its rule, span, and rewrite. The span is file:line, or a quote of the text. Name the rule by its section and a few words: "Necessary: restated claim". Write no praise and no summary.
3. To revise: save the original to a scratch file. Fix each violation. Then run the mechanical pass again. Stop after two reruns.
4. If the text is in a file, edit the file. Otherwise, output the revised text. Then output one line per rule that you applied. Add your questions for the writer, including each hit that still breaks a rule.
