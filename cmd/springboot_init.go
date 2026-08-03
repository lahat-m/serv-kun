package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	survey "github.com/AlecAivazis/survey/v2"
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

	files, err := scaffold.Generate(scaffold.Params{
		AppName:      answers.AppName,
		AppPort:      answers.AppPort,
		MaxBodySize:  answers.MaxBodySize,
		IncludeCerts: includeCerts,
	})
	if err != nil {
		return err
	}

	for _, f := range files {
		fullPath := filepath.Join(proj.Dir, f.Path)

		if _, err := os.Stat(fullPath); err == nil {
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
