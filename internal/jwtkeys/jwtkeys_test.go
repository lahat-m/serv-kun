package jwtkeys

import (
	"os"
	"path/filepath"
	"runtime"
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

// Handing the keys to the container's user must not also hand them to the rest
// of the host: ownership moves, permissions stay where Generate put them.
func TestGrantUserKeepsPermissions(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("chown and Unix file modes aren't meaningful on Windows")
	}

	dir := t.TempDir()
	if err := os.Chmod(dir, 0700); err != nil {
		t.Fatalf("Chmod dir: %v", err)
	}
	priv := filepath.Join(dir, "private.pem")
	pub := filepath.Join(dir, "public.pem")
	if err := os.WriteFile(priv, []byte("k"), 0600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	if err := os.WriteFile(pub, []byte("k"), 0644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	// Only root can hand a file to a different user, so the test re-grants it
	// to the current one — enough to prove the chown runs and the modes survive.
	if err := GrantUser(dir, os.Getuid(), os.Getgid()); err != nil {
		t.Fatalf("GrantUser: %v", err)
	}

	for path, want := range map[string]os.FileMode{dir: 0700, priv: 0600, pub: 0644} {
		info, err := os.Stat(path)
		if err != nil {
			t.Fatalf("Stat %s: %v", path, err)
		}
		if got := info.Mode().Perm(); got != want {
			t.Errorf("%s mode = %o after GrantUser, want %o", filepath.Base(path), got, want)
		}
	}
}

// GrantUser has nothing to hand over if a key is missing, and says so rather
// than reporting success on a directory the container still can't use.
func TestGrantUserFailsOnMissingKey(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("chown isn't meaningful on Windows")
	}

	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "private.pem"), []byte("k"), 0600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	if err := GrantUser(dir, os.Getuid(), os.Getgid()); err == nil {
		t.Error("GrantUser() = nil with public.pem missing")
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
