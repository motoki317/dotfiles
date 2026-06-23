---
name: codex-advisor
description: Independent cross-model second opinion from OpenAI Codex — an outside view that catches blind spots a Claude-only review (including the built-in advisor) shares. Invoke on your own judgment when a decision is genuinely uncertain or hard to reverse, when you diverge from the built-in advisor, or before finishing non-trivial work; also on explicit request ("ask codex", "codex review", "get a second opinion"). Read-only by default; full-access only when asked to apply changes.
---

# Codex advisor / reviewer

`codex-consult` (on `PATH`) consults OpenAI Codex as an **independent, cross-model** reviewer. Codex is a different model family from Claude, so it catches blind spots a Claude-only review — including the built-in `advisor` tool — shares. The two are complementary: `advisor` is transcript-aware, Codex starts cold and is repo-aware.

## When to invoke

On your own judgment, not only on request:

- Approaches split, or your analysis diverges from the built-in advisor.
- A high-stakes or hard-to-reverse design/architecture/risk decision, before committing.
- You are stuck — an approach is not converging.
- A non-trivial task is about to be declared done — review the finished diff.

Skip mechanical or low-stakes work (renames, one-line fixes, formatting); each call spends quota and time.

## Running it

`codex-consult` runs **read-only**, anchored to a working directory, and prints **only the final verdict** to stdout (log path and token usage go to stderr).

```bash
codex-consult "<short prompt>"            # anchors to the current directory
codex-consult -C <dir> "<short prompt>"   # or target another directory
codex-consult -C <dir> < brief.md         # non-trivial prompt: pipe a file (see "Writing the prompt")
```

- Keep the anchored root the **project**, not a broad parent like `$HOME`: it bounds Codex's exploration and, in a write sandbox, what it may change. The tool warns if the root is `$HOME`.
- Options (`--help`): `-s workspace-write` lets Codex run tests or apply edits inside the repo, `-m <model>` overrides the model, `-v` streams the full run to stderr.

## Writing the prompt

**Task only.** The wrapper's preamble already gives Codex its role and tells it to end with a verdict — supply context and the question, never instructions on how to behave. For a code review, pipe `git diff` or just ask it to review the uncommitted changes (it reads the repo). For a design question, write a self-contained brief: the question, what is decided or tried, and which files matter.

**Avoid shell corruption.** A double-quoted argument is shell-parsed first, so backticks, `$(…)`, `$VAR`, and `!` expand and silently corrupt code-heavy prompts. For anything non-trivial, write the brief with the Write tool and pipe it via stdin — stdin is never shell-interpreted. Pass inline `"…"` only for a short one-liner with no backtick, `$`, or `!`.

```bash
codex-consult -C <repo> < /tmp/brief.md                          # brief from a file
cat /tmp/brief.md <(git diff main) | codex-consult -C <repo>     # metachar-heavy instruction + a diff
```

## Applying changes (explicit request only)

A consult is read-only by design. Only when the user asks Codex to *apply* changes, raise the sandbox: `codex-consult -s workspace-write -C <repo> "..."` keeps edits inside the repo; the user's `cdy` alias (`codex --dangerously-bypass-approvals-and-sandbox`: no sandbox, no approvals) grants full access for an interactive session.

## Handling the result

stdout is the verdict — relay it; the stderr log path holds Codex's full reasoning when you want it. Weigh findings as you weigh the built-in advisor's: evaluate, don't apply blindly. When Codex and the advisor disagree, surface the conflict to the user rather than silently picking a side.

The wrapper is a thin shim around the self-building Go source at `~/.claude/skills/codex-advisor/scripts/codex-consult.go`; edit that to change the helper itself.
