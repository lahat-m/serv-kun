package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/lahat-m/serv-kun/internal/compose"
	"github.com/spf13/cobra"
)

var springbootCleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Delete the stack's containers, images, and build cache, keeping every volume",
	Long: `clean takes the stack down and deletes everything it can be rebuilt from:
the containers, the images behind every service, and the builder cache.

It never deletes a volume. That is the whole point of the command existing
separately from "docker compose down -v": the volumes are the only thing on
the host that cannot be recovered by pulling again, and -v takes
postgres-data, postgres-backups, and certbot-certs with it — the mistake the
README warns about and that "springboot tls" exists to partially undo.

Reach for this when the host has run out of disk, or when a pull has left
the content store in a state that fails on a blob that should already be
cached ("failed commit on ref ... no such file or directory"). Deleting the
images forces the next deploy to fetch them cleanly.

This stops the site. Nothing is serving between this command and the next
"springboot deploy", so the images have to come back down over the network
before the app is up again.

The builder cache and dangling images are host-wide, not scoped to this
project: on a box running more than this one stack, that reclaims space
belonging to the others too.`,
	Args: cobra.NoArgs,
	RunE: runSpringbootClean,
}

func init() {
	springbootCleanCmd.Flags().BoolP("yes", "y", false,
		"skip the confirmation prompt (required when stdin is not a terminal)")
}

func runSpringbootClean(cmd *cobra.Command, _ []string) error {
	requireRoot(cmd)
	printBanner("springboot clean")

	proj := projectFromFlags(cmd)
	warnProjectNameCollision(proj)

	// Read the volume names out of the compose file rather than naming them in
	// a literal here: a project that added a volume this tool has never heard
	// of should still see it listed among what survives, and a wrong list is
	// worse than no list on a prompt whose only job is to be believed.
	volumes, err := survivingVolumes(proj)
	if err != nil {
		return err
	}

	fmt.Println("Deletes:")
	fmt.Println("  - every container in this project")
	fmt.Println("  - every image its services run (app, postgres, nginx, certbot)")
	fmt.Println("  - the builder cache and dangling images (host-wide)")
	fmt.Println()
	fmt.Println("Keeps:")
	if len(volumes) == 0 {
		// Not reassuring, so don't dress it up: either the file declares no
		// volumes or compose could not be asked, and the operator should know
		// which before answering.
		fmt.Println("  - (no named volumes declared in " + proj.File + ")")
	}
	for _, v := range volumes {
		fmt.Printf("  - %s\n", v)
	}
	fmt.Println()
	fmt.Println("The site is down from here until the next 'springboot deploy'.")
	fmt.Println()

	yes, _ := cmd.Flags().GetBool("yes")
	if !yes {
		if !stdinIsTerminal() {
			return fmt.Errorf("refusing to delete without confirmation and stdin is not a terminal\n\n"+
				"Re-run with --yes to confirm non-interactively:\n  sudo %s --yes", cmd.CommandPath())
		}
		var confirmed bool
		prompt := &survey.Confirm{
			Message: "Delete this project's containers and images?",
			Default: false,
		}
		if err := survey.AskOne(prompt, &confirmed); err != nil {
			return err
		}
		if !confirmed {
			fmt.Println("\nNothing deleted.")
			return nil
		}
		fmt.Println()
	}

	// --remove-orphans so a service that was renamed or dropped out of the
	// compose file since the last "up" does not survive the one command whose
	// job is to leave nothing running.
	fmt.Println("Taking the stack down and removing its images...")
	if err := proj.Run("down", "--rmi", "all", "--remove-orphans"); err != nil {
		return fmt.Errorf("compose down: %w", err)
	}

	// Pruned after the down, not before: the images the services were using
	// are only removable once nothing references them.
	fmt.Println("\nPruning the builder cache...")
	if err := runDocker("builder", "prune", "-af"); err != nil {
		return fmt.Errorf("builder prune: %w", err)
	}

	fmt.Println("\nPruning dangling images...")
	if err := runDocker("image", "prune", "-f"); err != nil {
		return fmt.Errorf("image prune: %w", err)
	}

	fmt.Println()
	printRule()
	fmt.Println("Containers and images deleted. Volumes untouched.")
	fmt.Println()
	fmt.Println("Bring it back with:")
	fmt.Printf("  sudo serv-kun springboot deploy --dir %s\n", proj.Dir)
	fmt.Println()
	fmt.Println("The first deploy after a clean re-pulls every image, so it takes")
	fmt.Println("markedly longer than usual and needs the registry to be reachable.")
	return nil
}

// survivingVolumes lists the named volumes the compose file declares — the
// things clean deliberately leaves behind.
//
// A compose file that cannot be parsed is a hard error rather than an empty
// list: every later step runs against that same file, so failing here reports
// the problem while nothing has been deleted yet.
func survivingVolumes(proj compose.Project) ([]string, error) {
	out, err := proj.Output("config", "--volumes")
	if err != nil {
		return nil, fmt.Errorf("could not read the volumes from %s: %w\n%s", proj.File, err, out)
	}
	var volumes []string
	for _, line := range strings.Split(out, "\n") {
		if v := strings.TrimSpace(line); v != "" {
			volumes = append(volumes, v)
		}
	}
	return volumes, nil
}

// runDocker runs a plain docker command, streaming it to the operator. The
// prune subcommands are not compose commands, so they cannot go through
// compose.Project, which prefixes every invocation with the project's -f and
// --env-file.
func runDocker(args ...string) error {
	cmd := exec.Command("docker", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// stdinIsTerminal reports whether there is a human on the other end to answer
// a prompt. Without this, clean run from the CI deploy step over SSH would
// block on survey forever instead of failing with something an operator can
// read in the job log.
func stdinIsTerminal() bool {
	info, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeCharDevice != 0
}
