// Package firewall configures the host firewall — ufw on Debian-family
// systems, firewalld on RHEL — with a default-deny inbound policy and an
// explicit allow-list. SSH is always permitted before default-deny takes
// effect, so enabling the firewall can't strand the session that ran it.
package firewall

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/lahat-m/serv-kun/internal/system"
)

type Backend string

const (
	UFW       Backend = "ufw"
	Firewalld Backend = "firewalld"
)

// Rule is one inbound allow.
type Rule struct {
	Port    string // "22", "80", "8080"
	Proto   string // "tcp" or "udp"
	Comment string // shown in `ufw status`; ignored by firewalld
}

func (r Rule) spec() string { return r.Port + "/" + r.Proto }

// Detect picks the firewall that fits the OS family.
func Detect(osInfo *system.Info) Backend {
	if osInfo.Family == system.RHEL {
		return Firewalld
	}
	return UFW
}

// IsInstalled reports whether the backend's CLI is present.
func IsInstalled(b Backend) bool {
	_, err := exec.LookPath(command(b))
	return err == nil
}

func command(b Backend) string {
	if b == Firewalld {
		return "firewall-cmd"
	}
	return "ufw"
}

// Install fetches the backend's package.
func Install(b Backend, osInfo *system.Info) error {
	switch osInfo.Family {
	case system.Debian:
		if err := system.RunCmd("apt-get", "update", "-qq"); err != nil {
			return fmt.Errorf("apt-get update: %w", err)
		}
		return system.RunCmd("apt-get", "install", "-y", "ufw")

	case system.RHEL:
		return system.RunCmd("dnf", "install", "-y", "firewalld")
	}
	return fmt.Errorf("unsupported OS family")
}

// DefaultRules is the baseline for a web server: SSH on whichever port it
// actually listens on, plus HTTP and HTTPS.
func DefaultRules(sshPort string) []Rule {
	return []Rule{
		{Port: sshPort, Proto: "tcp", Comment: "SSH"},
		{Port: "80", Proto: "tcp", Comment: "HTTP"},
		{Port: "443", Proto: "tcp", Comment: "HTTPS"},
	}
}

// Configure applies the allow-list and turns the firewall on. Step order is
// backend-specific and load-bearing — see each implementation.
func Configure(b Backend, rules []Rule) error {
	switch b {
	case UFW:
		return configureUFW(rules)
	case Firewalld:
		return configureFirewalld(rules)
	}
	return fmt.Errorf("unsupported firewall backend %q", b)
}

// configureUFW adds the allow rules first, then sets the default policy, then
// enables.
//
// The order matters in both directions. `ufw default deny incoming` takes
// effect immediately on an already-active firewall, so setting the policy
// first opens a window where new inbound connections — including SSH — are
// refused; an interrupt in that window leaves deny-all with no way back in.
// (An in-flight session survives either way, via ufw's RELATED,ESTABLISHED
// rule in before.rules.) And `ufw enable` applies the policy the moment it
// runs, so it has to come last, after the allows are in place.
func configureUFW(rules []Rule) error {
	for _, r := range rules {
		args := []string{"allow", r.spec()}
		if r.Comment != "" {
			args = append(args, "comment", r.Comment)
		}
		if err := system.RunCmd("ufw", args...); err != nil {
			return fmt.Errorf("ufw allow %s: %w", r.spec(), err)
		}
	}

	if err := system.RunCmd("ufw", "default", "deny", "incoming"); err != nil {
		return fmt.Errorf("ufw default deny incoming: %w", err)
	}
	if err := system.RunCmd("ufw", "default", "allow", "outgoing"); err != nil {
		return fmt.Errorf("ufw default allow outgoing: %w", err)
	}

	// --force skips ufw's interactive "may disrupt existing ssh connections"
	// prompt, which we've already made safe with the allow rules above.
	return system.RunCmd("ufw", "--force", "enable")
}

// configureFirewalld starts the daemon first — firewall-cmd won't accept
// changes otherwise. That's safe here because firewalld's default public
// zone already permits SSH.
func configureFirewalld(rules []Rule) error {
	if err := system.RunCmd("systemctl", "enable", "--now", "firewalld"); err != nil {
		return fmt.Errorf("start firewalld: %w", err)
	}

	for _, r := range rules {
		if err := system.RunCmd("firewall-cmd", "--permanent", "--add-port="+r.spec()); err != nil {
			return fmt.Errorf("add port %s: %w", r.spec(), err)
		}
	}

	if err := system.RunCmd("firewall-cmd", "--reload"); err != nil {
		return fmt.Errorf("firewalld reload: %w", err)
	}
	return nil
}

// Status returns a human-readable dump of the active rules.
func Status(b Backend) (string, error) {
	var out []byte
	var err error

	switch b {
	case UFW:
		out, err = exec.Command("ufw", "status", "verbose").CombinedOutput()
	case Firewalld:
		out, err = exec.Command("firewall-cmd", "--list-all").CombinedOutput()
	default:
		return "", fmt.Errorf("unsupported firewall backend %q", b)
	}

	if err != nil {
		return "", fmt.Errorf("read firewall status:\n%s", strings.TrimSpace(string(out)))
	}
	return string(out), nil
}
