package firewall

import (
	"strings"
	"testing"
)

// ufwSequence mirrors the command order configureUFW emits, so the ordering
// invariant can be asserted without a live ufw. It calls the same arg builder
// the real path does, so a change to rule rendering can't silently drift away
// from what these tests check.
func ufwSequence(rules []Rule) []string {
	var seq []string
	for _, r := range rules {
		for _, args := range ufwArgs(r) {
			seq = append(seq, strings.Join(args, " "))
		}
	}
	return append(seq, "default deny incoming", "default allow outgoing", "enable")
}

// isAllowStep reports whether a step opens a port. Both ufw verbs do: `limit`
// is `allow` plus a rate cap, so both must land before the policy change.
func isAllowStep(step string) bool {
	return strings.HasPrefix(step, "allow ") || strings.HasPrefix(step, "limit ")
}

// `ufw default deny incoming` applies immediately on an already-active
// firewall, so the allow rules — SSH above all — have to land first. An
// interrupt between the two would otherwise leave deny-all with no way in.
func TestAllowRulesPrecedeDefaultDeny(t *testing.T) {
	seq := ufwSequence(DefaultRules(SSHOptions{Port: "22"}))
	joined := strings.Join(seq, " | ")

	ssh := indexOfPrefix(seq, "allow 22/tcp")
	deny := indexOf(seq, "default deny incoming")
	enable := indexOf(seq, "enable")

	if ssh == -1 || deny == -1 || enable == -1 {
		t.Fatalf("sequence missing an expected step: %s", joined)
	}
	if ssh > deny {
		t.Errorf("SSH allow (%d) comes after default-deny (%d): %s", ssh, deny, joined)
	}
	if deny > enable {
		t.Errorf("default-deny (%d) comes after enable (%d): %s", deny, enable, joined)
	}
	if enable != len(seq)-1 {
		t.Errorf("enable must be last, got position %d of %d: %s", enable, len(seq), joined)
	}
}

// Every allow must precede the policy change, not just the SSH one.
func TestAllRulesPrecedeDefaultDeny(t *testing.T) {
	rules := append(DefaultRules(SSHOptions{Port: "2222"}), Rule{Port: "8080", Proto: "tcp", Comment: "custom"})
	seq := ufwSequence(rules)
	deny := indexOf(seq, "default deny incoming")

	for i, step := range seq {
		if isAllowStep(step) && i > deny {
			t.Errorf("%q at position %d comes after default-deny at %d", step, i, deny)
		}
	}
}

// A rate-limited SSH rule is still an allow: `ufw limit` opens the port as
// well as capping it, so it has to land before default-deny for the same
// reason the plain allow does.
func TestRateLimitedSSHStillPrecedesDefaultDeny(t *testing.T) {
	seq := ufwSequence(DefaultRules(SSHOptions{Port: "22", RateLimit: true}))
	joined := strings.Join(seq, " | ")

	ssh := indexOfPrefix(seq, "limit 22/tcp")
	deny := indexOf(seq, "default deny incoming")

	if ssh == -1 {
		t.Fatalf("rate-limited SSH rule did not emit `limit 22/tcp`: %s", joined)
	}
	if ssh > deny {
		t.Errorf("SSH limit (%d) comes after default-deny (%d): %s", ssh, deny, joined)
	}
}

// Same invariant for a source-restricted rule, which renders in the longer
// `from X to any port N` form and so wouldn't be caught by the checks above.
func TestSourceRestrictedSSHPrecedesDefaultDeny(t *testing.T) {
	seq := ufwSequence(DefaultRules(SSHOptions{
		Port: "22",
		From: []string{"203.0.113.0/24", "198.51.100.7"},
	}))
	deny := indexOf(seq, "default deny incoming")

	var sshRules int
	for i, step := range seq {
		if !isAllowStep(step) {
			continue
		}
		if i > deny {
			t.Errorf("%q at position %d comes after default-deny at %d", step, i, deny)
		}
		if strings.Contains(step, "port 22") {
			sshRules++
		}
	}

	// One rule per source: ufw accepts a single source per rule, so a
	// two-address list that emitted one command would be silently allowing
	// only the first.
	if sshRules != 2 {
		t.Errorf("got %d SSH rules for 2 sources, want 2: %s", sshRules, strings.Join(seq, " | "))
	}
}

func indexOf(list []string, want string) int {
	for i, s := range list {
		if s == want {
			return i
		}
	}
	return -1
}

func indexOfPrefix(list []string, prefix string) int {
	for i, s := range list {
		if strings.HasPrefix(s, prefix) {
			return i
		}
	}
	return -1
}
