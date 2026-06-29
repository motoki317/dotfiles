# Interaction Capability — ISO/IEC 25010:2023

> Degree to which a product or system can be interacted with by specified users to exchange information via the user interface to complete specific tasks in a variety of contexts of use.

(Renamed from "Usability" in the 2011 edition; the 2011 "Accessibility" sub-characteristic was split into Inclusivity and User assistance. Applies beyond GUIs — to CLIs, public APIs, error messages, and logs.)

## Official sub-characteristic definitions (ISO/IEC 25010:2023)
- **Appropriateness recognizability** — Degree to which users can recognize whether a product or system is appropriate for their needs.
- **Learnability** — Degree to which the functions of a product or system can be learnt to be used by specified users within a specified amount of time.
- **Operability** — Degree to which a product or system has attributes that make it easy to operate and control.
- **User error protection** — Degree to which a system prevents users against operation errors.
- **User engagement** — Degree to which a user interface presents functions and information in an inviting and motivating manner encouraging continued interaction.
- **Inclusivity** — Degree to which a product or system can be used by people of various backgrounds (such as people of various ages, abilities, cultures, ethnicities, languages, genders, economic situations, etc.).
- **User assistance** — Degree to which a product can be used by people with the widest range of characteristics and capabilities to achieve specified goals in a specified context of use.
- **Self-descriptiveness** — Degree to which a product presents appropriate information, where needed by the user, to make its capabilities and use immediately obvious to the user without excessive interactions with a product or other resources (such as user documentation, help desks or other users).

## What to look for (review guidance — not part of ISO/IEC 25010)
- **Appropriateness recognizability** — command/API names and signatures that obscure what they do.
- **Learnability** — required steps that aren't discoverable, missing `--help`/usage, undocumented preconditions.
- **Operability** — missing sane defaults, flags that are awkward to combine, behaviour that's hard to script.
- **User error protection** — no input validation or confirmation, easy-to-misuse destructive defaults, no dry-run.
- **User engagement** — terse or unhelpful feedback, no progress indication on long operations.
- **Inclusivity** — locale/timezone assumptions, colour-only signals, output that breaks screen readers or non-UTF-8 terminals.
- **User assistance** — error messages with no remediation, missing docs for new behaviour.
- **Self-descriptiveness** — cryptic errors, unlabeled fields, magic numbers/strings surfaced to the user.
