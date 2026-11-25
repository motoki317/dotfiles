# Conventions

- GitHub: `gh`. A shell wrapper picks the account from the git origin of the cwd, not from `-R`: attmcojp → toki-attm, other GitHub → motoki317, no origin → the active account. `gh auth status` bypasses it. To reach a non-attmcojp repo from outside its clone, run `gh` from `~`.
- HTTP and web pages: `ax`, never `curl` — modes and recovery paths: `~/.claude/refs/web.md`.
- In-place edits from the shell: `perl -i`, never `sed -i` or `awk`. With `-CSD` and non-ASCII text in the code, also pass `-Mutf8`. Missing command: retry with `nix run nixpkgs#<command>`.
