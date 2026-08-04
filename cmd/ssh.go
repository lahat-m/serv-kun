package cmd

import "github.com/spf13/cobra"

var sshCmd = &cobra.Command{
	Use:   "ssh",
	Short: "Harden the OpenSSH server — disable root login and password authentication",
}
