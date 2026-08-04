package cmd

import "github.com/spf13/cobra"

var nginxCmd = &cobra.Command{
	Use:   "nginx",
	Short: "Set up nginx as a reverse proxy and configure SSL via Let's Encrypt",
	Args:  cobra.NoArgs, // reject mistyped subcommands — see showHelp in root.go
	RunE:  showHelp,
}
