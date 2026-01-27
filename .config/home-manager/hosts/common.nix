{ config, pkgs, lib, username, homeDirectory, ... }:
{
  nixpkgs.config.allowUnfreePredicate = pkg: builtins.elem (lib.getName pkg) [
    "1password-cli"
  ];

  imports = [
    ../common/tmux.nix
  ];

  programs.starship = {
    enable = true;
    settings = {
      gcloud.disabled = true;
      directory = {
        truncation_length = 0;
      };
      git_status = {
        modified = "~";
      };
    };
  };

  programs.fzf = {
    enable = true;
  };

  programs.zoxide = {
    enable = true;
    enableZshIntegration = true;
    enableBashIntegration = true;
  };

  programs.direnv = {
    enable = true;
    enableZshIntegration = true;
    enableBashIntegration = true;
    nix-direnv.enable = true;
    silent = true;
  };
}
