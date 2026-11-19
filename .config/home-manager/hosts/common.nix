{ config, pkgs, lib, username, homeDirectory, inputs, ... }:
let
  # gh keeps one active account per host, so it ignores the attmcojp/ URL split
  # that url.insteadOf already applies to SSH keys. GH_TOKEN overrides per call.
  # A command on PATH, not a shell function, so that it also reaches shells that
  # never read the rc files below: scripts, agents' shells, and zsh on WSL, where
  # home-manager manages only bash.
  ghAccountSwitch = lib.hiPrio (pkgs.writeShellScriptBin "gh" ''
    gh=${pkgs.gh}/bin/gh
    [[ $1 == auth ]] && exec "$gh" "$@"
    case "$(git remote get-url origin 2>/dev/null)" in
      *[:/]attmcojp/*) account=toki-attm ;;
      *github*)        account=motoki317 ;;
      *) exec "$gh" "$@" ;;
    esac
    # Stop rather than fall back to the active account, which may be the other one.
    GH_TOKEN=$("$gh" auth token --user "$account") || exit
    export GH_TOKEN
    exec "$gh" "$@"
  '');
in
{
  nixpkgs.config.allowUnfreePredicate = pkg: builtins.elem (lib.getName pkg) [
    "1password-cli"
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

  home.packages = [ ghAccountSwitch ];

  # ha - git worktree manager (sourced as shell functions)
  programs.zsh.initContent = ''
    source ${inputs.ha}/ha.sh
    compdef _ha ha
  '';
  programs.bash.initExtra = ''
    source ${inputs.ha}/ha.sh
  '';
}
