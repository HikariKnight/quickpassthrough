{
  description = "Quick ";

  inputs = {
    nixpkgs.url = "nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    gomod2nix.url = "github:nix-community/gomod2nix";
    gomod2nix.inputs.nixpkgs.follows = "nixpkgs";
    gomod2nix.inputs.flake-utils.follows = "flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils, gomod2nix }:
  flake-utils.lib.eachDefaultSystem (system: let
      pkgs = import nixpkgs { inherit system; }; 
    in {
      packages.default = pkgs.callPackage ./. {
        inherit (gomod2nix.legacyPackages.${system}) buildGoApplication;
      };

      packages.updater = pkgs.writeShellScriptBin "updater" ''
	set -euo pipefail
	echo "WARNING: This package is for updating this nix package! You should also update default.nix release variable if you wanna update this too"
	${pkgs.nix}/bin/nix run github:nix-community/gomod2nix
      '';

      devShells.default = pkgs.mkShell {
        packages = with pkgs; [ nil go gopls ]; 
      };
    }
  );
}
