package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/lahat-m/serv-kun/internal/compose"
	"github.com/lahat-m/serv-kun/internal/dns"
	"github.com/lahat-m/serv-kun/internal/envfile"
	"github.com/lahat-m/serv-kun/internal/system"
	"github.com/spf13/cobra"
)

var springbootTLSCmd = &cobra.Command{
	Use:   "tls",
	Short: "Issue or reissue the Let's Encrypt certificate nginx serves",
	Long: `tls runs the certificate bootstrap on its own, against the APP_DOMAIN and
CERTBOT_EMAIL already in .env.prod.

The same bootstrap runs as part of "springboot setup", but setup is a wizard
that re-prompts for every secret and rewrites .env.prod. That is the wrong tool
for a stack whose configuration is fine and whose certificate is not — which is
what you have after "docker compose down -v", because -v destroys the
certbot-certs volume.

Nothing else recovers from that. The certbot service runs "certbot renew", which
reads /etc/letsencrypt/renewal/*.conf and does nothing at all when that
directory is empty, so it never issues a first certificate. nginx mounts the
same volume read-only and its config names the certificate at startup, so it
cannot start to serve the ACME challenge that would create one. Breaking that
deadlock needs a throwaway self-signed certificate first, which is this command.

Certificates already installed are left alone unless --force is passed:
Let's Encrypt allows five duplicate certificates per domain per 168 hours, and
spending those on a stack that did not need reissuing locks the domain out for
a week.`,
	Args: cobra.NoArgs,
	RunE: runSpringbootTLS,
}

func init() {
	springbootTLSCmd.Flags().Bool("force", false,
		"reissue even if a valid certificate is already installed (spends a Let's Encrypt rate-limit slot)")
	springbootTLSCmd.Flags().Bool("skip-dns-check", false,
		"issue without first confirming the domain resolves to this server")
}

func runSpringbootTLS(cmd *cobra.Command, _ []string) error {
	requireRoot(cmd)
	printBanner("springboot tls")

	proj := projectFromFlags(cmd)
	warnProjectNameCollision(proj)

	env, err := envfile.Read(filepath.Join(proj.Dir, proj.EnvFile))
	if err != nil {
		return fmt.Errorf("%w — run 'serv-kun springboot setup' first", err)
	}
	domain := env["APP_DOMAIN"]
	if domain == "" {
		return fmt.Errorf("%s does not set APP_DOMAIN", proj.EnvFile)
	}
	// Empty is allowed: bootstrapTLS falls back to
	// --register-unsafely-without-email rather than failing.
	email := env["CERTBOT_EMAIL"]

	fmt.Printf("domain: %s\n", domain)
	if email != "" {
		fmt.Printf("email:  %s\n", email)
	} else {
		fmt.Println("email:  (none — registering without one, so no expiry warnings)")
	}
	fmt.Println()

	force, _ := cmd.Flags().GetBool("force")
	state, detail, err := certificateState(proj, domain)
	if err != nil {
		return err
	}
	switch state {
	case certLive:
		if !force {
			return fmt.Errorf("a valid certificate for %s is already installed (expires %s)\n\n"+
				"Reissuing spends one of five duplicate certificates allowed per 168 hours.\n"+
				"If nginx is serving the old one, reload it instead:\n"+
				"  docker compose -f %s exec nginx nginx -s reload\n\n"+
				"To reissue anyway: sudo %s --force",
				domain, detail, proj.File, cmd.CommandPath())
		}
		fmt.Printf("! A valid certificate is already installed (expires %s) — reissuing because --force was passed.\n\n", detail)
	case certPlaceholder:
		fmt.Println("Found a self-signed placeholder — a previous bootstrap got as far as")
		fmt.Println("starting nginx but never completed issuance. Replacing it.")
		fmt.Println()
	case certExpired:
		fmt.Println("Found an expired certificate. Reissuing.")
		fmt.Println()
	case certNone:
		fmt.Println("No certificate for this domain yet.")
		fmt.Println()
	case certUnknown:
		fmt.Println("! Could not classify the existing certificate; continuing anyway.")
		fmt.Println()
	}

	skipDNS, _ := cmd.Flags().GetBool("skip-dns-check")
	if skipDNS {
		fmt.Println("! Skipping the DNS check — issuance will fail if the domain does not resolve here.")
		fmt.Println()
	} else if err := verifyDomainResolvesHere(domain); err != nil {
		return err
	}

	if err := bootstrapTLS(proj, domain, email); err != nil {
		return fmt.Errorf("TLS bootstrap: %w", err)
	}

	fmt.Printf("\nnginx is now serving the real certificate for %s.\n", domain)
	fmt.Println("Renewal is automatic — the certbot service re-runs 'certbot renew' every 12h,")
	fmt.Println("and nginx reloads itself every 6h to pick up what it renews.")
	return nil
}

// verifyDomainResolvesHere refuses issuance when the domain doesn't point at
// this server. setup asks whether to proceed anyway; this command errors
// instead, because it exists to be run from recovery scripts where a prompt
// would hang forever with nobody watching — and because a failed validation is
// not free: Let's Encrypt allows five per hostname per hour.
func verifyDomainResolvesHere(domain string) error {
	fmt.Println("Checking that the domain points at this server...")
	serverIP, err := system.GetPublicIP()
	if err != nil {
		return fmt.Errorf("could not determine this server's public IP: %w\n"+
			"Re-run with --skip-dns-check to attempt issuance regardless", err)
	}
	fmt.Printf("  server IP: %s\n", serverIP)

	results := dns.CheckDomains([]string{domain}, serverIP)
	for _, r := range results {
		switch {
		case r.Proxied:
			fmt.Printf("  ✓ %s → Cloudflare edge (proxied — challenge still reaches this server)\n", r.Domain)
		case r.OK:
			fmt.Printf("  ✓ %s → %s\n", r.Domain, r.MatchedIP)
		default:
			fmt.Printf("  ✗ %s → %v (%s)\n", r.Domain, r.ResolvedIPs, r.Err)
		}
	}
	fmt.Println()

	if !dns.AllOK(results) {
		return fmt.Errorf("%s does not resolve to this server, so the HTTP-01 challenge cannot reach it\n\n"+
			"Point the A record here and re-run, or pass --skip-dns-check to try anyway.\n"+
			"Also confirm port 80 is open — the challenge is served over HTTP, not HTTPS:\n"+
			"  sudo serv-kun firewall", domain)
	}
	return nil
}

// certState classifies what is sitting in the certbot-certs volume for a
// domain. The distinction that matters is certLive vs certPlaceholder: a failed
// bootstrap leaves behind the throwaway self-signed certificate it wrote to get
// nginx started, and refusing to reissue over *that* would strand the operator
// on an invalid certificate with no way forward short of --force.
type certState int

const (
	certUnknown certState = iota
	certNone
	certPlaceholder
	certExpired
	certLive
)

// certStateScript classifies the certificate for $1 and prints the verdict
// behind a sentinel prefix. The prefix is not decoration: compose writes
// container lifecycle lines onto the same combined stream, so the answer has to
// be findable rather than simply read off as the whole output.
//
// A renewal config is what separates a real certificate from the placeholder.
// certbot writes /etc/letsencrypt/renewal/<domain>.conf on issuance; the
// placeholder is written by openssl straight into live/ and never has one.
const certStateScript = `set -e
L="/etc/letsencrypt/live/$1/fullchain.pem"
R="/etc/letsencrypt/renewal/$1.conf"
if [ ! -s "$L" ]; then
  echo "SERVKUN-CERT none"
elif [ ! -f "$R" ]; then
  echo "SERVKUN-CERT placeholder"
elif openssl x509 -in "$L" -noout -checkend 0 >/dev/null 2>&1; then
  echo "SERVKUN-CERT live $(openssl x509 -in "$L" -noout -enddate | cut -d= -f2)"
else
  echo "SERVKUN-CERT expired"
fi`

func certificateState(proj compose.Project, domain string) (certState, string, error) {
	out, err := proj.Output("run", "--rm", "-T", "--entrypoint", "sh", "certbot", "-c", certStateScript, "sh", domain)
	if err != nil {
		return certUnknown, "", fmt.Errorf("could not inspect the existing certificate: %w\n%s", err, out)
	}
	state, detail := parseCertState(out)
	return state, detail, nil
}

func parseCertState(out string) (certState, string) {
	for _, line := range strings.Split(out, "\n") {
		rest, ok := strings.CutPrefix(strings.TrimSpace(line), "SERVKUN-CERT ")
		if !ok {
			continue
		}
		kind, detail, _ := strings.Cut(rest, " ")
		switch kind {
		case "none":
			return certNone, ""
		case "placeholder":
			return certPlaceholder, ""
		case "expired":
			return certExpired, ""
		case "live":
			return certLive, strings.TrimSpace(detail)
		}
	}
	return certUnknown, ""
}

// bootstrapTLS mirrors scripts/init-letsencrypt.sh: nginx's config expects
// /etc/letsencrypt/live/<domain>/{fullchain,privkey}.pem to already exist, so
// on a brand-new box it can't start clean. This writes a throwaway
// self-signed cert so nginx *can* start, deletes it, asks certbot for the
// real one over the HTTP-01 challenge, then reloads nginx with it.
func bootstrapTLS(proj compose.Project, domain, email string) error {
	fmt.Println("\nBootstrapping TLS certificate...")

	// The domain is passed as a positional argument rather than interpolated
	// into the script text, so it can never word-split or be read as shell
	// syntax. That matters most for the rm below: a domain with a stray space
	// would otherwise expand into extra operands and wipe every certificate
	// on the host.
	const dummyScript = `set -e
mkdir -p "/etc/letsencrypt/live/$1"
openssl req -x509 -nodes -newkey rsa:2048 -days 1 \
  -keyout "/etc/letsencrypt/live/$1/privkey.pem" \
  -out "/etc/letsencrypt/live/$1/fullchain.pem" -subj '/CN=localhost'`

	if err := proj.Run("run", "--rm", "-T", "--entrypoint", "sh", "certbot", "-c", dummyScript, "sh", domain); err != nil {
		return fmt.Errorf("create dummy certificate: %w", err)
	}

	if err := proj.Run("up", "-d", "nginx"); err != nil {
		return fmt.Errorf("start nginx: %w", err)
	}

	const cleanupScript = `set -e
rm -rf "/etc/letsencrypt/live/$1" "/etc/letsencrypt/archive/$1" "/etc/letsencrypt/renewal/$1.conf"`

	if err := proj.Run("run", "--rm", "-T", "--entrypoint", "sh", "certbot", "-c", cleanupScript, "sh", domain); err != nil {
		return fmt.Errorf("remove dummy certificate: %w", err)
	}

	// -T disables compose's TTY allocation and --non-interactive stops certbot
	// prompting: between them, anything certbot would have asked about becomes
	// an error we can report instead of a wizard that hangs forever waiting on
	// input nobody is watching for.
	certbotArgs := []string{
		"run", "--rm", "-T", "--entrypoint", "certbot", "certbot", "certonly",
		"--non-interactive",
		"--webroot", "-w", "/var/www/certbot",
		"-d", domain,
		"--rsa-key-size", "2048", "--agree-tos", "--no-eff-email",
	}
	if email != "" {
		certbotArgs = append(certbotArgs, "--email", email)
	} else {
		certbotArgs = append(certbotArgs, "--register-unsafely-without-email")
	}
	if err := proj.Run(certbotArgs...); err != nil {
		// The dummy certificate was deleted to make way for the real one, so
		// a failure here leaves nginx with no certificate at all — it would
		// fail to start on the next restart, taking the site down over and
		// above the failed issuance. Put a throwaway one back.
		fmt.Println("\nCertificate issuance failed — restoring a placeholder so nginx can still start.")
		if restoreErr := proj.Run("run", "--rm", "-T", "--entrypoint", "sh", "certbot", "-c", dummyScript, "sh", domain); restoreErr != nil {
			return fmt.Errorf("%w\n\nnginx now has no certificate for %s and will not start; "+
				"re-run 'serv-kun springboot setup' to rebuild it (restore also failed: %v)",
				err, domain, restoreErr)
		}
		fmt.Println("✓ Placeholder restored — nginx will start, but serves an invalid certificate until issuance succeeds.")
		return fmt.Errorf("request certificate: %w", err)
	}

	if err := proj.Run("exec", "nginx", "nginx", "-s", "reload"); err != nil {
		return fmt.Errorf("reload nginx: %w", err)
	}

	fmt.Println("✓ TLS certificate issued")
	return nil
}
