---
name: cold-read
description: A context-free reviewer reads durable prose as its real audience. The writer applies the surviving cuts. Use at task completion (Orchestrator Accept) whenever docs, README, comment blocks, or PR/issue bodies were written or revised.
---

# Cold read

The writer cannot measure its own prose, so a fresh reader measures it instead. Its report is evidence for the editor, not an edit script.

## Procedure

1. Collect the task's durable prose. Copy each non-file artifact (a PR body, its diff) verbatim into its own scratch file.
2. Give the reviewer what the artifact's real reader will have. For docs and comments, give the working tree, never the diff or the plan. A comment that is clear only next to the diff is the defect that this test catches. For a PR body, give the body plus its diff. Give nothing from the session.
3. Spawn one fresh `reviewer` agent with the template verbatim. Fill only {paths} and the audience line, which names the artifact's real reader ("a Go developer new to this repo", or "the human reviewing this PR"). Anything more is context that the real reader will not have, and the test cannot detect a gap that this context fills. If your harness has no fresh-context subagents, list the artifacts in your run report instead. The Orchestrator then runs this skill at Accept.
4. Apply each cut that drops no reader decision, precondition, or warning. To keep a span, name the loss that its cut causes: "adds nuance" is the writer's bias, not a loss. Rewrite each confusing span so that it refers to a fact that the reader can see. Add prose only for a failed probe or for a fact that the reader can reach nowhere else. If you cut most of the text, rewrite it from the survivors. For style repairs, follow `$HOME/.claude/skills/tech-writing/SKILL.md`.
5. For load-bearing docs (README, spec, onboarding), spawn three reviewers. A span that all of them cut is dead weight. Rewrite the span of each confusion that any of them raises, as step 4 describes.

## Reviewer template

    This material is new to you. You know nothing about recent changes or
    conversations. Act as this audience: {audience line}.

    Read: {paths}. You can read other repository files to verify claims. Do
    not read git history or any diff that is not listed above.

    Report three lists with file:line spans:
    1. Probe — what is each artifact for? What would you do differently
       because you read it? Cite only what you read.
    2. Confusions — every place you stopped, reread, guessed, or hit a
       referent you cannot resolve ("the above fix", "now").
    3. Cuts — every span whose deletion loses nothing that your audience
       needs. A span belongs in the text exactly when it changes what the
       reader does or believes. Propose no additions. Name gaps under Confusions.

    Report to an editor, not a person: no praise, no hedging, no summary.
