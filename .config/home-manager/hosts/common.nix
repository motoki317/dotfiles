{ config, pkgs, lib, username, homeDirectory, inputs, ... }:
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

    # gh keeps one active account per host, so it ignores the attmcojp/ URL split
    # that url.insteadOf already applies to SSH keys. GH_TOKEN overrides per call.
    gh() {
      [[ $1 == auth ]] && { command gh "$@"; return }
      local account
      case "$(command git remote get-url origin 2>/dev/null)" in
        *[:/]attmcojp/*) account=toki-attm ;;
        *github*)        account=motoki317 ;;
        *) command gh "$@"; return ;;
      esac
      GH_TOKEN=$(command gh auth token --user $account) command gh "$@"
    }
  '';
  programs.bash.initExtra = ''
    source ${inputs.ha}/ha.sh
  '';
}
