# Conventions

## Tools
- GitHub: `gh` CLI for all operations.
- Missing command: retry with `nix run nixpkgs#<command>`.
- Text processing: `perl`, never `sed`/`awk`.
- Long output: `<cmd> | tee /tmp/<name>.log`, then grep the file.
- Run independent tasks in parallel — batch the independent tool calls into one message; prefer symbol-level edits over reading whole files.
- Output in English.

## Prose
- Use the `japanese-tech-writing` skill for all human-facing prose — comments in any language, reports, PR/issue text — and when revising existing prose.
