package jwtkeys

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExistsRequiresBothHalves(t *testing.T) {
	dir := t.TempDir()

	if Exists(dir) {
		t.Error("Exists() true for an empty directory")
	}

	if err := os.WriteFile(filepath.Join(dir, "private.pem"), []byte("k"), 0600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	if Exists(dir) {
		t.Error("Exists() true with only private.pem — a half-generated pair must not count as present")
	}

	if err := os.WriteFile(filepath.Join(dir, "public.pem"), []byte("k"), 0644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	if !Exists(dir) {
		t.Error("Exists() false with both halves present")
	}
}

// Generate cleans up after itself on failure. Exists() requires both files, so
// a stray private.pem would make the next run mint a fresh key rather than
// reuse this one — silently invalidating every token already signed with it.
func TestGenerateLeavesNothingBehindOnFailure(t *testing.T) {
	dir := t.TempDir()
	certs := filepath.Join(dir, "certs")

	// A directory where private.pem must go makes openssl's write fail.
	if err := os.MkdirAll(filepath.Join(certs, "private.pem"), 0700); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}

	if err := Generate(certs); err == nil {
		t.Skip("openssl unexpectedly succeeded; nothing to assert")
	}

	if Exists(certs) {
		t.Error("Exists() true after a failed Generate")
	}
	if _, err := os.Stat(filepath.Join(certs, "public.pem")); err == nil {
		t.Error("public.pem left behind after a failed Generate")
	}
}
