# Conventions

- GitHub: `gh`. A shell wrapper picks the account from the cwd's git origin, never `-R`: attmcojp → toki-attm, other GitHub → motoki317, no origin → the active account. `gh auth status` bypasses it. For a non-attmcojp repo, run `gh` from `~`.
- HTTP and web pages: `ax`, never `curl` — modes and recovery paths: `~/.claude/refs/web.md`.
- In-place edits from the shell: `perl -i`, never `sed -i` or `awk`. Missing command: retry with `nix run nixpkgs#<command>`.
- Long output: `<cmd> | tee <log file>`, then grep the file.
- Batch independent tool calls into one message. Read and edit at symbol level, not whole files.
- Output in English.
