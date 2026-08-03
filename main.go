package main

import "github.com/lahat-m/serv-kun/cmd"

// version is set at build time via -ldflags "-X main.version=...";
// see .github/workflows/release.yml.
var version = "dev"

func main() {
	cmd.Execute(version)
}
