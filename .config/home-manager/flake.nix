{
  description = "My home-manager configurations for multiple devices";

  inputs = {
    # Specify the source of Home Manager and Nixpkgs.
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    home-manager = {
      url = "github:nix-community/home-manager";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs =
    { nixpkgs, home-manager, ... }:
    {
      # Windows Desktop (WSL2 Ubuntu)
      homeConfigurations."moto" = home-manager.lib.homeManagerConfiguration {
        pkgs = nixpkgs.legacyPackages.x86_64-linux;
        modules = [
          ./hosts/wsl.nix
          ./hosts/common.nix
        ];
        extraSpecialArgs = {
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
          username = "toki";
          homeDirectory = "/Users/toki";
        };
      };
    };
}
