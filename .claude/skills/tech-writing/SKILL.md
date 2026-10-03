---
name: tech-writing
description: Writing norms for human-facing prose in any language — sentence mechanics, argument rigor, redundancy, filler, plus a check mode. Load before you write or revise reports, PR/issue text, docs, or comments. Also load it for any request to make text read less machine-written. 日本語の技術文書、記事、書籍原稿の執筆と推敲、「AIっぽさを消して」「読みやすくして」という依頼にも使用する。
---

# Technical writing

Your reader has none of your context and cannot ask you a question. Every sentence must pass two checks:

1. **Complete**: on one read, a reader can tell who does what, under which condition, and how sure the writer is.
2. **Necessary**: the sentence adds information that the reader does not have yet.

When the checks conflict, cut whole points, and keep the grammar, actor, and condition of the points that remain. If no rule below covers a case, apply the two checks directly.

Before you write, read the reference for the text's language. For another language, apply this file alone.

- English: `$HOME/.claude/skills/tech-writing/references/english.md`
- Japanese: `$HOME/.claude/skills/tech-writing/references/japanese.md`
- Japanese book chapters and long articles: also `$HOME/.claude/skills/tech-writing/references/manuscript.md`, which relaxes some bans of this file for such text. `$HOME/.claude/skills/cognitive-rhythm-writing/SKILL.md` sets where a short sentence can be a deliberate beat.

## Complete

- Classify each passage (a section or a paragraph), and do not mix the two kinds in one passage. Procedural text tells the reader what to do: imperative mood, one instruction per sentence, and the exact command, code, or setting. Descriptive text explains what a thing is or does, with no imperative. In Japanese, the document's stance sets the endings of both kinds: japanese.md「文書の立場と文末」.
- If a sentence carries two instructions or two separate points, split it.
- Put the condition before the command: "If the build fails, read the log." Write the condition and its result in place of "as needed": "If the server returns 404, retry up to three times."
- Name the actor, and describe the action with a literal verb.
  - For a decision, name who decided: "It was decided" → "The team decided".
  - Use a verb, not a nominalization: "compress", not "perform compression".
  - Give a verb of intent or feeling only to a person. "The server rejects the request" and "The results show" stay.
  - Replace a figurative verb with the operation or state change that a reader can observe: "The cache fosters faster page loads" → "The cache makes pages load faster".
- Make every reference resolvable. A pronoun, a demonstrative, or a connective must point to something that the reader can identify. Replace "the other one" with its noun, and make "however" match the actual relation. Use one name per concept, and do not rotate synonyms (check/verify/confirm). Define a term before its first use. Keep its definition, and any categories that you sort things into, fixed across the document.
- State each claim at the strength of its evidence.
  - A requirement is "must". A capability is "can". Do not write "should" in instructions. State a recommendation as fact with its reason ("X avoids Y"), or mark it "recommended" if the choice stays with the reader.
  - Keep a hedge that carries epistemic content: an unverified fact, an inference from logs, a reader's likely doubt, a counterfactual, or the writer's impression in a personal account. If evidence in the text grounds a claim, assert it without a timid predicate such as "one option is".
  - Do not promise detection, guarantee, or resolution absolutely. State the condition under which it holds. Narrow every claim to what its example supports.
- Keep the parts of an argument distinct.
  - Do not lump distinct decisions, causes, or problems under one label. Name each one, then state how they relate. Before you group concepts under one umbrella term, state in one sentence why they reduce to the same thing.
  - Map each cause of a multi-cause event to the part that it explains.
  - When you claim causality, give the mechanism in one sentence: not "splitting by step makes changes ripple", but "each step shares the hand-off format, so a format change affects every step".
  - Resolve each deferred point ("covered in the next section") in its target section. After a concession ("however"), return to the main claim. Do not end on the concession.
- Warnings: the command or condition first, then the risk. Error messages: what happened (past tense), the cause if known, then the fix as an imperative.

## Necessary

A sentence fails this check when its deletion costs the reader no information. Delete such a sentence or word. Do not swap it for a synonym.

- Delete words that add stance but no information: announcements and self-labels, wrap-ups that only restate, unearned endings (a payoff, lesson, or call to action that the text does not support), generic openers, empty intensifiers, grandiose abstractions, words that only sound concrete or insightful, undefined coinages, empty verbs, and connectives that mark no real turn. The language references list examples. If an intensifier stands for a measurable property, write the property: "robust" → "retries three times, then stops".
- State one claim once. Do not restate it in other words or in parentheses: "raw output (the unmodified output)". A parenthesis that identifies a referent stays.
- Do not summarize what you just showed (an example, a log, a scene). Add only the one sentence that gives it meaning.
- Merge parallel facts that share one logical role into one sentence: together they make one point. Skip the steps that a reader can infer from the text.
- Do not stage a dialogue with an imagined reader, frame the text meta-textually ("a natural continuation of this is…"), add author's disclaimers, or add a punchline sentence for emphasis ("Order is everything."). A question that readers actually ask can stay a question.
- Negate only a misreading that the reader is likely to make, and give the reason in one sentence. "It is not just a cache, it is a contract" → "It is a contract".
- When you shorten a text, keep the grammar and the context that a reader needs. Keep articles and "that" in English, and particles in Japanese: "Ensure file exists before running" → "Make sure that the file exists before you run the command." Keep scope, comparison axis, and open questions.

## Length and format

- Keep each sentence within the length rule of its language: english.md "Word caps", japanese.md「文の長さと読点」. To meet it, split at a clause boundary where a connective at the start of the second sentence keeps the relation (reason, means, order, contrast). Do not split a purpose clause from the conclusion that it governs.
- Write reasoning as paragraphs, one topic per paragraph. Use a list only for parallel items: parameters, options, steps, a mapping. If the items of a "**Label**: text" list carry an argument, write a paragraph.
- Bold a term where you define it. Bold at most one or two other spans per section, each at a logical pivot: a negation that prevents a misreading, or the section's conclusion. Do not bold a list label that is not a defined term.
- Do not use emoji. Use the punctuation and spacing of the text's language: japanese.md「記号と空白」.

## Revising another writer's text

Leave exact: code, identifiers, commands, flags, file paths, quoted errors and logs, product and proper names, and boilerplate that you must reproduce. Exception: a quoted message is editable when that text itself is the artifact that you were asked to revise.

Keep four properties of each sentence. The rules in this section win over "Complete" and "Necessary".

1. Claim: the point and its logical relations.
2. Weight: what the sentence puts first, and what it negates or plays down. Do not turn a negated item into an addition ("A as well as B"), and do not rank items that the source does not rank.
3. Strength: an assertion stays an assertion, and a hedge stays a hedge, even where evidence grounds the claim.
4. Function: evaluation, explanation, request, plan, or impression.

Add no actor, object, referent, number, setting, cause, effect, example, or term that the source does not state, and do not narrow a vague word to a specific fact. If the surrounding text names an omitted actor or referent, you can restore it. If a rewrite needs a missing fact, keep the statement as broad as the source and ask the writer.

Delete a label or intensifier that only decorates. If it carries the sentence's evaluation, move the evaluation into the predicate: "The key point is speed" → "Speed is the key point". Keep a negation that corrects a likely misreading, and add no reason that the source lacks. When you replace a figurative word, keep the attitude that it carried, such as regret or understatement.

## Check mode

If the task is to check text rather than write it:

1. Run the language's mechanical pass: english.md "Self-check", japanese.md「点検」. Classify each hit before you rewrite it: the patterns also match rule-following text, for example epistemic hedges.
2. Check "Complete", "Necessary", and "Length and format".
3. Write each rewrite from the sentence's actor, action, and object. Do not swap the flagged word for a synonym: "robust" → "solid" keeps the defect. For another writer's text, follow "Revising another writer's text".
4. Report each violation as: device name — offending span (file:line) — compliant rewrite. The device name is the heading of the section that states the rule, plus a few words of the rule: "Necessary: empty intensifiers".
5. Report to an editor: no praise, no hedge without epistemic content, no summary.

If the task is to revise rather than report, apply the rewrites. Then run the mechanical pass again on the result. Stop after two reruns. Do not change a sentence only to clear a hit. Output the revised text, then one line for each device that you applied, then your questions for the writer.

After you finish a durable text, have a reader with no session context read it as its real audience: `$HOME/.claude/skills/cold-read/SKILL.md`.

## Sources

- AminBlg/SimpleEnglish (MIT): the sentence mechanics and `references/english.md`. <https://github.com/AminBlg/SimpleEnglish>
- k16shikano's Japanese writing norms: the argument and redundancy devices, and `references/manuscript.md`. <https://gist.github.com/k16shikano/fd287c3133457c4fd8f5601d34aa817d>
- nanaism/yomiyasu (MIT, commit 8d5abeeb): the Complete check (`gemini-syntax.md` 原則1), the revision rules (`SKILL.md`「意味の保持」「情報の不増補」「セルフラベリングと否定対比の整理」), and `references/japanese.md`. <https://github.com/nanaism/yomiyasu>
