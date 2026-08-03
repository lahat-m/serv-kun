package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	survey "github.com/AlecAivazis/survey/v2"
	pgpkg "github.com/lahat-m/serv-kun/internal/postgres"
	"github.com/spf13/cobra"
)

var postgresBackupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Back up a PostgreSQL database to a file",
	RunE:  runPostgresBackup,
}

func runPostgresBackup(cmd *cobra.Command, _ []string) error {
	requireRoot(cmd)

	printBanner("postgres backup")

	if !pgpkg.IsInstalled() {
		return fmt.Errorf("PostgreSQL is not installed — run: sudo serv-kun postgres setup")
	}

	databases, err := pgpkg.ListDatabases()
	if err != nil {
		return fmt.Errorf("list databases: %w", err)
	}
	if len(databases) == 0 {
		return fmt.Errorf("no databases found")
	}

	var database string
	if err := survey.AskOne(&survey.Select{
		Message: "Database to back up:",
		Options: databases,
	}, &database); err != nil {
		return err
	}

	defaultPath := filepath.Join("/var/backups/postgres", fmt.Sprintf("%s-%s.sql.gz", database, time.Now().Format("20060102-150405")))

	var destPath string
	if err := survey.AskOne(&survey.Input{
		Message: "Backup file path:",
		Default: defaultPath,
		Help:    "Gzip-compressed plain SQL — restore it with 'serv-kun postgres restore'. Not encrypted: handle it as sensitive as the database itself",
	}, &destPath, survey.WithValidator(survey.Required)); err != nil {
		return err
	}

	fmt.Printf("\nBacking up %q to %s...\n", database, destPath)
	if err := pgpkg.Backup(database, destPath); err != nil {
		return fmt.Errorf("backup: %w", err)
	}

	size := ""
	if info, err := os.Stat(destPath); err == nil {
		size = fmt.Sprintf(" (%.1f MB)", float64(info.Size())/1024/1024)
	}

	fmt.Printf("\n✓ Backup written → %s%s\n", destPath, size)
	return nil
}
