package cmd

import "github.com/spf13/cobra"

var postgresCmd = &cobra.Command{
	Use:   "postgres",
	Short: "Install PostgreSQL, manage a deployment role/database, and back up/restore",
}
