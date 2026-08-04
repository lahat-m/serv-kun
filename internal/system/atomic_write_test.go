package system

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func skipOnWindows(t *testing.T) {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("Unix file modes aren't meaningful on Windows")
	}
}

func TestAtomicWriteFileCreatesWithMode(t *testing.T) {
	skipOnWindows(t)

	path := filepath.Join(t.TempDir(), "new.conf")
	if err := AtomicWriteFile(path, []byte("hello"), 0600); err != nil {
		t.Fatalf("AtomicWriteFile: %v", err)
	}

	fi, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat: %v", err)
	}
	if got := fi.Mode().Perm(); got != 0600 {
		t.Errorf("mode = %v, want %v", got, os.FileMode(0600))
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if string(data) != "hello" {
		t.Errorf("content = %q, want %q", data, "hello")
	}
}

func TestAtomicWriteFileLeavesNoTempFiles(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "cfg")
	if err := AtomicWriteFile(path, []byte("x"), 0644); err != nil {
		t.Fatalf("AtomicWriteFile: %v", err)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("ReadDir: %v", err)
	}
	if len(entries) != 1 || entries[0].Name() != "cfg" {
		var names []string
		for _, e := range entries {
			names = append(names, e.Name())
		}
		t.Errorf("directory contains %v, want just [cfg]", names)
	}
}

func TestModeOrDefault(t *testing.T) {
	skipOnWindows(t)

	dir := t.TempDir()
	existing := filepath.Join(dir, "existing")
	if err := os.WriteFile(existing, []byte("x"), 0600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	if got := ModeOrDefault(existing, 0644); got != 0600 {
		t.Errorf("ModeOrDefault(existing) = %v, want 0600", got)
	}
	if got := ModeOrDefault(filepath.Join(dir, "absent"), 0644); got != 0644 {
		t.Errorf("ModeOrDefault(absent) = %v, want the 0644 default", got)
	}
}

// The regression that shipped in v1.2.0: AtomicWriteFile replaces the inode
// and applies whatever mode it's given, unlike os.WriteFile which ignores its
// perm argument for an existing file. Rewriting a 0600 config with a
// hardcoded 0644 therefore widened it. Pairing with ModeOrDefault is what
// keeps a hardening command from loosening permissions.
func TestRewriteWithModeOrDefaultPreservesPermissions(t *testing.T) {
	skipOnWindows(t)

	path := filepath.Join(t.TempDir(), "sshd_config")
	if err := os.WriteFile(path, []byte("PermitRootLogin yes\n"), 0600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	if err := AtomicWriteFile(path, []byte("PermitRootLogin no\n"), ModeOrDefault(path, 0644)); err != nil {
		t.Fatalf("AtomicWriteFile: %v", err)
	}

	fi, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat: %v", err)
	}
	if got := fi.Mode().Perm(); got != 0600 {
		t.Errorf("mode = %v after rewrite, want 0600 preserved", got)
	}
}

// A hardcoded mode is exactly what caused the regression — assert the
// widening actually happens without ModeOrDefault, so this test fails loudly
// if someone reintroduces the pattern believing it to be harmless.
func TestHardcodedModeWidensPermissions(t *testing.T) {
	skipOnWindows(t)

	path := filepath.Join(t.TempDir(), "sshd_config")
	if err := os.WriteFile(path, []byte("a\n"), 0600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	if err := AtomicWriteFile(path, []byte("b\n"), 0644); err != nil {
		t.Fatalf("AtomicWriteFile: %v", err)
	}

	fi, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat: %v", err)
	}
	if got := fi.Mode().Perm(); got != 0644 {
		t.Fatalf("expected the hardcoded mode to win (%v), got %v — "+
			"if AtomicWriteFile now preserves modes itself, ModeOrDefault callers can be simplified",
			os.FileMode(0644), got)
	}
}
