with (import <nixpkgs> {});
mkShell {
  buildInputs = [
    go
    gosec
    go-tools
  ];
}
