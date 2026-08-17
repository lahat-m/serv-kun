package sshd

import (
	"strings"
	"testing"
)

// Getting this wrong is expensive in one direction: a peer that goes
// unrecognised is a peer the lockout check can't protect, and the operator
// applies an allow-list that drops their own session.
func TestParseSSPeers(t *testing.T) {
	out := strings.Join([]string{
		"0 0   10.0.0.5:22   203.0.113.7:54321",
		"0 64  10.0.0.5:22   203.0.113.7:54322", // same operator, two windows
		"0 0   [2001:db8::5]:22   [2001:db8::99]:40000",
		"0 0   ::ffff:10.0.0.5:22   ::ffff:198.51.100.4:40001",
		"", // ss ends with a blank line
	}, "\n")

	got := parseSSPeers(out)
	want := []string{"198.51.100.4", "2001:db8::99", "203.0.113.7"} // sorted

	if strings.Join(got, " ") != strings.Join(want, " ") {
		t.Errorf("parseSSPeers() = %v, want %v", got, want)
	}
}

// ss keeps the State column when no state filter is given and drops it when
// one is. Both forms must yield the peer, never the server's own address —
// an allow-list built from the local address covers nobody.
func TestParseSSPeersIgnoresColumnShift(t *testing.T) {
	withState := "ESTAB 0 0 10.0.0.5:22 203.0.113.7:54321"
	withoutState := "0 0 10.0.0.5:22 203.0.113.7:54321"

	for _, line := range []string{withState, withoutState} {
		got := parseSSPeers(line)
		if len(got) != 1 || got[0] != "203.0.113.7" {
			t.Errorf("parseSSPeers(%q) = %v, want [203.0.113.7]", line, got)
		}
	}
}

func TestParseSSPeersIgnoresMalformedLines(t *testing.T) {
	out := "ESTAB 0 0\nnot a table row at all\n\n"

	if got := parseSSPeers(out); len(got) != 0 {
		t.Errorf("parseSSPeers() = %v, want none", got)
	}
}

func TestStripPort(t *testing.T) {
	cases := []struct{ in, want, why string }{
		{in: "203.0.113.7:54321", want: "203.0.113.7", why: "plain IPv4"},
		{in: "[2001:db8::99]:40000", want: "2001:db8::99", why: "bracketed IPv6"},
		{in: "2001:db8::99:40000", want: "2001:db8::99", why: "some ss builds omit the brackets"},
		{in: "::ffff:198.51.100.4:4000", want: "198.51.100.4", why: "unmapped so it matches an IPv4 CIDR"},
		{in: "nonsense", want: "", why: "not an address"},
		{in: "ESTAB", want: "", why: "the State column"},
		{in: "0", want: "", why: "a queue count"},
		{in: `users:(("sshd",pid=1,fd=4))`, want: "", why: "process info"},
		{in: "0.0.0.0:*", want: "", why: "a listening socket has no peer"},
		{in: "[::]:*", want: "", why: "same, IPv6"},
	}

	for _, tc := range cases {
		if got := stripPort(tc.in); got != tc.want {
			t.Errorf("stripPort(%q) = %q, want %q (%s)", tc.in, got, tc.want, tc.why)
		}
	}
}
