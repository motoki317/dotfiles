{ config, pkgs, username, homeDirectory, ... }:

{
  home.username = username;
  home.homeDirectory = homeDirectory;
  home.stateVersion = "25.05";

  # tmux clipboard tool consumed by common/tmux.nix
  _module.args.clipboardCommand = "pbcopy";

  home.packages = import ../common/packages.nix { inherit pkgs; }
   ++
   (with pkgs; [
    lima
   ]);

  home.file = {
    # Not needed with Docker Desktop
    # ".docker/cli-plugins/docker-buildx".source = "${pkgs.docker-buildx}/bin/docker-buildx";
    # ".docker/cli-plugins/docker-compose".source = "${pkgs.docker-compose}/bin/docker-compose";
  };
  home.sessionVariables = {};
  # Let Home Manager install and manage itself.
  programs.home-manager.enable = true;

  programs.zsh = {
    enable = true;
    initContent = ''
      # Re-add Homebrew to PATH. home-manager owns the zsh dotfiles, dropping the
      # `brew shellenv` line the Homebrew installer wrote to ~/.zprofile on Apple
      # Silicon. It lives in initContent (~/.zshrc), NOT profileExtra (~/.zprofile),
      # because tmux here runs non-login shells (common/tmux.nix: default-command
      # $SHELL) that never source ~/.zprofile; ~/.zshrc is sourced by every
      # interactive shell. Prefix is fixed at /opt/homebrew on Apple Silicon.
      eval "$(/opt/homebrew/bin/brew shellenv)"

      # Load common
      if [ -f ~/.profile.common ]; then
        . ~/.profile.common
      fi

      # Load local
      if [ -f ~/.zshenv.local ]; then
        . ~/.zshenv.local
      fi
    '';
  };
}
