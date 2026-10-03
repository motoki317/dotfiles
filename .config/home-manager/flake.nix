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
    # Declarative OpenClaw (first-party): package overlay + Home Manager module. WSL host only.
    # Pinned so a blanket `nix flake update` cannot move the gateway away from the
    # imperatively installed Slack plugin; upgrade per openclaw/README.md.
    nix-openclaw = {
      url = "github:openclaw/nix-openclaw/f62d33f760bcbdbc6a52ac589eae22bf99201f90";
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
