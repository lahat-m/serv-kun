package cmd

import (
	"fmt"
	"os"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/lahat-m/serv-kun/internal/sshd"
	"github.com/lahat-m/serv-kun/internal/sudouser"
	"github.com/lahat-m/serv-kun/internal/system"
	"github.com/spf13/cobra"
)

var sshHardenCmd = &cobra.Command{
	Use:   "harden",
	Short: "Disable root login and password authentication over SSH",
	Long: `Disables root login and (optionally) password authentication in sshd.

Before changing anything, this verifies that a non-root sudo user has a
working SSH key installed — otherwise these settings would lock everyone
out of the server. The new config is validated with sshd's own check and
rolled back automatically if it doesn't pass.`,
	RunE: runSSHHarden,
}

func runSSHHarden(cmd *cobra.Command, _ []string) error {
	requireRoot(cmd)

	printBanner("ssh harden")

	osInfo, err := system.Detect()
	if err != nil {
		return err
	}
	fmt.Printf("Detected OS: %s\n\n", osInfo.Pretty)

	// ── current state ──────────────────────────────────────────────────────
	settings, err := sshd.EffectiveSettings()
	if err != nil {
		return fmt.Errorf("read current SSH config: %w", err)
	}
	fmt.Println("Current effective settings:")
	fmt.Printf("  Port                    %s\n", settings["port"])
	fmt.Printf("  PermitRootLogin         %s\n", settings["permitrootlogin"])
	fmt.Printf("  PasswordAuthentication  %s\n", settings["passwordauthentication"])
	fmt.Println()

	// ── lockout safety check ───────────────────────────────────────────────
	var username string
	if err := survey.AskOne(&survey.Input{
		Message: "Admin user that must keep access:",
		Default: os.Getenv("SUDO_USER"),
		Help:    "this user needs sudo rights and an SSH key, since root and passwords are about to go away",
	}, &username, survey.WithValidator(survey.Required)); err != nil {
		return err
	}

	safe, reason := checkLockoutSafety(username, osInfo)
	fmt.Println()
	if !safe {
		fmt.Printf("⚠  %s\n\n", reason)
		fmt.Println("Applying these settings now would leave no way back into this server.")
		fmt.Println("Fix it first with:  sudo serv-kun user")
		fmt.Println()

		var proceed bool
		if err := survey.AskOne(&survey.Confirm{
			Message: "Continue anyway (only if you have console/rescue access)?",
			Default: false,
		}, &proceed); err != nil {
			return err
		}
		if !proceed {
			fmt.Println("Aborted — no changes made.")
			return nil
		}
	} else {
		fmt.Printf("✓ %s\n", reason)
	}
	fmt.Println()

	// ── what to change ─────────────────────────────────────────────────────
	var opts sshd.Options
	if err := survey.AskOne(&survey.Confirm{
		Message: "Disable root login over SSH?",
		Default: true,
	}, &opts.DisableRootLogin); err != nil {
		return err
	}
	if err := survey.AskOne(&survey.Confirm{
		Message: "Disable password authentication (keys only)?",
		Default: true,
	}, &opts.DisablePasswordAuth); err != nil {
		return err
	}

	if !opts.Any() {
		fmt.Println("\nNothing selected — no changes made.")
		return nil
	}

	// ── preview ────────────────────────────────────────────────────────────
	useDropIn, err := sshd.UsesIncludeDir()
	if err != nil {
		return err
	}
	target := sshd.ConfigPath
	if useDropIn {
		target = sshd.IncludeDir + "/" + sshd.DropInName
	}

	fmt.Printf("\nWill write → %s\n", target)
	if useDropIn {
		fmt.Println("(a drop-in that sorts first, so it wins over cloud-init's defaults)")
	}
	fmt.Println()
	printRule()
	fmt.Println(sshd.Render(opts))
	printRule()
	fmt.Println()

	var apply bool
	if err := survey.AskOne(&survey.Confirm{
		Message: "Apply this configuration?",
		Default: true,
	}, &apply); err != nil {
		return err
	}
	if !apply {
		fmt.Println("Aborted — no changes made.")
		return nil
	}

	// ── apply, validate, roll back on failure ──────────────────────────────
	written, backup, err := sshd.Apply(opts)
	if err != nil {
		return err
	}
	if backup != "" {
		fmt.Printf("  Backed up existing config → %s\n", backup)
	}
	fmt.Printf("  Config written → %s\n", written)

	fmt.Println("\nValidating SSH configuration...")
	if err := sshd.Test(); err != nil {
		fmt.Println("✗ Validation failed — rolling back.")
		if restoreErr := sshd.Restore(written, backup); restoreErr != nil {
			return fmt.Errorf("%w\n\nROLLBACK ALSO FAILED: %v\nRestore %s by hand before reloading sshd", err, restoreErr, written)
		}
		fmt.Println("✓ Rolled back — sshd was never reloaded, so nothing changed.")
		return err
	}
	fmt.Println("✓ Configuration is valid")

	if err := sshd.Reload(osInfo); err != nil {
		return err
	}
	fmt.Println("✓ sshd reloaded")

	// ── confirm it actually took effect ────────────────────────────────────
	if after, err := sshd.EffectiveSettings(); err == nil {
		fmt.Println("\nNew effective settings:")
		fmt.Printf("  PermitRootLogin         %s\n", after["permitrootlogin"])
		fmt.Printf("  PasswordAuthentication  %s\n", after["passwordauthentication"])
	}

	fmt.Println()
	printRule()
	fmt.Println("KEEP THIS SESSION OPEN.")
	fmt.Printf("Open a new terminal and confirm you can still log in:\n\n    ssh %s@<server>\n\n", username)
	fmt.Println("If that fails, this session is still connected and can undo the change.")
	printRule()

	return nil
}

// checkLockoutSafety verifies someone will still be able to get in and
// administer the box once root login and passwords are disabled.
func checkLockoutSafety(username string, osInfo *system.Info) (bool, string) {
	if username == "root" {
		return false, "That's the root account — disabling root login would lock it out. Pick a non-root sudo user."
	}

	keys, err := sshd.AuthorizedKeys(username)
	if err != nil {
		return false, fmt.Sprintf("Could not check %q: %v", username, err)
	}
	if len(keys) == 0 {
		return false, fmt.Sprintf("User %q has no SSH keys in ~/.ssh/authorized_keys.", username)
	}

	canSudo, reason := sudouser.CanSudo(username, osInfo)
	if !canSudo {
		return false, fmt.Sprintf("User %q has SSH keys but cannot run sudo — %s. You'd get in with no way to administer the server.",
			username, reason)
	}

	return true, fmt.Sprintf("User %q has %d SSH key(s) and is %s.", username, len(keys), reason)
}
