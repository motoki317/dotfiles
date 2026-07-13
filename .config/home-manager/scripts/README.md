# scripts

Small Go helpers, one module per directory, built and installed on PATH by
home-manager — see the `goBin` helper in [`../common/packages.nix`](../common/packages.nix).

| Directory             | Binary               | Used by                                              |
| --------------------- | -------------------- | ---------------------------------------------------- |
| `claude-statusline/`  | `claude-statusline`  | Claude Code status line (`~/.claude/settings.json`)  |
| `codex-run/`          | `codex-run`          | the `codex-advisor` and `execute` skills             |
| `session-transcript/` | `session-transcript` | `codex-run … --context` (invokes it by PATH name)    |
| `sync-claude-skills/` | `sync-claude-skills` | shared Claude/Codex skill registry                   |

## Editing these

They are compiled at `home-manager switch`, **not** at runtime (they used to
self-build on every call; that was dropped for one uniform build path). So:

> **After editing any file here, run `home-manager switch` — a plain save changes nothing.**

```sh
# iterate in-dir without a switch (buildGoModule runs these same tests):
cd claude-statusline && go test ./...

# apply your changes to the live binaries:
home-manager switch --flake ~/.config/home-manager#<host>   # <host>: toki (macOS), moto (WSL)
```

`buildGoModule` runs each module's tests in its check phase, so a failing test
fails the switch. Binaries are per-host: switch on each machine you use.
