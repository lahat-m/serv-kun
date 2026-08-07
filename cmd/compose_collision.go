package cmd

import (
	"fmt"
	"strings"

	"github.com/lahat-m/serv-kun/internal/compose"
)

// warnProjectNameCollision reports another compose file in the project
// directory that resolves to the same compose project name as this one.
//
// It is worth a warning of its own because the failure is invisible until it
// is expensive. Compose matches containers by project *and service* name, not
// by container_name, so a dev docker-compose.yml sitting beside the
// production file — same directory, therefore the same default project name —
// makes an ordinary `docker compose up -d` recreate the production database
// container as if it were the dev one, and `--remove-orphans` delete the rest
// of the stack outright. Nothing in docker's output says what it just did.
//
// It only ever prints. The obvious repair — renaming this project — is the
// one an operator must not be nudged into blindly, because the project name
// prefixes every volume name: renaming it strands the data and brings the
// stack back up empty. That trade-off is the operator's to make, so this
// spells it out and stops.
func warnProjectNameCollision(proj compose.Project) {
	collisions := compose.Collisions(proj.Dir, proj.File)
	if len(collisions) == 0 {
		return
	}
	name := compose.ProjectNameOf(proj.Dir, proj.File)

	for _, c := range collisions {
		fmt.Printf("⚠ %s and %s both resolve to the compose project %q — they share\n", proj.File, c.File, name)
		fmt.Println("  one container namespace.")
		fmt.Println()
		if len(c.Services) > 0 {
			fmt.Printf("  Both declare: %s\n", strings.Join(c.Services, ", "))
			fmt.Println()
		}
		fmt.Printf("  Compose finds containers by project and service name, not by container_name,\n")
		fmt.Printf("  so `docker compose -f %s up -d` recreates the containers for those\n", c.File)
		fmt.Println("  services out from under this stack, and --remove-orphans removes the rest.")
		fmt.Println()
		fmt.Printf("  Fix it on the other stack: give %s a distinct top-level `name:`.\n", c.File)
		fmt.Printf("  Renaming *this* project instead would repoint its volume prefix — %s_*\n", name)
		fmt.Println("  volumes would be left behind and recreated empty on the next start.")
		fmt.Println()
		fmt.Printf("  (If both files are meant to run together as one stack, this is expected.)\n")
		fmt.Println()
	}
}
