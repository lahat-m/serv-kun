package cmd

import "github.com/spf13/cobra"

var postgresCmd = &cobra.Command{
	Use:   "postgres",
	Short: "Install PostgreSQL, manage a deployment role/database, and back up/restore",
	// Without these a mistyped subcommand prints help and exits 0, so
	// `serv-kun postgres backupp` would "succeed" under set -e having backed
	// up nothing. Both are needed — see showHelp in root.go.
	Args: cobra.NoArgs,
	RunE: showHelp,
}
