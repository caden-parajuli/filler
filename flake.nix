{
  inputs = {
    flake-utils.url = "github:numtide/flake-utils";
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    process-compose-flake.url = "github:Platonic-Systems/process-compose-flake";
    services-flake.url = "github:juspay/services-flake";
  };

  outputs =
    inputs@{
    flake-parts,
    nixpkgs,
    flake-utils,
    ...
    }:
    flake-parts.lib.mkFlake { inherit inputs; } {
      imports = [
        inputs.process-compose-flake.flakeModule
      ];
      systems = [ "x86_64-linux" "aarch64-linux" "aarch64-darwin" "x86_64-darwin" ];

      perSystem = { config, self', inputs', pkgs, system, ... }: {
        devShells.default = with pkgs; mkShell rec {
          nativeBuildInputs = [
            emmet-language-server
            prettierd
            typescript

            gopls
            gcc # for cgo
            go
          ];

          buildInputs = [
            sqlite
          ];

          LD_LIBRARY_PATH = lib.makeLibraryPath buildInputs;
          CGO_ENABLED = 1;
        };
        packages = {
          # backend = pkgs.stdenv.mkDerivation {
          backend = pkgs.buildGoModule {
            name = "backend";
            src = ./server;
            vendorHash = "sha256-QBWyMxEa/orGV5j8oFf6meY5pVYaOB3ym4GdNCEtWkU=";
          };
        };
        process-compose."myservices" = {
          imports = [
            inputs.services-flake.processComposeModules.default
          ];
        };
      };
    };
}
