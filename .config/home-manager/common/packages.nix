{ pkgs, inputs }:
let
  # Personal Go helpers under ../scripts, built onto PATH. buildGoModule runs each package's
  # tests at build, so a broken helper fails the switch; stdlib-only, hence vendorHash = null.
  goBin = name: pkgs.buildGoModule {
    pname = name;
    version = "0";
    src = ../scripts/${name};
    vendorHash = null;
  };
  # headless browser automation CLI for AI agents (https://github.com/vercel-labs/agent-browser).
  # In nixpkgs but lagging (0.27 vs 0.31), so pin the latest prebuilt binary like ax/hunk; this
  # let binding shadows pkgs.agent-browser in the list. Bump: ./fetch-prebuilt-hashes.sh agent-browser <ver>.
  agent-browser =
    let
      version = "0.31.1";
      bin = {
        x86_64-linux = {
          suffix = "linux-x64";
          hash = "sha256-csE7z9L9axiDJb3SPGRtBsppoalkqc2qs35P+PR6pcY=";
        };
        aarch64-darwin = {
          suffix = "darwin-arm64";
          hash = "sha256-/XrNF7MHH/f3WgPB7NMFAZWdnC0GO9qgWttvd6vyp78=";
        };
      }.${pkgs.system};
    in
    pkgs.stdenvNoCC.mkDerivation {
      pname = "agent-browser";
      inherit version;
      src = pkgs.fetchurl {
        url = "https://github.com/vercel-labs/agent-browser/releases/download/v${version}/agent-browser-${bin.suffix}";
        inherit (bin) hash;
      };
      dontUnpack = true;
      installPhase = ''
        install -Dm755 $src $out/bin/agent-browser
      '';
    };
  # terminal UI for viewing coding-agent (Claude Code / Codex) logs and costs
  # (https://github.com/motoki317/agtlog). Prebuilt release tarball, not upstream's flake, for
  # the same x86_64-darwin-drop reason as hunk. Bump: ./fetch-prebuilt-hashes.sh agtlog <ver>.
  agtlog =
    let
      version = "0.2.4";
      bin = {
        x86_64-linux = {
          suffix = "linux-amd64";
          hash = "sha256-gKcw6LM9SuMSnT8gl7lc3tWzy385lBN1gYX9Z4JGZJI=";
        };
        aarch64-darwin = {
          suffix = "darwin-arm64";
          hash = "sha256-fm+k9yswMIPzExBeSJkusdi3yJoFGu7UVUJftljvZRs=";
        };
      }.${pkgs.system};
    in
    pkgs.stdenvNoCC.mkDerivation {
      pname = "agtlog";
      inherit version;
      src = pkgs.fetchurl {
        url = "https://github.com/motoki317/agtlog/releases/download/v${version}/agtlog-${version}-${bin.suffix}.tar.gz";
        inherit (bin) hash;
      };
      # release tarball extracts loose files (no wrapping dir), so unpackPhase needs "." as root.
      sourceRoot = ".";
      installPhase = ''
        runHook preInstall
        install -Dm755 agtlog $out/bin/agtlog
        runHook postInstall
      '';
    };
  # HTTP client for AI agents (https://github.com/yusukebe/ax), preferred over curl. Not
  # packaged; upstream ships only prebuilt binaries. Bump: ./fetch-prebuilt-hashes.sh ax <ver>.
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
  # review-first terminal diff viewer (https://github.com/modem-dev/hunk), driven from nvim.
  # Prebuilt, not upstream's flake: that builds via bun2nix, whose flake-parts eval trips
  # nixpkgs >=26.11's dropped-x86_64-darwin error. Bump: ./fetch-prebuilt-hashes.sh hunk <ver>.
  hunk =
    let
      version = "0.17.0";
      bin = {
        x86_64-linux = {
          suffix = "linux-x64";
          hash = "sha256-DGJvemaHqYJjBOod9pbaXUnt+EJx7M31f//1g0KJ4OI=";
        };
        aarch64-darwin = {
          suffix = "darwin-arm64";
          hash = "sha256-cAIhZppRt4yDWYW0nKZ+Wku6r9tDSmygP50mL3qCaT4=";
        };
      }.${pkgs.system};
    in
    pkgs.stdenvNoCC.mkDerivation {
      pname = "hunk";
      inherit version;
      src = pkgs.fetchurl {
        url = "https://github.com/modem-dev/hunk/releases/download/v${version}/hunkdiff-${bin.suffix}.tar.gz";
        inherit (bin) hash;
      };
      installPhase = ''
        runHook preInstall
        install -Dm755 hunk $out/bin/hunk
        runHook postInstall
      '';
    };
in
with pkgs; [
  _1password-cli
  age
  agent-browser
  agtlog
  awscli2
  ax
  buf
  cachix
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
  hunk
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
