{
  description = "A basic gomod2nix flake";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    gomod2nix.url = "github:tweag/gomod2nix";

    # dev
    nix-filter.url = "github:ilkecan/nix-filter/add-_assertPathIsDirectory";
    devshell.url = "github:numtide/devshell";
    flake-utils.url = "github:numtide/flake-utils";
    flake-compat = {
      url = "github:edolstra/flake-compat";
      flake = false;
    };
  };

  outputs = {
    self,
    nix-filter,
    flake-utils,
    ...
  } @ inputs: (
    flake-utils.lib.eachDefaultSystem
    (system: let
      inherit (pkgs) lib;
      pkgs = import inputs.nixpkgs {
        inherit system;
        overlays = [
          inputs.gomod2nix.overlays.default
          inputs.devshell.overlay
        ];
      };
    in {
      packages.yamlfmt = pkgs.buildGoApplication {
        pname = "yamlfmt";
        version = "0.1.0";
        src = with nix-filter.lib;
          filter {
            root = ./.;
            include = [
              "go.mod"
              "go.sum"
              "gomod2nix.toml"
              "yamlfmt.go"
              (inDirectory "cmd")
              (inDirectory "command")
              (inDirectory "engine")
              (inDirectory "formatters")
            ];
          };
        modules = ./gomod2nix.toml;
      };
      packages.default = self.packages.${system}.yamlfmt;

      apps.yamlfmt = flake-utils.lib.mkApp {
        drv = self.packages.${system}.yamlfmt;
        name = "yamlfmt";
      };
      apps.default = self.apps.${system}.yamlfmt;

      devShells.default = pkgs.devshell.mkShell {
        packages = with pkgs; [
          (mkGoEnv {pwd = ./.;})
          gomod2nix
          golangci-lint
          alejandra
          taplo-cli
          treefmt
          self.packages.${system}.yamlfmt
        ];
      };
    })
  );
}
