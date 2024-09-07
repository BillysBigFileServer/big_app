{
  description = "A basic gomod2nix flake";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  inputs.flake-utils.url = "github:numtide/flake-utils";
  inputs.gomod2nix.url = "github:nix-community/gomod2nix";
  inputs.gomod2nix.inputs.nixpkgs.follows = "nixpkgs";
  inputs.gomod2nix.inputs.flake-utils.follows = "flake-utils";

  outputs = { self, nixpkgs, flake-utils, gomod2nix }:
    (flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
        buildInputs = with pkgs; [
          pkg-config
          mesa
          xorg.libX11
          xorg.libXinerama
          xorg.libXrandr
          xorg.libXcursor
          xorg.libXi
          mesa.dev
          mesa
          libglvnd.dev
          xorg.libXxf86vm
          wayland
          vulkan-tools
          vulkan-headers
          vulkan-loader
          vulkan-validation-layers
          libxkbcommon
          libxkbcommon.dev
        ];
        nativeBuildInputs = buildInputs;
        buildGoApp = gomod2nix.legacyPackages.${system}.buildGoApplication {
          inherit (gomod2nix.legacyPackages.${system})
          ;
          inherit buildInputs;
          inherit nativeBuildInputs;
        };

        # The current default sdk for macOS fails to compile go projects, so we use a newer one for now.
        # This has no effect on other platforms.
        callPackage =
          pkgs.darwin.apple_sdk_11_0.callPackage or pkgs.callPackage;
      in {
        packages.default = callPackage ./. {
          inherit (gomod2nix.legacyPackages.${system}) buildGoApp;
        };
        devShells.default = callPackage ./shell.nix {
          inherit (gomod2nix.legacyPackages.${system}) mkGoEnv gomod2nix;
        };
      }));
}
