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

	dummyScript := fmt.Sprintf(
		"mkdir -p /etc/letsencrypt/live/%s && "+
			"openssl req -x509 -nodes -newkey rsa:2048 -days 1 "+
			"-keyout /etc/letsencrypt/live/%s/privkey.pem "+
			"-out /etc/letsencrypt/live/%s/fullchain.pem -subj '/CN=localhost'",
		domain, domain, domain,
	)
	if err := proj.Run("run", "--rm", "--entrypoint", "sh", "certbot", "-c", dummyScript); err != nil {
		return fmt.Errorf("create dummy certificate: %w", err)
	}

	if err := proj.Run("up", "-d", "nginx"); err != nil {
		return fmt.Errorf("start nginx: %w", err)
	}

	cleanupScript := fmt.Sprintf(
		"rm -rf /etc/letsencrypt/live/%s /etc/letsencrypt/archive/%s /etc/letsencrypt/renewal/%s.conf",
		domain, domain, domain,
	)
	if err := proj.Run("run", "--rm", "--entrypoint", "sh", "certbot", "-c", cleanupScript); err != nil {
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
		return fmt.Errorf("request certificate: %w", err)
	}

	if err := proj.Run("exec", "nginx", "nginx", "-s", "reload"); err != nil {
		return fmt.Errorf("reload nginx: %w", err)
	}

	fmt.Println("✓ TLS certificate issued")
	return nil
}
