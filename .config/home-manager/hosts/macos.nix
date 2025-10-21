{ config, pkgs, username, homeDirectory, ... }:

{
  home.username = username;
  home.homeDirectory = homeDirectory;
  home.stateVersion = "25.05";

  home.packages = import ../common/packages.nix { inherit pkgs; };

  home.file = {};
  home.sessionVariables = {};
  # Let Home Manager install and manage itself.
  programs.home-manager.enable = true;

  programs.zsh = {
    enable = true;
    initContent = ''
      # Load common
      if [ -f ~/.profile.common ]; then
        . ~/.profile.common
      fi

      # Load local
      if [ -f ~/.zshenv.local ]; then
        . ~/.zshenv.local
      fi

      # pure (zsh theme)
      fpath+=("$HOME/.zsh/pure")
      autoload -U promptinit
      promptinit
      prompt pure

      # Initialize fnm (Fast Node Manager)
      eval "$(fnm env --use-on-cd --shell zsh)"
    '';
  };
}
