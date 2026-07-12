# Conventions

## Tools
- GitHub: `gh` CLI for all operations.
- HTTP requests: `ax`, not `curl` — curl flags work; every fetch prints a `{status, ok, ms, headers, body}` report. Details: `ax --help`.
- Web pages: `ax <url> --md --budget <n>` gives the page's exact text, cut from the top when over budget (a typical docs page is ~2k tokens whole). WebFetch answers through a sub-model summary instead — use it for a single fact in a huge page, or for claude.ai URLs only it can reach; when exact wording matters, ax.
- Missing command: retry with `nix run nixpkgs#<command>`.
- Text processing: `perl`, never `sed`/`awk`.
- Long output: `<cmd> | tee /tmp/<name>.log`, then grep the file.
- Run independent tasks in parallel — batch the independent tool calls into one message; prefer symbol-level edits over reading whole files.
- Output in English.

## Prose
- Use the `japanese-tech-writing` skill for all human-facing prose — comments in any language, reports, PR/issue text — and when revising existing prose.
