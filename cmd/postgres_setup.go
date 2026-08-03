package cmd

import (
	"fmt"

	survey "github.com/AlecAivazis/survey/v2"
	pgpkg "github.com/lahat-m/serv-kun/internal/postgres"
	"github.com/lahat-m/serv-kun/internal/system"
	"github.com/spf13/cobra"
)

var postgresSetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Install PostgreSQL and create a deployment role and database",
	RunE:  runPostgresSetup,
}

func runPostgresSetup(cmd *cobra.Command, _ []string) error {
	requireRoot(cmd)

	printBanner("postgres setup")

	osInfo, err := system.Detect()
	if err != nil {
		return err
	}
	fmt.Printf("Detected OS: %s\n\n", osInfo.Pretty)

	if !pgpkg.IsInstalled() {
		var install bool
		if err := survey.AskOne(&survey.Confirm{
			Message: "PostgreSQL is not installed. Install it now?",
			Default: true,
		}, &install); err != nil {
			return err
		}
		if !install {
			return fmt.Errorf("PostgreSQL is required to continue")
		}
		fmt.Println("Installing PostgreSQL...")
		if err := pgpkg.Install(osInfo); err != nil {
			return fmt.Errorf("install postgresql: %w", err)
		}
		if err := pgpkg.Enable(); err != nil {
			return fmt.Errorf("enable postgresql: %w", err)
		}
		fmt.Println("✓ PostgreSQL installed and started")
	} else {
		fmt.Println("✓ PostgreSQL is already installed")
	}
	fmt.Println()

	var answers struct {
		Role     string
		Password string
		Database string
	}

	qs := []*survey.Question{
		{
			Name: "role",
			Prompt: &survey.Input{
				Message: "Role (username) to create for deployment:",
				Default: "deploy",
			},
			Validate: survey.Required,
		},
		{
			Name: "password",
			Prompt: &survey.Password{
				Message: "Password for this role:",
			},
			Validate: survey.Required,
		},
		{
			Name: "database",
			Prompt: &survey.Input{
				Message: "Database name to create (owned by this role):",
			},
			Validate: survey.Required,
		},
	}
	if err := survey.Ask(qs, &answers); err != nil {
		return err
	}

	fmt.Printf("\nCreating role %q...\n", answers.Role)
	if err := pgpkg.CreateRole(answers.Role, answers.Password); err != nil {
		return fmt.Errorf("create role: %w", err)
	}
	fmt.Println("✓ Role ready")

	fmt.Printf("Creating database %q...\n", answers.Database)
	if err := pgpkg.CreateDatabase(answers.Database, answers.Role); err != nil {
		return fmt.Errorf("create database: %w", err)
	}
	fmt.Println("✓ Database ready")

	fmt.Println("\nEnabling password auth for local TCP connections...")
	if err := pgpkg.EnableLocalPasswordAuth(osInfo); err != nil {
		return fmt.Errorf("configure pg_hba.conf: %w", err)
	}
	if err := pgpkg.Restart(); err != nil {
		return fmt.Errorf("restart postgresql: %w", err)
	}
	fmt.Println("✓ PostgreSQL restarted")

	fmt.Printf("\n✓ Connect with: postgresql://%s:<password>@127.0.0.1:5432/%s\n", answers.Role, answers.Database)
	return nil
}
