package sshd

// Authentication-log auditing. Everything in this file reads logs and
// computes; nothing here changes the system.
//
// The parsing is deliberately tolerant. sshd's messages differ between
// OpenSSH releases and distro patches, and syslog and journalctl disagree
// about how a line begins — so a line that doesn't match any known shape is
// skipped rather than treated as an error. An audit that reports nine of ten
// attackers is useful; one that aborts on an unfamiliar log line is not.

import (
	"fmt"
	"net/netip"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Kind classifies what an authentication log line represents.
type Kind int

const (
	// FailedAuth: a credential was offered for a real account and rejected.
	// The only kind that could have succeeded with a luckier guess.
	FailedAuth Kind = iota

	// ProbedUser: sshd rejected the username outright because no such
	// account exists here. Nothing offered for it could ever have worked,
	// so these are scanning, not a threat to any particular account.
	ProbedUser

	// Probe: a connection that opened and hung up before authenticating.
	// Port scanners and banner grabbers look exactly like this.
	Probe

	// Accepted: someone logged in. The lines that actually matter.
	Accepted
)

// Event is one parsed log line.
type Event struct {
	Time    time.Time // zero when the line carried no parseable timestamp
	IP      string
	User    string // "" when the line names no account
	Method  string // "password", "publickey", "" when the line doesn't say
	Invalid bool   // the account named does not exist on this host
	Kind    Kind
}

// Source aggregates every event from one address.
type Source struct {
	IP       string
	Failed   int
	Probed   int // ProbedUser + Probe
	Accepted int
	Users    []string  // distinct accounts named, most-seen first
	Last     time.Time // newest event from this address; zero if none had a readable timestamp
}

// Attempts is the total unwanted activity from this address.
func (s Source) Attempts() int { return s.Failed + s.Probed }

// GotIn reports an address that failed at least once and then succeeded —
// the shape of a brute-force that worked, and the single most important
// thing this audit can surface.
func (s Source) GotIn() bool { return s.Accepted > 0 && s.Failed > 0 }

// UserCount is one targeted account: how often it was named, and by how many
// distinct addresses.
//
// The second number is the one that matters. An account tried 600 times from
// one address is a single bot that a rate limit will stop; 600 tries spread
// over 300 addresses is the same attack rebuilt so that no per-source control
// can see it.
type UserCount struct {
	Name    string
	Count   int
	Sources int
}

// Subnet groups source addresses by network block — /24 for IPv4, /48 for
// IPv6. Botnet nodes and rented ranges cluster, so a block holding forty
// addresses is usually one operator rather than forty.
type Subnet struct {
	CIDR     string
	Sources  int
	Attempts int
}

// Intensity is attempts per address: how hard the block's hosts are working,
// rather than how many of them there are.
//
// Address count alone cannot rank these. Two blocks holding seven addresses
// each can differ by three orders of magnitude in what they did — one
// operator running a sustained attack from seven hosts, versus an
// internet-wide sweep that touched seven addresses in the block once apiece
// on its way through every address on the internet. Ranking by address count
// printed those as neighbours; only intensity separates them, and only the
// first is worth blocking, since the sweep's hosts already moved on.
func (s Subnet) Intensity() float64 {
	if s.Sources == 0 {
		return 0
	}
	return float64(s.Attempts) / float64(s.Sources)
}

// Report is the finished audit.
type Report struct {
	LogSource string // where the lines came from, for the header
	Lines     int    // log lines read
	Events    int    // lines recognised as authentication activity

	Failed   int
	Probed   int
	Accepted int
	RootHits int // events naming "root", successful or not

	Sources   []Source    // most active first
	Subnets   []Subnet    // blocks holding more than one source, most populated first
	Logins    []Event     // Accepted events, oldest first
	Usernames []UserCount // accounts targeted, most-tried first
}

// Distributed reports the activity a per-source limit would never have
// caught: addresses that each stayed under threshold attempts, and how many
// attempts they made between them.
//
// This is the shape of a distributed brute force, and the reason rate
// limiting is not an answer to one. Every address looks harmless on its own;
// only the total is alarming, and no control that decides per address can see
// the total. The threshold is in attempts while ufw's limit counts
// connections, so this is an estimate of what would slip through, not a
// simulation of it.
func (r *Report) Distributed(threshold int) (sources, attempts int) {
	for _, s := range r.Sources {
		if n := s.Attempts(); n > 0 && n < threshold {
			sources++
			attempts += n
		}
	}
	return sources, attempts
}

// Sprayed returns the account tried from the greatest number of distinct
// addresses. One account named by many hosts is credential spraying — a
// coordinated attempt on a known login rather than untargeted scanning.
func (r *Report) Sprayed() (UserCount, bool) {
	var worst UserCount
	for _, u := range r.Usernames {
		if u.Sources > worst.Sources {
			worst = u
		}
	}
	return worst, worst.Sources > 0
}

// ParseWindow accepts a Go duration ("90m", "12h") and additionally a day
// suffix ("7d"), which time.ParseDuration rejects but is the unit people
// reach for when looking back over logs.
func ParseWindow(s string) (time.Duration, error) {
	s = strings.TrimSpace(s)
	invalid := fmt.Errorf("invalid window %q — try 24h, 90m or 7d", s)

	if rest, ok := strings.CutSuffix(s, "d"); ok {
		days, err := strconv.Atoi(rest)
		if err != nil || days <= 0 {
			return 0, invalid
		}
		return time.Duration(days) * 24 * time.Hour, nil
	}

	d, err := time.ParseDuration(s)
	if err != nil || d <= 0 {
		return 0, invalid
	}
	return d, nil
}

// logFiles are the syslog destinations sshd writes to, by distro family.
var logFiles = []string{"/var/log/auth.log", "/var/log/secure"}

// Collect gathers authentication log lines from the last window.
//
// journalctl comes first because it filters by time itself and survives log
// rotation, but a plain syslog file is used when there's no journal, when the
// journal has no ssh unit, or when the caller names a file outright.
func Collect(window time.Duration, file string) ([]string, string, error) {
	if file != "" {
		lines, err := readLogFile(file)
		return lines, file, err
	}

	if lines, ok := journalLines(window); ok {
		return lines, "journalctl -u ssh -u sshd", nil
	}

	for _, path := range logFiles {
		if _, err := os.Stat(path); err != nil {
			continue
		}
		lines, err := readLogFile(path)
		return lines, path, err
	}

	return nil, "", fmt.Errorf(
		"no SSH authentication log found — tried journalctl and %s\n"+
			"point at one directly with --file", strings.Join(logFiles, ", "))
}

func journalLines(window time.Duration) ([]string, bool) {
	if _, err := exec.LookPath("journalctl"); err != nil {
		return nil, false
	}

	// An absolute timestamp rather than a relative one: systemd's relative
	// forms vary across versions, and a --since journalctl doesn't understand
	// is rejected outright rather than widened.
	since := time.Now().Add(-window).Format("2006-01-02 15:04:05")

	out, err := exec.Command("journalctl",
		"-u", "ssh", "-u", "sshd",
		"--since", since,
		"--no-pager", "-o", "short-iso").Output()
	if err != nil {
		return nil, false
	}

	lines := splitLines(string(out))
	// A host whose sshd logs only to syslog has the units but no entries.
	// Reporting "0 attempts" there would be a lie by omission, so fall back.
	if len(lines) == 0 {
		return nil, false
	}
	return lines, true
}

func readLogFile(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	return splitLines(string(data)), nil
}

func splitLines(s string) []string {
	var out []string
	for _, line := range strings.Split(s, "\n") {
		line = strings.TrimSpace(line)
		// journalctl's own "-- Logs begin at ... --" banners aren't log lines.
		if line == "" || strings.HasPrefix(line, "-- ") {
			continue
		}
		out = append(out, line)
	}
	return out
}

// The message shapes OpenSSH emits around authentication. Anchored at the
// start of the message so a username or hostname containing "Failed" can't
// forge a match.
var (
	reAccepted = regexp.MustCompile(`^Accepted (\S+) for (\S+) from (\S+) port \d+`)
	reFailed   = regexp.MustCompile(`^Failed (\S+) for (invalid user )?(\S+) from (\S+) port \d+`)
	reMaxTries = regexp.MustCompile(`^maximum authentication attempts exceeded for (invalid user )?(\S+) from (\S+) port \d+`)
	reInvalid  = regexp.MustCompile(`^Invalid user (\S*) from (\S+)`)
	rePreauth  = regexp.MustCompile(`^(?:Connection closed by|Disconnected from) (?:(invalid|authenticating) user (\S+) )?(\S+) port \d+`)
)

// classify turns an sshd message into an Event, or reports that the line
// isn't about authentication.
//
// A single failed login produces several lines — "Invalid user bob", then a
// "Failed password" per guess, then "Connection closed by invalid user bob".
// They are counted on different axes on purpose: Failed measures credentials
// offered, Probed measures accounts and connections tried. Neither is a
// connection count, and the report says so.
func classify(msg string) (Event, bool) {
	for _, prefix := range []string{"error: ", "fatal: ", "PAM: "} {
		msg = strings.TrimPrefix(msg, prefix)
	}

	if m := reAccepted.FindStringSubmatch(msg); m != nil {
		return Event{Kind: Accepted, Method: m[1], User: m[2], IP: m[3]}, true
	}
	if m := reFailed.FindStringSubmatch(msg); m != nil {
		return Event{Kind: FailedAuth, Method: m[1], User: m[3], IP: m[4], Invalid: m[2] != ""}, true
	}
	if m := reMaxTries.FindStringSubmatch(msg); m != nil {
		return Event{Kind: FailedAuth, User: m[2], IP: m[3], Invalid: m[1] != ""}, true
	}
	if m := reInvalid.FindStringSubmatch(msg); m != nil {
		return Event{Kind: ProbedUser, User: m[1], IP: m[2], Invalid: true}, true
	}
	if m := rePreauth.FindStringSubmatch(msg); m != nil {
		return Event{Kind: Probe, User: m[2], IP: m[3], Invalid: m[1] == "invalid"}, true
	}
	return Event{}, false
}

// splitMessage separates a log line's prefix (timestamp, host, tag) from the
// message the daemon wrote.
func splitMessage(line string) (prefix, msg string, ok bool) {
	// "host sshd[1234]: Failed password ..." — the first "]: " is the tag's,
	// since a message body can only contribute a later one.
	if i := strings.Index(line, "]: "); i >= 0 {
		return line[:i], line[i+len("]: "):], true
	}
	// Some builds log without a pid: "host sshd: Failed password ...".
	if i := strings.Index(line, "sshd: "); i >= 0 {
		return line[:i+len("sshd")], line[i+len("sshd: "):], true
	}
	return "", "", false
}

// parseTimestamp reads the leading timestamp in whichever of the two formats
// the line uses. A zero return means the time is unknown, which callers treat
// as "include it" rather than dropping the event.
func parseTimestamp(line string, now time.Time) time.Time {
	fields := strings.Fields(line)
	if len(fields) == 0 {
		return time.Time{}
	}

	// journalctl -o short-iso: 2026-08-16T09:14:23+0300
	for _, layout := range []string{"2006-01-02T15:04:05-0700", time.RFC3339} {
		if t, err := time.Parse(layout, fields[0]); err == nil {
			return t
		}
	}

	// Classic syslog: "Aug 16 09:14:23" — no year at all. Assume the current
	// one, and step back if that lands in the future, which is what a
	// December log read in January looks like.
	if len(fields) >= 3 {
		stamp := strings.Join(fields[:3], " ")
		if t, err := time.ParseInLocation("Jan 2 15:04:05", stamp, time.Local); err == nil {
			t = t.AddDate(now.Year(), 0, 0)
			if t.Sub(now) > 25*time.Hour {
				t = t.AddDate(-1, 0, 0)
			}
			return t
		}
	}

	return time.Time{}
}

// Parse turns raw log lines into events, dropping anything older than cutoff.
// Events whose timestamp couldn't be read are kept: a missing timestamp is a
// parsing gap, and silently discarding activity would understate the report.
func Parse(lines []string, cutoff time.Time) []Event {
	now := time.Now()
	var events []Event

	for _, line := range lines {
		prefix, msg, ok := splitMessage(line)
		if !ok || !strings.Contains(prefix, "sshd") {
			continue
		}
		event, ok := classify(msg)
		if !ok {
			continue
		}
		event.Time = parseTimestamp(line, now)
		if !event.Time.IsZero() && event.Time.Before(cutoff) {
			continue
		}
		events = append(events, event)
	}

	return events
}

// Summarize aggregates events by source address and by targeted account.
func Summarize(events []Event) *Report {
	r := &Report{Events: len(events)}

	sources := map[string]*Source{}
	userSeen := map[string]map[string]int{} // ip → account → times named
	targeted := map[string]int{}
	// Both of these are counted across addresses rather than within one,
	// because that's the only level at which a distributed attack exists.
	targetedBy := map[string]map[string]bool{} // account → addresses that named it
	blocks := map[string]*Subnet{}
	blockSeen := map[string]map[string]bool{} // block → addresses seen in it

	for _, e := range events {
		switch e.Kind {
		case Accepted:
			r.Accepted++
			r.Logins = append(r.Logins, e)
		case FailedAuth:
			r.Failed++
		default:
			r.Probed++
		}

		if e.User == "root" {
			r.RootHits++
		}
		if e.User != "" && e.Kind != Accepted {
			targeted[e.User]++
			if e.IP != "" {
				if targetedBy[e.User] == nil {
					targetedBy[e.User] = map[string]bool{}
				}
				targetedBy[e.User][e.IP] = true
			}
		}
		if e.IP == "" {
			continue
		}

		s := sources[e.IP]
		if s == nil {
			s = &Source{IP: e.IP}
			sources[e.IP] = s
			userSeen[e.IP] = map[string]int{}
		}

		switch e.Kind {
		case Accepted:
			s.Accepted++
		case FailedAuth:
			s.Failed++
		default:
			s.Probed++
		}
		if e.User != "" {
			userSeen[e.IP][e.User]++
		}
		if e.Time.After(s.Last) {
			s.Last = e.Time
		}

		if cidr := subnetOf(e.IP); cidr != "" && e.Kind != Accepted {
			b := blocks[cidr]
			if b == nil {
				b = &Subnet{CIDR: cidr}
				blocks[cidr] = b
				blockSeen[cidr] = map[string]bool{}
			}
			b.Attempts++
			if !blockSeen[cidr][e.IP] {
				blockSeen[cidr][e.IP] = true
				b.Sources++
			}
		}
	}

	for ip, s := range sources {
		s.Users = rankKeys(userSeen[ip])
		r.Sources = append(r.Sources, *s)
	}

	// A block holding one address says nothing that the address itself
	// didn't; only clustering is new information.
	for _, b := range blocks {
		if b.Sources > 1 {
			r.Subnets = append(r.Subnets, *b)
		}
	}
	// Hardest-working first, not most-populated first — see Intensity.
	sort.Slice(r.Subnets, func(i, j int) bool {
		a, b := r.Subnets[i], r.Subnets[j]
		if a.Intensity() != b.Intensity() {
			return a.Intensity() > b.Intensity()
		}
		if a.Attempts != b.Attempts {
			return a.Attempts > b.Attempts
		}
		if a.Sources != b.Sources {
			return a.Sources > b.Sources
		}
		return a.CIDR < b.CIDR
	})

	// Busiest first. An address that got in outranks a noisier one that
	// didn't, since a successful login is the finding worth reading.
	sort.Slice(r.Sources, func(i, j int) bool {
		a, b := r.Sources[i], r.Sources[j]
		if a.GotIn() != b.GotIn() {
			return a.GotIn()
		}
		if a.Attempts() != b.Attempts() {
			return a.Attempts() > b.Attempts()
		}
		return a.IP < b.IP
	})

	for _, name := range rankKeys(targeted) {
		r.Usernames = append(r.Usernames, UserCount{
			Name:    name,
			Count:   targeted[name],
			Sources: len(targetedBy[name]),
		})
	}

	sort.Slice(r.Logins, func(i, j int) bool { return r.Logins[i].Time.Before(r.Logins[j].Time) })

	return r
}

// subnetOf returns the network block an address belongs to — /24 for IPv4,
// /48 for IPv6, the smallest units usually assigned to one operator. Returns
// "" for anything that isn't an address, which is what sshd logs when
// `UseDNS yes` resolves the client to a hostname.
func subnetOf(ip string) string {
	addr, err := netip.ParseAddr(ip)
	if err != nil {
		return ""
	}
	addr = addr.Unmap()

	bits := 48
	if addr.Is4() {
		bits = 24
	}
	prefix, err := addr.Prefix(bits)
	if err != nil {
		return ""
	}
	return prefix.String()
}

// rankKeys orders map keys by count descending, then by name so that equal
// counts don't reshuffle between runs.
func rankKeys(counts map[string]int) []string {
	keys := make([]string, 0, len(counts))
	for k := range counts {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		if counts[keys[i]] != counts[keys[j]] {
			return counts[keys[i]] > counts[keys[j]]
		}
		return keys[i] < keys[j]
	})
	return keys
}

// Audit collects and summarizes SSH authentication activity over window.
// Pass file to read a specific log instead of letting Collect choose.
func Audit(window time.Duration, file string) (*Report, error) {
	lines, source, err := Collect(window, file)
	if err != nil {
		return nil, err
	}

	report := Summarize(Parse(lines, time.Now().Add(-window)))
	report.LogSource = source
	report.Lines = len(lines)
	return report, nil
}
