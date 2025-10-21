{ config, pkgs, username, homeDirectory, ... }:

{
  home.username = username;
  home.homeDirectory = homeDirectory;
  home.stateVersion = "25.05";

  home.packages = import ../common/packages.nix { inherit pkgs; };

  home.file = {
    ".docker/cli-plugins/docker-buildx".source = "${pkgs.docker-buildx}/bin/docker-buildx";
    ".docker/cli-plugins/docker-compose".source = "${pkgs.docker-compose}/bin/docker-compose";
  };
  home.sessionVariables = {};
  # Let Home Manager install and manage itself.
  programs.home-manager.enable = true;

  programs.bash = {
    enable = true;
    # Load for both non-interactive and login shells
    bashrcExtra = ''
      # Load common
      if [ -f ~/.profile.common ]; then
        . ~/.profile.common
      fi

      # Load local
      if [ -f ~/.bashrc.local ]; then
        . ~/.bashrc.local
      fi

      # Initialize fnm (Fast Node Manager)
      eval "$(fnm env --use-on-cd --shell bash)"
    '';
    # Load default color profiles and so for login shells
    initExtra = builtins.readFile ./.bashrc;
  };
}
