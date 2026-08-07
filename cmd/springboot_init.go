package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/lahat-m/serv-kun/internal/compose"
	"github.com/lahat-m/serv-kun/internal/scaffold"
	"github.com/spf13/cobra"
)

var springbootInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Scaffold docker-compose.prod.yml, nginx config, backup script, and CI/CD for a new project",
	RunE:  runSpringbootInit,
}

func runSpringbootInit(cmd *cobra.Command, _ []string) error {
	printBanner("springboot init")

	proj := projectFromFlags(cmd)

	var answers struct {
		AppName     string
		AppPort     string
		MaxBodySize string
	}
	if err := survey.Ask([]*survey.Question{
		{
			Name:     "appName",
			Prompt:   &survey.Input{Message: "App name:", Help: "lowercase, used for container names (e.g. \"myapp\" -> myapp-app, myapp-db, ...)"},
			Validate: survey.Required,
		},
		{
			Name:     "appPort",
			Prompt:   &survey.Input{Message: "Port the app listens on inside its container:", Default: "8080"},
			Validate: survey.Required,
		},
		{
			Name:   "maxBodySize",
			Prompt: &survey.Input{Message: "Max upload size (nginx client_max_body_size):", Default: "50M", Help: "should match the app's multipart max-request-size"},
		},
	}, &answers); err != nil {
		return err
	}

	var includeCerts bool
	if err := survey.AskOne(&survey.Confirm{
		Message: "Mount certs/ into the app container (for RS256 JWT keys)?",
		Default: true,
	}, &includeCerts); err != nil {
		return err
	}
	fmt.Println()

	projectName := initProjectName(proj.Dir, answers.AppName)

	files, err := scaffold.Generate(scaffold.Params{
		AppName:      answers.AppName,
		AppPort:      answers.AppPort,
		MaxBodySize:  answers.MaxBodySize,
		IncludeCerts: includeCerts,
		ProjectName:  projectName,
	})
	if err != nil {
		return err
	}

	fmt.Printf("Compose project: %s — volumes will be named %s_*\n\n", projectName, projectName)

	for _, f := range files {
		fullPath := filepath.Join(proj.Dir, f.Path)

		if _, err := os.Stat(fullPath); err == nil {
			// Overwriting a compose file that is already running under a
			// different project name renames the namespace its volumes live
			// in. The stack comes back up on empty ones and the real data sits
			// orphaned under the old prefix, which looks like data loss and is
			// not obviously reversible to whoever is staring at it at 2am.
			if f.Path == scaffold.ComposeFile {
				if current := compose.ProjectNameOf(proj.Dir, f.Path); current != "" && current != projectName {
					fmt.Printf("\n⚠ The existing %s runs as compose project %q; this one would be %q.\n", f.Path, current, projectName)
					fmt.Printf("  Overwriting renames the volume prefix: %s_postgres-data and every other\n", current)
					fmt.Printf("  %s_* volume would be left behind, and the stack would start on empty\n", current)
					fmt.Println("  ones. If this project is already deployed, skip this file, or keep the")
					fmt.Printf("  old name by editing `name:` back to %q before the next up.\n\n", current)
				}
			}

			var overwrite bool
			if err := survey.AskOne(&survey.Confirm{
				Message: fmt.Sprintf("%s already exists — overwrite (a timestamped backup is kept)?", f.Path),
				Default: false,
			}, &overwrite); err != nil {
				return err
			}
			if !overwrite {
				fmt.Printf("- skipped %s\n", f.Path)
				continue
			}
			backup := fullPath + ".bak." + time.Now().Format("20060102-150405")
			if err := os.Rename(fullPath, backup); err != nil {
				return fmt.Errorf("back up existing %s: %w", f.Path, err)
			}
		}

		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			return fmt.Errorf("create directory for %s: %w", f.Path, err)
		}

		mode := os.FileMode(0644)
		if f.Exec {
			mode = 0755
		}
		if err := os.WriteFile(fullPath, []byte(f.Content), mode); err != nil {
			return fmt.Errorf("write %s: %w", f.Path, err)
		}
		fmt.Printf("✓ %s\n", f.Path)
	}

	fmt.Println("\nNext: add a Dockerfile if this project doesn't have one yet, then run  sudo serv-kun springboot setup")
	return nil
}

// initProjectName picks the compose project name the generated file declares.
//
// It is the directory name with "-prod" appended, because the directory name
// is exactly what compose would otherwise default *both* this stack and a dev
// compose file in the same directory to. Appending the suffix is what makes
// the two namespaces distinct, and doing it in the file rather than through
// -p means a raw "docker compose -f docker-compose.prod.yml" gets it too —
// which matters, since the commands that cause the damage are the ones typed
// by hand, not the ones this tool runs.
func initProjectName(dir, appName string) string {
	name := compose.DefaultProjectName(dir)
	if name == "" {
		// A directory whose name survives none of the normalization (say, one
		// named in a script with no ASCII in it) still needs a project name.
		name = compose.NormalizeProjectName(appName)
	}
	if name == "" {
		return "app-prod"
	}
	return name + "-prod"
}
