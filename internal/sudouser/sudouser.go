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

// PasswordState reports whether username has a usable login password.
// `passwd -S` prints a status field: P/PS means a password is set, L/LK that
// it's locked, NP that there is none. Accounts created with only an SSH key
// land in L — useradd locks the password by default.
func PasswordState(username string) (usable bool, status string, err error) {
	out, err := exec.Command("passwd", "-S", username).Output()
	if err != nil {
		return false, "", fmt.Errorf("passwd -S %s: %w", username, err)
	}

	fields := strings.Fields(string(out))
	if len(fields) < 2 {
		return false, "", fmt.Errorf("unexpected passwd -S output: %q", strings.TrimSpace(string(out)))
	}

	status = fields[1]
	return status == "P" || status == "PS", status, nil
}

// HasNoPasswdSudo reports whether sudoers grants username at least one
// command without a password prompt.
func HasNoPasswdSudo(username string) bool {
	out, err := exec.Command("sudo", "-n", "-l", "-U", username).CombinedOutput()
	if err != nil {
		return false
	}
	return strings.Contains(string(out), "NOPASSWD")
}

// CanSudo reports whether username can actually run sudo, and why not when
// it can't. Group membership alone doesn't answer this: the default sudoers
// policy authenticates with the user's password, so a key-only account sits
// in the sudo group and still gets "sudo: a password is required" for every
// command — logged in, but unable to administer anything.
func CanSudo(username string, osInfo *system.Info) (bool, string) {
	inGroup, err := HasSudo(username, osInfo)
	if err != nil {
		return false, fmt.Sprintf("could not check group membership: %v", err)
	}
	if !inGroup {
		return false, fmt.Sprintf("not in the %q group", SudoGroup(osInfo))
	}

	if HasNoPasswdSudo(username) {
		return true, "in the sudo group with a NOPASSWD sudoers rule"
	}

	usable, status, err := PasswordState(username)
	if err != nil {
		// Don't block on an unreadable password database — report the group
		// membership we did confirm and let the operator judge.
		return true, fmt.Sprintf("in the %q group (password state unknown: %v)", SudoGroup(osInfo), err)
	}
	if !usable {
		return false, fmt.Sprintf("in the %q group but has no usable password (passwd -S reports %q), "+
			"so sudo will reject every command", SudoGroup(osInfo), status)
	}

	return true, fmt.Sprintf("in the %q group with a usable password", SudoGroup(osInfo))
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
