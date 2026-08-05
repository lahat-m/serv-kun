package nginx

import (
	"os"
	"path/filepath"
	"testing"
)

// An invalid config left on disk stops nginx starting at the next reload or
// reboot, long after the command that wrote it finished — so a run that fails
// validation must leave nothing behind.
func TestRestoreConfigRemovesANewlyWrittenFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "example.com.conf")

	if err := os.WriteFile(path, []byte("server { broken\n"), 0644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	// Backup empty == the config did not exist before this run.
	if err := RestoreConfig(Written{Path: path}); err != nil {
		t.Fatalf("RestoreConfig: %v", err)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Errorf("config still present after rollback (err=%v)", err)
	}
}

func TestRestoreConfigPutsBackPreviousContents(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "example.com.conf")
	backup := path + ".bak.20260101-000000"

	const original = "server { listen 80; }\n"
	if err := os.WriteFile(backup, []byte(original), 0640); err != nil {
		t.Fatalf("WriteFile backup: %v", err)
	}
	if err := os.WriteFile(path, []byte("server { broken\n"), 0640); err != nil {
		t.Fatalf("WriteFile config: %v", err)
	}

	if err := RestoreConfig(Written{Path: path, Backup: backup}); err != nil {
		t.Fatalf("RestoreConfig: %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if string(got) != original {
		t.Errorf("content = %q, want %q", got, original)
	}
	if fi, err := os.Stat(path); err == nil {
		if fi.Mode().Perm() != 0640 {
			t.Errorf("mode = %v, want 0640 preserved", fi.Mode().Perm())
		}
	}
}

func TestRestoreConfigRemovesASymlinkItCreated(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "example.com.conf")
	link := filepath.Join(dir, "enabled.conf")

	if err := os.WriteFile(target, []byte("x"), 0644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	if err := os.Symlink(target, link); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}

	if err := RestoreConfig(Written{Path: target, Link: link}); err != nil {
		t.Fatalf("RestoreConfig: %v", err)
	}
	if _, err := os.Lstat(link); !os.IsNotExist(err) {
		t.Errorf("symlink still present after rollback (err=%v)", err)
	}
}

// enableSite replaces any symlink already there, so rollback has to put the
// old one back rather than just deleting ours.
func TestRestoreConfigReinstatesAPreviousSymlink(t *testing.T) {
	dir := t.TempDir()
	oldTarget := filepath.Join(dir, "old.conf")
	newTarget := filepath.Join(dir, "new.conf")
	link := filepath.Join(dir, "enabled.conf")

	for _, p := range []string{oldTarget, newTarget} {
		if err := os.WriteFile(p, []byte("x"), 0644); err != nil {
			t.Fatalf("WriteFile: %v", err)
		}
	}
	if err := os.Symlink(newTarget, link); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}

	err := RestoreConfig(Written{Path: newTarget, Link: link, LinkTarget: oldTarget})
	if err != nil {
		t.Fatalf("RestoreConfig: %v", err)
	}

	got, err := os.Readlink(link)
	if err != nil {
		t.Fatalf("Readlink: %v", err)
	}
	if got != oldTarget {
		t.Errorf("symlink points at %q, want the previous target %q", got, oldTarget)
	}
}

func TestRestoreConfigIsSafeOnAnEmptyRecord(t *testing.T) {
	if err := RestoreConfig(Written{}); err != nil {
		t.Errorf("RestoreConfig(zero) = %v, want nil", err)
	}
}
