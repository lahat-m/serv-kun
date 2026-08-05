// Package jwtkeys generates the RSA keypair Kuniel-shaped apps use for
// signing JWTs (jwt.public-key / jwt.private-key), matching the exact
// openssl invocations documented in DEPLOY.md.
package jwtkeys

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/lahat-m/serv-kun/internal/system"
)

// Exists reports whether both private.pem and public.pem are already present in dir.
func Exists(dir string) bool {
	_, privErr := os.Stat(filepath.Join(dir, "private.pem"))
	_, pubErr := os.Stat(filepath.Join(dir, "public.pem"))
	return privErr == nil && pubErr == nil
}

// Generate creates dir/private.pem and dir/public.pem via openssl.
func Generate(dir string) (err error) {
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("create %s: %w", dir, err)
	}

	priv := filepath.Join(dir, "private.pem")
	pub := filepath.Join(dir, "public.pem")

	// Clear a half-generated pair. Exists() requires both files, so leaving
	// only private.pem behind would make the next run generate a fresh key
	// rather than reuse this one — silently invalidating every token already
	// signed with it.
	defer func() {
		if err != nil {
			_ = os.Remove(priv)
			_ = os.Remove(pub)
		}
	}()

	if err = system.RunCmd("openssl", "genrsa", "-out", priv, "2048"); err != nil {
		return fmt.Errorf("generate private key: %w", err)
	}
	if err = os.Chmod(priv, 0600); err != nil {
		return fmt.Errorf("chmod private key: %w", err)
	}
	if err = system.RunCmd("openssl", "rsa", "-in", priv, "-pubout", "-out", pub); err != nil {
		return fmt.Errorf("derive public key: %w", err)
	}
	return nil
}
