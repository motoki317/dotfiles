---
name: tech-writing
description: Writing norms and a check mode for human-facing prose in any language. Load before you write or revise docs, code comments, PR/issue text, or reports. Also load it for any request to make text read less machine-written. 日本語の技術文書、記事、書籍原稿の執筆と推敲、「AIっぽさを消して」「読みやすくして」という依頼にも使用する。
---

# Technical writing

Read the reference for the text's language. Hold the three principles below while you write, and run the check before you post or save the text.

- Japanese: `$HOME/.claude/skills/tech-writing/references/japanese.md`
- English: `$HOME/.claude/skills/tech-writing/references/english.md`
- Japanese book chapters and published articles (a blog, Zenn): also `$HOME/.claude/skills/tech-writing/references/manuscript.md`

## 1. Choose what to say, then say it in full

List what the reader must do or believe after reading. Cut every other point whole, such as the story of how you found the result. Do not force the kept points short. Write each one as you explain it aloud to a colleague: who or what does what, when, and under which condition.

| Shortened to a label | Written in full |
|---|---|
| Cache invalidation is eventual. | After an update, other servers can return the old value for up to 60 seconds. |
| Restart as needed. | If memory use exceeds 90%, restart the worker. |
| 権限を一元管理する | 権限はIAMグループにだけ付与し、個人には付与しない |

## 2. Use the words that the reader uses

Use the name that the code, the UI, or the reader already uses, and one name for each concept.

## 3. Keep the meaning of the source

The source is the notes, the data, the code, or the text that you revise. Add no cause, number, actor, or guarantee that the source does not state. Keep a hedge as a hedge and an assertion as an assertion. Keep code, identifiers, commands, paths, quoted errors, and names exact.

- Source: 「原因はキャッシュの期限切れだと思われます」 → 「原因はキャッシュの期限切れと考えられます」, not 「原因はキャッシュの期限切れです」.
- Source: "no timeouts in a one-hour load test on staging" → "A one-hour load test on staging had no timeouts", not "This fixes the timeouts."

When you revise another writer's text, also keep what each sentence stresses or plays down, and whether it states a finished change, current behavior, or a plan. If the source does not state the event behind a label, keep the label and ask the writer. These revision rules do not bind your own earlier draft.

## Check

Before you post or save any text, a chat message included, run the check of its language: japanese.md「点検」 or english.md "Self-check". Write a draft that has no file to a scratch file first. Fix a hit only if its sentence breaks a principle or a rule of the reference.

- When you are asked to check a text, report each violation with its rule, its span (file:line or a quote), and a rewrite. Write no praise and no summary.
- To revise a text, save the original to a scratch file first. If the text is in a file, edit the file. Otherwise, output the revised text. Then list each point that you cut and your questions for the writer.
