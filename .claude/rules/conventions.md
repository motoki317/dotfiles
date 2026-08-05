# Conventions

- GitHub: `gh` CLI for all operations.
- HTTP and web pages: `ax`, never `curl` — modes and recovery paths: `~/.claude/refs/web.md`.
- Text processing: `perl`, never `sed`/`awk`. Missing command: retry with `nix run nixpkgs#<command>`.
- Long output: `<cmd> | tee /tmp/<name>.log`, then grep the file.
- Batch independent tool calls into one message; read and edit at symbol level, not whole files.
- Output in English.
