# home-manager

Standalone [home-manager](https://github.com/nix-community/home-manager) flake for toki's machines — deliberately **not** integrated into nix-darwin.

- Hosts: `toki` (macOS, `hosts/macos.nix` + `hosts/common.nix`) and `moto` (WSL, `hosts/wsl.nix`).
- Apply: `home-manager switch --flake ~/.config/home-manager#toki`.
- nix-darwin lives separately at `/etc/nix-darwin` (config attr `darwinConfigurations."toki"`; Nix via the open-source DetSys installer): it owns `/etc/zsh*`, Nix itself, and the Linux builder, and does not reference home-manager.

## Ownership

- home-manager owns `~/.zshenv`, `~/.zprofile`, `~/.zshrc` as read-only nix-store symlinks — hand-editing fails or gets wiped; change `macos.nix`/`common.nix` instead.
- Non-declarative escape hatches: `~/.profile.common` (a real hand-maintained file, sourced from zsh `initContent`) and optional `~/.zshenv.local`.

## PATH mechanics (macOS)

Homebrew (`/opt/homebrew`, Apple Silicon) enters PATH declaratively via `programs.zsh.initContent` — `eval "$(/opt/homebrew/bin/brew shellenv)"` in `macos.nix`. It must live in `initContent` (`~/.zshrc`), **not** `profileExtra` (`~/.zprofile`): tmux here runs non-login shells (`common/tmux.nix`: `default-command $SHELL`), which never source `~/.zprofile`. Non-login is intentional — macOS `/etc/zprofile` runs `path_helper`, which reorders PATH and shadows nix binaries. So anything that must be on interactive PATH belongs in `initContent`; macOS-only settings go in `macos.nix`, not `common.nix`. (The brew line was added 2026-06-23 after home-manager's takeover of `~/.zprofile` dropped the original one and hid all brew tools.)

## Gotchas

### nixpkgs ≥ 26.11 dropped x86_64-darwin

`nixos-unstable` at 26.11 dropped x86_64-darwin: evaluating `pkgs` for that platform throws. Since this flake tracks nixos-unstable, any flake input whose transitive evaluation instantiates x86_64-darwin fails `home-manager switch` on the aarch64-darwin host — even though no host targets that platform. flake-parts flakes are the trap: they enumerate every system in their `systems` input and can strictly evaluate per-system outputs (`formatter`/`checks`) regardless of which attr you request (hit adding `modem-dev/hunk`, whose `bun2nix` input lists x86_64-darwin).

**Rule:** for a tool not in nixpkgs, install its prebuilt release binary instead of consuming its build-from-source flake — the `ax`/`hunk` pattern in `common/packages.nix` (fetchurl the per-`pkgs.system` asset, `stdenvNoCC.mkDerivation`, install the binary). Works on both hosts and sidesteps flake-parts eval entirely.

### Apple M4: nix-darwin's `nix.linux-builder` (qemu) crash-loops

On Apple Silicon M4 (confirmed on M4 Max / `Mac16,6`), nix-darwin's `nix.linux-builder.enable` does not work: the qemu darwin-builder VM aborts at vCPU init with `hvf_arch_init_vcpu: assertion failed (HV_SYS_REG_SMCR_EL1 == …)` — a qemu/Hypervisor.framework mismatch on the M4's SME registers (`SMCR_EL1`, qemu 11.0.1; `accel=hvf:tcg` aborts rather than falling back). The launchd daemon `org.nixos.linux-builder` then crash-loops (`last exit code = 134`, SIGABRT). The hand-rolled nixpkgs `darwin-builder` (qemu, ssh `localhost:31022`) has the identical bug.

**Fix:** use an Apple Virtualization.framework (vfkit/vz) builder — Apple's own hypervisor handles M4/SME natively. Options: `virby-nix-darwin` (`github:quinneden/virby-nix-darwin`, turnkey vfkit module with on-demand + Rosetta) or a Lima/vz builder wired in via `nix.buildMachines`.

Diagnose a dead builder with `sudo launchctl print system/org.nixos.linux-builder` (check `last exit code` / `runs`), then run the start script in the foreground to capture qemu's stderr.
