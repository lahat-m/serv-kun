package cmd

import (
	"fmt"

	"github.com/lahat-m/serv-kun/internal/compose"
)

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

	if err := proj.Run("run", "--rm", "--entrypoint", "sh", "certbot", "-c", dummyScript, "sh", domain); err != nil {
		return fmt.Errorf("create dummy certificate: %w", err)
	}

	if err := proj.Run("up", "-d", "nginx"); err != nil {
		return fmt.Errorf("start nginx: %w", err)
	}

	const cleanupScript = `set -e
rm -rf "/etc/letsencrypt/live/$1" "/etc/letsencrypt/archive/$1" "/etc/letsencrypt/renewal/$1.conf"`

	if err := proj.Run("run", "--rm", "--entrypoint", "sh", "certbot", "-c", cleanupScript, "sh", domain); err != nil {
		return fmt.Errorf("remove dummy certificate: %w", err)
	}

	certbotArgs := []string{
		"run", "--rm", "--entrypoint", "certbot", "certbot", "certonly",
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
		if restoreErr := proj.Run("run", "--rm", "--entrypoint", "sh", "certbot", "-c", dummyScript, "sh", domain); restoreErr != nil {
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
