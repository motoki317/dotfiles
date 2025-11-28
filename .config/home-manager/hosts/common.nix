{ config, pkgs, username, homeDirectory, ... }:
{
  programs.starship = {
    enable = true;
    settings = {
      gcloud.disabled = true;
      git_status = {
        modified = "~";
      };
    };
  };
}
