package nginx

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/lahat-m/serv-kun/internal/system"
)

// Config holds the values needed to generate an nginx reverse proxy config.
type Config struct {
	Domains     []string
	Upstream    string
	MaxBodySize string
	OSInfo      *system.Info
}

// configTemplate is the nginx server block we generate.
// proxy_read_timeout is set to 86400s to support long-lived WebSocket connections.
const configTemplate = `server {
    listen 80;
    server_name {{.ServerNames}};

    client_max_body_size {{.MaxBodySize}};

    location / {
        proxy_pass {{.Upstream}};
        proxy_http_version 1.1;

        # WebSocket support
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";

        # Forwarding headers
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        # Timeouts
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 86400s;

        proxy_buffering off;
    }
}
`

type templateData struct {
	ServerNames string
	Upstream    string
	MaxBodySize string
}

func IsInstalled() (bool, error) {
	_, err := exec.LookPath("nginx")
	if err != nil {
		return false, nil
	}
	return true, nil
}

func Install(osInfo *system.Info) error {
	switch osInfo.Family {
	case system.Debian:
		if err := system.RunCmd("apt-get", "update", "-qq"); err != nil {
			return fmt.Errorf("apt-get update: %w", err)
		}
		return system.RunCmd("apt-get", "install", "-y", "nginx")

	case system.RHEL:
		return system.RunCmd("dnf", "install", "-y", "nginx")
	}

	return fmt.Errorf("unsupported OS family")
}

// Enable starts nginx and enables it on boot.
func Enable() error {
	return system.RunCmd("systemctl", "enable", "--now", "nginx")
}

// GenerateConfig renders the nginx config template.
func GenerateConfig(cfg Config) (string, error) {
	tmpl, err := template.New("nginx").Parse(configTemplate)
	if err != nil {
		return "", err
	}

	data := templateData{
		ServerNames: strings.Join(cfg.Domains, " "),
		Upstream:    cfg.Upstream,
		MaxBodySize: cfg.MaxBodySize,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// ConfigPath returns where the config file should live for the given OS and primary domain.
func ConfigPath(osInfo *system.Info, primaryDomain string) string {
	name := primaryDomain + ".conf"
	if osInfo.Family == system.Debian {
		return filepath.Join("/etc/nginx/sites-available", name)
	}
	return filepath.Join("/etc/nginx/conf.d", name)
}

// BackupFile copies path to path.bak.<timestamp> if it exists.
func BackupFile(path string) (string, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read for backup: %w", err)
	}

	backup := path + ".bak." + time.Now().Format("20060102-150405")
	if err := system.AtomicWriteFile(backup, data, system.ModeOrDefault(path, 0644)); err != nil {
		return "", fmt.Errorf("write backup: %w", err)
	}
	return backup, nil
}

// Written records everything WriteConfig changed, so RestoreConfig can undo
// it if the result turns out not to validate.
type Written struct {
	Path       string // config file that was written
	Backup     string // previous contents; empty if the config file was new
	Link       string // sites-enabled symlink touched; empty if none
	LinkTarget string // what that symlink pointed at before; empty if it was new
}

// WriteConfig writes content to the appropriate config path, backs up any
// existing file, and (on Debian-family systems) creates the sites-enabled
// symlink. The returned Written lets the caller roll all of that back — an
// invalid config left in place would stop nginx starting on the next reload
// or reboot, long after the command that wrote it.
func WriteConfig(cfg Config, content string) (Written, error) {
	path := ConfigPath(cfg.OSInfo, cfg.Domains[0])
	w := Written{Path: path}

	backup, err := BackupFile(path)
	if err != nil {
		return w, err
	}
	w.Backup = backup
	if backup != "" {
		fmt.Printf("  Backed up existing config → %s\n", backup)
	}

	// Atomic, so a reload can never pick up a half-written server block, and
	// mode-preserving, so rewriting a config never changes its permissions.
	if err := system.AtomicWriteFile(path, []byte(content), system.ModeOrDefault(path, 0644)); err != nil {
		return w, fmt.Errorf("write config: %w", err)
	}
	fmt.Printf("  Config written → %s\n", path)

	if cfg.OSInfo.Family == system.Debian {
		link, prevTarget, err := enableSite(path, cfg.Domains[0])
		w.Link, w.LinkTarget = link, prevTarget
		if err != nil {
			return w, err
		}
	}

	return w, nil
}

// RestoreConfig undoes WriteConfig. Used when `nginx -t` rejects the result,
// so a run that fails validation leaves nothing behind.
func RestoreConfig(w Written) error {
	if w.Link != "" {
		if err := os.Remove(w.Link); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("remove symlink %s: %w", w.Link, err)
		}
		// Put back whatever the symlink pointed at before, if anything.
		if w.LinkTarget != "" {
			if err := os.Symlink(w.LinkTarget, w.Link); err != nil {
				return fmt.Errorf("restore symlink %s: %w", w.Link, err)
			}
		}
	}

	if w.Path == "" {
		return nil
	}

	// No backup means the config file didn't exist before this run.
	if w.Backup == "" {
		if err := os.Remove(w.Path); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("remove %s: %w", w.Path, err)
		}
		return nil
	}

	data, err := os.ReadFile(w.Backup)
	if err != nil {
		return fmt.Errorf("read backup %s: %w", w.Backup, err)
	}
	if err := system.AtomicWriteFile(w.Path, data, system.ModeOrDefault(w.Backup, 0644)); err != nil {
		return fmt.Errorf("restore %s: %w", w.Path, err)
	}
	return nil
}

// enableSite points sites-enabled at the config, returning the link path and
// whatever it pointed at before so the change can be reversed.
func enableSite(configPath, domain string) (link, prevTarget string, err error) {
	link = filepath.Join("/etc/nginx/sites-enabled", domain+".conf")

	if target, err := os.Readlink(link); err == nil {
		prevTarget = target
	}
	if err := os.Remove(link); err != nil && !os.IsNotExist(err) {
		return link, prevTarget, fmt.Errorf("replace symlink %s: %w", link, err)
	}

	if err := os.Symlink(configPath, link); err != nil {
		return link, prevTarget, fmt.Errorf("enable site (symlink): %w", err)
	}
	fmt.Printf("  Site enabled → %s\n", link)
	return link, prevTarget, nil
}

// Test runs `nginx -t` and returns the combined output on failure.
func Test() error {
	out, err := exec.Command("nginx", "-t").CombinedOutput()
	if err != nil {
		return fmt.Errorf("nginx config test failed:\n%s", string(out))
	}
	return nil
}

// Reload reloads nginx without dropping connections.
func Reload() error {
	out, err := exec.Command("systemctl", "reload", "nginx").CombinedOutput()
	if err != nil {
		return fmt.Errorf("nginx reload failed:\n%s", string(out))
	}
	return nil
}
