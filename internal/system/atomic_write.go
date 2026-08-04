package system

import (
	"os"
	"path/filepath"
)

// ModeOrDefault returns path's current permission bits, or def if the file
// doesn't exist.
//
// Pair this with AtomicWriteFile whenever you're replacing a file that may
// already exist. os.WriteFile *ignores* its perm argument for an existing
// file — the mode only applies at creation — so code converted from it will
// silently start changing permissions, since AtomicWriteFile replaces the
// inode and applies whatever mode it is handed. That regression widened a
// 0600 sshd_config to 0644.
func ModeOrDefault(path string, def os.FileMode) os.FileMode {
	if fi, err := os.Stat(path); err == nil {
		return fi.Mode().Perm()
	}
	return def
}

// AtomicWriteFile writes data to path atomically by creating a temp file in the
// same directory and renaming it into place. The perm argument sets the
// permissions on the temporary file before rename.
func AtomicWriteFile(path string, data []byte, perm os.FileMode) error {
	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, ".tmp-*")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()
	defer func() {
		_ = tmp.Close()
		_ = os.Remove(tmpPath)
	}()

	// Attempt to set the requested permissions on the temp file. Ignore error on
	// systems where this isn't supported.
	_ = tmp.Chmod(perm)

	if _, err := tmp.Write(data); err != nil {
		return err
	}
	if err := tmp.Sync(); err != nil {
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmpPath, path)
}
