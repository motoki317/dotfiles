# OpenClaw agents on WSL2 with Nix

A Slack assistant that runs as the `moto` Linux user, with that user's tools and logins, on the Claude Max subscription. One gateway, a systemd user service, serves every agent in this directory's `agents.nix`. `default.nix` holds the config.

## Where things live

| Path | Owner | Contents |
|---|---|---|
| `~/.openclaw/openclaw.json` | Nix | Read-only symlink into the Nix store |
| `~/.openclaw/` (everything else) | OpenClaw | Sessions, SQLite state, `logs/` |
| `~/.claude/projects/-home-moto-projects-openclaw-<id>/` | Claude Code | The agent's Claude conversations: prompts, tool calls, replies, token usage |
| `~/projects/openclaw/<id>/` | agent | Working directory: `MEMORY.md`, `memory/`, scratch files |
| `~/.secrets/openclaw/secrets.json` | you, not in git | Gateway and Slack tokens, and the OpenAI API key for `memory_search` embeddings. Mode 600 and a regular file: OpenClaw rejects symlinks and hard links, so sops-nix and agenix cannot provide it. |
| `~/.secrets/openclaw/agents.nix` | you, not in git | Each agent's Slack channel and member IDs |
| `~/.secrets/openclaw/workspace/`, `~/.secrets/openclaw/agents/<id>/` | you, not in git | Workspace instruction files |
| `~/.secrets/openclaw/slack-alert-prompt.md`, `~/.secrets/openclaw/grafana-api.nix`, `~/.secrets/openclaw/grafana/` | you, not in git | Alert-channel prompt, the `grafana-api` client and its tokens |

This repo is public. Keep workplace details out of it: channel names, alert definitions and handling, internal hosts and URLs. Put them under `~/.secrets/openclaw/`.

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

Nix ignores untracked files, so `git add` each new file in this directory first. `openclaw config set`, `onboard`, `doctor --fix`, and `plugins install` refuse to run in Nix mode. `openclaw-migrate-state` runs `doctor --fix` outside Nix mode, on a throwaway copy of the config. Add a plugin to `runtimePlugins` in `default.nix`.

### Upgrade OpenClaw

`flake.nix` pins a commit of the fork [motoki317/nix-openclaw](https://github.com/motoki317/nix-openclaw). Each repo pins its OpenClaw version in `nix/sources/openclaw-source.nix`. Neither repo can run 2026.10.7 here yet, as [OpenClaw 2026.10.7 does not build here](#openclaw-2026107-does-not-build-here) explains. If upstream nix-openclaw pins the release you want and that section no longer applies, set `github:openclaw/nix-openclaw/<commit>` in `flake.nix` and go to the backup and switch below.

Otherwise, pin the newest stable release in the fork's clone:

```bash
cd ~/projects/nix-openclaw
tag=$(gh release view -R openclaw/openclaw --json tagName -q .tagName)
sha=$(git ls-remote https://github.com/openclaw/openclaw.git "refs/tags/$tag^{}" | cut -f1)
GITHUB_ACTIONS=true scripts/update-pins.sh apply "$tag" "$sha" "" ""
```

Without `GITHUB_ACTIONS=true`, the script refuses to run outside GitHub Actions. The two empty arguments keep the macOS app pin, which this host does not use.

If the gateway build fails with `expected exactly one …` or `contract is unsupported or already patched`, the release changed code that `nix/scripts/patch-openclaw-npm-dist.mjs` matches exactly. Add the new form to the patcher and to its `.fixtures.mjs` and `.test.mjs`, as commit `7f57218` does for 2026.10. Commit that change on its own, then run the script again.

Commit the generated files and push them to the fork's `main`. Set that commit in `flake.nix`.

On either path, `openclaw-migrate-state` migrates state one way, so copy `~/.openclaw` before you switch. The build runs before the stop, so a failed build leaves the gateway running:

```bash
cd ~/.config/home-manager && nix flake lock
home-manager build --flake .#moto --impure
systemctl --user stop openclaw-gateway
cp -a ~/.openclaw ~/.openclaw.bak-$(date +%F)
home-manager switch --flake .#moto --impure && openclaw-migrate-state
```

`openclaw-migrate-state` stops the gateway, runs `openclaw doctor --fix` to migrate the state databases, and starts the gateway again. Until it runs, a newer gateway can exit with status 78 and log `uses schema version …; … run openclaw doctor --fix`. Doctor also merges each agent's `TOOLS.md` into its `AGENTS.md`, and the next switch reverts that merge.

To roll back, stop the gateway and replace `~/.openclaw` with the copy. Restore `flake.nix` and `flake.lock` from git to their state before the upgrade. Then switch and start the gateway:

```bash
home-manager switch --flake ~/.config/home-manager#moto --impure
systemctl --user start openclaw-gateway
```

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

The gateway's output also goes to `~/.openclaw/logs/openclaw-gateway.log`, which survives reboots. Each turn on the subscription logs `provider=claude-cli`, whatever the model. If the gateway's environment has `ANTHROPIC_API_KEY`, turns bill to the API instead of the subscription.

## Keep WSL running

By default, WSL stops the distro 15 s after the last terminal closes, even while systemd services such as the gateway run. `C:\Users\moto\.wslconfig` turns this shutdown off with `instanceIdleTimeout=-1` under `[general]`. WSL reads that file only when it starts.

Linger is on for `moto`, so the gateway starts whenever WSL starts. The `\Start WSL` scheduled task starts WSL at Windows logon with `wsl.exe echo hello`. After a reboot or `wsl --shutdown`, the bot stays offline until your next Windows logon or until you open a WSL terminal.

## Known limitations

### Slack loads as an unverified plugin

Every gateway start logs `OpenClaw can't verify where this plugin came from … reason=record-missing` for `slack`. Nix builds the plugin, so the state database holds no install record that proves it is the official package. In OpenClaw 2026.10.3, an unverified channel plugin cannot run hook agent turns (`dispatchHookAgentTurn`). Its senders get owner authority only from `commands.ownerAllowFrom`, which this config leaves unset. Owner authority unlocks owner-only chat commands and tools, such as `terminal` and `screen`. The Slack plugin does not use hook agent turns.

After an OpenClaw upgrade, check that the gateway log has no `only available for trusted plugins` or `cannot admit authenticated owner` line. Either line means that the release blocks Slack on verification.

Do not reinstall Slack from npm, as the message suggests. An npm install points at the exact gateway build that it was installed for. Any later gateway rebuild, even one from a nixpkgs update, then stops Slack until you reinstall it.

### OpenClaw 2026.10.7 does not build here

As of 2026-10-27, upstream nix-openclaw pins OpenClaw 2026.10.7, whose `package.json` requires `node >=24.22.1 <25 || >=26.12.1`. The `nixos-unstable` nixpkgs has Node 24.21.0, so upstream falls back to the official Node binary and patches it with `autoPatchelfHook`. The patchelf 0.19.2 in nixpkgs corrupts that binary ([NixOS/patchelf#674](https://github.com/NixOS/patchelf/issues/674)), and the build fails when the patched `node` segfaults in its install check. The fork has no such fallback: if it pins 2026.10.7, the gateway runs on a Node that the release does not support. This setup stays on the fork's 2026.10.3, which runs on the Node from nixpkgs.

Move to upstream after one of these reaches `nixos-unstable`: Node 24.22.1, or a patchelf release with [NixOS/patchelf#672](https://github.com/NixOS/patchelf/pull/672). To check, set the upstream commit in `flake.nix` and run `home-manager build --flake ~/.config/home-manager#moto --impure`. The build does not touch the running gateway. If it fails, restore `flake.nix` and `flake.lock` from git.

### acpx warning at startup

From 2026.10.3, the gateway logs `Bundled plugin acpx identity unavailable: plugin runtime artifact contains unsupported symlink entry`. The agents here run on claude-cli, not ACP, so they do not use that plugin.

### Follow-ups abort, not steer

A new message in the same thread or DM conversation stops the running turn, including the command it is running. The bot then starts on the new message. Even "thanks" does this. The exception is an unmentioned reply sent before the bot posts in the thread. Such a reply is a room event: it waits for the running turn to end instead of stopping it. OpenClaw 2026.10.3 and earlier cannot inject a message into a running claude-cli turn ([openclaw#124731](https://github.com/openclaw/openclaw/issues/124731)).

Once a release includes the fix ([openclaw#159816](https://github.com/openclaw/openclaw/pull/159816)), upgrade OpenClaw. Then delete `queue.byChannel.slack` from `default.nix`. Also delete the sentence that begins "A new message in the thread stops your current run" from `slack-channel-prompt.md`. The acknowledgement sentence before it exists because room events do not interrupt. Keep it unless that release lets room events steer.

### Usage and access

Each message from a listed member in a listed channel starts a turn, mention or not. These turns draw on the same Max limits as your interactive Claude Code sessions. Anthropic designs subscription logins "to support ordinary use of Claude Code" ([legal and compliance](https://code.claude.com/docs/en/legal-and-compliance)), and the terms do not clearly cover use by other people.

The agent can read and use every credential `moto` has, and the channel and member allowlists are the only gate. Content the agent fetches can still try to give it instructions. Slack content goes to Anthropic for inference.

## References

- [nix-openclaw](https://github.com/openclaw/nix-openclaw)
- OpenClaw: [Nix mode](https://docs.openclaw.ai/install/nix) · [Anthropic / Claude CLI](https://docs.openclaw.ai/providers/anthropic) · [CLI backends](https://docs.openclaw.ai/gateway/cli-backends) · [Multi-agent](https://docs.openclaw.ai/concepts/multi-agent) · [Secrets](https://docs.openclaw.ai/gateway/config-secrets-env) · [Slack setup](https://docs.openclaw.ai/channels/slack/setup)
- Microsoft: [WSL config](https://learn.microsoft.com/en-us/windows/wsl/wsl-config)
