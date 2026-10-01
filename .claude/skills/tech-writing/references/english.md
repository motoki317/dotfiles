# English surface rules

Adapted from ASD-STE100 Simplified Technical English via AminBlg/SimpleEnglish (MIT, see `../LICENSE-SimpleEnglish`). Unofficial. Structural rules only — this file does not reproduce the official dictionary.

## Grammar

- Simple tenses only: simple present, simple past, simple future. No perfect ("has been installed" → "was installed"), no progressive.
- Active voice. Use the passive only in descriptive text where the agent is unknown.
- Use a past participle only as an adjective ("the cached response") or in a passive that the rule above allows.
- Use an "-ing" form only inside a technical noun ("logging"). Never as a verb: ", making restarts unnecessary" → a new sentence with a real subject.
- No contractions.
- No semicolons — write two sentences. No "e.g." / "i.e." / "etc." — write "for example", "that is", or name the items.
- American spelling. An established repo convention wins.

## Modals

| You wrote | Write |
|---|---|
| should (requirement) | must |
| should (recommendation) | State it as fact with the reason, or write "recommended: X, because Y", or delete. |
| may / might / could (present capability, permission) | can |
| could (past capability: "could not connect") | Keep it. It is the simple past of "can". |
| may / might / could (epistemic uncertainty) | Keep the uncertainty — SKILL.md "Requirement vs uncertainty". |
| would (conditional) | Restructure: "If X occurs, Y occurs." |
| would (counterfactual — contrary to fact) | Keep it. It carries epistemic content. |

## Word caps

- Procedural sentence: 20 words. Descriptive sentence, and a note inside a procedure: 25. Paragraph: 6 sentences.
- Count as one word each: `code spans`, numbers with units, identifiers, quoted text, proper names.
- Break noun chains over three words with prepositions: "the connection pool timeout configuration value" → "the timeout value for the connection pool".

## Filler table

A row's replacement applies only when the word carries a fact. When it does not, delete the word.

| You wrote | Write |
|---|---|
| leverage, utilize | use |
| in order to / prior to / due to the fact that / in the event that / when it comes to | to / before / because / if / for |
| ensure | make sure that |
| it is worth noting that / it is important to / crucially | (delete — state the fact) |
| simply, just, easily, seamlessly, effortlessly | (delete) |
| robust, powerful, comprehensive, performant, blazingly fast, crucial, pivotal, vital, key (adjective), vibrant | (delete, or give the measurable property) |
| serves as, stands as, functions as, represents (meaning "is") | is |
| boasts, features, offers (meaning "has") | has |
| underscore, highlight, showcase (meaning "show") | show, or state the fact |
| foster, garner, bolster, enhance, align with | (name the effect: "retries failed uploads three times") |
| a testament to, tapestry, landscape (abstract), pivotal moment, indelible mark, deeply rooted | (delete — state the plain fact) |
| deep dive, valuable insights, interplay, intricate, meticulous | (say what was examined or found) |
| in today's fast-paced world, in an ever-evolving landscape | (delete) |
| I hope this helps, let me know if, happy coding | (delete) |
| Additionally, Furthermore, Moreover (sentence-initial) | (delete, unless it marks a real turn) |
| enables you to, allows you to | you can |
| is designed to, aims to | (say what it does) |
| facilitate / streamline | help / make simpler |
| delve into, dive into | read, examine |
| as needed, as necessary | (state the condition) |
| and/or | "X, or Y, or both" |
| gracefully handles | (say what it does: "retries three times, then stops") |
| out of the box / under the hood | by default / internally |
| plethora, myriad | many |
| functionality | function, feature |

## Rotations to collapse

One term per concept (pick one):

- check / verify / confirm / validate / make sure that — not "ensure" (filler table)
- config / configuration / settings / options
- delete / remove / drop / destroy — one per meaning
- error / issue / problem / failure — one per meaning (a message reports an error, an operation fails)
- run / execute / invoke / launch
- show / display / render / present

## Self-check (searchable)

Search the draft. A hit outside code blocks and quoted text is a candidate. Classify it against the rules above before you rewrite it.

| Search | Violation → fix |
|---|---|
| `'ll` `'re` `'ve` `n't` `it's` | contraction → expand |
| `has been` `have been` `had been`, has/have + participle | perfect tense → simple past or present |
| `should` `would` `may` `might` `could` | modal → the Modals table above |
| `is being` `are being` | progressive passive → active, simple tense |
| `, making` `, allowing` `, enabling` `, ensuring` | "-ing" clause → new sentence |
| `;` | semicolon → two sentences |
| ` if ` / ` when ` mid-sentence, in procedural text | trailing condition → move to the front, add a comma |
| `not just` `not only` `isn't just` | negative parallelism → SKILL.md "Redundancy" |
| `serves as` `stands as` `boasts` | copula avoidance → "is", "has" |
| `- **` at the start of a line | inline-header list whose items carry an argument → a paragraph |
| `**` | bold other than a defined term, more than two spans in a section → SKILL.md "Format" |
| emoji | delete |

Then count words in the three longest sentences against the caps. Search for the rotation synonyms that you did not pick.

`$HOME/.claude/skills/tech-writing/scripts/ste_lint.py` automates these counts for before/after measurement only. It reports aggregate totals without spans and exits 0 regardless — it is not the check-mode pass.
