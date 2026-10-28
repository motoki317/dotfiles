---
name: tech-writing
description: Writing norms and a check mode for human-facing prose in any language. Load before you write or revise docs, code comments, PR/issue text, or reports. Also load it for any request to make text read less machine-written. 日本語の技術文書、記事、書籍原稿の執筆と推敲、「AIっぽさを消して」「読みやすくして」という依頼にも使用する。
---

# Technical writing

The finished text must pass two checks:

1. **Complete**: a reader who has only this text can tell who does what, when, under which condition, and how sure the writer is.
2. **Necessary**: deleting any sentence or word costs that reader information.

Read the reference for the text's language before you write. Where a reference differs from this file, the reference wins. For another language, apply this file alone.

- English: `$HOME/.claude/skills/tech-writing/references/english.md`
- Japanese: `$HOME/.claude/skills/tech-writing/references/japanese.md`
- Japanese book chapters and published articles (a blog, Zenn): also `$HOME/.claude/skills/tech-writing/references/manuscript.md`. It wins over japanese.md.

## Write

1. List what the reader must do or believe after reading. Cut every other point whole. Cut a point that the reader can derive from another point or from the code. Cut the story of how you reached the result, such as a dead end. Keep a rejected alternative only if the reader can propose it. If the reader owns the fix, give the problem and the expected outcome, not a fix plan.
2. Write each point the way you explain it aloud to a colleague who has not seen your work: who or what does what, when, under which condition, and why. "Aloud" chooses the content, not the tone.
3. Delete each sentence and word whose loss the reader does not notice. Keep the grammar words: "Check config before deploy" → "Check the config before you deploy." If the text still exceeds a length that the task sets, cut a point, not a condition. Do not compress a point into a label or an abbreviation.
4. Read only the final text, as a reader who did not hear your explanation. Compare it with your step-2 draft, and restore each actor, time, condition, or reason that the cut lost. Then fix each term or referent ("this", "the fix") that the text does not identify. Fix each connective ("however", "so") that does not match the relation.
5. Before you deliver any text except a chat message, run the mechanical pass on it: english.md "Self-check" or japanese.md「点検」. Write a draft that has no file to a scratch file first. Fix a hit only if its sentence breaks a rule.

## Write the event, not a label

A label is a short phrase that stands for an event that only the writer can see. Write the event: the actor, the action or state change, and its condition.

| Label | Event |
|---|---|
| Cache invalidation is eventual. | After an update, other servers can return the old value for up to 60 seconds. |
| 権限を一元管理する | 権限は IAM グループにだけ付与し、個人には付与しない |
| a robust parser | The parser skips a malformed line and logs its line number. |
| Restart as needed. | If memory use exceeds 90%, restart the worker. |

- In a change report, such as a PR body or a status update, give the old and the new state: "The timeout is 10 seconds (previously 30)", not "Updated the timeout."
- Use the name that the code, the UI, or the reader already uses. Coin a term only if no such name exists, and define it before its first use. Use one name per concept.
- Assert a claim that you verified or that evidence in the text supports. Hedge only an unverified fact, an inference, or an impression. State the condition under which a guarantee holds.
- Write a recommendation as the fact that supports it: "X avoids Y", not "You should use X". If the reader must choose among options, recommend one with its reason.
- Write a procedure in the imperative, one instruction per sentence, with the exact command, code, or setting.

## Delete on sight

Do not swap these for a synonym.

- An announcement of a point ("Key takeaway:"), an opener about the subject area, and a lesson that no fact in the text supports.
- A punchline ("Order is everything.") and a staged dialogue ("You may wonder why.").
- An intensifier that stands for no measurable property.
- A second statement of one claim: a retold example or log, or a parenthesis that restates ("raw output (the unmodified output)").
- A negation that nothing before it invites: "This is not just a refactor, it is a bug fix" → "This is a bug fix". If the text before it invites the wrong reading, negate that reading and say why.

## Length and format

- Write reasoning as paragraphs, one topic each. Use a list only for parallel items: parameters, options, steps, or a mapping.
- Use bold only for a term where you define it, and for at most two other spans under each heading.

## Revising another writer's text

This section and the revision rules of a reference win over every other rule. Your own earlier draft is not another writer's text.

Leave code, identifiers, commands, paths, quoted errors and logs, and names exact. In each sentence that you keep, keep four properties:

1. Claim: the point and its logical relations.
2. Weight: what the sentence puts first or plays down. Keep a negation that the text before it invites. Rank no items that the source does not rank.
3. Strength: an assertion stays an assertion, and a hedge stays a hedge.
4. Function and time: evaluation, explanation, request, or impression, and a finished change, current behavior, or a plan.

Add no fact that the source does not state, such as an actor, a number, or a cause. If the source does not state the event behind a label, keep the label's breadth and add a question for the writer. If an announcement or intensifier carries the sentence's evaluation, move it into the predicate: "What matters most is speed" → "Speed matters most". When you replace a figurative word, keep its attitude, such as regret or understatement.

## Check or revise

1. Run the mechanical pass. Then test each sentence against this file and its language reference.
2. To check: report each violation with its rule, its span (file:line or a quote), and a rewrite. Name the rule by its section and a few words: "Delete on sight: restated claim". Write no praise and no summary.
3. To revise: save the original to a scratch file. Fix each violation, then run the mechanical pass again. Stop after two reruns. If the text is in a file, edit the file. Otherwise, output the revised text first. Then list each point that you cut and your questions for the writer, including each hit that still breaks a rule.
