// Package sshd hardens the OpenSSH server: disabling root login and
// password authentication. Every change is validated with sshd's own
// `sshd -t` before the daemon is reloaded, and rolled back from a
// timestamped backup if that validation fails — a bad sshd_config that
// gets loaded is how you lose a server.
package sshd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"github.com/lahat-m/serv-kun/internal/system"
)

const (
	// ConfigPath is the main sshd config file.
	ConfigPath = "/etc/ssh/sshd_config"

	// IncludeDir holds drop-in configs on distros whose sshd_config has an
	// `Include` line — Ubuntu 22.04+, Debian 12+.
	IncludeDir = "/etc/ssh/sshd_config.d"

	// DropInName deliberately sorts ahead of vendor drop-ins like
	// 50-cloud-init.conf. sshd takes the *first* value it obtains for a
	// keyword and Include expands in lexical order, so the 00- prefix is
	// the whole reason our settings win instead of cloud-init's.
	DropInName = "00-serv-kun-hardening.conf"
)

// Options selects which hardening changes to apply.
type Options struct {
	DisableRootLogin    bool
	DisablePasswordAuth bool
}

type directive struct{ Key, Value string }

// directives returns the keywords to set, in the order they're written.
func (o Options) directives() []directive {
	var ds []directive
	if o.DisableRootLogin {
		ds = append(ds, directive{"PermitRootLogin", "no"})
	}
	if o.DisablePasswordAuth {
		ds = append(ds,
			directive{"PasswordAuthentication", "no"},
			// Without this, PAM can still offer an interactive password
			// prompt and PasswordAuthentication no buys nothing.
			directive{"KbdInteractiveAuthentication", "no"},
			directive{"PubkeyAuthentication", "yes"},
		)
	}
	return ds
}

// Any reports whether any change was requested.
func (o Options) Any() bool { return len(o.directives()) > 0 }

// binary locates sshd, which lives in sbin and so is often off the PATH.
func binary() string {
	if p, err := exec.LookPath("sshd"); err == nil {
		return p
	}
	return "/usr/sbin/sshd"
}

// EffectiveSettings runs `sshd -T`, which dumps the configuration sshd
// would actually use. It's the only reliable way to read the real values:
// it resolves Include files and distro drop-ins that a naive parse of
// sshd_config would miss entirely. Keys come back lowercased.
func EffectiveSettings() (map[string]string, error) {
	out, err := exec.Command(binary(), "-T").Output()
	if err != nil {
		return nil, fmt.Errorf("run sshd -T: %w", err)
	}

	m := make(map[string]string)
	sc := bufio.NewScanner(strings.NewReader(string(out)))
	for sc.Scan() {
		k, v, ok := strings.Cut(strings.TrimSpace(sc.Text()), " ")
		if !ok {
			continue
		}
		m[strings.ToLower(k)] = v
	}
	return m, sc.Err()
}

// Port returns the port sshd actually listens on, so the firewall opens the
// right one. Falls back to 22 if the config can't be read.
func Port() string {
	settings, err := EffectiveSettings()
	if err != nil {
		return "22"
	}
	if p := settings["port"]; p != "" {
		return p
	}
	return "22"
}

// Test runs `sshd -t`, sshd's own config validation.
func Test() error {
	out, err := exec.Command(binary(), "-t").CombinedOutput()
	if err != nil {
		return fmt.Errorf("sshd config test failed:\n%s", strings.TrimSpace(string(out)))
	}
	return nil
}

var keyTypePrefixes = []string{
	"ssh-rsa", "ssh-ed25519", "ssh-dss", "ecdsa-sha2-", "sk-ssh-", "sk-ecdsa-",
}

func isKeyLine(line string) bool {
	// Options may precede the key type (e.g. no-port-forwarding ssh-ed25519 ...),
	// so scan every field rather than just the first.
	for _, f := range strings.Fields(line) {
		for _, p := range keyTypePrefixes {
			if strings.HasPrefix(f, p) {
				return true
			}
		}
	}
	return false
}

// AuthorizedKeys returns the key lines in username's authorized_keys file.
// This is the evidence that someone can still log in after passwords and
// root are taken away; a nil result means nobody can.
func AuthorizedKeys(username string) ([]string, error) {
	u, err := user.Lookup(username)
	if err != nil {
		return nil, fmt.Errorf("lookup user %q: %w", username, err)
	}

	path := filepath.Join(u.HomeDir, ".ssh", "authorized_keys")
	f, err := os.Open(path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	defer f.Close()

	var keys []string
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if isKeyLine(line) {
			keys = append(keys, line)
		}
	}
	return keys, sc.Err()
}

// UsesIncludeDir reports whether sshd_config pulls in
// /etc/ssh/sshd_config.d/*.conf. When it does, editing the main file is a
// trap: drop-ins there are read first and win, so a cloud-init file setting
// PasswordAuthentication yes would quietly undo our change.
func UsesIncludeDir() (bool, error) {
	f, err := os.Open(ConfigPath)
	if err != nil {
		return false, fmt.Errorf("read %s: %w", ConfigPath, err)
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		fields := strings.Fields(strings.TrimSpace(sc.Text()))
		if len(fields) >= 2 && !strings.HasPrefix(fields[0], "#") &&
			strings.EqualFold(fields[0], "Include") &&
			strings.Contains(fields[1], "sshd_config.d") {
			return true, nil
		}
	}
	return false, sc.Err()
}

func blockLines(ds []directive) []string {
	lines := []string{
		"# ── serv-kun SSH hardening ──────────────────────────────────────",
		"# sshd uses the first value it obtains for a keyword, so these",
		"# are placed ahead of anything that would otherwise override them.",
	}
	for _, d := range ds {
		lines = append(lines, d.Key+" "+d.Value)
	}
	return append(lines, "")
}

// Render returns the directive block Apply would install, for previewing
// the change before it's made.
func Render(opts Options) string {
	return strings.Join(blockLines(opts.directives()), "\n")
}

// rewriteConfig comments out top-level settings for the keys we manage and
// prepends our block. Occurrences inside Match blocks are left alone: those
// are deliberate per-connection overrides, and blanking them could change
// behavior well beyond what was asked for.
func rewriteConfig(content string, ds []directive) string {
	managed := make(map[string]bool, len(ds))
	for _, d := range ds {
		managed[strings.ToLower(d.Key)] = true
	}

	out := blockLines(ds)
	inMatch := false

	for _, line := range strings.Split(content, "\n") {
		trimmed := strings.TrimSpace(line)
		fields := strings.Fields(trimmed)

		if len(fields) > 0 && strings.EqualFold(fields[0], "Match") {
			inMatch = true
		}

		if !inMatch && len(fields) > 0 && !strings.HasPrefix(trimmed, "#") &&
			managed[strings.ToLower(fields[0])] {
			out = append(out, "# serv-kun: superseded above — "+line)
			continue
		}
		out = append(out, line)
	}

	return strings.Join(out, "\n")
}

// backupFile copies path to path.bak.<timestamp>, returning "" when there
// was no existing file to preserve.
func backupFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("read %s for backup: %w", path, err)
	}

	backup := path + ".bak." + time.Now().Format("20060102-150405")
	if err := os.WriteFile(backup, data, 0600); err != nil {
		return "", fmt.Errorf("write backup: %w", err)
	}
	return backup, nil
}

// Apply installs the hardening directives, returning the file it wrote and
// the backup it took (empty when the file was newly created).
func Apply(opts Options) (written, backup string, err error) {
	ds := opts.directives()
	if len(ds) == 0 {
		return "", "", fmt.Errorf("no changes requested")
	}

	useDropIn, err := UsesIncludeDir()
	if err != nil {
		return "", "", err
	}

	if useDropIn {
		if err := os.MkdirAll(IncludeDir, 0755); err != nil {
			return "", "", fmt.Errorf("create %s: %w", IncludeDir, err)
		}
		path := filepath.Join(IncludeDir, DropInName)
		backup, err = backupFile(path)
		if err != nil {
			return "", "", err
		}
		if err := os.WriteFile(path, []byte(strings.Join(blockLines(ds), "\n")), 0644); err != nil {
			return "", backup, fmt.Errorf("write %s: %w", path, err)
		}
		return path, backup, nil
	}

	data, err := os.ReadFile(ConfigPath)
	if err != nil {
		return "", "", fmt.Errorf("read %s: %w", ConfigPath, err)
	}
	backup, err = backupFile(ConfigPath)
	if err != nil {
		return "", "", err
	}
	if err := os.WriteFile(ConfigPath, []byte(rewriteConfig(string(data), ds)), 0644); err != nil {
		return "", backup, fmt.Errorf("write %s: %w", ConfigPath, err)
	}
	return ConfigPath, backup, nil
}

// Restore undoes Apply. With no backup the file didn't exist beforehand, so
// removing it is what returns the system to its prior state.
func Restore(path, backup string) error {
	if backup == "" {
		if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("remove %s: %w", path, err)
		}
		return nil
	}

	data, err := os.ReadFile(backup)
	if err != nil {
		return fmt.Errorf("read backup %s: %w", backup, err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("restore %s: %w", path, err)
	}
	return nil
}

// ServiceName is the systemd unit for the SSH daemon: `ssh` on the Debian
// family, `sshd` on RHEL.
func ServiceName(osInfo *system.Info) string {
	if osInfo.Family == system.RHEL {
		return "sshd"
	}
	return "ssh"
}

// Reload has sshd re-read its config. Reload rather than restart, so
// established sessions survive — including the one running this command.
func Reload(osInfo *system.Info) error {
	name := ServiceName(osInfo)
	out, err := exec.Command("systemctl", "reload", name).CombinedOutput()
	if err == nil {
		return nil
	}

	// Unit naming isn't fully consistent across distros; try the other name
	// before reporting failure.
	alt := "sshd"
	if name == "sshd" {
		alt = "ssh"
	}
	altOut, altErr := exec.Command("systemctl", "reload", alt).CombinedOutput()
	if altErr == nil {
		return nil
	}

	return fmt.Errorf("reload %s failed:\n%s\nreload %s failed:\n%s",
		name, strings.TrimSpace(string(out)), alt, strings.TrimSpace(string(altOut)))
}
