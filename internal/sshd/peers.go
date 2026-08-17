package sshd

// Who is connected right now.
//
// This exists for one reason: restricting SSH to a list of source addresses
// is the only firewall control that stops a distributed attack, and it is
// also the fastest way to lock yourself out of a server. Before writing such
// a list, something has to check it against the sessions that are actually
// open — including the one running the command.
//
// SSH_CLIENT and SSH_CONNECTION would be the obvious source, but sudo's
// default env_reset drops both, and every command that would use this runs
// under sudo. Asking the kernel is also better than asking one process: it
// covers a colleague's session as well as your own.

import (
	"fmt"
	"net"
	"net/netip"
	"os/exec"
	"sort"
	"strings"
)

// ActivePeers lists the remote addresses of established connections to the
// SSH port, newest state as the kernel sees it.
//
// An error here is not fatal to the caller — it means "this can't be
// verified", which is worth saying out loud rather than treating as "nobody
// is connected".
func ActivePeers(port string) ([]string, error) {
	if _, err := exec.LookPath("ss"); err != nil {
		return nil, fmt.Errorf("ss is not installed")
	}

	out, err := exec.Command("ss", "-Htn", "state", "established",
		"sport", "=", ":"+port).Output()
	if err != nil {
		return nil, fmt.Errorf("run ss: %w", err)
	}
	return parseSSPeers(string(out)), nil
}

// parseSSPeers pulls peer addresses out of `ss -Htn` output. Duplicates are
// collapsed: one operator with three windows open is one address to allow.
//
// The peer is taken as the *last* field on the line that parses as an
// address, not a fixed column. ss omits the State column when a state filter
// is given and keeps it otherwise, so a positional parse reads the local
// address instead on whichever form it wasn't written for — and an allow-list
// built from the server's own address covers nobody.
func parseSSPeers(out string) []string {
	seen := map[string]bool{}
	var peers []string

	for _, line := range strings.Split(out, "\n") {
		var host string
		for _, field := range strings.Fields(line) {
			if h := stripPort(field); h != "" {
				host = h
			}
		}
		if host == "" || seen[host] {
			continue
		}
		seen[host] = true
		peers = append(peers, host)
	}

	sort.Strings(peers)
	return peers
}

// stripPort turns one `address:port` field into a bare address, or "" if the
// field isn't one. Returning "" for anything unparseable is what lets the
// caller scan every field without picking up process names or queue counts.
//
// The IPv4-mapped form (`::ffff:203.0.113.7`) a dual-stack listener reports
// is normalised down to plain IPv4; left as-is it would never match an IPv4
// CIDR, and the operator would be told they aren't covered when they are.
func stripPort(addrPort string) string {
	host, _, err := net.SplitHostPort(addrPort)
	if err != nil {
		// SplitHostPort rejects the unbracketed `::ffff:1.2.3.4:22` that ss
		// prints, so fall back to cutting at the last colon.
		i := strings.LastIndex(addrPort, ":")
		if i <= 0 {
			return ""
		}
		host = addrPort[:i]
	}
	host = strings.Trim(host, "[]")

	// IsUnspecified drops the `0.0.0.0:*` and `[::]:*` a listening socket
	// carries in the peer column, which is not somebody to allow.
	addr, err := netip.ParseAddr(host)
	if err != nil || addr.IsUnspecified() {
		return ""
	}
	return addr.Unmap().String()
}
