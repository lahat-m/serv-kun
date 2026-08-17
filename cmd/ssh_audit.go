package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/lahat-m/serv-kun/internal/sshd"
	"github.com/spf13/cobra"
)

// Thresholds for turning counts into findings. A distributed attack is
// defined by what it isn't — no single address busy enough to notice — so it
// takes the last three to describe rather than the first.
const (
	noisyThreshold     = 50 // one address this busy is worth rate-limiting
	quietThreshold     = 10 // below this, no per-source limit would ever react
	distributedSources = 20 // that many quiet addresses is coordination, not coincidence
	sprayThreshold     = 10 // one account named from this many addresses
	clusterThreshold   = 5  // addresses from a single /24 or /48
)

var sshAuditCmd = &cobra.Command{
	Use:   "audit",
	Short: "Report SSH authentication attempts from the system logs",
	Long: `Reads the SSH authentication log and summarises who has been trying to
log in: which addresses, which accounts they named, and who actually got in.

Read-only — this command changes nothing. It reads the systemd journal where
there is one, falling back to /var/log/auth.log or /var/log/secure.

Findings are weighed against the running sshd config, since thousands of
rejected passwords mean something very different on a server that accepts
passwords than on one that doesn't.

Activity is also totalled across addresses, not only per address. A run split
over enough hosts keeps every one of them below any rate limit, so a per-source
view of it looks like nothing at all.`,
	Args: cobra.NoArgs,
	RunE: runSSHAudit,
}

func init() {
	sshAuditCmd.Flags().String("since", "24h", "how far back to look (e.g. 90m, 24h, 7d)")
	sshAuditCmd.Flags().Int("top", 10, "how many source addresses to list")
	sshAuditCmd.Flags().String("file", "", "read this log file instead of auto-detecting")
}

func runSSHAudit(cmd *cobra.Command, _ []string) error {
	// Both the journal and /var/log/auth.log are root-only on a stock
	// install, so a non-root run would report zero attempts rather than fail
	// — the most misleading possible output for a security check.
	requireRoot(cmd)

	sinceFlag, _ := cmd.Flags().GetString("since")
	top, _ := cmd.Flags().GetInt("top")
	file, _ := cmd.Flags().GetString("file")

	window, err := sshd.ParseWindow(sinceFlag)
	if err != nil {
		return err
	}

	printBanner("ssh audit")

	report, err := sshd.Audit(window, file)
	if err != nil {
		return err
	}

	fmt.Printf("Source: %s\n", report.LogSource)
	fmt.Printf("Window: last %s\n", sinceFlag)
	fmt.Printf("Read:   %d log lines, %d of them authentication events\n\n", report.Lines, report.Events)

	// The sshd config decides what the numbers below mean, so read it before
	// reporting them. A failure here is not fatal: the log summary stands on
	// its own, only the assessment needs the settings.
	settings, settingsErr := sshd.EffectiveSettings()

	printSummary(report)
	printSources(report, top)
	printSubnets(report)
	printTargetedUsers(report)
	printLogins(report)
	printAssessment(report, settings, settingsErr)

	return nil
}

func printSummary(r *sshd.Report) {
	printRule()
	fmt.Println("Summary")
	fmt.Printf("  Rejected credentials    %d\n", r.Failed)
	fmt.Printf("  Probes and scans        %d\n", r.Probed)
	fmt.Printf("  Successful logins       %d\n", r.Accepted)
	fmt.Printf("  Events naming root      %d\n", r.RootHits)
	fmt.Printf("  Distinct addresses      %d\n", len(r.Sources))
	printRule()
	fmt.Println()

	if r.Events == 0 {
		return
	}
	// Without this, a reader divides the totals by each other and concludes
	// the wrong thing: one login attempt writes several lines, and the two
	// counters above measure different axes of the same connection.
	fmt.Println("One attempt writes several log lines — read these as relative volume,")
	fmt.Println("not as a count of connections.")
	fmt.Println()
}

func printSources(r *sshd.Report, top int) {
	if len(r.Sources) == 0 {
		return
	}

	shown := r.Sources
	if top > 0 && len(shown) > top {
		shown = shown[:top]
	}

	fmt.Printf("Busiest sources (%d of %d)\n", len(shown), len(r.Sources))
	fmt.Printf("  %-40s %8s %8s %6s  %s\n", "ADDRESS", "REJECTED", "PROBES", "IN", "ACCOUNTS NAMED")
	for _, s := range shown {
		marker := " "
		if s.GotIn() {
			marker = "!"
		}
		fmt.Printf("%s %-40s %8d %8d %6d  %s\n",
			marker, s.IP, s.Failed, s.Probed, s.Accepted, summarizeUsers(s.Users, 5))
	}

	// Volume alone doesn't say whether this is still going on. A run that
	// stopped yesterday and one still in progress read identically in the
	// table above, and they call for different responses.
	if latest, ip := mostRecent(r.Sources); !latest.IsZero() {
		fmt.Printf("\n  Most recent activity: %s (%s)\n", latest.Format(time.RFC3339), ip)
	}
	fmt.Println()
}

// mostRecent finds the newest activity across all sources, and where it came
// from. A zero time means no source carried a readable timestamp.
func mostRecent(sources []sshd.Source) (time.Time, string) {
	var newest time.Time
	var ip string
	for _, s := range sources {
		if s.Last.After(newest) {
			newest, ip = s.Last, s.IP
		}
	}
	return newest, ip
}

// summarizeUsers renders the first few account names an address tried, with a
// count of the rest — a scanner naming 200 accounts should not print 200.
func summarizeUsers(users []string, limit int) string {
	if len(users) == 0 {
		return "—"
	}
	if len(users) <= limit {
		return strings.Join(users, ", ")
	}
	return fmt.Sprintf("%s (+%d more)", strings.Join(users[:limit], ", "), len(users)-limit)
}

// printSubnets shows network blocks holding several sources. One operator
// renting a range looks like forty attackers in the per-address table and
// like one line here, which is also the level at which it can be blocked.
func printSubnets(r *sshd.Report) {
	var clustered []sshd.Subnet
	for _, s := range r.Subnets {
		if s.Sources >= clusterThreshold {
			clustered = append(clustered, s)
		}
	}
	if len(clustered) == 0 {
		return
	}
	if len(clustered) > 5 {
		clustered = clustered[:5]
	}

	fmt.Println("Clustered network blocks")
	fmt.Printf("  %-24s %9s %10s\n", "BLOCK", "ADDRESSES", "ATTEMPTS")
	for _, s := range clustered {
		fmt.Printf("  %-24s %9d %10d\n", s.CIDR, s.Sources, s.Attempts)
	}
	fmt.Println()
}

func printTargetedUsers(r *sshd.Report) {
	if len(r.Usernames) == 0 {
		return
	}

	shown := r.Usernames
	if len(shown) > 8 {
		shown = shown[:8]
	}

	// The address count is the column that separates one persistent bot from
	// a spray coordinated across hosts, so it gets a column of its own rather
	// than being folded into the attempt total.
	fmt.Println("Most-targeted accounts")
	fmt.Printf("  %-24s %9s %10s\n", "ACCOUNT", "ATTEMPTS", "ADDRESSES")
	for _, u := range shown {
		name := u.Name
		if name == "" {
			name = "(none named)"
		}
		fmt.Printf("  %-24s %9d %10d\n", name, u.Count, u.Sources)
	}
	fmt.Println()
}

func printLogins(r *sshd.Report) {
	if len(r.Logins) == 0 {
		fmt.Println("Successful logins: none in this window.")
		fmt.Println()
		return
	}

	fmt.Printf("Successful logins (%d)\n", len(r.Logins))
	for _, e := range r.Logins {
		when := "unknown time"
		if !e.Time.IsZero() {
			when = e.Time.Format(time.RFC3339)
		}
		method := e.Method
		if method == "" {
			method = "unknown method"
		}
		fmt.Printf("  %-25s %-16s %-12s from %s\n", when, e.User, method, e.IP)
	}
	fmt.Println()
}

// printAssessment turns the counts into the two or three sentences someone
// would actually act on.
func printAssessment(r *sshd.Report, settings map[string]string, settingsErr error) {
	printRule()
	fmt.Println("Assessment")

	if settingsErr != nil {
		fmt.Printf("  ?  Could not read the running sshd config (%v) — the counts above\n", settingsErr)
		fmt.Println("     are reported without it.")
	} else {
		passwordAuth := settings["passwordauthentication"]
		rootLogin := settings["permitrootlogin"]

		// An empty value means sshd -T didn't report the keyword at all, so
		// there's nothing to judge — say nothing rather than warn about "".
		switch {
		case passwordAuth == "":
		case passwordAuth == "no" && r.Failed > 0:
			fmt.Printf("  ✓  PasswordAuthentication is off, so none of those %d rejected\n", r.Failed)
			fmt.Println("     credentials could have succeeded. This is background noise.")
		case passwordAuth != "no" && r.Failed > 0:
			fmt.Printf("  ⚠  PasswordAuthentication is %q and %d credentials were rejected —\n", passwordAuth, r.Failed)
			fmt.Println("     these are guesses that could eventually land. Turn passwords off:")
			fmt.Println("         sudo serv-kun ssh harden")
		case passwordAuth != "no":
			fmt.Printf("  ⚠  PasswordAuthentication is %q. Nothing has tried it in this window,\n", passwordAuth)
			fmt.Println("     but keys-only is the safer default:  sudo serv-kun ssh harden")
		}

		switch {
		case rootLogin == "" || r.RootHits == 0:
		case rootLogin != "no":
			fmt.Printf("  ⚠  PermitRootLogin is %q and %d events named root. Disable it:\n", rootLogin, r.RootHits)
			fmt.Println("         sudo serv-kun ssh harden")
		default:
			fmt.Printf("  ✓  %d events named root, but PermitRootLogin is \"no\" — all refused.\n", r.RootHits)
		}
	}

	// Anyone who failed and then succeeded gets called out by name. This is
	// the one finding here that can mean an active compromise rather than
	// ambient scanning, so it is stated plainly and not buried in a count.
	var brokeIn []string
	for _, s := range r.Sources {
		if s.GotIn() {
			brokeIn = append(brokeIn, fmt.Sprintf("%s (%d rejected, then %d login(s))", s.IP, s.Failed, s.Accepted))
		}
	}
	if len(brokeIn) > 0 {
		fmt.Println("  !!  These addresses failed to authenticate and then succeeded:")
		for _, line := range brokeIn {
			fmt.Printf("         %s\n", line)
		}
		fmt.Println("      Confirm each was you. If not, that is a live intrusion — rotate the")
		fmt.Println("      account's keys and check ~/.ssh/authorized_keys for additions.")
	}

	var noisy int
	for _, s := range r.Sources {
		if s.Attempts() >= noisyThreshold {
			noisy++
		}
	}
	if noisy > 0 {
		fmt.Printf("  ⚠  %d address(es) made %d+ attempts each. Rate-limit SSH so a single\n", noisy, noisyThreshold)
		fmt.Println("     source is throttled after 6 connections in 30 seconds:")
		fmt.Println("         sudo serv-kun firewall")
	}

	printDistributedFindings(r)

	if r.Events == 0 {
		fmt.Println("  ✓  No authentication activity in this window.")
		fmt.Println("     If that looks wrong, sshd may log elsewhere — try --file or a wider --since.")
	}

	printRule()
}

// printDistributedFindings reports the attacks that per-address controls
// cannot see.
//
// This exists because the tables above can look reassuring during exactly the
// attack you should worry about: spread the same run over two thousand hosts
// and every row reads "3 attempts", no rate limit ever fires, and nothing in
// a per-source view is out of the ordinary. Only the totals give it away.
func printDistributedFindings(r *sshd.Report) {
	quietSources, quietAttempts := r.Distributed(quietThreshold)
	if quietSources >= distributedSources {
		fmt.Printf("  ⚠  %d addresses made %d attempts between them, none of them busy\n", quietSources, quietAttempts)
		fmt.Printf("     enough on its own for a rate limit to react (under %d each).\n", quietThreshold)
		fmt.Println("     That is a distributed attack, and no per-source control can see it.")
		fmt.Println("     What does work, in order of effect:")
		fmt.Println("       1. keys-only authentication — a distributed guess at a key is")
		fmt.Println("          not feasible:  sudo serv-kun ssh harden")
		fmt.Println("       2. restricting SSH to known source ranges, which refuses every")
		fmt.Println("          address not on the list:  sudo serv-kun firewall")
	}

	if sprayed, ok := r.Sprayed(); ok && sprayed.Sources >= sprayThreshold {
		fmt.Printf("  ⚠  %q was tried from %d different addresses (%d attempts) — one\n",
			sprayed.Name, sprayed.Sources, sprayed.Count)
		fmt.Println("     account, many hosts. That is credential spraying against a login")
		fmt.Println("     someone believes exists here, not untargeted scanning.")
		if sprayed.Name != "root" {
			fmt.Printf("     Confirm %q is an account you meant to have.\n", sprayed.Name)
		}
	}

	// Blocking a whole range is a blunt instrument, so it's a suggestion to
	// run by hand rather than something the firewall wizard applies.
	//
	// `prepend` is load-bearing: ufw takes the first matching rule, so a
	// plain `ufw deny from` is appended after the SSH allow and never
	// reached. `insert 1` puts it in the same place but fails outright on an
	// empty rule list ("Invalid position '1'"), so prepend is the form that
	// works whether or not the firewall has been set up yet.
	if len(r.Subnets) > 0 && r.Subnets[0].Sources >= clusterThreshold {
		top := r.Subnets[0]
		fmt.Printf("  ⚠  %d of those addresses sit in %s alone — one operator, not %d.\n",
			top.Sources, top.CIDR, top.Sources)
		fmt.Println("     Denying the block beats chasing its addresses:")
		fmt.Printf("         sudo ufw prepend deny from %s\n", top.CIDR)
		fmt.Println("     (prepend, not deny — an appended rule sits after the SSH allow")
		fmt.Println("      and is never reached.)")
	}
}
