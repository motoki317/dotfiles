{ pkgs }:
let
  # Personal Go helpers under ../scripts, built and installed on PATH. buildGoModule runs each
  # package's tests in its check phase, so a broken helper fails `home-manager switch` instead
  # of shipping. Stdlib-only, hence vendorHash = null. Editing one now takes effect on the next
  # switch — they previously self-built from ~/.claude on every call.
  goBin = name: pkgs.buildGoModule {
    pname = name;
    version = "0";
    src = ../scripts/${name};
    vendorHash = null;
  };
in
with pkgs; [
  _1password-cli
  age
  agent-browser
  awscli2
  buf
  cachix
  # claude-code
  (goBin "claude-statusline")
  (goBin "codex-consult")
  (goBin "session-transcript")
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
