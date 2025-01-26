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
        devShells.default =
          with pkgs;
          mkShell rec {
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
          backend = pkgs.buildGoModule {
            name = "backend";
            src = ./server;
            vendorHash = "sha256-QBWyMxEa/orGV5j8oFf6meY5pVYaOB3ym4GdNCEtWkU=";
          };

          esbuild-script = pkgs.buildGoModule {
            name = "esbuild-script";
            src = ./web/esbuild;
            vendorHash = "sha256-uQuzfDwae3XK7QyrrrcR9cb/q8lcX3/+fcRFkaJ5PT4=";
          };

          frontend = pkgs.stdenv.mkDerivation {
            name = "frontend";
            src = ./web;
            buildInputs = [
              pkgs.go
              self'.packages.esbuild-script
            ];
            buildPhase = ''
                # see https://github.com/NixOS/nix/issues/670
                export HOME=$(pwd)
                make ESBUILD_SCRIPT=${self'.packages.esbuild-script}/bin/esbuild
            '';
            installPhase = ''
                runHook preInstall
                export HOME=$(pwd)
                make install ESBUILD_SCRIPT=${self'.packages.esbuild-script}/bin/esbuild DEST_DIR=$out
                runHook postInstall
            '';
          };
        };

        process-compose."default" = {
          imports = [
            inputs.services-flake.processComposeModules.default
          ];

          settings.processes.backend = {
            command = pkgs.writeShellApplication {
              name = "backend";
              runtimeInputs = [
                self'.packages.backend
                pkgs.sqlite
              ];
              text = ''
                ${self'.packages.backend}/bin/server
              '';
            };
          };

          services.nginx."nginx1" = rec {
            enable = true;
            port = 80;
            dataDir = "./nginx_data";

            configFile = pkgs.writeText "nginx.conf" ''
              # pid ./data/nginx1/nginx/nginx.pid;
              pid ${dataDir}/nginx/nginx.pid;

              error_log stderr debug;
              daemon off;
              
              user nginx nginx;
              
              events {
                
              }
              
              http {
                access_log off;
                client_body_temp_path ${dataDir}/nginx/;
                proxy_temp_path ${dataDir}/nginx/;
                fastcgi_temp_path ${dataDir}/nginx/;
                scgi_temp_path ${dataDir}/nginx/;
                uwsgi_temp_path ${dataDir}/nginx/; 

                include ${pkgs.mailcap}/etc/nginx/mime.types;
                # include /nix/store/3cd6mykxz9s07v786aq7nxp1aipmnz67-mailcap-2.1.54/etc/nginx/mime.types;

                default_type  application/octet-stream;
                map $http_upgrade $connection_upgrade {
                  \'\' close;
                  default upgrade;
                }
              
                upstream websocket {
                  server 127.0.0.1:42069;
                }
              
                server {
                  listen 80;
                  server_name localhost;
              
                  root ${self'.packages.frontend}/var/www/filler;
              
                  location / {
              
                  }
              
                  location /ws/ {
                      proxy_pass http://websocket;
                      proxy_http_version 1.1;
                      proxy_set_header Upgrade $http_upgrade;
                      proxy_set_header Connection $connection_upgrade;
                      proxy_set_header Host $host;
                  }
                }
              
              }
            '';
          };
        };
      };
    };
}
