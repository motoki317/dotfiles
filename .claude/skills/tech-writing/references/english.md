# English surface rules

## Grammar

- Use simple tenses only: simple present, simple past, simple future. "has failed" → "failed". Exception: a counterfactual keeps "would have": "Without the cap, the job would have run forever."
- Use the active voice. Use the passive only in descriptive text where the agent is unknown.
- Do not use a participle clause: ", making restarts unnecessary" → a new sentence with a real subject. A participle or gerund that acts as an adjective or a noun is correct: "cached response", "existing users", "logging".
- No contractions. No semicolons: write two sentences.
- Write "for example", "that is", or the items in place of "e.g.", "i.e.", and "etc."

## Modals

| You wrote | Write |
|---|---|
| should (requirement) | must |
| should (recommendation) | SKILL.md "Complete": recommendation |
| may / might / could (present capability, permission) | can |
| could (past capability: "could not connect") | Keep it. |
| may / might / could (epistemic uncertainty) | Keep it. |
| would (conditional) | Restructure: "If X occurs, Y occurs." |
| would (counterfactual) | Keep it. |

## Word caps

- Procedural sentence: 20 words. Descriptive sentence, and a note inside a procedure: 25.
- Count as one word each: `code spans`, numbers with units, identifiers, quoted text, proper names.
- Break noun chains over three words with prepositions: "the connection pool timeout configuration value" → "the timeout value for the connection pool".

## Filler table

Keep a word in a sense that its row does not cover: "the job just finished".

| You wrote | Write |
|---|---|
| leverage, utilize | use |
| in order to / prior to / due to the fact that / in the event that / when it comes to | to / before / because / if / for |
| ensure | make sure that (a person), or say what the thing does |
| it is worth noting that / it is important to / crucially | (delete) |
| simply, just, easily, seamlessly, effortlessly | (delete) |
| robust, powerful, comprehensive, performant, crucial, pivotal, vital, key (adjective) | (delete, or give the measurable property) |
| serves as, stands as, functions as, represents (meaning "is") | is |
| features, offers (meaning "has") | has |
| underscore, highlight, showcase (meaning "show") | show, or state the fact |
| foster, garner, bolster, enhance, align with | (name the effect) |
| deep dive, dive into, delve into, valuable insights, interplay, intricate, meticulous | (say what you examined or found) |
| I hope this helps, let me know if, happy coding | (delete) |
| Additionally, Furthermore, Moreover (sentence-initial) | (delete) |
| enables you to, allows you to | you can |
| is designed to, aims to, gracefully handles | (say what it does) |
| facilitate / streamline | help / make simpler |
| and/or | "X, or Y, or both" |

## One term per set

- check / verify / confirm / validate
- config / configuration / settings / options
- delete / remove / drop / destroy, one per meaning
- error / issue / problem / failure, one per meaning: a message reports an error, and an operation fails
- run / execute / invoke / launch
- show / display / render / present

## Self-check

Search the draft outside code blocks and quoted text.

| Search | Violation → fix |
|---|---|
| `'ll` `'re` `'ve` `'d` `n't` `it's` `that's` `there's` `let's` | contraction → expand |
| `has been` `have been` `had been`, has/have + participle | perfect tense → simple past or present |
| `should` `would` `may` `might` `could` | modal → "Modals" |
| `is being` `are being` | progressive passive → active, simple tense |
| `, making` `, allowing` `, enabling` `, ensuring` | participle clause → new sentence |
| `;` | semicolon → two sentences |
| `e.g.` `i.e.` `etc.` | abbreviation → "Grammar" |
| ` if ` / ` when ` mid-sentence, in procedural text | trailing condition → move it to the front |
| `not just` `not only` | uninvited negation → SKILL.md "Necessary" |
| each word in "Filler table" | filler → "Filler table" |
| `- **` at the start of a line | bold-label list of reasoning → paragraphs (SKILL.md "Length and format") |
| `**` | bold beyond SKILL.md "Length and format" |

Then count the words of the three longest procedural sentences and the three longest descriptive sentences against "Word caps". Search for the words of each "One term per set" set that you did not pick.
