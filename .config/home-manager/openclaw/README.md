# OpenClaw agents on WSL2 with Nix

A Slack assistant that runs as the `moto` Linux user, with that user's tools and logins, on the Claude Max subscription. One gateway, a systemd user service, serves every agent in this directory's `agents.nix`. `default.nix` holds the config.

## Where things live

| Path | Owner | Contents |
|---|---|---|
| `~/.openclaw/openclaw.json` | Nix | Read-only symlink into the Nix store |
| `~/.openclaw/` (everything else) | OpenClaw | Sessions, SQLite state, the Slack plugin, `logs/` |
| `~/.claude/projects/-home-moto-projects-openclaw-<id>/` | Claude Code | The agent's Claude conversations: prompts, tool calls, replies, token usage |
| `~/projects/openclaw/<id>/` | agent | Working directory: `MEMORY.md`, `memory/`, scratch files |
| `~/.secrets/openclaw/secrets.json` | you, not in git | Gateway and Slack tokens, and the OpenAI API key for `memory_search` embeddings. Mode 600 and a regular file: OpenClaw rejects symlinks and hard links, so sops-nix and agenix cannot provide it. |
| `~/.secrets/openclaw/agents.nix` | you, not in git | Each agent's Slack channel and member IDs |
| `~/.secrets/openclaw/workspace/`, `~/.secrets/openclaw/agents/<id>/` | you, not in git | Workspace instruction files |

Every switch copies `AGENTS.md`, `SOUL.md`, `TOOLS.md`, `IDENTITY.md`, and `USER.md` into each agent's working directory. Each file comes from `~/.secrets/openclaw/agents/<id>/` if it exists there, else from `~/.secrets/openclaw/workspace/`. `IDENTITY.md` must exist in the agent's own directory. Edit the files under `~/.secrets/openclaw/`: a switch overwrites edits in the working directory.

Agents run your native, self-updating Claude Code (`~/.local/bin/claude`) with your own `~/.claude`: its login, `settings.json`, instructions, skills, and subagents. OpenClaw replaces your MCP servers and effort level with its own.

## Operations

```bash
systemctl --user restart openclaw-gateway
systemctl --user stop openclaw-gateway
```

Stopping the gateway also stops every process an agent started under it. A process that an agent starts in a tmux server that was already running, such as yours, survives.

### Change the config

After you edit `secrets.json`, restart. After you edit this directory, the private `agents.nix`, or the workspace files, switch and restart:

```bash
home-manager switch --flake ~/.config/home-manager#moto --impure && systemctl --user restart openclaw-gateway
```

Nix ignores untracked files, so `git add` each new file in this directory first. `openclaw config set`, `onboard`, `doctor --fix`, and `plugins install` refuse to run in Nix mode.

### Upgrade OpenClaw

Change the pinned nix-openclaw revision in `flake.nix`, then run:

```bash
cd ~/.config/home-manager && nix flake lock
home-manager switch --flake .#moto --impure && openclaw-install-slack
```

`openclaw-install-slack` stops the gateway, installs the Slack plugin, and starts the gateway again. Run it after every switch that changes the OpenClaw version, rollbacks included.

### Change the Claude model

Change `model.primary` in `default.nix`. A new model can need a newer Claude Code. `claude --version` shows the installed release, and `claude update` installs the latest.

## Add an agent

1. Add a block named `<id>` to this directory's `agents.nix`. The optional `entry` overrides that agent's config, for example its model. Every Anthropic model runs through the subscription.
2. Add its Slack IDs to `~/.secrets/openclaw/agents.nix`:
   ```nix
   <id>.slack = {
     # Slack channel ID -> member IDs whose messages the agent acts on there.
     channels = { C0123456789 = [ "U0123456789" ]; };
     # Member IDs who may DM the agent; [ ] turns DMs off.
     dm = [ "U0123456789" ];
   };
   ```
3. Create `~/.secrets/openclaw/agents/<id>/IDENTITY.md` with the agent's name and character.
4. Copy `slack-app-manifest.json`. In the copy, change `display_information.name` and `features.bot_user.display_name`.
5. Create a Slack app from the copy.
6. Install the app to the Slack workspace.
7. Under Basic Information → App-Level Tokens, generate a token with `connections:write` (`xapp-…`).
8. Put that token and the Bot User OAuth Token (`xoxb-…`) under `slack.<id>` in `secrets.json`, as `appToken` and `botToken`.
9. Invite the bot only to that agent's channels.
10. Switch and restart, as in [Change the config](#change-the-config).

All agents share one gateway, one `~/.claude`, and one Claude login and allowance. Separating them needs a second gateway, which this module does not set up.

## Check health

```bash
systemctl --user status openclaw-gateway
openclaw gateway status --require-rpc
openclaw channels status --probe    # each agent: "running, connected, … works"
openclaw logs --follow              # the daily log in /tmp/openclaw/, kept 24 h
journalctl --user -u openclaw-gateway -g 'Started|Stopping'   # each gateway start and stop
claude auth status --text           # must show the Max account
tr '\0' '\n' < /proc/$(systemctl --user show -p MainPID --value openclaw-gateway)/environ | grep ANTHROPIC_   # must print nothing
```

The gateway's output also goes to `~/.openclaw/logs/openclaw-gateway.log`, which survives reboots. Turns on the subscription log `provider=claude-cli`, whatever the model. If the gateway's environment has `ANTHROPIC_API_KEY`, turns bill to the API instead of the subscription.

## Keep WSL running

By default, WSL stops the distro 15 s after the last terminal closes, even while systemd services such as the gateway run. `C:\Users\moto\.wslconfig` turns this shutdown off with `instanceIdleTimeout=-1` under `[general]`. WSL reads that file only when it starts.

Linger is on for `moto`, so the gateway starts whenever WSL starts. The `\Start WSL` scheduled task starts WSL at Windows logon with `wsl.exe echo hello`. After a reboot or `wsl --shutdown`, the bot stays offline until your next Windows logon or until you open a WSL terminal.

## Known limitations

### Slack plugin outside Nix

nix-openclaw's `runtimePlugins = [ "slack" ]` fails OpenClaw 2026.9.5's trust gate: `openChannelIngressQueue is only available for trusted plugins … reason=record-missing`. [nix-openclaw#158](https://github.com/openclaw/nix-openclaw/issues/158) reports the same gate failure for Discord. `openclaw-install-slack` installs the official npm package instead.

Once `runtimePlugins` passes the gate, move Slack back into Nix. Add `runtimePlugins = [ "slack" ]` to `default.nix`. Delete `installSlackPlugin`, its `home.packages` entry, and the `environment` and `home.sessionVariables` entries that set `OPENCLAW_DISABLE_PERSISTED_PLUGIN_REGISTRY`.

### Follow-ups abort, not steer

A new message in the same thread or DM conversation stops the running turn, including the command it is running. The bot then starts on the new message. Even "thanks" does this. The exception is an unmentioned reply sent before the bot posts in the thread. It is a room event, which does not stop the running turn and waits for it to end. OpenClaw cannot inject a message into a running claude-cli turn yet ([openclaw#124731](https://github.com/openclaw/openclaw/issues/124731)).

Once a release includes the fix ([openclaw#159816](https://github.com/openclaw/openclaw/pull/159816)), upgrade OpenClaw. Then delete `queue.byChannel.slack` from `default.nix`. Also delete the sentence that begins "A new message in the thread stops your current run" from `slack-channel-prompt.md`.

### Usage and access

Each message from a listed member in a listed channel starts a turn, mention or not. These turns draw on the same Max limits as your interactive Claude Code sessions. Anthropic designs subscription logins "to support ordinary use of Claude Code" ([legal and compliance](https://code.claude.com/docs/en/legal-and-compliance)), and the terms do not clearly cover use by other people.

The agent can read and use every credential `moto` has, and the channel and member allowlists are the only gate. Content the agent fetches can still try to give it instructions. Slack content goes to Anthropic for inference.

## References

- [nix-openclaw](https://github.com/openclaw/nix-openclaw)
- OpenClaw: [Nix mode](https://docs.openclaw.ai/install/nix) · [Anthropic / Claude CLI](https://docs.openclaw.ai/providers/anthropic) · [CLI backends](https://docs.openclaw.ai/gateway/cli-backends) · [Multi-agent](https://docs.openclaw.ai/concepts/multi-agent) · [Secrets](https://docs.openclaw.ai/gateway/config-secrets-env) · [Slack setup](https://docs.openclaw.ai/channels/slack/setup)
- Microsoft: [WSL config](https://learn.microsoft.com/en-us/windows/wsl/wsl-config)
