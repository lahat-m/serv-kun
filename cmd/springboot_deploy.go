package cmd

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/lahat-m/serv-kun/internal/compose"
	"github.com/lahat-m/serv-kun/internal/envfile"
	"github.com/lahat-m/serv-kun/internal/healthcheck"
	"github.com/spf13/cobra"
)

var springbootDeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Pull the latest app image and redeploy, rolling back automatically on a failed health check",
	RunE:  runSpringbootDeploy,
}

const (
	deployHealthTimeout  = 90 * time.Second
	deployHealthInterval = 3 * time.Second
)

func runSpringbootDeploy(cmd *cobra.Command, _ []string) error {
	requireRoot(cmd)
	printBanner("springboot deploy")

	proj := projectFromFlags(cmd)

	envPath := filepath.Join(proj.Dir, proj.EnvFile)
	env, err := envfile.Read(envPath)
	if err != nil {
		return fmt.Errorf("%w — run 'serv-kun springboot setup' first", err)
	}
	domain := env["APP_DOMAIN"]
	appImage := env["APP_IMAGE"]
	if domain == "" || appImage == "" {
		return fmt.Errorf("%s is missing APP_DOMAIN or APP_IMAGE", proj.EnvFile)
	}
	healthURL := "https://" + domain + "/actuator/health"

	// Capture the currently-running image ID *before* pulling — pulling a
	// floating tag (e.g. :latest) repoints it, but the old image blob stays
	// on disk under this ID until something prunes it, which is what makes
	// rollback possible.
	var oldImageID string
	if containerID, err := proj.ContainerID("app"); err == nil {
		if id, err := compose.ImageID(containerID); err == nil {
			oldImageID = id
		}
	}
	if oldImageID == "" {
		fmt.Println("⚠ No previous image found — this deploy has no automatic rollback safety net.")
	}

	fmt.Println("Pulling latest image...")
	if err := proj.Run("pull", "app"); err != nil {
		return fmt.Errorf("pull: %w", err)
	}

	fmt.Println("Recreating app container...")
	if err := proj.Run("up", "-d", "--no-deps", "app"); err != nil {
		return fmt.Errorf("up: %w", err)
	}

	fmt.Printf("Waiting for %s to become healthy...\n", healthURL)
	err = healthcheck.WaitHealthy(healthURL, deployHealthTimeout, deployHealthInterval)
	if err == nil {
		fmt.Println("\n✓ Deployed and healthy")
		return nil
	}

	fmt.Printf("\n✗ New deploy failed health check: %s\n", err)
	fmt.Println("\nLast 50 lines of app logs:")
	printRule()
	if err := proj.Run("logs", "--tail=50", "app"); err != nil {
		fmt.Printf("(could not fetch logs: %s)\n", err)
	}
	printRule()

	if oldImageID == "" {
		return fmt.Errorf("no previous image recorded — cannot roll back automatically; investigate with: serv-kun springboot logs")
	}

	fmt.Println("Rolling back to the previous image...")
	if err := compose.Tag(oldImageID, appImage); err != nil {
		return fmt.Errorf("retag previous image: %w", err)
	}
	if err := proj.Run("up", "-d", "--no-deps", "--force-recreate", "app"); err != nil {
		return fmt.Errorf("recreate with previous image: %w", err)
	}

	fmt.Printf("Waiting for %s to become healthy again...\n", healthURL)
	if err := healthcheck.WaitHealthy(healthURL, deployHealthTimeout, deployHealthInterval); err != nil {
		return fmt.Errorf("rollback also failed health check: %w — needs manual investigation", err)
	}

	fmt.Println("\n✓ Rolled back — previous version is healthy again")
	return fmt.Errorf("deploy of %s failed and was rolled back", appImage)
}
