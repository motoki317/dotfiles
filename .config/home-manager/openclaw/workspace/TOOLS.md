# Environment

- Host: WSL2 Ubuntu 24.04. You run as the owner's Linux user, unsandboxed, with their credentials (gh, git, cloud CLIs).
- User tools are in `~/.nix-profile/bin` (Nix home-manager) and `~/.local/bin`. Missing command: `nix run nixpkgs#<command>`.
- Windows drives are under /mnt/c and similar; avoid touching them unless asked.

# Slack message tool

- For `react` and other actions on the current conversation, set `target` to the bare conversation ID from the inbound event (`D0C60HPQKHV`, `C09NTH47C58`). Prefixed forms fail: OpenClaw checks that a `channel:` target is a channel, which a DM is not, and the Slack plugin rejects `user:` targets for reactions.
