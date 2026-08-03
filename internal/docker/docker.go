// Package docker installs and checks for Docker Engine + the Compose plugin.
package docker

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/lahat-m/serv-kun/internal/system"
)

// IsInstalled reports whether both the docker binary and the compose plugin are present.
func IsInstalled() bool {
	if _, err := exec.LookPath("docker"); err != nil {
		return false
	}
	return exec.Command("docker", "compose", "version").Run() == nil
}

// Install runs Docker's official convenience script, which detects the distro
// itself and installs Engine + the Compose/Buildx plugins in one step — unlike
// nginx/postgres, there's no need to branch on OS family here.
func Install() error {
	// Only bounds fetching the (small) script — the install it then runs can
	// legitimately take a couple of minutes and is left unbounded.
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get("https://get.docker.com")
	if err != nil {
		return fmt.Errorf("download get.docker.com: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download get.docker.com: unexpected status %s", resp.Status)
	}

	cmd := exec.Command("sh")
	cmd.Stdin = resp.Body
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("run get.docker.com script: %w", err)
	}
	return nil
}

// Enable starts the Docker daemon and enables it on boot.
func Enable() error {
	return system.RunCmd("systemctl", "enable", "--now", "docker")
}
