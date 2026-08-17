package firewall

import (
	"strings"
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
	rules := DefaultRules(SSHOptions{Port: "2222"})

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

// Rate limiting is for SSH alone: `ufw limit` on 80 or 443 would throttle
// ordinary browsers, which open several connections per page load.
func TestOnlySSHIsRateLimited(t *testing.T) {
	for _, r := range DefaultRules(SSHOptions{Port: "22", RateLimit: true}) {
		wantLimit := r.Comment == "SSH"
		if r.Limit != wantLimit {
			t.Errorf("rule %+v: Limit = %v, want %v", r, r.Limit, wantLimit)
		}
		wantVerb := "allow"
		if wantLimit {
			wantVerb = "limit"
		}
		if got := r.verb(); got != wantVerb {
			t.Errorf("rule %+v: verb() = %q, want %q", r, got, wantVerb)
		}
	}

	for _, r := range DefaultRules(SSHOptions{Port: "22"}) {
		if r.Limit {
			t.Errorf("rule %+v: rate-limiting was declined but Limit is set", r)
		}
	}
}

// Source restriction applies to SSH alone — putting HTTP behind an
// allow-list would take the site off the internet.
func TestOnlySSHIsSourceRestricted(t *testing.T) {
	for _, r := range DefaultRules(SSHOptions{Port: "22", From: []string{"203.0.113.0/24"}}) {
		if r.Comment == "SSH" {
			if len(r.From) != 1 {
				t.Errorf("SSH rule %+v: want one source", r)
			}
			continue
		}
		if len(r.From) != 0 {
			t.Errorf("rule %+v: %s must stay reachable from anywhere", r, r.Comment)
		}
	}
}

// ufw takes one source per rule, and a source-restricted rule can't use the
// `22/tcp` shorthand — that form accepts no source at all, so emitting it
// would open the port to everyone.
func TestUFWArgsRenderOneRulePerSource(t *testing.T) {
	r := Rule{
		Port:    "22",
		Proto:   "tcp",
		Comment: "SSH",
		Limit:   true,
		From:    []string{"203.0.113.0/24", "198.51.100.7"},
	}

	got := ufwArgs(r)
	if len(got) != 2 {
		t.Fatalf("ufwArgs() returned %d commands for 2 sources, want 2: %v", len(got), got)
	}

	want := [][]string{
		{"limit", "from", "203.0.113.0/24", "to", "any", "port", "22", "proto", "tcp", "comment", "SSH"},
		{"limit", "from", "198.51.100.7", "to", "any", "port", "22", "proto", "tcp", "comment", "SSH"},
	}
	for i := range want {
		if strings.Join(got[i], " ") != strings.Join(want[i], " ") {
			t.Errorf("ufwArgs()[%d]\n got %v\nwant %v", i, got[i], want[i])
		}
	}
}

func TestUFWArgsWithoutSourcesUsesShorthand(t *testing.T) {
	got := ufwArgs(Rule{Port: "443", Proto: "tcp", Comment: "HTTPS"})

	if len(got) != 1 {
		t.Fatalf("ufwArgs() returned %d commands, want 1: %v", len(got), got)
	}
	if want := "allow 443/tcp comment HTTPS"; strings.Join(got[0], " ") != want {
		t.Errorf("ufwArgs() = %q, want %q", strings.Join(got[0], " "), want)
	}
}

// firewalld's --add-port opens a port to every source, so a source-restricted
// rule must become a rich rule *instead of* the port, never as well as it.
func TestFirewalldArgsUseRichRulesForSources(t *testing.T) {
	got := firewalldArgs(Rule{
		Port: "22", Proto: "tcp", Comment: "SSH",
		From: []string{"203.0.113.0/24", "2001:db8::/48"},
	})

	if len(got) != 2 {
		t.Fatalf("firewalldArgs() returned %d commands, want 2: %v", len(got), got)
	}
	for _, args := range got {
		joined := strings.Join(args, " ")
		if strings.Contains(joined, "--add-port") {
			t.Errorf("%q opens the port to every source alongside the rich rule", joined)
		}
		if !strings.Contains(joined, "--add-rich-rule=") {
			t.Errorf("%q is not a rich rule", joined)
		}
	}

	if !strings.Contains(strings.Join(got[0], " "), `family="ipv4"`) {
		t.Errorf("IPv4 source did not render as family ipv4: %v", got[0])
	}
	if !strings.Contains(strings.Join(got[1], " "), `family="ipv6"`) {
		t.Errorf("IPv6 source did not render as family ipv6: %v", got[1])
	}
}

// The lockout check runs on this, so a wrong answer costs the operator their
// access. Unparseable input must read as "not covered", never as covered.
func TestCovers(t *testing.T) {
	cases := []struct {
		sources []string
		addr    string
		want    bool
	}{
		{sources: []string{"203.0.113.0/24"}, addr: "203.0.113.7", want: true},
		{sources: []string{"203.0.113.0/24"}, addr: "203.0.114.7"},
		{sources: []string{"198.51.100.7"}, addr: "198.51.100.7", want: true},
		{sources: []string{"198.51.100.7"}, addr: "198.51.100.8"},
		{sources: []string{"10.0.0.0/8", "203.0.113.0/24"}, addr: "203.0.113.9", want: true},
		// A /24 written with host bits set still describes that network.
		{sources: []string{"203.0.113.7/24"}, addr: "203.0.113.200", want: true},
		{sources: []string{"2001:db8::/48"}, addr: "2001:db8:0:1::5", want: true},
		{sources: []string{"2001:db8::/48"}, addr: "2001:db9::5"},
		// An IPv4 address is not inside an IPv6 prefix, or the reverse.
		{sources: []string{"2001:db8::/48"}, addr: "203.0.113.7"},
		{sources: []string{"203.0.113.0/24"}, addr: "2001:db8::1"},
		// The dual-stack listener reports IPv4 peers in mapped form.
		{sources: []string{"203.0.113.0/24"}, addr: "::ffff:203.0.113.7", want: true},
		// Garbage in must not read as a match.
		{sources: []string{"not-a-cidr"}, addr: "203.0.113.7"},
		{sources: []string{"203.0.113.0/24"}, addr: "not-an-address"},
		{sources: nil, addr: "203.0.113.7"},
	}

	for _, tc := range cases {
		if got := Covers(tc.sources, tc.addr); got != tc.want {
			t.Errorf("Covers(%v, %q) = %v, want %v", tc.sources, tc.addr, got, tc.want)
		}
	}
}

// firewalld's rich-rule limit counts per rule rather than per source address,
// so honouring Rule.Limit there would let one attacker starve every other
// address of the budget — the operator included.
func TestOnlyUFWSupportsLimit(t *testing.T) {
	if !SupportsLimit(UFW) {
		t.Error("SupportsLimit(UFW) = false, want true")
	}
	if SupportsLimit(Firewalld) {
		t.Error("SupportsLimit(Firewalld) = true — rich-rule limits are not per-source")
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
