package cmd

import (
	"fmt"
	"os"
	"strings"

	survey "github.com/AlecAivazis/survey/v2"
	nodepkg "github.com/lahat-m/serv-kun/internal/node"
	"github.com/lahat-m/serv-kun/internal/system"
	"github.com/spf13/cobra"
)

var nodeCmd = &cobra.Command{
	Use:   "node",
	Short: "Install Node.js via nvm for a Linux user, and pick a default version",
	RunE:  runNode,
}

func runNode(cmd *cobra.Command, _ []string) error {
	requireRoot(cmd)

	printBanner("serv-kun node")

	osInfo, err := system.Detect()
	if err != nil {
		return err
	}

	var username string
	if err := survey.AskOne(&survey.Input{
		Message: "Linux user to install Node.js for:",
		Default: os.Getenv("SUDO_USER"),
		Help:    "nvm and Node will be installed into this user's home directory",
	}, &username, survey.WithValidator(survey.Required)); err != nil {
		return err
	}

	installed, err := nodepkg.IsInstalled(username)
	if err != nil {
		return err
	}

	if !installed {
		fmt.Printf("\nInstalling nvm for %q...\n", username)
		if err := nodepkg.InstallNVM(osInfo, username); err != nil {
			return fmt.Errorf("install nvm: %w", err)
		}
		fmt.Println("✓ nvm installed")
	} else {
		fmt.Println("\n✓ nvm is already installed for this user")
	}
	fmt.Println()

	var versions []string
	message := "Node.js version to install:"
	def := "--lts"
	for {
		var v string
		if err := survey.AskOne(&survey.Input{
			Message: message,
			Default: def,
			Help:    "e.g. --lts, 20, 22.4.0, node (latest)",
		}, &v); err != nil {
			return err
		}
		v = strings.TrimSpace(v)
		if v == "" {
			if len(versions) == 0 {
				return fmt.Errorf("at least one Node.js version is required")
			}
			break
		}

		fmt.Printf("\nInstalling Node.js %s...\n", v)
		if err := nodepkg.InstallVersion(username, v); err != nil {
			return fmt.Errorf("install node %s: %w", v, err)
		}
		versions = append(versions, v)

		message = "Install another version? (leave blank to finish):"
		def = ""
	}

	fmt.Println()
	if out, err := nodepkg.Versions(username); err == nil {
		fmt.Println("Installed versions:")
		fmt.Println(strings.TrimRight(out, "\n"))
		fmt.Println()
	}

	defaultVersion := versions[len(versions)-1]
	if len(versions) > 1 {
		if err := survey.AskOne(&survey.Select{
			Message: "Which version should be the default?",
			Options: versions,
			Default: defaultVersion,
		}, &defaultVersion); err != nil {
			return err
		}
	}

	if err := nodepkg.SetDefault(username, defaultVersion); err != nil {
		return fmt.Errorf("set default version: %w", err)
	}
	fmt.Printf("\n✓ Node.js %s is set as default for %q\n", defaultVersion, username)

	if err := nodepkg.LinkSystemBinaries(username); err != nil {
		return fmt.Errorf("link node/npm/npx into /usr/local/bin: %w", err)
	}
	fmt.Println("✓ node/npm/npx linked into /usr/local/bin — usable right away, no new login needed")

	return nil
}
