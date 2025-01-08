{
  inputs = {
    flake-utils.url = "github:numtide/flake-utils";
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs =
    { nixpkgs, flake-utils, ... }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs { inherit system; };
      in
      with pkgs;
      {
        devShell = mkShell rec {
          nativeBuildInputs = [ ];
          buildInputs = [
            emmet-language-server
            prettierd
            typescript
            nodejs_22

            go
            gopls
            gcc # for cgo

            sqlite
          ];
          LD_LIBRARY_PATH = lib.makeLibraryPath buildInputs;
          CGO_ENABLED = 1;
        };
      }
    );
}
