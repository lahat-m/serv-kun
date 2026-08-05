package firewall

import (
	"strings"
	"testing"
)

// ufwSequence mirrors the command order configureUFW emits, so the ordering
// invariant can be asserted without a live ufw.
func ufwSequence(rules []Rule) []string {
	var seq []string
	for _, r := range rules {
		seq = append(seq, "allow "+r.spec())
	}
	seq = append(seq, "default deny incoming", "default allow outgoing", "enable")
	return seq
}

// `ufw default deny incoming` applies immediately on an already-active
// firewall, so the allow rules — SSH above all — have to land first. An
// interrupt between the two would otherwise leave deny-all with no way in.
func TestAllowRulesPrecedeDefaultDeny(t *testing.T) {
	seq := ufwSequence(DefaultRules("22"))
	joined := strings.Join(seq, " | ")

	ssh := indexOf(seq, "allow 22/tcp")
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
	rules := append(DefaultRules("2222"), Rule{Port: "8080", Proto: "tcp", Comment: "custom"})
	seq := ufwSequence(rules)
	deny := indexOf(seq, "default deny incoming")

	for i, step := range seq {
		if strings.HasPrefix(step, "allow ") && i > deny {
			t.Errorf("%q at position %d comes after default-deny at %d", step, i, deny)
		}
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
