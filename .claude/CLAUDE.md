# Critical Rules

- Use `perl` for text processing (never `sed` or `awk`)
- Execute independent tasks in parallel (multiple Task calls in one message)
- Use symbol-level operations over reading entire files
- Output in English

# Coding

- Codes should explain itself, HOW things are done
- Test codes should explain WHAT the code does
- Code comments and commit messages should explain WHY a code exists and WHY NOT a feature / an alternative was not selected

# Tool Usage

- **GitHub**: Use `gh` CLI for all operations (PRs, issues, repos)
- **Docs**: Use Context7 MCP to verify latest library documentation
- **Missing commands**: Retry with `nix run nixpkgs#<command>`
- **Long output**: Save with `<cmd> | tee /tmp/<name>.log` first, then grep the file. Lets you refine filters or inspect surrounding context without re-running the source command.

# Workflow

- Delegate detailed work to sub-agents; focus on orchestration
- Check existing code/patterns before implementing new features
- Only perform Git operations when explicitly requested
- Always verify the valid tool/library usage or configuration values by checking official docs. Don't guess.

# Error Handling

- Sub-agent fails → retry with alternative agent
- Memory not found → document gap, investigate
- Conflicting outputs → flag uncertainty to user
- Security/destructive risk → BLOCK and require acknowledgment
