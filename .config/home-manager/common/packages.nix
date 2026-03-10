{ pkgs }:
let
  version = "0.27.2";

  platforms = {
    "aarch64-darwin" = "aarch64-apple-darwin";
    "x86_64-linux" = "x86_64-unknown-linux-musl";
  };

  hashes = {
    "aarch64-darwin" = "sha256-p2XccLukutDpFu0v2ssDz9nZBqGulPxrXuyXFR3MGhc=";
    "x86_64-linux" = "sha256-hi+4QwyeNsBOUhQQ8sRZVF/BoGkLpxaNMzXugETfkIQ=";
  };

  rtk = pkgs.stdenv.mkDerivation {
    pname = "rtk";
    inherit version;

    src = pkgs.fetchurl {
      url = "https://github.com/rtk-ai/rtk/releases/download/v${version}/rtk-${platforms.${pkgs.stdenv.hostPlatform.system}}.tar.gz";
      hash = hashes.${pkgs.stdenv.hostPlatform.system};
    };

    sourceRoot = ".";

    installPhase = ''
      install -Dm755 rtk $out/bin/rtk
    '';

    meta = {
      description = "Token-optimized CLI proxy for Claude Code";
      homepage = "https://github.com/rtk-ai/rtk";
      license = pkgs.lib.licenses.mit;
      platforms = builtins.attrNames platforms;
    };
  };
in
with pkgs; [
  _1password-cli
  age
  awscli2
  buf
  cachix
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
  nodejs_24
  omekasy
  opentofu
  postgresql
  ripgrep
  rsync
  rtk
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
