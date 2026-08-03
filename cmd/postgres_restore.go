package cmd

import (
	"errors"
	"fmt"
	"os"

	survey "github.com/AlecAivazis/survey/v2"
	pgpkg "github.com/lahat-m/serv-kun/internal/postgres"
	"github.com/spf13/cobra"
)

var postgresRestoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restore a PostgreSQL database from a backup file",
	RunE:  runPostgresRestore,
}

const newDatabaseOption = "<create a new database>"

func runPostgresRestore(cmd *cobra.Command, _ []string) error {
	requireRoot(cmd)

	printBanner("postgres restore")

	if !pgpkg.IsInstalled() {
		return fmt.Errorf("PostgreSQL is not installed — run: sudo serv-kun postgres setup")
	}

	var srcPath string
	if err := survey.AskOne(&survey.Input{
		Message: "Path to the backup file:",
		Help:    "A gzip-compressed plain-SQL file, e.g. one written by 'serv-kun postgres backup'",
	}, &srcPath, survey.WithValidator(validateBackupPath)); err != nil {
		return err
	}

	databases, err := pgpkg.ListDatabases()
	if err != nil {
		return fmt.Errorf("list databases: %w", err)
	}
	options := append([]string{newDatabaseOption}, databases...)

	var choice string
	if err := survey.AskOne(&survey.Select{
		Message: "Restore into which database?",
		Options: options,
	}, &choice); err != nil {
		return err
	}

	database := choice
	if choice == newDatabaseOption {
		var answers struct {
			Database string
			Owner    string
		}
		qs := []*survey.Question{
			{
				Name:     "database",
				Prompt:   &survey.Input{Message: "New database name:"},
				Validate: survey.Required,
			},
			{
				Name: "owner",
				Prompt: &survey.Input{
					Message: "Owner role for the new database:",
					Default: "postgres",
				},
				Validate: survey.Required,
			},
		}
		if err := survey.Ask(qs, &answers); err != nil {
			return err
		}

		fmt.Printf("\nCreating database %q...\n", answers.Database)
		if err := pgpkg.CreateDatabase(answers.Database, answers.Owner); err != nil {
			return fmt.Errorf("create database: %w", err)
		}
		database = answers.Database
	} else {
		var confirm bool
		if err := survey.AskOne(&survey.Confirm{
			Message: fmt.Sprintf("This will overwrite existing data in database %q. Continue?", database),
			Default: false,
		}, &confirm); err != nil {
			return err
		}
		if !confirm {
			return fmt.Errorf("aborted — no changes made")
		}
	}

	fmt.Printf("\nRestoring %s into %q...\n", srcPath, database)
	if err := pgpkg.Restore(database, srcPath); err != nil {
		return fmt.Errorf("restore: %w", err)
	}

	fmt.Printf("\n✓ Restored into %q\n", database)
	return nil
}

func validateBackupPath(ans interface{}) error {
	s, _ := ans.(string)
	if s == "" {
		return errors.New("path is required")
	}
	info, err := os.Stat(s)
	if err != nil {
		return fmt.Errorf("cannot access %q: %w", s, err)
	}
	if info.IsDir() {
		return fmt.Errorf("%q is a directory, not a file", s)
	}
	return nil
}
