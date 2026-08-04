package sudouser

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/lahat-m/serv-kun/internal/system"
)

// Exists reports whether a Linux user account already exists.
func Exists(username string) bool {
	_, err := user.Lookup(username)
	return err == nil
}

// Create adds a new Linux user with a home directory and bash shell.
func Create(username string) error {
	return system.RunCmd("useradd", "-m", "-s", "/bin/bash", username)
}

// SudoGroup returns the group that grants root privileges via sudo for the OS family.
func SudoGroup(osInfo *system.Info) string {
	if osInfo.Family == system.RHEL {
		return "wheel"
	}
	return "sudo"
}

// GrantSudo adds username to the OS's sudo/wheel group.
func GrantSudo(username string, osInfo *system.Info) error {
	return system.RunCmd("usermod", "-aG", SudoGroup(osInfo), username)
}

// HasSudo reports whether username is already in the OS's sudo/wheel group.
func HasSudo(username string, osInfo *system.Info) (bool, error) {
	u, err := user.Lookup(username)
	if err != nil {
		return false, fmt.Errorf("lookup user %q: %w", username, err)
	}

	group, err := user.LookupGroup(SudoGroup(osInfo))
	if err != nil {
		return false, fmt.Errorf("lookup group %q: %w", SudoGroup(osInfo), err)
	}

	// GroupIds covers the primary group too, not just supplementary ones.
	gids, err := u.GroupIds()
	if err != nil {
		return false, fmt.Errorf("list groups for %q: %w", username, err)
	}
	for _, gid := range gids {
		if gid == group.Gid {
			return true, nil
		}
	}
	return false, nil
}

// SetPassword sets the Unix login password for username via chpasswd.
func SetPassword(username, password string) error {
	cmd := exec.Command("chpasswd")
	cmd.Stdin = strings.NewReader(fmt.Sprintf("%s:%s\n", username, password))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// AddSSHKey appends pubKey to username's authorized_keys file, creating
// ~/.ssh if needed and fixing ownership/permissions so sshd will accept it.
func AddSSHKey(username, pubKey string) error {
	u, err := user.Lookup(username)
	if err != nil {
		return fmt.Errorf("lookup user: %w", err)
	}

	sshDir := filepath.Join(u.HomeDir, ".ssh")
	if err := os.MkdirAll(sshDir, 0700); err != nil {
		return fmt.Errorf("create .ssh dir: %w", err)
	}

	authFile := filepath.Join(sshDir, "authorized_keys")
	f, err := os.OpenFile(authFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("open authorized_keys: %w", err)
	}
	defer f.Close()

	if _, err := f.WriteString(strings.TrimSpace(pubKey) + "\n"); err != nil {
		return fmt.Errorf("write authorized_keys: %w", err)
	}

	uid, err := strconv.Atoi(u.Uid)
	if err != nil {
		return fmt.Errorf("parse uid: %w", err)
	}
	gid, err := strconv.Atoi(u.Gid)
	if err != nil {
		return fmt.Errorf("parse gid: %w", err)
	}

	if err := os.Chown(sshDir, uid, gid); err != nil {
		return fmt.Errorf("chown .ssh dir: %w", err)
	}
	if err := os.Chown(authFile, uid, gid); err != nil {
		return fmt.Errorf("chown authorized_keys: %w", err)
	}

	return nil
}
