package cmd

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/lahat-m/serv-kun/internal/envfile"
	"github.com/spf13/cobra"
)

var springbootStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show container status and check the app's health endpoint",
	RunE:  runSpringbootStatus,
}

func runSpringbootStatus(cmd *cobra.Command, _ []string) error {
	requireRoot(cmd)
	proj := projectFromFlags(cmd)
	printBanner("springboot status")

	if err := proj.Run("ps"); err != nil {
		return fmt.Errorf("compose ps: %w", err)
	}
	fmt.Println()

	env, err := envfile.Read(filepath.Join(proj.Dir, proj.EnvFile))
	if err != nil || env["APP_DOMAIN"] == "" {
		fmt.Println("(no APP_DOMAIN in env file — skipping health check)")
		return nil
	}

	healthURL := "https://" + env["APP_DOMAIN"] + "/actuator/health"
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(healthURL)
	if err != nil {
		fmt.Printf("✗ %s unreachable: %s\n", healthURL, err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Printf("✓ %s → 200\n", healthURL)
	} else {
		fmt.Printf("✗ %s → %d\n", healthURL, resp.StatusCode)
	}
	return nil
}
