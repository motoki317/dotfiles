{ config, pkgs, lib, username, homeDirectory, inputs, ... }:
let
  # gh keeps one active account per host, so it ignores the attmcojp/ URL split
  # that url.insteadOf already applies to SSH keys. GH_TOKEN overrides per call.
  # Sourced by both zsh and bash, so it sticks to syntax both accept.
  ghAccountSwitch = ''
    gh() {
      [[ $1 == auth ]] && { command gh "$@"; return; }
      local account
      case "$(command git remote get-url origin 2>/dev/null)" in
        *[:/]attmcojp/*) account=toki-attm ;;
        *github*)        account=motoki317 ;;
        *) command gh "$@"; return ;;
      esac
      GH_TOKEN=$(command gh auth token --user "$account") command gh "$@"
    }
  '';
in
{
  nixpkgs.config.allowUnfreePredicate = pkg: builtins.elem (lib.getName pkg) [
    "1password-cli"
    "claude-code"
    "ngrok"
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

  # ha - git worktree manager (sourced as shell functions)
  programs.zsh.initContent = ''
    source ${inputs.ha}/ha.sh
    compdef _ha ha

    ${ghAccountSwitch}
  '';
  programs.bash.initExtra = ''
    source ${inputs.ha}/ha.sh

    ${ghAccountSwitch}
  '';
}
