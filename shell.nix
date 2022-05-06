let
  pkgs = import <nixpkgs> { };
in pkgs.mkShell {
  buildInputs = with pkgs; [
    go_1_18
    gnumake
    sqlite

    # mattn/go-sqlite3 requires gcc to be installed, at least when running `go get`
    gcc

    # pprof visualizer calls graphviz
    graphviz
  ];

  NIX_LD_LIBRARY_PATH = pkgs.lib.makeLibraryPath [
    pkgs.stdenv.cc.cc
  ];

  NIX_LD = builtins.readFile "${pkgs.stdenv.cc}/nix-support/dynamic-linker";
}
