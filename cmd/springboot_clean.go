package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/core"
	"github.com/lahat-m/serv-kun/internal/compose"
	"github.com/spf13/cobra"
)

var springbootCleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Delete the stack's containers, images, and build cache, keeping every volume",
	Long: `clean takes the stack down and deletes everything it can be rebuilt from:
the containers, the images behind every service, and the builder cache.

Run interactively, it asks which of those to delete, and lists every named
volume alongside them. Nothing is selected for a volume unless you select it:
the volumes are the only thing on the host that cannot be recovered by
pulling again, which is why this command exists separately from
"docker compose down -v" — -v takes postgres-data, postgres-backups, and
certbot-certs with it, the mistake the README warns about and that
"springboot tls" exists to partially undo. Selecting a volume asks again,
by name, before anything is deleted.

With --yes there is no selection: it deletes the containers, images,
builder cache, and dangling images, and never a volume.

Reach for this when the host has run out of disk, or when a pull has left
the content store in a state that fails on a blob that should already be
cached ("failed commit on ref ... no such file or directory"). Deleting the
images forces the next deploy to fetch them cleanly.

Deleting the containers stops the site, and it stays down until the stack is
brought back with "docker compose up -d" — not "springboot deploy", which
recreates only the app service. clean prints the exact commands when it
finishes. The images have to come back down over the network before the
app is up again. Pruning only the builder
cache and dangling images leaves the site running.

The builder cache and dangling images are host-wide, not scoped to this
project: on a box running more than this one stack, that reclaims space
belonging to the others too.`,
	Args: cobra.NoArgs,
	RunE: runSpringbootClean,
}

func init() {
	springbootCleanCmd.Flags().BoolP("yes", "y", false,
		"skip the prompts and delete everything but the volumes (required when stdin is not a terminal)")
}

// composeVolume is a named volume as compose resolves it: Key is what the
// compose file calls it, Name is what docker calls it on the host (normally
// "<project>_<key>"), and External marks one compose neither created nor
// owns.
type composeVolume struct {
	Key      string
	Name     string
	External bool
}

// cleanPlan is what one run of clean deletes. Everything not in it is kept.
type cleanPlan struct {
	Containers   bool
	Images       bool
	BuilderCache bool
	Dangling     bool
	Volumes      []composeVolume
}

func (p cleanPlan) empty() bool {
	return !p.Containers && !p.Images && !p.BuilderCache && !p.Dangling && len(p.Volumes) == 0
}

const (
	cleanOptContainers = "containers       every container in this project (takes the site down)"
	cleanOptImages     = "images           every image its services run (app, postgres, nginx, certbot)"
	cleanOptCache      = "builder cache    (host-wide)"
	cleanOptDangling   = "dangling images  (host-wide)"
	cleanOptVolumePfx  = "volume           "
)

func volumeOption(v composeVolume) string {
	return cleanOptVolumePfx + v.Key + "  (data — cannot be pulled back)"
}

// cleanOptions returns the choices offered and the ones selected by default.
// The default is exactly what --yes deletes, so pressing enter on the prompt
// and passing --yes do the same thing. External volumes are not offered: they
// were made outside this project, and clean deleting one would surprise
// whoever made it.
func cleanOptions(volumes []composeVolume) (options, defaults []string) {
	options = []string{cleanOptContainers, cleanOptImages, cleanOptCache, cleanOptDangling}
	defaults = append([]string(nil), options...)
	for _, v := range volumes {
		if !v.External {
			options = append(options, volumeOption(v))
		}
	}
	return options, defaults
}

// planFromSelection turns the prompt's answers back into a plan.
func planFromSelection(selected []string, volumes []composeVolume) cleanPlan {
	chosen := make(map[string]bool, len(selected))
	for _, s := range selected {
		chosen[s] = true
	}
	plan := cleanPlan{
		Containers:   chosen[cleanOptContainers],
		Images:       chosen[cleanOptImages],
		BuilderCache: chosen[cleanOptCache],
		Dangling:     chosen[cleanOptDangling],
	}
	for _, v := range volumes {
		if !v.External && chosen[volumeOption(v)] {
			plan.Volumes = append(plan.Volumes, v)
		}
	}
	return plan
}

// validate rejects a plan docker would refuse halfway through: an image or a
// volume still in use by a container cannot be removed, so either one needs
// the containers gone first.
func (p cleanPlan) validate() error {
	if p.Containers {
		return nil
	}
	if p.Images {
		return errors.New("images can only be deleted along with the containers that run them")
	}
	if len(p.Volumes) > 0 {
		return errors.New("volumes can only be deleted along with the containers that mount them")
	}
	return nil
}

func runSpringbootClean(cmd *cobra.Command, _ []string) error {
	requireRoot(cmd)
	printBanner("springboot clean")

	proj := projectFromFlags(cmd)
	warnProjectNameCollision(proj)

	// Read the volumes out of the compose file rather than naming them in a
	// literal here: a project that added a volume this tool has never heard of
	// should still see it listed, and a wrong list is worse than no list on a
	// prompt whose only job is to be believed.
	volumes, err := projectVolumes(proj)
	if err != nil {
		return err
	}

	_, defaults := cleanOptions(volumes)
	plan := planFromSelection(defaults, volumes)

	yes, _ := cmd.Flags().GetBool("yes")
	if !yes {
		if !stdinIsTerminal() {
			return fmt.Errorf("refusing to delete without confirmation and stdin is not a terminal\n\n"+
				"Re-run with --yes to confirm non-interactively:\n  sudo %s --yes", cmd.CommandPath())
		}
		plan, err = askCleanPlan(volumes)
		if err != nil {
			return err
		}
		if plan.empty() {
			fmt.Println("\nNothing selected. Nothing deleted.")
			return nil
		}
	}

	printCleanPlan(plan, volumes, proj.File)

	if !yes {
		confirmed, err := confirmCleanPlan(plan)
		if err != nil {
			return err
		}
		if !confirmed {
			fmt.Println("\nNothing deleted.")
			return nil
		}
		fmt.Println()
	}

	return executeCleanPlan(plan, proj)
}

func askCleanPlan(volumes []composeVolume) (cleanPlan, error) {
	options, defaults := cleanOptions(volumes)
	var selected []string
	err := survey.AskOne(&survey.MultiSelect{
		Message: "Select what to delete (anything left unselected is kept):",
		Options: options,
		Default: defaults,
		Help:    "space toggles, enter confirms. Volumes are unselected unless you select them.",
	}, &selected, survey.WithValidator(func(ans interface{}) error {
		answers, _ := ans.([]core.OptionAnswer)
		names := make([]string, len(answers))
		for i, a := range answers {
			names[i] = a.Value
		}
		return planFromSelection(names, volumes).validate()
	}), survey.WithPageSize(len(options)))
	if err != nil {
		return cleanPlan{}, err
	}
	return planFromSelection(selected, volumes), nil
}

// confirmCleanPlan asks once for a plan that only deletes what can be pulled
// or rebuilt again, and makes the operator type the volume names back for one
// that deletes data: a y pressed out of habit should not be enough to lose a
// database.
func confirmCleanPlan(plan cleanPlan) (bool, error) {
	if len(plan.Volumes) == 0 {
		var confirmed bool
		err := survey.AskOne(&survey.Confirm{
			Message: "Delete the selected items?",
			Default: false,
		}, &confirmed)
		return confirmed, err
	}

	keys := make([]string, len(plan.Volumes))
	for i, v := range plan.Volumes {
		keys[i] = v.Key
	}
	want := strings.Join(keys, " ")
	var typed string
	err := survey.AskOne(&survey.Input{
		Message: fmt.Sprintf("This permanently deletes the data in %s. Type %q to confirm:", want, want),
	}, &typed)
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(typed) == want, nil
}

func printCleanPlan(plan cleanPlan, volumes []composeVolume, composeFile string) {
	var deletes, keeps []string
	sortInto := func(selected bool, label string) {
		if selected {
			deletes = append(deletes, label)
		} else {
			keeps = append(keeps, label)
		}
	}
	sortInto(plan.Containers, "every container in this project")
	sortInto(plan.Images, "every image its services run (app, postgres, nginx, certbot)")
	sortInto(plan.BuilderCache, "the builder cache (host-wide)")
	sortInto(plan.Dangling, "dangling images (host-wide)")

	deleted := make(map[string]bool, len(plan.Volumes))
	for _, v := range plan.Volumes {
		deleted[v.Key] = true
	}
	for _, v := range volumes {
		label := v.Key
		if v.External {
			label += " (external)"
		}
		sortInto(deleted[v.Key], label)
	}

	fmt.Println()
	fmt.Println("Deletes:")
	for _, d := range deletes {
		fmt.Printf("  - %s\n", d)
	}
	fmt.Println()
	fmt.Println("Keeps:")
	if len(volumes) == 0 {
		// Not reassuring, so don't dress it up: either the file declares no
		// volumes or compose could not be asked, and the operator should know
		// which before answering.
		fmt.Println("  - (no named volumes declared in " + composeFile + ")")
	}
	for _, k := range keeps {
		fmt.Printf("  - %s\n", k)
	}
	fmt.Println()
	if plan.Containers {
		fmt.Println("The site is down from here until the stack is brought back up.")
	} else {
		fmt.Println("The containers are kept, so the site stays up.")
	}
	fmt.Println()
}

func executeCleanPlan(plan cleanPlan, proj compose.Project) error {
	if plan.Containers {
		// --remove-orphans so a service that was renamed or dropped out of the
		// compose file since the last "up" does not keep a volume mounted, or
		// an image referenced, that this run was asked to delete.
		args := []string{"down", "--remove-orphans"}
		if plan.Images {
			fmt.Println("Taking the stack down and removing its images...")
			args = append(args, "--rmi", "all")
		} else {
			fmt.Println("Taking the stack down...")
		}
		if err := proj.Run(args...); err != nil {
			return fmt.Errorf("compose down: %w", err)
		}
	}

	// Removed by name, one at a time, rather than with "down -v": -v takes
	// every volume in the file, and the point of selecting is that it doesn't.
	for _, v := range plan.Volumes {
		if exec.Command("docker", "volume", "inspect", v.Name).Run() != nil {
			fmt.Printf("\nVolume %s (%s) does not exist on this host — skipping.\n", v.Key, v.Name)
			continue
		}
		fmt.Printf("\nDeleting volume %s (%s)...\n", v.Key, v.Name)
		if err := runDocker("volume", "rm", v.Name); err != nil {
			return fmt.Errorf("volume rm %s: %w", v.Name, err)
		}
	}

	// Pruned after the down, not before: the images the services were using
	// are only removable once nothing references them.
	if plan.BuilderCache {
		fmt.Println("\nPruning the builder cache...")
		if err := runDocker("builder", "prune", "-af"); err != nil {
			return fmt.Errorf("builder prune: %w", err)
		}
	}

	if plan.Dangling {
		fmt.Println("\nPruning dangling images...")
		if err := runDocker("image", "prune", "-f"); err != nil {
			return fmt.Errorf("image prune: %w", err)
		}
	}

	fmt.Println()
	printRule()
	if len(plan.Volumes) == 0 {
		fmt.Println("Selected items deleted. Volumes untouched.")
	} else {
		fmt.Println("Selected items deleted, including the selected volumes.")
	}
	if !plan.Containers {
		return nil
	}
	fmt.Println()
	fmt.Println("Bring it back with:")
	for _, step := range cleanNextSteps(plan, proj) {
		fmt.Printf("  %s\n", step)
	}
	if plan.Images {
		fmt.Println()
		fmt.Println("The first 'up' after a clean re-pulls every image, so it takes")
		fmt.Println("markedly longer than usual and needs the registry to be reachable.")
	}
	if plan.deletes("postgres-data") {
		fmt.Println()
		fmt.Println("postgres-data was deleted, so postgres starts with an empty database.")
	}
	return nil
}

// deletes reports whether the plan removes the volume the compose file calls
// key.
func (p cleanPlan) deletes(key string) bool {
	for _, v := range p.Volumes {
		if v.Key == key {
			return true
		}
	}
	return false
}

// cleanNextSteps is what to run after a clean that took the containers down.
//
// Not "springboot deploy": that recreates only the app service (--no-deps),
// so postgres, nginx, and certbot would stay gone. A plain "up -d" recreates
// every service and pulls whatever images are missing.
//
// "springboot tls" goes first when the certificates were deleted, because
// nginx refuses to start without the certificate its config names, and tls
// is what gets nginx up on a placeholder and then issues the real one.
func cleanNextSteps(plan cleanPlan, proj compose.Project) []string {
	// Absolute, because --dir defaults to "." and "cd ." tells the operator
	// nothing once they have moved on to another directory.
	dir := proj.Dir
	if abs, err := filepath.Abs(dir); err == nil {
		dir = abs
	}
	// The serv-kun steps need the same project this clean ran against, so
	// non-default file names are carried over rather than silently dropped.
	flags := "--dir " + dir
	if proj.File != "docker-compose.prod.yml" {
		flags += " --compose-file " + proj.File
	}
	if proj.EnvFile != ".env.prod" {
		flags += " --env-file " + proj.EnvFile
	}

	var steps []string
	if plan.deletes("certbot-certs") {
		steps = append(steps, "sudo serv-kun springboot tls "+flags)
	}
	return append(steps,
		"cd "+dir,
		fmt.Sprintf("sudo docker compose -f %s --env-file %s up -d", proj.File, proj.EnvFile),
		"sudo serv-kun springboot status "+flags,
	)
}

// projectVolumes lists the named volumes the compose file declares, with the
// names docker knows them by.
//
// A compose file that cannot be parsed is a hard error rather than an empty
// list: every later step runs against that same file, so failing here reports
// the problem while nothing has been deleted yet.
func projectVolumes(proj compose.Project) ([]composeVolume, error) {
	out, err := proj.Output("config", "--format", "json")
	if err != nil {
		return nil, fmt.Errorf("could not read the volumes from %s: %w\n%s", proj.File, err, out)
	}
	return parseComposeVolumes(out)
}

// parseComposeVolumes reads the top-level volumes out of "compose config
// --format json". Compose fills in each volume's host name there, so no
// guessing at the project-name prefix is needed.
func parseComposeVolumes(configJSON string) ([]composeVolume, error) {
	var cfg struct {
		Volumes map[string]struct {
			Name     string `json:"name"`
			External bool   `json:"external"`
		} `json:"volumes"`
	}
	if err := json.Unmarshal([]byte(configJSON), &cfg); err != nil {
		return nil, fmt.Errorf("could not parse compose config: %w", err)
	}
	volumes := make([]composeVolume, 0, len(cfg.Volumes))
	for key, v := range cfg.Volumes {
		name := v.Name
		if name == "" {
			name = key
		}
		volumes = append(volumes, composeVolume{Key: key, Name: name, External: v.External})
	}
	sort.Slice(volumes, func(i, j int) bool { return volumes[i].Key < volumes[j].Key })
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
