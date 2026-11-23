{ config, pkgs, inputs, username, homeDirectory, ... }:

{
  home.username = username;
  home.homeDirectory = homeDirectory;
  home.stateVersion = "25.05";

  # tmux clipboard tool consumed by common/tmux.nix
  _module.args.clipboardCommand = "clip.exe";

  home.packages = import ../common/packages.nix { inherit pkgs inputs; };

  home.file = {
    ".docker/cli-plugins/docker-buildx".source = "${pkgs.docker-buildx}/bin/docker-buildx";
    ".docker/cli-plugins/docker-compose".source = "${pkgs.docker-compose}/bin/docker-compose";
  };
  home.sessionVariables = {};
  # Let Home Manager install and manage itself.
  programs.home-manager.enable = true;

  programs.tmux.shell = "${config.home.profileDirectory}/bin/bash";
}
