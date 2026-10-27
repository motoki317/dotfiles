# One OpenClaw gateway (systemd user service) serving every agent in ./agents.nix.
# Agents run unsandboxed as this user, on the Claude subscription through the
# native, self-updating Claude Code in ~/.local/bin. Runbook: ./README.md.
{
  config,
  lib,
  pkgs,
  inputs,
  ...
}:
let
  home = config.home.homeDirectory;
  # Slack IDs, the owner (USER.md), and each agent's identity stay out of git,
  # so evaluating this module needs --impure.
  privateDir = "${home}/.secrets/openclaw";
  privateAgentsFile = "${privateDir}/agents.nix";
  privateAgents = import privateAgentsFile;
  agents = lib.mapAttrs (id: agent: agent // { slack = privateAgents.${id}.slack or { }; }) (
    import ./agents.nix
  );
  # Gateway-wide duties (heartbeat, system agent, auth source, talk) need one named owner.
  owner = if agents ? main then "main" else lib.head (lib.attrNames agents);
  # Workspaces hold runtime files (memory, scratch); git-tracked config stays here.
  workspaceRoot = "${home}/projects/openclaw";
  stateDir = "${home}/.openclaw";
  # Private 0600 JSON outside git; OpenClaw rejects symlinked or hardlinked secret files.
  secretsFile = "${privateDir}/secrets.json";
  secret = pointer: {
    source = "file";
    provider = "local";
    id = pointer;
  };
  # Internal hosts and query notes stay out of this public repo.
  grafanaApi = import "${privateDir}/grafana-api.nix" { inherit pkgs privateDir; };

  # OpenClaw refuses workspace files that symlink outside the workspace, so
  # each is copied from the first that exists. The files are read at
  # activation and never enter the Nix store.
  workspaceFileNames = [
    "AGENTS.md"
    "SOUL.md"
    "TOOLS.md"
    "IDENTITY.md"
    "USER.md"
  ];
  workspaceFile =
    id: name:
    lib.findFirst builtins.pathExists
      (throw "openclaw: create ${privateDir}/workspace/${name} or ${privateDir}/agents/${id}/${name}.")
      [
        "${privateDir}/agents/${id}/${name}"
        "${privateDir}/workspace/${name}"
      ];

  gateway = config.programs.openclaw.instances.default.package;

  # A new OpenClaw can refuse to start until `doctor --fix` migrates its state
  # databases, and doctor refuses Nix mode. Its config edits go to a throwaway
  # copy. Outside Nix mode it also rewrites and starts the managed gateway unit
  # (docs/cli/doctor/recovery.md), so it gets a throwaway HOME and no user bus.
  # Outside Nix mode, OpenClaw also rejects the hardlinks that store optimisation
  # leaves in Nix-built plugins, and doctor then reinstalls each one from npm.
  # So the throwaway config loads plain copies of the plugins.
  migrateState = pkgs.writeShellScriptBin "openclaw-migrate-state" ''
    set -euo pipefail
    jq=${lib.getExe pkgs.jq}
    config=${lib.escapeShellArg "${stateDir}/openclaw.json"}
    tmp=$(mktemp -d)
    trap 'rm -rf "$tmp"' EXIT
    mkdir "$tmp/plugins"
    "$jq" -r '.plugins.load.paths // [] | .[]' "$config" | while read -r plugin; do
      cp -r --no-preserve=mode "$plugin" "$tmp/plugins/"
    done
    (umask 077 && "$jq" --arg dir "$tmp/plugins" \
      'if .plugins.load.paths then .plugins.load.paths |= map("\($dir)/\(split("/") | last)") else . end' \
      "$config" > "$tmp/openclaw.json")
    systemctl --user stop openclaw-gateway
    env -u DBUS_SESSION_BUS_ADDRESS -u XDG_RUNTIME_DIR HOME="$tmp" \
      OPENCLAW_NIX_MODE=0 OPENCLAW_CONFIG_PATH="$tmp/openclaw.json" \
      OPENCLAW_STATE_DIR=${lib.escapeShellArg stateDir} \
      ${gateway}/bin/openclaw doctor --fix --non-interactive
    systemctl --user start openclaw-gateway
  '';

  slackAccount =
    id: agent:
    let
      dmUsers = agent.slack.dm or [ ];
    in
    {
      # DMs only from the listed members; an empty list keeps DMs off.
      dmPolicy = if dmUsers == [ ] then "disabled" else "allowlist";
      allowFrom = dmUsers;
      appToken = secret "/slack/${id}/appToken";
      botToken = secret "/slack/${id}/botToken";
    }
    # Reads go through the user token, so the agent reads channels the bot has not
    # joined; posts stay on the bot token (userTokenReadOnly defaults to true).
    // lib.optionalAttrs (agent.slack.userToken or false) {
      userToken = secret "/slack/${id}/userToken";
    }
    // {
      channels = lib.mapAttrs (
        channelId: users:
        let
          alertSenders = agent.slack.botChannels.${channelId} or null;
          botChannel = alertSenders != null;
        in
        {
          enabled = true;
          requireMention = false;
          ignoreOtherMentions = false;
          replyToMode = "all";
          # Overrides channels.slack.allowBots = false for alert channels only.
          allowBots = botChannel;
          # The member allowlist gates bot senders too, so alert bots must be listed.
          users = users ++ lib.optionals botChannel alertSenders;
          systemPrompt =
            builtins.readFile ./slack-channel-prompt.md
            + lib.optionalString botChannel (builtins.readFile "${privateDir}/slack-alert-prompt.md");
        }
      ) (agent.slack.channels or { });
    };
in
{
  imports = [ inputs.nix-openclaw.homeManagerModules.openclaw ];
  nixpkgs.overlays = [ inputs.nix-openclaw.overlays.default ];

  assertions = lib.concatLists (
    lib.mapAttrsToList (id: agent: [
      {
        assertion = privateAgents ? ${id};
        message = "openclaw: add agent ${id}'s Slack IDs to ${privateAgentsFile}.";
      }
      {
        # A shared fallback would give every agent the same identity.
        assertion = builtins.pathExists "${privateDir}/agents/${id}/IDENTITY.md";
        message = "openclaw: create ${privateDir}/agents/${id}/IDENTITY.md.";
      }
      {
        # The Slack plugin treats an empty member list as "every member".
        assertion = lib.all (users: users != [ ]) (lib.attrValues (agent.slack.channels or { }));
        message = "openclaw: agent ${id} has a Slack channel with no member IDs.";
      }
    ]) agents
  );

  programs.openclaw = {
    # Agents use this user's tools (the unit PATH below). The bundled toolchain
    # would shadow them on the gateway's PATH and adds ~2.8 GB.
    toolNames = [ ];
    instances.default = {
      inherit stateDir;
      workspaceDir = workspaceRoot;
      # nix-openclaw builds these against its overlay's gateway, which this instance
      # runs while `package` keeps its default. A gateway rebuild rebuilds them too.
      runtimePlugins = [ "slack" ];
      # The default under /tmp is emptied at boot, and only activation recreates it.
      logPath = "${stateDir}/logs/openclaw-gateway.log";
      # Agents share this user's ~/.claude (login, settings, instructions). No
      # CLAUDE_CONFIG_DIR: OpenClaw 2026.9.5 resumes a conversation only if it
      # finds the transcript under ~/.claude, so a separate directory makes the
      # agent forget once its Claude process idles (openclaw/openclaw#145309).
      config = {
        secrets.providers.local = {
          source = "file";
          path = secretsFile;
          mode = "json";
        };
        # Embeddings for memory_search only; chat stays on the Claude subscription.
        models.providers.openai.apiKey = secret "/openai/apiKey";
        # One shared ClickUp login (toki's): `openclaw mcp login clickup`.
        # Every allowed Slack member reads and writes ClickUp as toki, by toki's choice.
        mcp.servers.clickup = {
          url = "https://mcp.clickup.com/mcp";
          transport = "streamable-http";
          auth = "oauth";
        };
        # Same model for Slite: one shared login (toki's) with read and write,
        # by toki's choice. `openclaw mcp login slite`.
        mcp.servers.slite = {
          url = "https://api.slite.com/mcp";
          transport = "streamable-http";
          auth = "oauth";
        };
        # Same model for Sentry, for alert investigation. `openclaw mcp login sentry`.
        mcp.servers.sentry = {
          url = "https://mcp.sentry.dev/mcp";
          transport = "streamable-http";
          auth = "oauth";
        };
        # Every Slack thread and heartbeat run is its own session, and by default its
        # MCP runtime lives until the gateway stops. After a week they hit the
        # gateway's 256-runtime cap and new turns failed (2026-10-26).
        mcp.sessionIdleTtlMs = 3600000;
        gateway = {
          mode = "local";
          bind = "loopback";
          auth.token = secret "/gateway/token";
        };
        agents = {
          ownership = "explicit";
          defaults = {
            workspace = workspaceRoot;
            # Workspace files come from git; OpenClaw must not seed its templates.
            skipBootstrap = true;
            sandbox.mode = "off";
            heartbeat =
              let
                reportChannel = agents.${owner}.slack.heartbeatChannel or null;
              in
              {
                # Explicit ownership needs a named owner even with periodic turns off.
                agentId = owner;
              }
              // (
                if reportChannel == null then
                  { every = "0m"; }
                else
                  {
                    # Each turn draws on the Claude subscription, so hourly. Around the
                    # clock: an urgent item must not wait for morning.
                    every = "1h";
                    # The checklist lives in memory/todo.md, so prior chat history only costs tokens.
                    isolatedSession = true;
                    target = "slack";
                    accountId = owner;
                    to = "channel:${reportChannel}";
                    prompt = builtins.readFile ./heartbeat-prompt.md;
                  }
              );
            systemAgent.agentId = owner;
            authInheritance.agentId = owner;
            model.primary = "anthropic/claude-opus-5-5";
            # Subscription billing: every Anthropic model, including per-agent
            # overrides, runs through Claude Code, never the paid API.
            models."anthropic/*".agentRuntime.id = "claude-cli";
          };
          entries = lib.mapAttrs (
            id: agent: { workspace = "${workspaceRoot}/${id}"; } // (agent.entry or { })
          ) agents;
        };
        bindings = lib.mapAttrsToList (id: _: {
          agentId = id;
          match = {
            channel = "slack";
            accountId = id;
          };
        }) agents;
        talk.agentId = owner;
        tools = {
          profile = "full";
          fs.workspaceOnly = false;
          exec = {
            host = "gateway";
            mode = "full";
            strictInlineEval = false;
            applyPatch.workspaceOnly = false;
          };
        };
        messages = {
          ackReactionScope = "off";
          # A Slack follow-up aborts the running turn and starts on the new message: the
          # default `steer` cannot reach a claude-cli turn and waits for it to end
          # (openclaw/openclaw#124731). Once a release ships the fix
          # (openclaw/openclaw#159816), upgrade and delete this line to steer instead.
          # Room events never interrupt, so slack-channel-prompt.md has the agent post
          # in the thread first.
          queue.byChannel.slack = "interrupt";
          groupChat = {
            unmentionedInbound = "room_event";
            visibleReplies = "message_tool";
            historyLimit = 50;
          };
        };
        channels.slack = {
          enabled = true;
          mode = "socket";
          groupPolicy = "allowlist";
          allowBots = false;
          joinIntro = false;
          accounts = lib.mapAttrs slackAccount agents;
        };
      };
    };
  };

  systemd.user.services.openclaw-gateway = {
    # nix-openclaw's unit has no [Install] section, so it would not start at boot.
    Install.WantedBy = [ "default.target" ];
    # The claude-cli runtime finds `claude` here, and agent commands inherit it
    # through OpenClaw's exec tool. Same order as the login shell; systemd's
    # default lacks the Nix and ~/.local/bin tools.
    Service.Environment = [
      "PATH=${home}/.nix-profile/bin:/nix/var/nix/profiles/default/bin:${home}/.local/bin:/usr/local/bin:/usr/bin:/bin"
    ];
  };

  home.packages = [
    migrateState
    grafanaApi
  ];

  home.activation.openclawAgentWorkspaces = lib.hm.dag.entryAfter [ "writeBoundary" ] (
    lib.concatStrings (
      lib.mapAttrsToList (
        id: _:
        lib.concatMapStrings (name: ''
          run install -Dm644 ${lib.escapeShellArg "${workspaceFile id name}"} ${lib.escapeShellArg "${workspaceRoot}/${id}/${name}"}
        '') workspaceFileNames
      ) agents
    )
  );
}
