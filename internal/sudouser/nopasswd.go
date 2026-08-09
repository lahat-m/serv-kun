package sudouser

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// SudoersFile is where serv-kun installs its NOPASSWD rule.
//
// The base name deliberately contains no ".": sudo's #includedir skips every
// file in /etc/sudoers.d whose name has a dot in it or ends in "~". A rule
// written to "serv-kun.conf" would be syntactically perfect, readable, owned
// by root — and silently never consulted.
const SudoersFile = "/etc/sudoers.d/serv-kun"

// Command is one command a NOPASSWD rule permits.
//
// Args pins the arguments sudo will accept. Leaving it empty permits any
// arguments, which is right for a wizard like serv-kun whose subcommands are
// the point, and wrong for a general-purpose tool like docker, where "any
// arguments" is a root shell in one hop: `docker run -v /:/host`.
type Command struct {
	Path string
	Args []string
}

func (c Command) String() string {
	if len(c.Args) == 0 {
		return c.Path
	}
	return c.Path + " " + strings.Join(c.Args, " ")
}

// NoPasswdRule renders the sudoers file granting username the listed commands
// without a password prompt.
func NoPasswdRule(username string, commands []Command) string {
	specs := make([]string, len(commands))
	for i, c := range commands {
		specs[i] = c.String()
	}

	var b strings.Builder
	b.WriteString("# Written by serv-kun.\n")
	b.WriteString("#\n")
	b.WriteString("# CI deploys over an SSH session with no terminal attached, so sudo has\n")
	b.WriteString("# nowhere to prompt for a password and fails with \"a terminal is required\n")
	b.WriteString("# to authenticate\". This grants exactly the commands the deploy runs, and\n")
	b.WriteString("# nothing else — every other sudo still asks for a password as usual.\n")
	b.WriteString("#\n")
	b.WriteString("# Re-run `sudo serv-kun ci` to regenerate. Remove the file to revoke.\n")
	fmt.Fprintf(&b, "%s ALL=(ALL) NOPASSWD: %s\n", username, strings.Join(specs, ", "))
	return b.String()
}

// InstallNoPasswd writes the rule to SudoersFile, validating it before it can
// take effect.
//
// The staging file is created inside /etc/sudoers.d under a name containing a
// dot, so sudo ignores it for exactly as long as it might be malformed. That
// ordering is the whole point: a bad file that sudo *does* read breaks sudo
// for every user on the machine at once, and recovering from that needs a root
// shell you may no longer be able to open. Only once visudo has accepted the
// contents is it renamed into place — a rename within one directory, so no
// window exists where the real path holds a partial file.
func InstallNoPasswd(username string, commands []Command) error {
	if username == "" {
		return errors.New("a username is required")
	}
	if len(commands) == 0 {
		return errors.New("no commands to grant")
	}

	dir := filepath.Dir(SudoersFile)
	if err := os.MkdirAll(dir, 0750); err != nil {
		return fmt.Errorf("create %s: %w", dir, err)
	}

	tmp, err := os.CreateTemp(dir, ".serv-kun-*.tmp")
	if err != nil {
		return fmt.Errorf("stage sudoers file: %w", err)
	}
	tmpPath := tmp.Name()
	defer func() {
		_ = tmp.Close()
		_ = os.Remove(tmpPath)
	}()

	if _, err := tmp.WriteString(NoPasswdRule(username, commands)); err != nil {
		return fmt.Errorf("write sudoers file: %w", err)
	}
	// Both are what sudo demands before it will read a file at all: mode no
	// wider than 0440, owned by root. A file it refuses to read is
	// indistinguishable, from the operator's side, from one that was never
	// written.
	if err := tmp.Chmod(0440); err != nil {
		return fmt.Errorf("chmod sudoers file: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("close sudoers file: %w", err)
	}
	if err := os.Chown(tmpPath, 0, 0); err != nil {
		return fmt.Errorf("chown sudoers file to root: %w", err)
	}

	if err := CheckSudoersSyntax(tmpPath); err != nil {
		return err
	}

	if err := os.Rename(tmpPath, SudoersFile); err != nil {
		return fmt.Errorf("install %s: %w", SudoersFile, err)
	}
	return nil
}

// CheckSudoersSyntax runs visudo's syntax check over a single file.
func CheckSudoersSyntax(path string) error {
	out, err := exec.Command(visudoPath(), "-c", "-f", path).CombinedOutput()
	if err != nil {
		return fmt.Errorf("visudo rejected the generated rule: %w\n%s", err, strings.TrimSpace(string(out)))
	}
	return nil
}

// visudoPath locates visudo, which lives in /usr/sbin and so is missing from
// PATH in some non-login root shells.
func visudoPath() string {
	if p, err := exec.LookPath("visudo"); err == nil {
		return p
	}
	return "/usr/sbin/visudo"
}
