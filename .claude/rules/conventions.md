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
- Use the `japanese-tech-writing` skill for all human-facing prose — comments in any language, reports, PR/issue text — and when revising existing prose.
- At task completion, run `/cold-read` over the task's durable prose — docs, comment blocks, PR/issue text — and apply the report.
