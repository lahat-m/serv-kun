package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2/terminal"
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
  ci                Grant the deploy user passwordless sudo for the commands CI runs
  ssh harden        Disable root login and password auth, with a lockout safety check
  ssh audit         Report SSH login attempts from the logs — who tried, who got in
  firewall          Default-deny inbound firewall allowing SSH, HTTP and HTTPS
  postgres setup    Install PostgreSQL and create a deployment role and database
  postgres backup   Back up a PostgreSQL database to a file
  postgres restore  Restore a PostgreSQL database from a backup file
  node              Install Node.js via nvm and pick a default version
  springboot init   Scaffold docker-compose.prod.yml, nginx config, backup script, and CI/CD
  springboot setup  Bootstrap docker, .env.prod, JWT keys, TLS, and bring the stack up
  springboot deploy Pull the latest app image and redeploy, with automatic rollback
  springboot clean  Delete the stack's containers, images, and build cache, keeping volumes
  springboot status Show container status and check the app's health endpoint
  springboot logs   Tail logs for a compose service`,
}

// Execute runs the CLI. version is reported via --version (set at build
// time in main.go's version var, "dev" for a local `go build`/`go run`).
func Execute(version string) {
	rootCmd.Version = version

	// Cobra prints the error itself; leaving that on as well as printing it
	// here reported every failure twice.
	rootCmd.SilenceErrors = true

	// PersistentPreRun happens after flags and arguments are validated but
	// before RunE, so a bad flag still gets a usage dump while a runtime
	// failure — a failed useradd, an interrupted prompt — does not. A wall of
	// flag documentation is no help when the problem isn't a flag.
	rootCmd.PersistentPreRun = func(cmd *cobra.Command, _ []string) {
		cmd.SilenceUsage = true
	}

	err := rootCmd.Execute()
	if err == nil {
		return
	}

	// Ctrl-C is a choice, not a failure: say so plainly and use the shell's
	// conventional 128+SIGINT status.
	if errors.Is(err, terminal.InterruptErr) {
		fmt.Fprintln(os.Stderr, "\nAborted.")
		os.Exit(130)
	}

	fmt.Fprintln(os.Stderr, "Error:", err)
	os.Exit(1)
}

// parentCommands exist only to group subcommands. Each must set both
// Args: cobra.NoArgs and RunE: showHelp — see showHelp for why one without
// the other silently does nothing. root_test.go enforces it.
func parentCommands() []*cobra.Command {
	return []*cobra.Command{nginxCmd, sshCmd, postgresCmd, springbootCmd}
}

// showHelp is the RunE for grouping commands, printing usage for a bare
// `serv-kun <parent>`.
//
// It also exists to make the parent *runnable*, which is load-bearing: cobra
// checks Runnable() before it validates arguments, so a parent with no Run at
// all returns "print help" and exits 0 — its Args validator is never reached
// and `serv-kun postgres backupp` succeeds having backed up nothing. Setting
// cobra.NoArgs alone does not fix that; the command has to be runnable for
// cobra to get as far as checking the arguments.
func showHelp(cmd *cobra.Command, _ []string) error {
	return cmd.Help()
}

func init() {
	nginxCmd.AddCommand(nginxSetupCmd, nginxSSLCmd)
	sshCmd.AddCommand(sshHardenCmd, sshAuditCmd)
	postgresCmd.AddCommand(postgresSetupCmd, postgresBackupCmd, postgresRestoreCmd)
	springbootCmd.AddCommand(springbootInitCmd, springbootSetupCmd, springbootDeployCmd, springbootCleanCmd, springbootStatusCmd, springbootLogsCmd)
	rootCmd.AddCommand(nginxCmd, userCmd, ciCmd, sshCmd, firewallCmd, postgresCmd, nodeCmd, springbootCmd)
}

// requireRoot exits with an actionable error if the current process isn't root.
func requireRoot(cmd *cobra.Command) {
	if os.Getuid() != 0 {
		fmt.Fprintf(os.Stderr, "error: must be run as root — try: sudo %s\n", cmd.CommandPath())
		os.Exit(1)
	}
}
