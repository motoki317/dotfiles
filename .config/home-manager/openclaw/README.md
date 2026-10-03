# OpenClaw agents on WSL2 with Nix

An always-on Slack assistant that runs as the `moto` Linux user, with that user's tools and logins, on the Claude Max subscription with Opus 5.5. The git-tracked home-manager flake owns the configuration. `openclaw/agents.nix` lists the agents.

Tested with nix-openclaw `f62d33f` (OpenClaw 2026.9.5). Agents run your native, self-updating Claude Code (`~/.local/bin/claude`); Opus 5.5 needs 2.1.280 or newer.

## Where things live

| Path | Owner | Contents |
|---|---|---|
| `~/.config/home-manager/openclaw/` | git | Nix module, agent list, workspace instruction files, Slack manifest, this runbook |
| `~/.openclaw/openclaw.json` | Nix | Read-only symlink into the Nix store |
| `~/.openclaw/` (everything else) | OpenClaw | Sessions, SQLite state, the Slack plugin, `logs/` |
| `~/.claude/projects/-home-moto-projects-openclaw-<agent>/` | Claude Code | The agent's Claude conversations: prompts, tool calls, replies, token usage |
| `~/projects/openclaw/<agent>/` | agent | Working directory: `MEMORY.md`, `memory/`, scratch files |
| `~/.secrets/openclaw/secrets.json` | you, not in git | Gateway and Slack tokens. Mode 600, a regular file with a single hard link: OpenClaw rejects symlinks, so sops-nix or agenix cannot provide it. |
| `~/.secrets/openclaw/agents.nix` | you, not in git | Each agent's Slack channel and member IDs, merged into `openclaw/agents.nix` at evaluation. Reading files in `~/.secrets/openclaw` is why every switch needs `--impure`. |
| `~/.secrets/openclaw/workspace/USER.md`, `~/.secrets/openclaw/agents/<agent>/IDENTITY.md` | you, not in git | Who you are, and each agent's name and character |

Every switch copies `AGENTS.md`, `SOUL.md`, `TOOLS.md`, `IDENTITY.md`, and `USER.md` into each workspace. Each comes from the first of these that has it: `~/.secrets/openclaw/agents/<id>/`, `~/.secrets/openclaw/workspace/`, `openclaw/workspace/`. Edit them there; edits in the workspace are overwritten.

## Before you start

- Slack workspace admin approval.
- The channel IDs (`C…`) and member IDs (`U…`) each agent acts on.
- Interactive access to the Claude Max login.
- `sudo` on WSL, and an elevated PowerShell on Windows.

WSL already runs systemd (`/etc/wsl.conf` has `systemd=true`).

## First install

`flake.nix` imports `./openclaw` for the `moto` (WSL) host only. The home directory is the git root, and Nix ignores untracked files, so `git -C ~ add` every new file under `openclaw/`.

1. Create the private files, which stay out of git:
   ```bash
   install -d -m700 ~/.secrets/openclaw/workspace ~/.secrets/openclaw/agents/main
   ```
   - `workspace/USER.md`: who you are, for example your name, timezone, and languages.
   - `agents/main/IDENTITY.md`: the agent's name and character.
   - `agents.nix`: each agent's Slack IDs.
   ```nix
   {
     main.slack = {
       # Slack channel ID -> member IDs whose messages the agent acts on there.
       channels = { C0123456789 = [ "U0123456789" ]; };
       # Member IDs who may DM the agent; [ ] turns DMs off.
       dm = [ "U0123456789" ];
     };
   }
   ```
2. Create the secrets file. Add the Slack tokens after step 4.
   ```bash
   ( umask 077; cat > ~/.secrets/openclaw/secrets.json <<EOF
   {
     "gateway": { "token": "$(head -c 32 /dev/urandom | base64 | tr -d '/+=')" },
     "slack": { "main": { "appToken": "xapp-…", "botToken": "xoxb-…" } }
   }
   EOF
   )
   ```
3. Check the Claude Code login. Agents use your own `~/.claude`: its login, `settings.json`, instructions, skills, and subagents. OpenClaw passes its own effort level and MCP servers.
   ```bash
   claude auth status --text                            # must show the Max account
   systemctl --user show-environment | grep ANTHROPIC_  # must print nothing
   ```
   An `ANTHROPIC_API_KEY` in the service environment would bill turns to the API instead of the subscription.
4. Create one Slack app per agent:
   1. Create the app from `openclaw/slack-app-manifest.json`. For each agent, change both `display_information.name` and `bot_user.display_name`. The manifest covers channel messages, DMs, and replies: no files or slash commands. DMs reach only the members in the agent's `slack.dm` list.
   2. Install it to the workspace and copy the **Bot User OAuth Token** (`xoxb-…`).
   3. Under Basic Information → App-Level Tokens, generate a token with `connections:write` (`xapp-…`).
   4. Put both tokens under `slack.<agent id>` in `secrets.json`.
   5. Invite the bot only to that agent's channels.

   Socket Mode connects outbound, so no port forwarding is needed.
5. Switch, then install the Slack plugin. The installer stops the gateway, installs, and starts it again.
   ```bash
   home-manager switch --flake ~/.config/home-manager#moto --impure
   openclaw-install-slack
   ```
   Then open a new WSL terminal, not a pane of an already running tmux server. Shells started before the switch lack `OPENCLAW_DISABLE_PERSISTED_PLUGIN_REGISTRY=0`, so their `openclaw` does not see the Slack plugin.
6. Keep it running across reboots:
   ```bash
   sudo loginctl enable-linger moto
   ```
   Systemd alone does not keep WSL running. In an elevated PowerShell, as the Windows user who owns the distro (not SYSTEM):
   ```powershell
   $Distro = "Ubuntu"       # this machine's name in `wsl --list --verbose`
   schtasks /Create /TN "OpenClaw WSL Boot" `
     /TR "wsl.exe -d $Distro -u moto --exec dbus-launch true" `
     /SC ONSTART /RU "$env:USERDOMAIN\$env:USERNAME" /F
   schtasks /Run /TN "OpenClaw WSL Boot"
   ```
   `schtasks` asks for your Windows password and stores it, so the task runs at boot without a logon. This is OpenClaw's documented workaround, not a guarantee for every WSL build. Keep Windows awake and online.

## Add an agent

1. Add a block to `openclaw/agents.nix`, and its Slack channel and member IDs to `~/.secrets/openclaw/agents.nix`. The optional `entry` overrides that agent's settings, for example its model. Every Anthropic model runs through the subscription.
2. Create `~/.secrets/openclaw/agents/<id>/IDENTITY.md`. The build fails without it. Other files there override the shared ones for that agent.
3. `git -C ~ add .config/home-manager/openclaw`
4. Create its Slack app (First install, step 4) and add `slack.<id>` to `secrets.json`.
5. `home-manager switch --flake ~/.config/home-manager#moto --impure && systemctl --user restart openclaw-gateway`

All agents share one gateway, one Claude login and allowance, and one tool policy. Separating them needs a second gateway, which this module does not set up.

## Acceptance tests

```bash
systemctl --user status openclaw-gateway
openclaw gateway status --require-rpc
openclaw channels status --probe
openclaw models status
openclaw logs --follow        # kept 24 h; persistent copy: ~/.openclaw/logs/openclaw-gateway.log
```

| Test | Required evidence |
|---|---|
| Authorized mention | Mention the bot and ask it to run `whoami`, `pwd`, and `command -v gh git codex`. It replies in the thread with `moto`, `~/projects/openclaw/<id>`, and three paths. |
| Subscription route | The log shows `provider=claude-cli model=claude-opus-5-5` for that turn. |
| Raw tools | It creates, reads, and deletes a file under `/tmp` without an approval prompt. |
| Ambient context | A listed member posts without a mention. The bot stays quiet, then recalls the message when mentioned. |
| Direct message | A member in the agent's `slack.dm` list DMs the bot and gets a reply. |
| Access filters | Messages from an unlisted member, in an unlisted channel, or a DM from someone outside `slack.dm` get no reply, and no `provider=claude-cli` log line follows them. |
| Persistence | Reboot Windows and test Slack **before logging in**. |

## Operations

- **Restart:** `systemctl --user restart openclaw-gateway`.
- **Stop:** `systemctl --user stop openclaw-gateway` stops the gateway and every process an agent started under it. A process handed to an already running server, such as your tmux, survives. `wsl --terminate <distro>` in PowerShell stops everything.
- **Config change:** edit the files in `openclaw/` (`git add` new ones) and switch, or edit `secrets.json` or the private `agents.nix` (then switch). Then restart: the gateway reads its config and secrets only at start. `openclaw config set`, `onboard`, `doctor --fix`, and `plugins install` refuse to run in Nix mode.
- **Upgrade OpenClaw:** change the pinned nix-openclaw revision in `flake.nix`, then:
  ```bash
  cd ~/.config/home-manager && nix flake lock
  home-manager switch --flake .#moto --impure && openclaw-install-slack
  ```
  The Slack plugin lives in `~/.openclaw`, outside the generation. After any switch that changes the OpenClaw version, rollbacks included, run `openclaw-install-slack` again.
- **New Claude model:** change `model.primary` in `default.nix`. Claude Code updates itself; check `claude --version` if the model needs a newer release.

## Known limitations

- **Slack plugin outside Nix.** nix-openclaw's `runtimePlugins = [ "slack" ]` fails OpenClaw 2026.9.5's trust gate: `openChannelIngressQueue is only available for trusted plugins … reason=record-missing`. [nix-openclaw#158](https://github.com/openclaw/nix-openclaw/issues/158) reports the same gate failure for Discord. `openclaw-install-slack` installs the official npm package instead. Once load-path plugins pass the gate, switch back: add `runtimePlugins`, and delete `installSlackPlugin`, its `home.packages` entry, and both `OPENCLAW_DISABLE_PERSISTED_PLUGIN_REGISTRY` settings from `default.nix`.
- **Subscription terms.** Anthropic designs subscription logins "to support ordinary use of Claude Code" and states that Max limits "assume ordinary, individual usage". OpenClaw runs the real Claude Code with its own login and never handles the token. Every member you list drives turns on your subscription, and the terms do not clearly cover use by other people. The template lists only you.
- **Every listed message is a turn.** Each message from a listed member in a listed channel starts an Opus turn, mention or not. The bot stays quiet because only its explicit message-tool calls post to Slack. These turns draw on the same Max limits as your interactive Claude Code sessions.
- **Raw access.** The agent can read and use every credential `moto` has. The channel and member allowlists are the only gate. Content the agent fetches can still try to steer it. Slack content goes to Anthropic for inference.

## References

- [nix-openclaw](https://github.com/openclaw/nix-openclaw)
- OpenClaw: [Nix mode](https://docs.openclaw.ai/install/nix) · [Anthropic / Claude CLI](https://docs.openclaw.ai/providers/anthropic) · [CLI backends](https://docs.openclaw.ai/gateway/cli-backends) · [Multi-agent](https://docs.openclaw.ai/concepts/multi-agent) · [Secrets](https://docs.openclaw.ai/gateway/config-secrets-env) · [Exec approvals](https://docs.openclaw.ai/tools/exec-approvals) · [Slack setup](https://docs.openclaw.ai/channels/slack/setup) · [Windows/WSL](https://docs.openclaw.ai/platforms/windows)
- Anthropic: [Claude Code legal and compliance](https://code.claude.com/docs/en/legal-and-compliance)
