{
  description = "Tools for development environment";

  # Flake inputs
  inputs = {
    nixpkgs.url = "nixpkgs/nixos-23.11";
  };

  # Flake outputs
  outputs = { self, nixpkgs }:
    let
      # Systems supported
      allSystems = [
        "x86_64-linux"   # 64-bit Intel/AMD Linux
        "aarch64-linux"  # 64-bit ARM linux
        "x86_64-darwin"  # 64-bit Intel macOS
        "aarch64-darwin" # 64-bit ARM macOS
      ];

      # Helper to provide system-specific attributes
      forAllSystems = f: nixpkgs.lib.genAttrs allSystems (system: f {
        pkgs = import nixpkgs { inherit system; };
      });
    in
    {
      # Development environment output
      devShells = forAllSystems ({ pkgs }: {
        default = pkgs.mkShell {
          # The Nix packages provided in the environment
          packages = with pkgs; [
            go
            # https://goreleaser.com/
            goreleaser
            # - https://golangci-lint.run/
            golangci-lint
          ];
        };
      });
    };
}
