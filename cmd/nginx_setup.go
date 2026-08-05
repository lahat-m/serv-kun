package cmd

import (
	"fmt"
	"strings"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/lahat-m/serv-kun/internal/certbot"
	nginxpkg "github.com/lahat-m/serv-kun/internal/nginx"
	"github.com/lahat-m/serv-kun/internal/system"
	"github.com/spf13/cobra"
)

var nginxSetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Full wizard: install nginx, write reverse-proxy config, optionally add SSL",
	RunE:  runNginxSetup,
}

func runNginxSetup(cmd *cobra.Command, _ []string) error {
	requireRoot(cmd)

	printBanner("nginx setup")

	osInfo, err := system.Detect()
	if err != nil {
		return err
	}
	fmt.Printf("Detected OS: %s\n\n", osInfo.Pretty)

	// ── nginx install ──────────────────────────────────────────────────────
	installed, err := nginxpkg.IsInstalled()
	if err != nil {
		return err
	}

	if !installed {
		var install bool
		if err := survey.AskOne(&survey.Confirm{
			Message: "nginx is not installed. Install it now?",
			Default: true,
		}, &install); err != nil {
			return err
		}
		if !install {
			return fmt.Errorf("nginx is required to continue")
		}
		fmt.Println("Installing nginx...")
		if err := nginxpkg.Install(osInfo); err != nil {
			return fmt.Errorf("install nginx: %w", err)
		}
		if err := nginxpkg.Enable(); err != nil {
			return fmt.Errorf("enable nginx: %w", err)
		}
		fmt.Println("✓ nginx installed and started")
	} else {
		fmt.Println("✓ nginx is already installed")
	}
	fmt.Println()

	// ── gather proxy config ────────────────────────────────────────────────
	var answers struct {
		Upstream    string
		Domains     string
		MaxBodySize string
	}

	qs := []*survey.Question{
		{
			Name: "upstream",
			Prompt: &survey.Input{
				Message: "Upstream URL:",
				Help:    "e.g. http://localhost:3000  or  https://localhost:8443",
			},
			Validate: survey.Required,
		},
		{
			Name: "domains",
			Prompt: &survey.Input{
				Message: "Domain(s):",
				Help:    "Space-separated — e.g. example.com www.example.com, or an IP for a proxy without TLS",
			},
			Validate: validateServerNameList,
		},
		{
			Name: "maxBodySize",
			Prompt: &survey.Input{
				Message: "Max upload size (client_max_body_size):",
				Default: "100M",
				Help:    "e.g. 100M, 1G",
			},
		},
	}

	if err := survey.Ask(qs, &answers); err != nil {
		return err
	}

	domains := strings.Fields(answers.Domains)

	cfg := nginxpkg.Config{
		Domains:     domains,
		Upstream:    answers.Upstream,
		MaxBodySize: answers.MaxBodySize,
		OSInfo:      osInfo,
	}

	content, err := nginxpkg.GenerateConfig(cfg)
	if err != nil {
		return err
	}

	fmt.Println()
	printRule()
	fmt.Print(content)
	printRule()
	fmt.Println()

	var apply bool
	if err := survey.AskOne(&survey.Confirm{
		Message: "Apply this configuration?",
		Default: true,
	}, &apply); err != nil {
		return err
	}
	if !apply {
		fmt.Println("Aborted — no changes made.")
		return nil
	}

	// ── write, test, reload ────────────────────────────────────────────────
	written, err := nginxpkg.WriteConfig(cfg, content)
	if err != nil {
		// Partially applied — undo whatever did land.
		if restoreErr := nginxpkg.RestoreConfig(written); restoreErr != nil {
			return fmt.Errorf("%w\n\nROLLBACK ALSO FAILED: %v", err, restoreErr)
		}
		return err
	}

	fmt.Println("\nTesting nginx configuration...")
	if err := nginxpkg.Test(); err != nil {
		// An invalid config left on disk would stop nginx starting at the
		// next reload or reboot, long after this command finished.
		fmt.Println("✗ Validation failed — rolling back.")
		if restoreErr := nginxpkg.RestoreConfig(written); restoreErr != nil {
			return fmt.Errorf("%w\n\nROLLBACK ALSO FAILED: %v\nFix %s by hand before reloading nginx",
				err, restoreErr, written.Path)
		}
		fmt.Println("✓ Rolled back — nginx was never reloaded, so nothing changed.")
		return err
	}

	if err := nginxpkg.Reload(); err != nil {
		return err
	}
	fmt.Println("✓ nginx configured and reloaded")
	fmt.Println()

	// ── optional SSL ───────────────────────────────────────────────────────
	// server_name accepts IPs and dotless hosts, but Let's Encrypt only issues
	// for public domain names — so don't offer SSL where certbot is certain to
	// fail. Better to say why than to let it error out mid-issuance.
	for _, d := range domains {
		if err := checkDomain(d); err != nil {
			fmt.Printf("\nDone. Skipping SSL: %s\n", err)
			fmt.Println("Let's Encrypt only issues certificates for public domain names.")
			return nil
		}
	}

	var doSSL bool
	if err := survey.AskOne(&survey.Confirm{
		Message: "Set up SSL with Let's Encrypt now?",
		Default: true,
	}, &doSSL); err != nil {
		return err
	}

	if !doSSL {
		fmt.Println("\nDone. Run  sudo serv-kun nginx ssl  later to add SSL.")
		return nil
	}

	return runSSLSetup(domains, osInfo, true)
}

// runSSLSetup is shared by the nginx setup and nginx ssl commands. Pass
// checkDNS=false when the caller has already resolved and reported DNS, so
// the operator isn't shown the same check twice.
func runSSLSetup(domains []string, osInfo *system.Info, checkDNS bool) error {
	printBanner("SSL setup")

	if checkDNS {
		fmt.Println("Resolving server public IP...")
		serverIP, err := system.GetPublicIP()
		if err != nil {
			return fmt.Errorf("could not determine server public IP: %w", err)
		}
		fmt.Printf("Server IP: %s\n\n", serverIP)

		printDNSResults(domains, serverIP)
	}

	// ── certbot install ────────────────────────────────────────────────────
	if !certbot.IsInstalled() {
		fmt.Println("certbot is not installed. Installing...")
		if err := certbot.Install(osInfo); err != nil {
			return fmt.Errorf("install certbot: %w", err)
		}
		fmt.Println("✓ certbot installed")
	} else {
		fmt.Println("✓ certbot is installed")
	}
	fmt.Println()

	var email string
	if err := survey.AskOne(&survey.Input{
		Message: "Email for Let's Encrypt renewal notices (leave blank to skip):",
	}, &email); err != nil {
		return err
	}

	fmt.Println("\nRunning certbot...")
	if err := certbot.Run(domains, email); err != nil {
		return fmt.Errorf("certbot: %w", err)
	}

	fmt.Printf("\n✓ SSL configured — your site is live at https://%s\n", domains[0])
	return nil
}
