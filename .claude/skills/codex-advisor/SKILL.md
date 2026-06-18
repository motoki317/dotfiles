---
name: codex-advisor
description: Independent cross-model second opinion from OpenAI Codex — an outside view that catches blind spots a Claude-only review (including the built-in advisor) shares. Invoke on your own judgment when a decision is genuinely uncertain or hard to reverse, when you diverge from the built-in advisor, or before finishing non-trivial work; also on explicit request ("ask codex", "codex review", "get a second opinion"). Read-only by default; full-access only when asked to apply changes.
---

# Codex advisor / reviewer

Consult OpenAI Codex as an **independent, cross-model** reviewer. Codex is a different model family from Claude, so it surfaces issues a Claude-only review — including the built-in `advisor` tool — tends to share a blind spot on.

## When to invoke (proactively)

Reach for Codex on your own judgment — do not wait to be asked:

- **Opinions split** — two or more approaches are genuinely viable, or your analysis diverges from the built-in advisor. Codex breaks the tie from outside the Claude model family.
- **A high-stakes or hard-to-reverse design / architecture / risk decision** is on the table, before committing to it.
- **You are stuck** — an approach is not converging, or results do not fit.
- **A non-trivial task is about to be declared done** — have Codex review the finished diff.

Skip it for mechanical or low-stakes work (renames, one-line fixes, formatting) and for short reactive tasks. Each call spends subscription quota and real time.

## How this differs from the built-in advisor

The built-in `advisor` auto-forwards the whole conversation to a Claude model; Codex starts cold and sees none of it. The helper's baked-in preamble already gives Codex its role and output discipline — independent cross-model critic, separate real defects from speculative risks, end with a clear verdict — so your prompt is the **task only**: context and the question, never instructions on how Codex should behave.

- For a **code review**, point Codex at the changes — pipe `git diff`, or just ask it to review the uncommitted changes; it reads the repo itself.
- For a **design or advice** question, write a self-contained brief: the question, what is already decided or tried, and which files matter.

Keep both reviewers — they are complementary. The Claude advisor is transcript-aware; Codex is repo-aware.

## Using it

`codex-consult` (on `PATH`) is the one interface. It runs **read-only**, anchors Codex to a working directory, and prints **only the final answer** to stdout — the full-transcript path and token usage go to stderr, so you never dig the verdict out of the banner-and-reasoning stream.

It defaults to the current directory, so running it from inside the repo just works; `-C <dir>` is an optional override to target a different directory without `cd`-ing there. Either way, keep the working root the project, not a broad parent like `$HOME`: Codex anchors its exploration there, and in a write sandbox that root bounds what it may change. The tool warns if the root is `$HOME`.

```bash
codex-consult "your question or review request"            # scopes to the current directory
codex-consult -C <dir> "your question or review request"   # or target a different directory
```

`--help` lists the options: `-s workspace-write` lets Codex run the test suite to verify a finding (writes stay inside the repo), `-m <model>` overrides the model, `-v` streams the full run to stderr.

`codex-consult` is a thin wrapper (installed on `PATH` via home-manager) around the self-building Go source at `~/.claude/skills/codex-advisor/scripts/codex-consult.go`. In the rare case you need to debug or change the helper itself, read or edit that file directly.

## Examples

Each prompt below is task-only — the preamble supplies the reviewer framing, so none of them tells Codex how to behave.

**Review a finished change before calling it done** — pipe the diff:

```bash
git diff main | codex-consult -C ~/projects/importer \
  "Review the bounded worker pool I just added. Most worried about goroutine leaks on early return, and whether context cancellation reaches in-flight jobs."
```

**Break a split decision** — a design question with no diff; Codex reads the named package itself:

```bash
codex-consult -C ~/projects/importer \
  "Bounded channel vs. golang.org/x/sync/semaphore for limiting concurrent uploads in internal/upload/? A WaitGroup + manual counter already raced under load, and the limit must be resizable at runtime. Recommend one, and give the main failure mode of the option you reject."
```

**Let Codex pull the diff itself** — run from inside the repo, so no `-C` and no piping:

```bash
codex-consult "Review the uncommitted changes for concurrency bugs and error-handling gaps."
```

## Applying changes (explicit request only)

A consult is read-only by design. Only when the user explicitly asks Codex to *apply* changes, raise the sandbox: `codex-consult -s workspace-write -C <repo> "fix the issues"` keeps Codex's edits inside the repo, while the user's full-access `cdy` alias (`codex --dangerously-bypass-approvals-and-sandbox`: no sandbox, no approvals) grants unrestricted access for an interactive session.

## Handling the result

stdout is already the verdict — relay it; the full-log path on stderr is there when you want Codex's reasoning. Weigh the findings the way you weigh the built-in advisor's — evaluate, do not apply blindly. When Codex and the Claude advisor disagree, surface the conflict to the user rather than silently picking a side.
