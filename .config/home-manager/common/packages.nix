{ pkgs, inputs }:
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
  # HTTP client for AI agents (https://github.com/yusukebe/ax), preferred over curl in agent
  # sessions. Not in nixpkgs and not on npm; upstream ships only prebuilt Bun binaries, so
  # install the release binary. Bump version and both hashes together (see checksums.txt on
  # the release).
  ax =
    let
      version = "0.1.10";
      bin = {
        x86_64-linux = {
          suffix = "linux-x64";
          hash = "sha256-dJWaZpyB58SmXv6SRMxyM8d1LSNSnM7wBHkFbjM04NA=";
        };
        aarch64-darwin = {
          suffix = "darwin-arm64";
          hash = "sha256-FCCrnigkNiCtCHsuzpskyZaxfJGGUdltOMu2A+WKvOk=";
        };
      }.${pkgs.system};
    in
    pkgs.stdenvNoCC.mkDerivation {
      pname = "ax";
      inherit version;
      src = pkgs.fetchurl {
        url = "https://github.com/yusukebe/ax/releases/download/v${version}/ax-${bin.suffix}";
        inherit (bin) hash;
      };
      dontUnpack = true;
      installPhase = ''
        install -Dm755 $src $out/bin/ax
      '';
    };
in
with pkgs; [
  _1password-cli
  age
  agent-browser
  awscli2
  ax
  buf
  cachix
  inputs.ccusage.packages.${pkgs.system}.default # statusline session cost (see flake.nix)
  # claude-code
  (goBin "claude-statusline")
  (goBin "codex-run")
  (goBin "session-transcript")
  (goBin "sync-claude-skills")
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
