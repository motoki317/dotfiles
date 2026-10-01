# Lens catalog

This catalog lists the lenses for `decision-analysis`. `/define` borrows them for requirement questions and NFR coverage. Pick the 3–5 that best expose *this* decision's blind spots — this is a menu, not a checklist. Each lens names what to ask and what a useful answer looks like.

## Core — consider for almost any decision

### Feasibility
- Asks: Can this be built with the available tech, skills, and time? What must be proven before you commit? What outright blocks it?
- Output: go / no-go / needs-spike, with the blocker or the unknown-to-prove named.

### Alternatives
- Asks: What other approaches solve the same need? Why this one over each? If you do not pick the runner-up, what is the opportunity cost?
- Output: the considered set, and the difference that makes this option win.

### Risk
- Asks: What can go wrong, how likely, how bad? Is the action reversible? What is the blast radius if it fails?
- Output: the dominant risks ranked, each with a mitigation or an accept-decision.

### Impact
- Asks: Who and what does this affect — users, operators, downstream systems, future work? Does it change a default, contract, or shared format others depend on?
- Output: the affected parties and the change each one sees.

### Cost / Effort
- Asks: What does it cost to build *and to carry*? Does the payoff justify the build and the ongoing maintenance?
- Output: a rough build-vs-maintenance estimate and whether the return justifies it.

## Conditional — include when the decision touches the area

### Security
- Include when: untrusted input, authn/authz, secrets or crypto, access-control changes, or sensitive data.
- Asks: What new attack surface or trust boundary does this open? What data can leak?
- Output: the exposure introduced and what closes it.

### Performance
- Include when: large/unbounded data, hot or real-time paths, DB queries, or latency/throughput-sensitive work.
- Asks: How does this behave at scale? Where is the bottleneck? What is the cost per request/item?
- Output: the scaling limit and the point where it first causes a problem.

### Maintainability
- Include when: the choice sets a pattern others will follow, or adds long-lived structure.
- Asks: What complexity does this commit the codebase to? How hard is it to change or remove later?
- Output: the carrying cost and any lock-in.

### Dependencies & Constraints
- Include when: the option relies on external systems, libraries, or hard limits (platform, quota, deadline).
- Asks: What must hold for this to work? What breaks if a dependency changes or a limit is hit?
- Output: the load-bearing assumptions and what each one rests on.

### Edge cases
- Include when: the decision hinges on boundary behavior — empty, maximal, concurrent, or malformed inputs.
- Asks: Where does the happy-path reasoning no longer hold? What unusual input or timing breaks it?
- Output: the boundaries that need explicit handling or an explicit decision to ignore.

### Precedent & prior art
- Include when: this problem, or one like it, likely came up before — in this codebase, this org, or a past attempt.
- Asks: Was this tried here before? Why is the current thing built the way it is? What happened last time, and what changed since then?
- Output: the relevant precedent and what it implies for this decision.

### Standard practice
- Include when: an established, idiomatic, or industry-standard way to solve this plausibly exists.
- Asks: What is the conventional approach? Do we deviate from it, and is the deviation justified?
- Output: the standard approach and a reasoned align-or-deviate call.
