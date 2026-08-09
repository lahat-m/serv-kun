package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/lahat-m/serv-kun/internal/compose"
	"github.com/lahat-m/serv-kun/internal/sudouser"
	"github.com/spf13/cobra"
)

var ciCmd = &cobra.Command{
	Use:   "ci",
	Short: "Grant the deploy user passwordless sudo for the commands CI runs",
	Long: `Write a sudoers rule letting the CI user run the deploy without a password.

The workflow scaffolded by ` + "`springboot init`" + ` SSHes in and runs
` + "`sudo serv-kun springboot deploy`" + `. That session has no terminal, so sudo
cannot prompt for a password and the deploy dies with "a terminal is required
to authenticate" — a message that names neither sudoers nor the user involved.

This grants exactly the commands that workflow runs, and nothing else. Every
other sudo still asks for a password.`,
	Args: cobra.NoArgs,
	RunE: runCI,
}

func runCI(cmd *cobra.Command, _ []string) error {
	requireRoot(cmd)
	printBanner("serv-kun ci")

	commands, err := ciCommands()
	if err != nil {
		return err
	}

	username, err := askCIUser()
	if err != nil {
		return err
	}

	fmt.Printf("\nThis writes %s:\n\n", sudouser.SudoersFile)
	printRule()
	fmt.Print(sudouser.NoPasswdRule(username, commands))
	printRule()
	fmt.Println()

	// Replacing the file silently would revoke whatever it already grants —
	// most plausibly a second deploy user added on an earlier run. Show it
	// rather than describe it: the operator can only judge what they can see.
	if existing, err := os.ReadFile(sudouser.SudoersFile); err == nil {
		fmt.Printf("⚠ %s already exists and will be replaced. It holds:\n\n", sudouser.SudoersFile)
		printRule()
		fmt.Println(strings.TrimRight(string(existing), "\n"))
		printRule()
		fmt.Println("\n  Anything above that isn't in the new rule is revoked. If it grants a")
		fmt.Println("  different user, copy it somewhere before continuing.")
		fmt.Println()
	}

	var confirm bool
	if err := survey.AskOne(&survey.Confirm{
		Message: fmt.Sprintf("Grant %q passwordless sudo for those commands?", username),
		Default: true,
	}, &confirm); err != nil {
		return err
	}
	if !confirm {
		fmt.Println("\nNothing written.")
		return nil
	}

	if err := sudouser.InstallNoPasswd(username, commands); err != nil {
		return err
	}
	fmt.Printf("\n✓ %s written (0440, root:root) and accepted by visudo\n", sudouser.SudoersFile)

	// Ask sudo rather than trust the write: a rule can be valid, installed,
	// and still not apply — an overriding entry later in /etc/sudoers, or a
	// username that resolves differently than it reads.
	if !sudouser.HasNoPasswdSudo(username) {
		return fmt.Errorf("the file was installed, but sudo still reports no NOPASSWD rule for %q — "+
			"check `sudo -l -U %s` and whether /etc/sudoers overrides it", username, username)
	}
	fmt.Printf("✓ sudo confirms %q can run them without a password\n\n", username)

	fmt.Println("Check it the way CI will, from your workstation:")
	fmt.Printf("  ssh %s@<server> 'sudo -n serv-kun springboot status'\n\n", username)
	fmt.Printf("The rule is per-user, so HETZNER_SSH_USER in the repo's GitHub secrets\n")
	fmt.Printf("has to be %q or it won't match.\n", username)
	return nil
}

// ciCommands is the set the scaffolded deploy workflow actually runs.
func ciCommands() ([]sudouser.Command, error) {
	servKun, err := servKunPath()
	if err != nil {
		return nil, err
	}

	// serv-kun is granted with no argument restriction: its subcommands are
	// the entire point, and pinning them would break on the next release that
	// adds one.
	commands := []sudouser.Command{{Path: servKun}}

	// docker is granted for one exact invocation — the workflow's trailing
	// `docker image prune -f`. Unrestricted docker would be a far larger grant
	// than passwordless serv-kun: `docker run -v /:/host` is a root shell.
	dockerPath, err := exec.LookPath("docker")
	if err != nil {
		fmt.Println("Note: docker isn't on PATH, so the image-prune rule is being left out.")
		fmt.Println("      Re-run this once docker is installed to add it.")
		fmt.Println()
		return commands, nil
	}
	return append(commands, sudouser.Command{
		Path: dockerPath,
		Args: []string{"image", "prune", "-f"},
	}), nil
}

// servKunPath resolves the path a sudoers rule has to name for `sudo serv-kun`
// to match it.
//
// LookPath comes first deliberately: it mirrors how sudo finds the command,
// searching secure_path and matching the entry it lands on without resolving
// symlinks. os.Executable() reads /proc/self/exe, which *is* symlink-resolved
// — the right answer to "where am I really", and the wrong one to "what will
// sudo compare against" if /usr/local/bin/serv-kun is a link to elsewhere.
func servKunPath() (string, error) {
	if p, err := exec.LookPath("serv-kun"); err == nil {
		// LookPath can return a relative path; sudoers requires a fully
		// qualified one.
		if abs, err := filepath.Abs(p); err == nil {
			return abs, nil
		}
	}
	if p, err := os.Executable(); err == nil {
		return p, nil
	}
	return "", errors.New("could not determine serv-kun's own path — install it at /usr/local/bin/serv-kun and re-run")
}

func askCIUser() (string, error) {
	// SUDO_USER is whoever invoked `sudo serv-kun ci`, which is nearly always
	// the account CI will SSH in as — this server was set up from it.
	var username string
	if err := survey.AskOne(&survey.Input{
		Message: "Which user does CI SSH in as?",
		Default: os.Getenv("SUDO_USER"),
		Help:    "Whatever the HETZNER_SSH_USER secret holds in the GitHub repo.",
	}, &username, survey.WithValidator(survey.Required)); err != nil {
		return "", err
	}
	username = strings.TrimSpace(username)

	if username == "root" {
		return "", errors.New("root runs everything without sudo already — no rule is needed")
	}
	if !sudouser.Exists(username) {
		return "", fmt.Errorf("no user %q on this server — create one with: sudo serv-kun user", username)
	}
	return username, nil
}

// deployWorkflowPath is the workflow `springboot init` scaffolds, relative to
// the project directory.
const deployWorkflowPath = ".github/workflows/deploy.yml"

// warnMissingCIAccess reports a scaffolded deploy workflow that cannot yet
// work, while the operator is still at the keyboard to fix it.
//
// The failure it heads off arrives days later, from a machine nobody is
// watching, as "sudo: a terminal is required to authenticate" — a complaint
// about terminals, on a host that has none, mentioning neither sudoers nor
// which user it means.
func warnMissingCIAccess(proj compose.Project) {
	if _, err := os.Stat(filepath.Join(proj.Dir, deployWorkflowPath)); err != nil {
		return
	}

	// Empty for a genuine root login, and root needs no rule either way.
	user := os.Getenv("SUDO_USER")
	if user == "" || user == "root" {
		return
	}
	if sudouser.HasNoPasswdSudo(user) {
		return
	}

	fmt.Printf("⚠ %s is here, but %q has no passwordless sudo rule.\n", deployWorkflowPath, user)
	fmt.Println("  That workflow runs `sudo serv-kun springboot deploy` over an SSH session")
	fmt.Println("  with no terminal attached, so sudo has nowhere to ask for a password and")
	fmt.Println("  the deploy fails with \"a terminal is required to authenticate\".")
	fmt.Println()
	fmt.Println("  Fix it before the first push:  sudo serv-kun ci")
	fmt.Println()
}
