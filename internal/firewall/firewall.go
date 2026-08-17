// Package firewall configures the host firewall — ufw on Debian-family
// systems, firewalld on RHEL — with a default-deny inbound policy and an
// explicit allow-list. SSH is permitted before default-deny takes effect, so
// the ordering alone can't strand the session that ran it.
//
// Restricting SSH to named sources still can, which is the one thing here
// that needs checking rather than ordering: see Covers, and the live-session
// check the firewall command runs before applying such a list.
package firewall

import (
	"fmt"
	"net/netip"
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
	Port    string   // "22", "80", "8080"
	Proto   string   // "tcp" or "udp"
	Comment string   // shown in `ufw status`; ignored by firewalld
	Limit   bool     // rate-limit per source IP — ufw only, see verb()
	From    []string // source CIDRs or addresses; empty means from anywhere
}

func (r Rule) spec() string { return r.Port + "/" + r.Proto }

// verb is the ufw command word for the rule. `limit` is `allow` plus a
// per-source-IP rate limit: ufw denies an address that opens 6 or more
// connections to the port within 30 seconds, which is what turns a
// credential-stuffing run into a handful of tries. It is per source, so one
// attacker hitting the limit cannot lock anyone else out.
//
// Only worth setting on SSH. Rate-limiting HTTP would throttle ordinary
// browsers, which open several connections to one page as a matter of course.
func (r Rule) verb() string {
	if r.Limit {
		return "limit"
	}
	return "allow"
}

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

// SSHOptions tunes the SSH rule. Both settings are decisions only the
// operator can make — one throttles their own automation, the other can
// strand them — so neither has a default worth guessing at here.
type SSHOptions struct {
	Port string

	// RateLimit selects ufw's `limit` verb: a per-source cap that slows a
	// single busy address. It does nothing about an attack spread across
	// many addresses, each staying under the cap.
	RateLimit bool

	// From restricts the rule to these source CIDRs or addresses. This is
	// the only firewall control that stops a distributed attack, because it
	// refuses everything not on the list however few attempts each source
	// makes. Empty means reachable from anywhere.
	From []string
}

// DefaultRules is the baseline for a web server: SSH on whichever port it
// actually listens on, plus HTTP and HTTPS.
func DefaultRules(ssh SSHOptions) []Rule {
	return []Rule{
		{Port: ssh.Port, Proto: "tcp", Comment: "SSH", Limit: ssh.RateLimit, From: ssh.From},
		{Port: "80", Proto: "tcp", Comment: "HTTP"},
		{Port: "443", Proto: "tcp", Comment: "HTTPS"},
	}
}

// Covers reports whether addr falls inside any of the given source
// specifications, each a CIDR ("203.0.113.0/24") or a bare address.
//
// The lockout check depends on this: an allow-list that doesn't cover a
// currently connected session will end it the moment the policy applies.
// Anything unparseable is deliberately not treated as a match — a typo must
// read as "not covered" and stop the operator, never as "covered" and let
// them proceed.
func Covers(sources []string, addr string) bool {
	target, err := netip.ParseAddr(addr)
	if err != nil {
		return false
	}
	target = target.Unmap()

	for _, src := range sources {
		if prefix, err := netip.ParsePrefix(src); err == nil {
			if prefix.Masked().Contains(target) {
				return true
			}
			continue
		}
		if one, err := netip.ParseAddr(src); err == nil && one.Unmap() == target {
			return true
		}
	}
	return false
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
		for _, args := range ufwArgs(r) {
			if err := system.RunCmd("ufw", args...); err != nil {
				return fmt.Errorf("ufw %s: %w", strings.Join(args, " "), err)
			}
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

// SupportsLimit reports whether the backend can rate-limit a port per source
// address. Only ufw can, so callers must not offer the option elsewhere.
//
// firewalld deliberately gets nothing here. Its rich-rule `limit` counts
// matches against the rule itself rather than per source, so a single attacker
// saturating the budget would have every other address — including the
// operator's — fall through to the default deny. A lockout is a worse outcome
// than the brute-force noise it would suppress.
func SupportsLimit(b Backend) bool { return b == UFW }

// configureFirewalld starts the daemon first — firewall-cmd won't accept
// changes otherwise. That's safe here because firewalld's default public
// zone already permits SSH.
//
// Rule.Limit is ignored — see SupportsLimit.
func configureFirewalld(rules []Rule) error {
	if err := system.RunCmd("systemctl", "enable", "--now", "firewalld"); err != nil {
		return fmt.Errorf("start firewalld: %w", err)
	}

	for _, r := range rules {
		for _, args := range firewalldArgs(r) {
			if err := system.RunCmd("firewall-cmd", args...); err != nil {
				return fmt.Errorf("firewall-cmd %s: %w", strings.Join(args, " "), err)
			}
		}
	}

	if err := system.RunCmd("firewall-cmd", "--reload"); err != nil {
		return fmt.Errorf("firewalld reload: %w", err)
	}
	return nil
}

// ufwArgs renders the ufw invocations for one rule — one per source address,
// since ufw takes a single source per rule. A source-restricted rule has to
// name the port explicitly (`from X to any port 22 proto tcp`); the shorthand
// `22/tcp` form accepts no source at all.
func ufwArgs(r Rule) [][]string {
	if len(r.From) == 0 {
		return [][]string{withComment([]string{r.verb(), r.spec()}, r)}
	}

	var cmds [][]string
	for _, src := range r.From {
		cmds = append(cmds, withComment(
			[]string{r.verb(), "from", src, "to", "any", "port", r.Port, "proto", r.Proto}, r))
	}
	return cmds
}

func withComment(args []string, r Rule) []string {
	if r.Comment == "" {
		return args
	}
	return append(args, "comment", r.Comment)
}

// firewalldArgs renders the firewall-cmd invocations for one rule.
//
// A source-restricted rule becomes a rich rule rather than `--add-port`,
// because the two are not additive in the direction you'd hope: an open port
// in the zone is reachable from anywhere regardless of any rich rule, so
// adding both would leave the port wide open. Rate limiting is still ignored
// here — see SupportsLimit.
func firewalldArgs(r Rule) [][]string {
	if len(r.From) == 0 {
		return [][]string{{"--permanent", "--add-port=" + r.spec()}}
	}

	var cmds [][]string
	for _, src := range r.From {
		family := "ipv4"
		if strings.Contains(src, ":") {
			family = "ipv6"
		}
		rule := fmt.Sprintf(
			`rule family="%s" source address="%s" port port="%s" protocol="%s" accept`,
			family, src, r.Port, r.Proto)
		cmds = append(cmds, []string{"--permanent", "--add-rich-rule=" + rule})
	}
	return cmds
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
