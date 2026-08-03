# Conventions

## Tools
- GitHub: `gh` CLI for all operations.
- HTTP requests: `ax`, not `curl` — curl flags work; every fetch prints a `{status, ok, ms, headers, body}` report.
- Web pages: `ax <url> --md --all --budget <n>` gives the page's exact text (`--all` matters: without it extract modes stop at 50 blocks; `--budget` still caps tokens, truncation is announced). Empty output = JS-rendered page: try appending `.md` to the URL (docs sites often serve raw markdown — read it with `--body`), else WebFetch. WebFetch answers through a sub-model summary — use it for a single fact in a huge page, or claude.ai URLs only it can reach.
- Missing command: retry with `nix run nixpkgs#<command>`.
- Text processing: `perl`, never `sed`/`awk`.
- Long output: `<cmd> | tee /tmp/<name>.log`, then grep the file.
- Run independent tasks in parallel — batch the independent tool calls into one message; prefer symbol-level edits over reading whole files.
- Output in English.

## Prose
- Load the `tech-writing` skill before you write or revise human-facing prose — reports, PR/issue text, docs, comment blocks — in any language.
- In any prose: one point per sentence, condition before command, one name per concept. State requirements as requirements ("must", not "should") and real uncertainty as uncertainty.
- Shorten by cutting points and sentences, never grammar or needed context.
- At task completion, run `/cold-read` over the task's durable prose and apply the report.
