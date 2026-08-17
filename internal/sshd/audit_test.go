package sshd

import (
	"fmt"
	"testing"
	"time"
)

// Real lines, copied in the shapes OpenSSH emits. The classifier's whole job
// is telling "someone guessed a password on a real account" apart from
// "someone scanned a port", so each case names which axis it lands on.
func TestClassify(t *testing.T) {
	cases := []struct {
		name    string
		msg     string
		want    Event
		matched bool
	}{
		{
			name:    "failed password on a real account",
			msg:     "Failed password for root from 203.0.113.7 port 40404 ssh2",
			want:    Event{Kind: FailedAuth, Method: "password", User: "root", IP: "203.0.113.7"},
			matched: true,
		},
		{
			name:    "failed password on an account that does not exist",
			msg:     "Failed password for invalid user admin from 203.0.113.7 port 40404 ssh2",
			want:    Event{Kind: FailedAuth, Method: "password", User: "admin", IP: "203.0.113.7", Invalid: true},
			matched: true,
		},
		{
			name:    "rejected key",
			msg:     "Failed publickey for deploy from 198.51.100.4 port 51000 ssh2: ED25519 SHA256:abc",
			want:    Event{Kind: FailedAuth, Method: "publickey", User: "deploy", IP: "198.51.100.4"},
			matched: true,
		},
		{
			name:    "successful key login",
			msg:     "Accepted publickey for deploy from 198.51.100.4 port 51000 ssh2: ED25519 SHA256:abc",
			want:    Event{Kind: Accepted, Method: "publickey", User: "deploy", IP: "198.51.100.4"},
			matched: true,
		},
		{
			name:    "successful password login",
			msg:     "Accepted password for deploy from 198.51.100.4 port 51000 ssh2",
			want:    Event{Kind: Accepted, Method: "password", User: "deploy", IP: "198.51.100.4"},
			matched: true,
		},
		{
			name:    "username rejected outright",
			msg:     "Invalid user oracle from 203.0.113.7 port 40404",
			want:    Event{Kind: ProbedUser, User: "oracle", IP: "203.0.113.7", Invalid: true},
			matched: true,
		},
		{
			name:    "max tries exceeded",
			msg:     "error: maximum authentication attempts exceeded for root from 203.0.113.7 port 40404 ssh2 [preauth]",
			want:    Event{Kind: FailedAuth, User: "root", IP: "203.0.113.7"},
			matched: true,
		},
		{
			name:    "max tries exceeded for a nonexistent account",
			msg:     "error: maximum authentication attempts exceeded for invalid user admin from 203.0.113.7 port 40404 ssh2 [preauth]",
			want:    Event{Kind: FailedAuth, User: "admin", IP: "203.0.113.7", Invalid: true},
			matched: true,
		},
		{
			name:    "hung up before authenticating",
			msg:     "Connection closed by 203.0.113.9 port 40404 [preauth]",
			want:    Event{Kind: Probe, IP: "203.0.113.9"},
			matched: true,
		},
		{
			name:    "hung up while authenticating as a real account",
			msg:     "Connection closed by authenticating user root 203.0.113.9 port 40404 [preauth]",
			want:    Event{Kind: Probe, User: "root", IP: "203.0.113.9"},
			matched: true,
		},
		{
			name:    "hung up while authenticating as a nonexistent account",
			msg:     "Disconnected from invalid user admin 203.0.113.9 port 40404 [preauth]",
			want:    Event{Kind: Probe, User: "admin", IP: "203.0.113.9", Invalid: true},
			matched: true,
		},
		{
			name:    "IPv6 source",
			msg:     "Failed password for root from 2001:db8::1 port 40404 ssh2",
			want:    Event{Kind: FailedAuth, Method: "password", User: "root", IP: "2001:db8::1"},
			matched: true,
		},
		{
			name: "unrelated sshd chatter",
			msg:  "Server listening on 0.0.0.0 port 22.",
		},
		{
			name: "session teardown is not an auth event",
			msg:  "pam_unix(sshd:session): session closed for user deploy",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := classify(tc.msg)
			if ok != tc.matched {
				t.Fatalf("classify(%q) matched = %v, want %v", tc.msg, ok, tc.matched)
			}
			if !tc.matched {
				return
			}
			if got != tc.want {
				t.Errorf("classify(%q)\n got %+v\nwant %+v", tc.msg, got, tc.want)
			}
		})
	}
}

// sshd logs the username exactly as the client sent it, so every pattern is
// anchored at the start of the message: a keyword appearing inside a username
// must not decide how the line is classified, or an attacker could write
// entries into this report attributed to an address of their choosing.
func TestClassifyIsAnchoredAtTheStartOfTheMessage(t *testing.T) {
	// An account literally named "Failed".
	got, ok := classify("Invalid user Failed from 203.0.113.7 port 40404")
	if !ok {
		t.Fatal("classify() skipped a line it should have matched")
	}
	if got.Kind != ProbedUser {
		t.Errorf("Kind = %v, want ProbedUser — the username's 'Failed' won", got.Kind)
	}
	if got.IP != "203.0.113.7" {
		t.Errorf("IP = %q, want 203.0.113.7", got.IP)
	}

	// An account literally named "Accepted" must not read as a login.
	got, ok = classify("Failed password for Accepted from 203.0.113.7 port 40404 ssh2")
	if !ok {
		t.Fatal("classify() skipped a line it should have matched")
	}
	if got.Kind != FailedAuth {
		t.Errorf("Kind = %v, want FailedAuth — a rejected login was read as a success", got.Kind)
	}
}

// The username is a single non-space token and the address is the token after
// the first " from ", so a crafted username cannot push the parser onto a
// later address and attribute its own attempts to someone else.
func TestClassifyTakesTheFirstAddressNotTheLast(t *testing.T) {
	msg := "Invalid user root from 10.0.0.1 port 22 ssh2 from 203.0.113.7 port 40404"

	got, ok := classify(msg)
	if !ok {
		t.Fatal("classify() skipped a line it should have matched")
	}
	if got.IP != "10.0.0.1" {
		t.Errorf("IP = %q, want 10.0.0.1 — the parser walked to a later address", got.IP)
	}
	if got.User != "root" {
		t.Errorf("User = %q, want root", got.User)
	}
}

func TestSplitMessage(t *testing.T) {
	cases := []struct {
		line       string
		wantMsg    string
		wantPrefix string
		wantOK     bool
	}{
		{
			line:       "2026-08-16T09:14:23+0300 web sshd[1234]: Failed password for root from 203.0.113.7 port 1 ssh2",
			wantPrefix: "2026-08-16T09:14:23+0300 web sshd[1234",
			wantMsg:    "Failed password for root from 203.0.113.7 port 1 ssh2",
			wantOK:     true,
		},
		{
			line:       "Aug 16 09:14:23 web sshd[1234]: Invalid user admin from 203.0.113.7",
			wantPrefix: "Aug 16 09:14:23 web sshd[1234",
			wantMsg:    "Invalid user admin from 203.0.113.7",
			wantOK:     true,
		},
		{
			// Some builds log without a pid.
			line:       "Aug 16 09:14:23 web sshd: Server listening on 0.0.0.0 port 22.",
			wantPrefix: "Aug 16 09:14:23 web sshd",
			wantMsg:    "Server listening on 0.0.0.0 port 22.",
			wantOK:     true,
		},
		{line: "a line with no tag at all"},
	}

	for _, tc := range cases {
		prefix, msg, ok := splitMessage(tc.line)
		if ok != tc.wantOK {
			t.Errorf("splitMessage(%q) ok = %v, want %v", tc.line, ok, tc.wantOK)
			continue
		}
		if !ok {
			continue
		}
		if prefix != tc.wantPrefix {
			t.Errorf("splitMessage(%q) prefix = %q, want %q", tc.line, prefix, tc.wantPrefix)
		}
		if msg != tc.wantMsg {
			t.Errorf("splitMessage(%q) msg = %q, want %q", tc.line, msg, tc.wantMsg)
		}
	}
}

// Only sshd's lines count. Parse filters on the tag, so a sudo or cron line
// mentioning a failure can't inflate the numbers.
func TestParseIgnoresOtherDaemons(t *testing.T) {
	lines := []string{
		"2026-08-16T09:14:23+0300 web sudo[99]: Failed password for root from 203.0.113.7 port 1 ssh2",
		"2026-08-16T09:14:24+0300 web sshd[1234]: Failed password for root from 203.0.113.7 port 1 ssh2",
	}

	events := Parse(lines, time.Time{})
	if len(events) != 1 {
		t.Fatalf("Parse() returned %d events, want 1 (the sudo line must be ignored)", len(events))
	}
}

func TestParseDropsEventsOlderThanTheCutoff(t *testing.T) {
	lines := []string{
		"2026-08-10T09:14:23+0300 web sshd[1]: Failed password for root from 203.0.113.7 port 1 ssh2",
		"2026-08-16T09:14:23+0300 web sshd[1]: Failed password for root from 203.0.113.8 port 1 ssh2",
	}
	cutoff := time.Date(2026, 8, 15, 0, 0, 0, 0, time.UTC)

	events := Parse(lines, cutoff)
	if len(events) != 1 {
		t.Fatalf("Parse() returned %d events, want 1", len(events))
	}
	if events[0].IP != "203.0.113.8" {
		t.Errorf("kept the wrong event: %+v", events[0])
	}
}

// A line whose timestamp can't be read is a gap in our parsing, not proof the
// event falls outside the window. Dropping it would understate the report.
func TestParseKeepsEventsWithUnreadableTimestamps(t *testing.T) {
	lines := []string{"garbled-timestamp web sshd[1]: Failed password for root from 203.0.113.7 port 1 ssh2"}

	events := Parse(lines, time.Now().Add(-time.Hour))
	if len(events) != 1 {
		t.Fatalf("Parse() returned %d events, want 1", len(events))
	}
	if !events[0].Time.IsZero() {
		t.Errorf("Time = %v, want zero to mark it unknown", events[0].Time)
	}
}

func TestParseTimestampFormats(t *testing.T) {
	now := time.Date(2026, 8, 16, 12, 0, 0, 0, time.UTC)

	iso := parseTimestamp("2026-08-16T09:14:23+0000 web sshd[1]: x", now)
	if iso.Year() != 2026 || iso.Month() != time.August || iso.Day() != 16 {
		t.Errorf("short-iso parsed as %v", iso)
	}

	// Syslog omits the year; the current one is the right assumption.
	syslog := parseTimestamp("Aug 16 09:14:23 web sshd[1]: x", now)
	if syslog.Year() != 2026 || syslog.Month() != time.August || syslog.Day() != 16 {
		t.Errorf("syslog stamp parsed as %v, want 2026-08-16", syslog)
	}

	// Single-digit days are space-padded to two columns in syslog.
	padded := parseTimestamp("Aug  6 09:14:23 web sshd[1]: x", now)
	if padded.Day() != 6 {
		t.Errorf("space-padded day parsed as %v", padded)
	}

	// December logs read in January would otherwise land a year in the future.
	january := time.Date(2027, 1, 2, 12, 0, 0, 0, time.UTC)
	rolled := parseTimestamp("Dec 28 09:14:23 web sshd[1]: x", january)
	if rolled.Year() != 2026 {
		t.Errorf("year-boundary stamp parsed as %v, want 2026", rolled)
	}
}

func TestSummarizeAggregatesBySource(t *testing.T) {
	events := []Event{
		{Kind: FailedAuth, IP: "203.0.113.7", User: "root"},
		{Kind: FailedAuth, IP: "203.0.113.7", User: "root"},
		{Kind: ProbedUser, IP: "203.0.113.7", User: "admin", Invalid: true},
		{Kind: Probe, IP: "203.0.113.9"},
		{Kind: Accepted, IP: "198.51.100.4", User: "deploy", Method: "publickey"},
	}

	r := Summarize(events)

	if r.Failed != 2 || r.Probed != 2 || r.Accepted != 1 {
		t.Errorf("totals: failed=%d probed=%d accepted=%d, want 2/2/1", r.Failed, r.Probed, r.Accepted)
	}
	if r.RootHits != 2 {
		t.Errorf("RootHits = %d, want 2", r.RootHits)
	}
	if len(r.Sources) != 3 {
		t.Fatalf("got %d sources, want 3", len(r.Sources))
	}
	if r.Sources[0].IP != "203.0.113.7" {
		t.Errorf("busiest source = %q, want 203.0.113.7", r.Sources[0].IP)
	}
	if got := r.Sources[0].Attempts(); got != 3 {
		t.Errorf("Attempts() = %d, want 3", got)
	}
	if len(r.Logins) != 1 || r.Logins[0].User != "deploy" {
		t.Errorf("Logins = %+v, want one login by deploy", r.Logins)
	}
	// Successful logins aren't attempts on an account, so they stay out of
	// the "most-targeted" ranking.
	for _, u := range r.Usernames {
		if u.Name == "deploy" {
			t.Error("a successful login was counted as a targeted account")
		}
	}
}

// Volume doesn't say whether an attack is still running, so each source
// carries the newest timestamp seen from it.
func TestSummarizeTracksMostRecentActivityPerSource(t *testing.T) {
	early := time.Date(2026, 8, 16, 9, 0, 0, 0, time.UTC)
	late := time.Date(2026, 8, 16, 11, 0, 0, 0, time.UTC)

	r := Summarize([]Event{
		{Kind: FailedAuth, IP: "203.0.113.7", User: "root", Time: late},
		{Kind: FailedAuth, IP: "203.0.113.7", User: "root", Time: early},
	})

	if !r.Sources[0].Last.Equal(late) {
		t.Errorf("Last = %v, want the newer event at %v", r.Sources[0].Last, late)
	}
}

// An address that failed and then succeeded is the one finding that can mean
// a live intrusion, so it must sort to the top however quiet it was.
func TestSummarizeRanksSuccessfulBruteForceFirst(t *testing.T) {
	events := []Event{
		{Kind: FailedAuth, IP: "203.0.113.7", User: "root"},
		{Kind: Accepted, IP: "203.0.113.7", User: "root"},
	}
	for range 500 {
		events = append(events, Event{Kind: Probe, IP: "198.51.100.99"})
	}

	r := Summarize(events)

	if r.Sources[0].IP != "203.0.113.7" {
		t.Errorf("top source = %q, want the address that got in", r.Sources[0].IP)
	}
	if !r.Sources[0].GotIn() {
		t.Error("GotIn() = false for an address that failed then succeeded")
	}
	if r.Sources[1].GotIn() {
		t.Error("GotIn() = true for an address that never authenticated")
	}
}

// The whole point of a distributed attack is that no single address looks
// worth acting on. Summarize has to surface it at the level it exists —
// across addresses — or the report reads as "nothing to see" during one.
func TestSummarizeCountsSpraySourcesPerAccount(t *testing.T) {
	var events []Event
	// 30 addresses, two tries at root each: no address near any rate limit.
	for i := range 30 {
		ip := fmt.Sprintf("203.0.113.%d", i+1)
		events = append(events,
			Event{Kind: FailedAuth, IP: ip, User: "root"},
			Event{Kind: FailedAuth, IP: ip, User: "root"},
		)
	}
	// One noisy address on a different account, for contrast.
	for range 40 {
		events = append(events, Event{Kind: FailedAuth, IP: "198.51.100.9", User: "admin"})
	}

	r := Summarize(events)

	sprayed, ok := r.Sprayed()
	if !ok {
		t.Fatal("Sprayed() found nothing")
	}
	if sprayed.Name != "root" {
		t.Errorf("Sprayed() = %q, want root — the busiest account by attempts won instead", sprayed.Name)
	}
	if sprayed.Sources != 30 {
		t.Errorf("root Sources = %d, want 30", sprayed.Sources)
	}

	// "admin" has more attempts from one host; it must not read as a spray.
	for _, u := range r.Usernames {
		if u.Name == "admin" && u.Sources != 1 {
			t.Errorf("admin Sources = %d, want 1", u.Sources)
		}
	}
}

// Distributed measures what a per-source limit would have missed: addresses
// individually below the threshold, and their combined volume.
func TestDistributedCountsOnlyQuietSources(t *testing.T) {
	var events []Event
	for i := range 25 {
		ip := fmt.Sprintf("203.0.113.%d", i+1)
		events = append(events, Event{Kind: FailedAuth, IP: ip, User: "root"})
		events = append(events, Event{Kind: FailedAuth, IP: ip, User: "root"})
	}
	// One loud address, which a rate limit *would* have caught, so it must
	// not be counted among the ones that slipped through.
	for range 60 {
		events = append(events, Event{Kind: FailedAuth, IP: "198.51.100.9", User: "root"})
	}

	sources, attempts := Summarize(events).Distributed(10)

	if sources != 25 {
		t.Errorf("quiet sources = %d, want 25 (the loud address must be excluded)", sources)
	}
	if attempts != 50 {
		t.Errorf("quiet attempts = %d, want 50", attempts)
	}
}

func TestSubnetOf(t *testing.T) {
	cases := []struct{ in, want string }{
		{in: "203.0.113.7", want: "203.0.113.0/24"},
		{in: "203.0.113.200", want: "203.0.113.0/24"},
		{in: "10.1.2.3", want: "10.1.2.0/24"},
		{in: "2001:db8:0:1::5", want: "2001:db8::/48"},
		{in: "2001:db8:1:1::5", want: "2001:db8:1::/48"},
		// A mapped IPv4 peer must group with its IPv4 block, not into ::/48.
		{in: "::ffff:203.0.113.7", want: "203.0.113.0/24"},
		// `UseDNS yes` logs a hostname, which belongs to no block we can name.
		{in: "not-an-address", want: ""},
		{in: "", want: ""},
	}

	for _, tc := range cases {
		if got := subnetOf(tc.in); got != tc.want {
			t.Errorf("subnetOf(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

// A block holding several addresses is one operator worth denying wholesale;
// a block holding one says nothing the address didn't.
func TestSummarizeClustersSourcesIntoSubnets(t *testing.T) {
	var events []Event
	for i := range 8 {
		events = append(events, Event{
			Kind: FailedAuth,
			IP:   fmt.Sprintf("203.0.113.%d", i+1),
			User: "root",
		})
	}
	events = append(events, Event{Kind: FailedAuth, IP: "198.51.100.9", User: "root"})

	r := Summarize(events)

	if len(r.Subnets) != 1 {
		t.Fatalf("got %d clustered subnets, want 1 (the lone address must not form one): %+v",
			len(r.Subnets), r.Subnets)
	}
	if r.Subnets[0].CIDR != "203.0.113.0/24" {
		t.Errorf("CIDR = %q, want 203.0.113.0/24", r.Subnets[0].CIDR)
	}
	if r.Subnets[0].Sources != 8 || r.Subnets[0].Attempts != 8 {
		t.Errorf("subnet = %+v, want 8 sources and 8 attempts", r.Subnets[0])
	}
}

// Two blocks can hold a similar number of addresses and differ by three
// orders of magnitude in what they did. Ranking by address count printed a
// sustained attack and a passing sweep as neighbours; only one is worth
// blocking, so intensity decides the order.
func TestSubnetsRankSustainedBlocksAboveSweeps(t *testing.T) {
	var events []Event

	// A sweep: 20 addresses, one probe each, on its way through the internet.
	for i := range 20 {
		events = append(events, Event{Kind: Probe, IP: fmt.Sprintf("198.51.100.%d", i+1)})
	}
	// A sustained run: 3 addresses, 200 attempts each.
	for i := range 3 {
		ip := fmt.Sprintf("203.0.113.%d", i+1)
		for range 200 {
			events = append(events, Event{Kind: FailedAuth, IP: ip, User: "root"})
		}
	}

	r := Summarize(events)

	if len(r.Subnets) != 2 {
		t.Fatalf("got %d subnets, want 2: %+v", len(r.Subnets), r.Subnets)
	}
	if r.Subnets[0].CIDR != "203.0.113.0/24" {
		t.Errorf("top block = %q, want 203.0.113.0/24 — a 20-address sweep outranked a sustained run",
			r.Subnets[0].CIDR)
	}
	if got := r.Subnets[0].Intensity(); got != 200 {
		t.Errorf("sustained Intensity() = %v, want 200", got)
	}
	if got := r.Subnets[1].Intensity(); got != 1 {
		t.Errorf("sweep Intensity() = %v, want 1", got)
	}
}

func TestSubnetIntensityHandlesEmptyBlock(t *testing.T) {
	if got := (Subnet{CIDR: "203.0.113.0/24"}).Intensity(); got != 0 {
		t.Errorf("Intensity() with no sources = %v, want 0", got)
	}
}

// A successful login is not an attack, so it must not inflate a block's
// numbers and get an operator's own range recommended for a deny rule.
func TestSubnetsIgnoreSuccessfulLogins(t *testing.T) {
	r := Summarize([]Event{
		{Kind: Accepted, IP: "203.0.113.1", User: "deploy"},
		{Kind: Accepted, IP: "203.0.113.2", User: "deploy"},
	})

	if len(r.Subnets) != 0 {
		t.Errorf("logins produced subnet clusters: %+v", r.Subnets)
	}
}

func TestParseWindow(t *testing.T) {
	cases := []struct {
		in   string
		want time.Duration
		ok   bool
	}{
		{in: "24h", want: 24 * time.Hour, ok: true},
		{in: "90m", want: 90 * time.Minute, ok: true},
		{in: "7d", want: 7 * 24 * time.Hour, ok: true},
		{in: " 1d ", want: 24 * time.Hour, ok: true},
		{in: "0h"},
		{in: "-2h"},
		{in: "0d"},
		{in: "week"},
		{in: ""},
	}

	for _, tc := range cases {
		got, err := ParseWindow(tc.in)
		if tc.ok {
			if err != nil {
				t.Errorf("ParseWindow(%q) errored: %v", tc.in, err)
			} else if got != tc.want {
				t.Errorf("ParseWindow(%q) = %v, want %v", tc.in, got, tc.want)
			}
			continue
		}
		if err == nil {
			t.Errorf("ParseWindow(%q) = %v, want an error", tc.in, got)
		}
	}
}
