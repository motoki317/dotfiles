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

  # OpenClaw refuses workspace files that symlink outside the workspace, so
  # each is copied from the first that exists. Private files are read at
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
        (./workspace + "/${name}")
      ];

  gateway = config.programs.openclaw.instances.default.package;

  # nix-openclaw's load-path plugins fail OpenClaw's trust gate for channel
  # ingress (nix-openclaw#158), so Slack comes from the official npm package,
  # whose install record in the state dir is trusted. The installer writes config
  # migrations, so it gets a throwaway copy of the Nix config. A running gateway
  # would take over the install and refuse it in Nix mode, so it is stopped first.
  installSlackPlugin = pkgs.writeShellScriptBin "openclaw-install-slack" ''
    set -eu
    # `plugins install` shells out to npm.
    export PATH=${lib.makeBinPath [ pkgs.nodejs_24 ]}:$PATH
    tmp=$(mktemp -d)
    trap 'rm -rf "$tmp"' EXIT
    install -m600 ${lib.escapeShellArg "${stateDir}/openclaw.json"} "$tmp/openclaw.json"
    systemctl --user stop openclaw-gateway
    OPENCLAW_NIX_MODE=0 OPENCLAW_CONFIG_PATH="$tmp/openclaw.json" \
      OPENCLAW_STATE_DIR=${lib.escapeShellArg stateDir} OPENCLAW_DISABLE_PERSISTED_PLUGIN_REGISTRY=0 \
      ${gateway}/bin/openclaw plugins install "@openclaw/slack@${gateway.version}" --force
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
      channels = lib.mapAttrs (_: users: {
        enabled = true;
        requireMention = false;
        ignoreOtherMentions = false;
        replyToMode = "all";
        inherit users;
        systemPrompt = builtins.readFile ./slack-channel-prompt.md;
      }) (agent.slack.channels or { });
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
      # The default under /tmp is emptied at boot, and only activation recreates it.
      logPath = "${stateDir}/logs/openclaw-gateway.log";
      # Agents share this user's ~/.claude (login, settings, instructions). No
      # CLAUDE_CONFIG_DIR: OpenClaw 2026.9.5 resumes a conversation only if it
      # finds the transcript under ~/.claude, so a separate directory makes the
      # agent forget once its Claude process idles (openclaw/openclaw#145309).
      environment = {
        # The package wrapper defaults this to 1, which hides the trusted Slack install record.
        OPENCLAW_DISABLE_PERSISTED_PLUGIN_REGISTRY = "0";
      };
      config = {
        plugins.entries.slack.enabled = true;
        secrets.providers.local = {
          source = "file";
          path = secretsFile;
          mode = "json";
        };
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
            heartbeat = {
              # No periodic turns; explicit ownership still needs a named owner.
              every = "0m";
              agentId = owner;
            };
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
    # The claude-cli runtime finds `claude` here, and Claude Code's native Bash
    # inherits it (tools.exec.pathPrepend does not reach it). Same order as the
    # login shell; systemd's default lacks the Nix and ~/.local/bin tools.
    Service.Environment = [
      "PATH=${home}/.nix-profile/bin:/nix/var/nix/profiles/default/bin:${home}/.local/bin:/usr/local/bin:/usr/bin:/bin"
    ];
  };

  home.packages = [ installSlackPlugin ];
  # Interactive `openclaw` sees the service's plugin registry (the wrapper defaults to hiding it).
  home.sessionVariables.OPENCLAW_DISABLE_PERSISTED_PLUGIN_REGISTRY = "0";

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
