{
  description = "My home-manager configurations for multiple devices";

  inputs = {
    # Specify the source of Home Manager and Nixpkgs.
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    home-manager = {
      url = "github:nix-community/home-manager";
      inputs.nixpkgs.follows = "nixpkgs";
    };
    ha = {
      url = "github:kawarimidoll/ha";
      flake = false;
    };
    # Declarative OpenClaw, from a fork while upstream nix-openclaw stays on 2026.9.5.
    # Pinned so a blanket `nix flake update` cannot move the gateway away from the
    # imperatively installed Slack plugin; upgrade per openclaw/README.md.
    nix-openclaw = {
      url = "github:motoki317/nix-openclaw/16ae435a1bc4c46765598db288944c4c9f05352c";
      inputs.nixpkgs.follows = "nixpkgs";
      inputs.home-manager.follows = "home-manager";
    };
  };

  outputs =
    { nixpkgs, home-manager, ... }@inputs:
    {
      # Windows Desktop (WSL2 Ubuntu)
      homeConfigurations."moto" = home-manager.lib.homeManagerConfiguration {
        pkgs = nixpkgs.legacyPackages.x86_64-linux;
        modules = [
          ./hosts/wsl.nix
          ./hosts/common.nix
          ./openclaw
        ];
        extraSpecialArgs = {
          inherit inputs;
          username = "moto";
          homeDirectory = "/home/moto";
        };
      };
      # Work MacBook Pro
      homeConfigurations."toki" = home-manager.lib.homeManagerConfiguration {
        pkgs = nixpkgs.legacyPackages.aarch64-darwin;
        modules = [
          ./hosts/macos.nix
          ./hosts/common.nix
        ];
        extraSpecialArgs = {
          inherit inputs;
          username = "toki";
          homeDirectory = "/Users/toki";
        };
      };
    };
}
