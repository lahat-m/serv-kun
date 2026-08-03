package node

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/lahat-m/serv-kun/internal/system"
)

// installScriptURL is pinned to a specific nvm release, per upstream's own
// install instructions, rather than following a mutable branch ref.
const installScriptURL = "https://raw.githubusercontent.com/nvm-sh/nvm/v0.40.1/install.sh"

// IsInstalled reports whether nvm is already installed for username.
func IsInstalled(username string) (bool, error) {
	u, err := user.Lookup(username)
	if err != nil {
		return false, fmt.Errorf("lookup user %q: %w", username, err)
	}

	_, err = os.Stat(filepath.Join(u.HomeDir, ".nvm", "nvm.sh"))
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// InstallNVM installs nvm into username's home directory.
func InstallNVM(osInfo *system.Info, username string) error {
	if err := ensureCurl(osInfo); err != nil {
		return fmt.Errorf("ensure curl: %w", err)
	}
	script := fmt.Sprintf("curl -o- %s | bash", installScriptURL)
	return runAsUserLogin(username, script)
}

// InstallVersion installs a Node.js version (e.g. "--lts", "20", "22.4.0") via nvm.
func InstallVersion(username, version string) error {
	script := fmt.Sprintf("%s nvm install %s", nvmInit, shellQuote(version))
	return runAsUserLogin(username, script)
}

// SetDefault aliases an already-installed version as username's default.
func SetDefault(username, version string) error {
	script := fmt.Sprintf("%s nvm alias default %s", nvmInit, shellQuote(version))
	return runAsUserLogin(username, script)
}

// Versions returns the output of `nvm ls` for username.
func Versions(username string) (string, error) {
	script := nvmInit + " nvm ls"
	return outputAsUserLogin(username, script)
}

// BinaryPath resolves the absolute path to a Node.js binary (e.g. "node",
// "npm", "npx") for the given nvm alias/version (e.g. "default", "20"), as
// installed for username.
func BinaryPath(username, alias, binary string) (string, error) {
	script := fmt.Sprintf("%s nvm which %s", nvmInit, shellQuote(alias))
	out, err := outputAsUserLogin(username, script)
	if err != nil {
		return "", fmt.Errorf("resolve nvm path: %w\n%s", err, out)
	}

	var nodePath string
	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "/") {
			nodePath = line
		}
	}
	if nodePath == "" {
		return "", fmt.Errorf("could not resolve nvm path for %q:\n%s", alias, out)
	}
	return filepath.Join(filepath.Dir(nodePath), binary), nil
}

// LinkSystemBinaries symlinks node/npm/npx for username's default nvm version
// into /usr/local/bin, which is on PATH for every shell by default — unlike
// ~/.nvm, which only gets wired up when a shell's rc file is sourced at
// login. This makes the binaries usable immediately, including in whatever
// session just ran the installer, without logging out and back in. Existing
// symlinks there are replaced, so this always reflects whichever version is
// currently aliased "default".
func LinkSystemBinaries(username string) error {
	nodePath, err := BinaryPath(username, "default", "node")
	if err != nil {
		return err
	}
	binDir := filepath.Dir(nodePath)

	for _, name := range []string{"node", "npm", "npx"} {
		src := filepath.Join(binDir, name)
		if _, err := os.Stat(src); err != nil {
			continue
		}
		dst := filepath.Join("/usr/local/bin", name)
		os.Remove(dst)
		if err := os.Symlink(src, dst); err != nil {
			return fmt.Errorf("symlink %s: %w", name, err)
		}
	}
	return nil
}

const nvmInit = `export NVM_DIR="$HOME/.nvm"; [ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh";`

func ensureCurl(osInfo *system.Info) error {
	if _, err := exec.LookPath("curl"); err == nil {
		return nil
	}
	switch osInfo.Family {
	case system.Debian:
		if err := system.RunCmd("apt-get", "update", "-qq"); err != nil {
			return err
		}
		return system.RunCmd("apt-get", "install", "-y", "curl")
	case system.RHEL:
		return system.RunCmd("dnf", "install", "-y", "curl")
	}
	return fmt.Errorf("unsupported OS family")
}

func runAsUserLogin(username, script string) error {
	cmd := exec.Command("su", "-l", username, "-c", script)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func outputAsUserLogin(username, script string) (string, error) {
	cmd := exec.Command("su", "-l", username, "-c", script)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func shellQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}
