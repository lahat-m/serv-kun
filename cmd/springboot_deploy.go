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

	// A container that is already on the image the pull just produced is not a
	// rollback target — retagging it would repoint the tag at itself and the
	// "rollback" would redeploy the identical broken build. That happens
	// routinely when CI has already run `up` and left the app crash-looping,
	// so the operator reaching for this command is looking at a container that
	// is *newer* than the tool assumes.
	var alreadyOnNewImage bool
	if oldImageID != "" {
		if newImageID, err := compose.LocalImageID(appImage); err == nil && newImageID == oldImageID {
			oldImageID = ""
			alreadyOnNewImage = true
			fmt.Println("⚠ The app is already running the image this pull produced — nothing to roll back to.")
		}
	}

	// The image that just arrived may not run as the user the one before it
	// did — a base image bump is enough to renumber it — and certs/ is still
	// owned by the old uid. Reassert ownership against the image about to run,
	// before it gets a chance to crash-loop on a key it cannot open.
	ensureCertsReadable(proj, appImage)

	// --force-recreate so this always replaces the container. Without it,
	// compose leaves a service whose image and config are both unchanged
	// alone and prints "Container <name> Running", which reads as a completed
	// deploy while nothing was deployed. That is not a hypothetical: when CI
	// has already run `up`, or when the container is crash-looping on the
	// image that was just pulled, this branch is the common one, and the
	// health check that follows would be judging a container this command
	// never touched.
	fmt.Println("Recreating app container...")
	if err := proj.Run("up", "-d", "--no-deps", "--force-recreate", "app"); err != nil {
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

	if alreadyOnNewImage {
		return fmt.Errorf("the app was already running %s before this deploy, so there is no earlier image to roll back to — "+
			"the fault is in the image or in %s, not in the upgrade; fix it and re-run this command",
			appImage, proj.EnvFile)
	}
	if oldImageID == "" {
		return fmt.Errorf("no previous image recorded — cannot roll back automatically; investigate with: serv-kun springboot logs")
	}

	fmt.Println("Rolling back to the previous image...")
	if err := compose.Tag(oldImageID, appImage); err != nil {
		return fmt.Errorf("retag previous image: %w", err)
	}
	// Hand certs/ back to the user the *old* image runs as. Skipping this
	// breaks rollback in precisely the case that needs it most: a deploy whose
	// new image renumbered its user, where the chown above has since pointed
	// the keys at a uid the image being restored doesn't have.
	ensureCertsReadable(proj, oldImageID)
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
