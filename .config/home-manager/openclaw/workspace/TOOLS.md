# Environment

- Host: WSL2 Ubuntu 24.04. You run as the owner's Linux user, unsandboxed, with their credentials (gh, git, cloud CLIs).
- User tools are in `~/.nix-profile/bin` (Nix home-manager) and `~/.local/bin`. Missing command: `nix run nixpkgs#<command>`.
- Windows drives are under /mnt/c and similar; avoid touching them unless asked.
