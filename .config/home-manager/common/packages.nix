{ pkgs }:
with pkgs; [
  _1password-cli
  age
  agent-browser
  awscli2
  buf
  cachix
  # claude-code
  # Short PATH alias for the self-building Codex review helper. It execs the source
  # in ~/.claude (synced across machines) rather than a nix-built binary, so edits
  # take effect on the next call without a home-manager rebuild.
  (writeShellScriptBin "codex-consult" ''
    exec sh "$HOME/.claude/skills/codex-advisor/scripts/codex-consult.go" "$@"
  '')
  curl
  dive
  docker-buildx
  docker-client
  docker-compose
  dyff
  fortune
  fzf
  gh
  git
  go
  go-task
  gonzo
  google-cloud-sdk
  hackgen-nf-font
  htop
  jq
  just
  k6
  k9s
  kubeconform
  kubectl
  kubernetes-helm
  kustomize
  kustomize-sops
  lazygit
  lua
  luajitPackages.luarocks
  neovim
  ngrok
  nodejs_24
  omekasy
  opentofu
  postgresql
  ripgrep
  rsync
  skaffold
  sops
  ssm-session-manager-plugin
  textlint
  tree
  tree-sitter
  treefmt
  unzip
  uv
  vegeta
  vim
  wget
  witr
  yq-go
  zip
]
