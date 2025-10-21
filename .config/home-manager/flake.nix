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
        modules = [ ./hosts/wsl.nix ];
        extraSpecialArgs = {
          username = "moto";
          homeDirectory = "/home/moto";
        };
      };
      # gmom M3 Pro
      homeConfigurations."usr0200788" = home-manager.lib.homeManagerConfiguration {
        pkgs = nixpkgs.legacyPackages.aarch64-darwin;
        modules = [ ./hosts/macos.nix ];
        extraSpecialArgs = {
          username = "usr0200788";
          homeDirectory = "/Users/usr0200788";
        };
      };
    };
}
