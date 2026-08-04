package firewall

import (
	"testing"

	"github.com/lahat-m/serv-kun/internal/system"
)

func TestDetectByOSFamily(t *testing.T) {
	if got := Detect(&system.Info{Family: system.Debian}); got != UFW {
		t.Errorf("Detect(debian) = %q, want %q", got, UFW)
	}
	if got := Detect(&system.Info{Family: system.RHEL}); got != Firewalld {
		t.Errorf("Detect(rhel) = %q, want %q", got, Firewalld)
	}
}

// The SSH rule must follow the running sshd's port, not a hardcoded 22.
func TestDefaultRulesUsesGivenSSHPort(t *testing.T) {
	rules := DefaultRules("2222")

	if len(rules) != 3 {
		t.Fatalf("DefaultRules() returned %d rules, want 3", len(rules))
	}
	if rules[0].Port != "2222" || rules[0].Comment != "SSH" {
		t.Errorf("first rule = %+v, want SSH on 2222", rules[0])
	}

	ports := map[string]bool{}
	for _, r := range rules {
		ports[r.Port] = true
		if r.Proto != "tcp" {
			t.Errorf("rule %+v: want tcp", r)
		}
	}
	for _, want := range []string{"2222", "80", "443"} {
		if !ports[want] {
			t.Errorf("DefaultRules() missing port %s", want)
		}
	}
}

func TestRuleSpec(t *testing.T) {
	if got := (Rule{Port: "443", Proto: "tcp"}).spec(); got != "443/tcp" {
		t.Errorf("spec() = %q, want %q", got, "443/tcp")
	}
}

func TestCommandByBackend(t *testing.T) {
	if got := command(UFW); got != "ufw" {
		t.Errorf("command(UFW) = %q", got)
	}
	if got := command(Firewalld); got != "firewall-cmd" {
		t.Errorf("command(Firewalld) = %q", got)
	}
}

func TestConfigureRejectsUnknownBackend(t *testing.T) {
	if err := Configure(Backend("nftables"), nil); err == nil {
		t.Error("Configure() with an unknown backend should error")
	}
}
