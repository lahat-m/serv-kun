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

	ssh := firewallpkg.SSHOptions{Port: sshPort}

	// Rate-limiting SSH is cheap brute-force protection, but it counts every
	// connection from an address, not just failed ones — so anything that
	// opens SSH sessions in bursts from one IP (a CI deploy, a parallel
	// rsync, a fleet behind one NAT gateway) can trip it. Ask rather than
	// impose it.
	if firewallpkg.SupportsLimit(backend) {
		if err := survey.AskOne(&survey.Confirm{
			Message: fmt.Sprintf("Rate-limit inbound SSH on port %s?", sshPort),
			Default: true,
			Help:    "blocks a source IP that opens 6+ connections in 30 seconds — beware if CI or a NAT'd office deploys over SSH",
		}, &ssh.RateLimit); err != nil {
			return err
		}
	}

	// A rate limit only helps against one busy address. An attack spread
	// across thousands of addresses, each trying twice, never reaches any
	// per-source cap — and an allow-list is the only firewall answer to it,
	// because it refuses everything not named however little each source
	// does. It is also the fastest way to lock yourself out, hence the
	// default of no and the check in askSSHSources.
	var restrict bool
	if err := survey.AskOne(&survey.Confirm{
		Message: "Restrict SSH to specific source addresses?",
		Default: false,
		Help:    "the only firewall control that stops a distributed attack — but leave an address out and you lose access from it",
	}, &restrict); err != nil {
		return err
	}
	if restrict {
		ssh.From, err = askSSHSources(sshPort)
		if err != nil {
			return err
		}
	}

	rules := firewallpkg.DefaultRules(ssh)

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
		note := r.Comment
		if len(r.From) > 0 {
			note += " from " + strings.Join(r.From, ", ")
		} else {
			note += " from anywhere"
		}
		if r.Limit {
			note += " (rate-limited: 6 connections per 30s, per source IP)"
		}
		fmt.Printf("  %-10s %s\n", r.Port+"/"+r.Proto, note)
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

// askSSHSources collects the addresses allowed to reach SSH, checked against
// the sessions that are open right now.
//
// Returning nil means "don't restrict after all" — falling back to an open
// SSH port is the right outcome when the operator can't confirm the list is
// safe, because a firewall nobody can get past is worse than one an attacker
// has to guess a key for.
func askSSHSources(sshPort string) ([]string, error) {
	// Ask the kernel who is connected rather than reading SSH_CLIENT: sudo's
	// default env_reset drops it, and this command always runs under sudo.
	// It also covers a colleague's session, not just this one.
	peers, err := sshd.ActivePeers(sshPort)
	switch {
	case err != nil:
		fmt.Printf("\n⚠  Could not read current SSH connections (%v).\n", err)
		fmt.Println("   Nothing here can check your list — get it wrong and you lose access.")
	case len(peers) == 0:
		fmt.Println("\n⚠  No established SSH connections found.")
		fmt.Println("   That's expected on a console session, but it means your list can't be checked.")
	default:
		fmt.Println("\nConnected over SSH right now:")
		for _, p := range peers {
			fmt.Printf("  %s\n", p)
		}
		fmt.Println("Every one of these must be covered, or its session ends when the policy applies.")
	}

	prompt := &survey.Input{
		Message: "Addresses allowed to reach SSH (space-separated CIDR or IP):",
		Help:    "e.g. 203.0.113.0/24 198.51.100.7 — include your office, your VPN, and any CI runner that deploys",
	}
	// Pre-filling with the live peers is a starting point, not an answer: a
	// single dynamic address will change, and CI won't be in the list at all.
	if len(peers) > 0 {
		prompt.Default = strings.Join(peers, " ")
	}

	var answer string
	if err := survey.AskOne(prompt, &answer,
		survey.WithValidator(validateSourceList)); err != nil {
		return nil, err
	}
	sources := strings.Fields(answer)

	var stranded []string
	for _, p := range peers {
		if !firewallpkg.Covers(sources, p) {
			stranded = append(stranded, p)
		}
	}
	if len(stranded) == 0 {
		return sources, nil
	}

	fmt.Printf("\n⚠  Not covered by that list: %s\n", strings.Join(stranded, ", "))
	fmt.Println("   Those sessions end the moment the firewall comes up — possibly this one.")

	var proceed bool
	if err := survey.AskOne(&survey.Confirm{
		Message: "Apply it anyway (only if you have console or rescue access)?",
		Default: false,
	}, &proceed); err != nil {
		return nil, err
	}
	if proceed {
		return sources, nil
	}

	fmt.Println("\nLeaving SSH reachable from anywhere. Re-run once you know the full list.")
	return nil, nil
}
