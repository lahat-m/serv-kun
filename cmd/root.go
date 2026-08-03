package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "serv-kun",
	Short: "Guided setup wizards for common Linux server provisioning tasks",
	Long: `serv-kun is a collection of guided setup wizards for provisioning a
Linux server — no config file editing required.

Commands:
  nginx setup       Install nginx, write a reverse-proxy config, optionally add SSL
  nginx ssl         Add SSL to an existing nginx config via Let's Encrypt
  user              Create a Linux user with sudo privileges
  postgres setup    Install PostgreSQL and create a deployment role and database
  postgres backup   Back up a PostgreSQL database to a file
  postgres restore  Restore a PostgreSQL database from a backup file
  node              Install Node.js via nvm and pick a default version
  springboot init   Scaffold docker-compose.prod.yml, nginx config, backup script, and CI/CD
  springboot setup  Bootstrap docker, .env.prod, JWT keys, TLS, and bring the stack up
  springboot deploy Pull the latest app image and redeploy, with automatic rollback
  springboot status Show container status and check the app's health endpoint
  springboot logs   Tail logs for a compose service`,
}

// Execute runs the CLI. version is reported via --version (set at build
// time in main.go's version var, "dev" for a local `go build`/`go run`).
func Execute(version string) {
	rootCmd.Version = version
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	nginxCmd.AddCommand(nginxSetupCmd, nginxSSLCmd)
	postgresCmd.AddCommand(postgresSetupCmd, postgresBackupCmd, postgresRestoreCmd)
	springbootCmd.AddCommand(springbootInitCmd, springbootSetupCmd, springbootDeployCmd, springbootStatusCmd, springbootLogsCmd)
	rootCmd.AddCommand(nginxCmd, userCmd, postgresCmd, nodeCmd, springbootCmd)
}

// requireRoot exits with an actionable error if the current process isn't root.
func requireRoot(cmd *cobra.Command) {
	if os.Getuid() != 0 {
		fmt.Fprintf(os.Stderr, "error: must be run as root — try: sudo %s\n", cmd.CommandPath())
		os.Exit(1)
	}
}
