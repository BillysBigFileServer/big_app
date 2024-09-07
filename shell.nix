{ pkgs ? (let inherit (builtins) fetchTree fromJSON readFile;
in import (fetchTree (fromJSON (readFile ./flake.lock)).nodes.nixpkgs.locked) {
  overlays = [ (import ./overlay.nix) ];
}), gomod2nix ? pkgs.gomod2nix, mkGoEnv ? pkgs.mkGoEnv }:

pkgs.mkShell {
  NIX_PATH = "nixpkgs=${builtins.toString pkgs.path}";
  nativeBuildInputs = with pkgs; [
    nixpkgs-fmt
    golangci-lint
    gomod2nix
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
    gcc
    gotools
    fyne

    (mkGoEnv { pwd = ./.; })
  ];
}
