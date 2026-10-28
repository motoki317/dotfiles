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
| should (recommendation) | SKILL.md "Write the event, not a label": recommendation |
| may / might / could (present capability, permission) | can |
| could (past capability: "could not connect") | Keep it. |
| may / might / could (epistemic uncertainty) | Keep it. |
| would (conditional) | Restructure: "If X occurs, Y occurs." |
| would (counterfactual) | Keep it. |

## Word caps

- Procedural sentence: 20 words. Descriptive sentence, and a note inside a procedure: 25.
- Split a sentence over its cap at a clause boundary: reason, means, order, or contrast. Keep a purpose clause with the conclusion that it governs.
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

Run this command on the file. It prints each hit outside code as `line: match`. A hit in quoted text is not a violation.

```sh
perl -CSD -Mutf8 -ne 'if(/^\s*```/){$c=!$c;next} next if $c; s/`[^`]*`//g; print "$.: $&\n" while /\b\w+(?:n[\x27’]t|[\x27’](?:ll|re|ve|d))\b|\b(?:it|that|there|let)[\x27’]s\b|\b(?:has|have|had) \w+(?:ed|en)\b|\b(?:should|would|may|might|could)\b|\b(?:is|are) being\b|, (?:making|allowing|enabling|ensuring)\b|;|\b(?:e\.g|i\.e|etc)\.|\bnot (?:just|only)\b|^\s*- \*\*|\*\*/gi' FILE
```

For a `not just` or `not only` hit, see SKILL.md "Delete on sight": negation. For a `**` hit, see SKILL.md "Length and format".

Then check these by hand:

- each word in "Filler table"
- the words that you did not pick in each "One term per set" set
- a trailing ` if ` or ` when ` in procedural text
- the three longest procedural and the three longest descriptive sentences, against "Word caps"
