package cmd

import (
	"fmt"
	"strings"

	survey "github.com/AlecAivazis/survey/v2"
	firewallpkg "github.com/lahat-m/serv-kun/internal/firewall"
	"github.com/lahat-m/serv-kun/internal/sshd"
	"github.com/lahat-m/serv-kun/internal/system"
	"github.com/spf13/cobra"
)

var firewallCmd = &cobra.Command{
	Use:   "firewall",
	Short: "Configure the host firewall — default-deny inbound, with SSH/HTTP/HTTPS allowed",
	Long: `Sets up ufw (Debian/Ubuntu) or firewalld (RHEL/Fedora) with a
default-deny inbound policy and an explicit allow-list.

SSH is detected from the running sshd config and allowed before default-deny
takes effect, so enabling the firewall won't cut off your session.`,
	RunE: runFirewall,
}

func runFirewall(cmd *cobra.Command, _ []string) error {
	requireRoot(cmd)

	printBanner("serv-kun firewall")

	osInfo, err := system.Detect()
	if err != nil {
		return err
	}
	fmt.Printf("Detected OS: %s\n", osInfo.Pretty)

	backend := firewallpkg.Detect(osInfo)
	fmt.Printf("Firewall backend: %s\n\n", backend)

	// ── install ────────────────────────────────────────────────────────────
	if !firewallpkg.IsInstalled(backend) {
		var install bool
		if err := survey.AskOne(&survey.Confirm{
			Message: fmt.Sprintf("%s is not installed. Install it now?", backend),
			Default: true,
		}, &install); err != nil {
			return err
		}
		if !install {
			return fmt.Errorf("%s is required to continue", backend)
		}
		fmt.Printf("Installing %s...\n", backend)
		if err := firewallpkg.Install(backend, osInfo); err != nil {
			return fmt.Errorf("install %s: %w", backend, err)
		}
		fmt.Printf("✓ %s installed\n", backend)
	} else {
		fmt.Printf("✓ %s is already installed\n", backend)
	}
	fmt.Println()

	// ── build the allow-list ───────────────────────────────────────────────
	// Read the real SSH port rather than assuming 22 — locking yourself out
	// of a server running sshd on a custom port is the classic way to lose it.
	sshPort, detected := sshd.Port()
	if detected {
		fmt.Printf("SSH is listening on port %s — it will be allowed.\n\n", sshPort)
	} else {
		// Don't quietly present the fallback as fact: if sshd is elsewhere,
		// allowing 22 and then enabling default-deny cuts off this session.
		fmt.Printf("⚠  Could not read the SSH port from sshd — assuming %s.\n", sshPort)
		fmt.Println("   Allowing the wrong port here would end your session when the firewall comes up.")
		if err := survey.AskOne(&survey.Input{
			Message: "SSH port to allow:",
			Default: sshPort,
			Help:    "check with: ss -tlnp | grep sshd",
		}, &sshPort, survey.WithValidator(validatePort)); err != nil {
			return err
		}
		fmt.Println()
	}

	rules := firewallpkg.DefaultRules(sshPort)

	var extra string
	if err := survey.AskOne(&survey.Input{
		Message: "Additional TCP ports to open (space-separated, blank for none):",
		Help:    "e.g. 8080 3000 — Postgres and other internal services should stay closed",
	}, &extra, survey.WithValidator(validateOptionalPortList)); err != nil {
		return err
	}
	for _, p := range strings.Fields(extra) {
		rules = append(rules, firewallpkg.Rule{Port: p, Proto: "tcp", Comment: "custom"})
	}

	// ── preview ────────────────────────────────────────────────────────────
	fmt.Println()
	printRule()
	fmt.Println("Default policy: deny incoming, allow outgoing")
	fmt.Println("Allowed inbound:")
	for _, r := range rules {
		fmt.Printf("  %-10s %s\n", r.Port+"/"+r.Proto, r.Comment)
	}
	printRule()
	fmt.Println()

	var apply bool
	if err := survey.AskOne(&survey.Confirm{
		Message: "Apply this firewall configuration?",
		Default: true,
	}, &apply); err != nil {
		return err
	}
	if !apply {
		fmt.Println("Aborted — no changes made.")
		return nil
	}

	// ── apply ──────────────────────────────────────────────────────────────
	fmt.Println()
	if err := firewallpkg.Configure(backend, rules); err != nil {
		return err
	}
	fmt.Printf("\n✓ Firewall enabled (%s)\n\n", backend)

	if status, err := firewallpkg.Status(backend); err == nil {
		printRule()
		fmt.Print(status)
		printRule()
	}

	return nil
}
