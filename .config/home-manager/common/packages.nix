{ pkgs }:
with pkgs; [
  _1password-cli
  age
  agent-browser
  awscli2
  buf
  cachix
  # claude-code
  # Short PATH aliases for the self-building codex-advisor helpers. Each execs its source
  # in ~/.claude (synced across machines) rather than a nix-built binary, so edits to the
  # source take effect on the next call; only adding a new alias needs a rebuild.
  (writeShellScriptBin "codex-consult" ''
    exec sh "$HOME/.claude/skills/codex-advisor/scripts/codex-consult.go" "$@"
  '')
  (writeShellScriptBin "session-transcript" ''
    exec sh "$HOME/.claude/skills/codex-advisor/scripts/session-transcript.go" "$@"
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
