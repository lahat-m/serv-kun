package certbot

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/lahat-m/serv-kun/internal/system"
)

func IsInstalled() bool {
	_, err := exec.LookPath("certbot")
	return err == nil
}

func Install(osInfo *system.Info) error {
	switch osInfo.Family {
	case system.Debian:
		if err := system.RunCmd("apt-get", "update", "-qq"); err != nil {
			return fmt.Errorf("apt-get update: %w", err)
		}
		return system.RunCmd("apt-get", "install", "-y", "certbot", "python3-certbot-nginx")

	case system.RHEL:
		// EPEL provides certbot on RHEL/CentOS — ignore error if already present
		_ = system.RunCmd("dnf", "install", "-y", "epel-release")
		return system.RunCmd("dnf", "install", "-y", "certbot", "python3-certbot-nginx")
	}

	return fmt.Errorf("unsupported OS family")
}

// Run calls certbot --nginx for the given domains.
// Pass a non-empty email for expiry notifications, or empty to skip registration.
func Run(domains []string, email string) error {
	if len(domains) == 0 {
		return fmt.Errorf("no domains provided")
	}

	args := []string{
		"--nginx",
		"--non-interactive",
		"--agree-tos",
		"--redirect", // auto-add HTTP→HTTPS redirect
	}

	if email != "" {
		args = append(args, "--email", email)
	} else {
		args = append(args, "--register-unsafely-without-email")
	}

	for _, d := range domains {
		args = append(args, "-d", d)
	}

	cmd := exec.Command("certbot", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
