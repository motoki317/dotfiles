---
name: cold-read
description: Acceptance test for durable prose — docs, README, comment blocks, PR/issue bodies. A reviewer with no session context reads the artifact as its real reader would; the writer applies the surviving cuts. Use after a task that wrote or materially revised such prose; not for chat replies or single-line comments.
---

# Cold read

The writer cannot measure its own prose, so a fresh reader measures it instead. Its report is evidence for the editor, not an edit script.

## Procedure

1. Collect the task's durable prose. A non-file artifact (a PR body, its diff) goes verbatim into a scratch file, one per artifact.
2. Give the reviewer what the artifact's real reader will have — for docs and comments the working tree, never the diff or plan (a comment clear only next to the diff is the defect this test catches); for a PR body, the body plus its diff. Nothing from the session.
3. Spawn one fresh `reviewer` agent with the template verbatim, filling only {paths} and the audience line — the artifact's real reader ("a Go developer new to this repo"; "the human reviewing this PR"). Anything more is context that reader will not have, and it blinds the test. A harness without fresh-context subagents instead lists the artifacts in its run report, and the Orchestrator runs this skill at Accept.
4. Apply: cuts land unless one drops a reader decision, precondition, or warning — keeping a span requires naming that loss, and "adds nuance" is the writer's bias, not a loss. Re-anchor each confusion to a fact the reader can see; add prose only for a failed probe or a fact the reader can reach nowhere else. If most is cut, rewrite from the survivors. Japanese style repairs: `$HOME/.claude/skills/japanese-tech-writing/SKILL.md`.
5. Load-bearing docs (README, spec, onboarding) get three reviewers: a span all of them cut is dead weight; a confusion any of them raises needs an anchor.

## Reviewer template

    You are reading this material for the first time. You know nothing about
    recent changes or conversations. Audience to embody: {audience line}.

    Read: {paths}. You may read other repository files to verify claims, but
    no git history and, unless a diff is listed above, no diffs.

    Report three lists with file:line spans:
    1. Probe — what is each artifact for, and what would you do differently
       because you read it? Cite only what you read.
    2. Confusions — every place you stopped, reread, guessed, or hit a
       referent you cannot resolve ("the above fix", "now").
    3. Cuts — every span whose deletion loses nothing your audience would
       miss; a span earns its place by changing what the reader does or
       believes. Propose no additions; name gaps under Confusions.

    Report to an editor, not a person: no praise, no hedging, no summary.
