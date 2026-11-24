{ config, pkgs, inputs, username, homeDirectory, ... }:

let
  # Gmail and Calendar CLI (https://github.com/openclaw/gogcli) for the agent on this host.
  # Its file keyring takes the password only from the environment, and the agent's shells and
  # systemd timers do not source shell rc files, so the wrapper loads it from a private file.
  gog = pkgs.writeShellScriptBin "gog" ''
    if [ -r "$HOME/.secrets/gog/env" ]; then
      set -a
      . "$HOME/.secrets/gog/env"
      set +a
    fi
    exec ${pkgs.gogcli}/bin/gog "$@"
  '';
in
{
  home.username = username;
  home.homeDirectory = homeDirectory;
  home.stateVersion = "25.05";

  # tmux clipboard tool consumed by common/tmux.nix
  _module.args.clipboardCommand = "clip.exe";

  home.packages = import ../common/packages.nix { inherit pkgs inputs; } ++ [ gog ];

  home.file = {
    ".docker/cli-plugins/docker-buildx".source = "${pkgs.docker-buildx}/bin/docker-buildx";
    ".docker/cli-plugins/docker-compose".source = "${pkgs.docker-compose}/bin/docker-compose";
  };
  home.sessionVariables = {};
  # Let Home Manager install and manage itself.
  programs.home-manager.enable = true;

  programs.tmux.shell = "${config.home.profileDirectory}/bin/bash";
}
