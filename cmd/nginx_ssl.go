package cmd

import (
	"fmt"
	"strings"

	survey "github.com/AlecAivazis/survey/v2"
	dnspkg "github.com/lahat-m/serv-kun/internal/dns"
	nginxpkg "github.com/lahat-m/serv-kun/internal/nginx"
	"github.com/lahat-m/serv-kun/internal/system"
	"github.com/spf13/cobra"
)

var nginxSSLCmd = &cobra.Command{
	Use:   "ssl",
	Short: "Add SSL to an existing nginx config via Let's Encrypt certbot",
	RunE:  runNginxSSL,
}

func runNginxSSL(cmd *cobra.Command, _ []string) error {
	requireRoot(cmd)

	printBanner("nginx ssl")

	osInfo, err := system.Detect()
	if err != nil {
		return err
	}
	fmt.Printf("Detected OS: %s\n\n", osInfo.Pretty)

	// default config path varies by OS family
	defaultConfig := "/etc/nginx/sites-available/default"
	if osInfo.Family == system.RHEL {
		defaultConfig = "/etc/nginx/conf.d/default.conf"
	}

	var answers struct {
		ConfigPath string
		Domains    string
	}

	qs := []*survey.Question{
		{
			Name: "configPath",
			Prompt: &survey.Input{
				Message: "Path to your nginx config file:",
				Default: defaultConfig,
				Help:    "This file will be backed up before certbot modifies it",
			},
		},
		{
			Name: "domains",
			Prompt: &survey.Input{
				Message: "Domain(s) to certify:",
				Help:    "Space-separated — must match the server_name in your config",
			},
			Validate: survey.Required,
		},
	}

	if err := survey.Ask(qs, &answers); err != nil {
		return err
	}

	domains := strings.Fields(answers.Domains)

	// ── DNS check ──────────────────────────────────────────────────────────
	fmt.Println("\nResolving server public IP...")
	serverIP, err := system.GetPublicIP()
	if err != nil {
		return fmt.Errorf("could not determine server public IP: %w", err)
	}
	fmt.Printf("Server IP: %s\n\n", serverIP)

	results := printDNSResults(domains, serverIP)
	allOK := dnspkg.AllOK(results)

	if !allOK {
		var proceed bool
		if err := survey.AskOne(&survey.Confirm{
			Message: "Some domains do not point to this server. certbot will fail for those. Proceed anyway?",
			Default: false,
		}, &proceed); err != nil {
			return err
		}
		if !proceed {
			return fmt.Errorf("aborted — fix your DNS A records and retry")
		}
	}

	// ── backup config ──────────────────────────────────────────────────────
	if answers.ConfigPath != "" {
		backup, err := nginxpkg.BackupFile(answers.ConfigPath)
		if err != nil {
			return err
		}
		if backup != "" {
			fmt.Printf("  Backed up config → %s\n\n", backup)
		}
	}

	return runSSLSetup(domains, osInfo)
}

// printDNSResults prints a check/cross for each domain and returns the results.
func printDNSResults(domains []string, serverIP string) []dnspkg.Result {
	results := dnspkg.CheckDomains(domains, serverIP)
	fmt.Println("DNS check:")
	for _, r := range results {
		if r.OK {
			fmt.Printf("  ✓ %-40s → %s\n", r.Domain, r.MatchedIP)
		} else if r.Err != nil {
			fmt.Printf("  ✗ %-40s — %s\n", r.Domain, r.Err)
		}
	}
	fmt.Println()
	return results
}
