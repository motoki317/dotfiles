# Critical Rules

- Use `perl` for text processing (never `sed` or `awk`)
- Execute independent tasks in parallel (multiple Task calls in one message)
- Use symbol-level operations over reading entire files
- Output in English

# Tool Usage

- **GitHub**: Use `gh` CLI for all operations (PRs, issues, repos)
- **Docs**: Use Context7 MCP to verify latest library documentation
- **Missing commands**: Retry with `nix run nixpkgs#<command>`

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
