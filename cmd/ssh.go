package cmd

import "github.com/spf13/cobra"

var sshCmd = &cobra.Command{
	Use:   "ssh",
	Short: "Harden and audit the OpenSSH server",
	Args:  cobra.NoArgs, // reject mistyped subcommands — see showHelp in root.go
	RunE:  showHelp,
}
