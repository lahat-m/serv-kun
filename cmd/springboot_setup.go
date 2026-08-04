package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/lahat-m/serv-kun/internal/dns"
	"github.com/lahat-m/serv-kun/internal/docker"
	"github.com/lahat-m/serv-kun/internal/envfile"
	"github.com/lahat-m/serv-kun/internal/healthcheck"
	"github.com/lahat-m/serv-kun/internal/jwtkeys"
	"github.com/lahat-m/serv-kun/internal/system"
	"github.com/spf13/cobra"
)

var springbootSetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Bootstrap docker, .env.prod, JWT keys, TLS, and bring the stack up",
	RunE:  runSpringbootSetup,
}

const (
	setupHealthTimeout  = 120 * time.Second
	setupHealthInterval = 3 * time.Second
)

// certsMountRe pulls the container-side path out of the app service's
// "./certs:/somewhere:ro" bind mount.
var certsMountRe = regexp.MustCompile(`\./certs:([^:\s]+)`)

// defaultCertsMount reads that path out of the compose file so the prompt's
// default matches what `springboot init` actually wrote. The old hardcoded
// /etc/app/certs was correct only for a project literally named "app";
// anywhere else it pointed JWT_PRIVATE_KEY at a path that doesn't exist in
// the container, and the app failed to start on a missing key file.
func defaultCertsMount(composePath string) string {
	const fallback = "/etc/app/certs"

	data, err := os.ReadFile(composePath)
	if err != nil {
		return fallback
	}
	if m := certsMountRe.FindSubmatch(data); m != nil {
		return string(m[1])
	}
	return fallback
}

func runSpringbootSetup(cmd *cobra.Command, _ []string) error {
	requireRoot(cmd)
	printBanner("springboot setup")

	proj := projectFromFlags(cmd)

	composePath := filepath.Join(proj.Dir, proj.File)
	if _, err := os.Stat(composePath); err != nil {
		return fmt.Errorf("%s not found — run this from the project directory, or pass --dir", composePath)
	}

	// ── Docker ──────────────────────────────────────────────────────────────
	if !docker.IsInstalled() {
		var install bool
		if err := survey.AskOne(&survey.Confirm{
			Message: "Docker is not installed. Install it now (via get.docker.com)?",
			Default: true,
		}, &install); err != nil {
			return err
		}
		if !install {
			return fmt.Errorf("docker is required to continue")
		}
		fmt.Println("Installing Docker...")
		if err := docker.Install(); err != nil {
			return fmt.Errorf("install docker: %w", err)
		}
		if err := docker.Enable(); err != nil {
			return fmt.Errorf("enable docker: %w", err)
		}
		fmt.Println("✓ Docker installed and started")
	} else {
		fmt.Println("✓ Docker is already installed")
	}
	fmt.Println()

	// ── domain / image (asked first — later defaults derive from these) ────
	var basics struct {
		AppDomain    string
		CertbotEmail string
		AppImage     string
	}
	if err := survey.Ask([]*survey.Question{
		{
			Name:     "appDomain",
			Prompt:   &survey.Input{Message: "Domain this app will be served on:", Help: "e.g. kuniel.com"},
			Validate: validateDomain,
		},
		{
			Name:     "certbotEmail",
			Prompt:   &survey.Input{Message: "Email for Let's Encrypt renewal notices:"},
			Validate: survey.Required,
		},
		{
			Name:     "appImage",
			Prompt:   &survey.Input{Message: "App image:", Help: "e.g. ghcr.io/you/app:latest"},
			Validate: survey.Required,
		},
	}, &basics); err != nil {
		return err
	}
	// survey stores the raw input; the validator only inspected the trimmed
	// form, so normalize before this flows into paths and container args.
	basics.AppDomain = strings.TrimSpace(basics.AppDomain)
	fmt.Println()

	// ── always required: Postgres + CORS ────────────────────────────────────
	var core struct {
		PostgresDB         string
		PostgresUser       string
		PostgresPassword   string
		BackupKeepDays     string
		CorsAllowedOrigins string
	}
	if err := survey.Ask([]*survey.Question{
		{Name: "postgresDB", Prompt: &survey.Input{Message: "Postgres database name:", Default: "app"}, Validate: survey.Required},
		{Name: "postgresUser", Prompt: &survey.Input{Message: "Postgres role name:", Default: "app"}, Validate: survey.Required},
		{Name: "postgresPassword", Prompt: &survey.Password{Message: "Postgres password:"}, Validate: survey.Required},
		{Name: "backupKeepDays", Prompt: &survey.Input{Message: "Days to keep nightly Postgres backups:", Default: "7"}},
		{Name: "corsAllowedOrigins", Prompt: &survey.Input{Message: "Allowed CORS origin(s):", Default: "https://" + basics.AppDomain}},
	}, &core); err != nil {
		return err
	}
	fmt.Println()

	answers := envfile.Answers{
		AppDomain:          basics.AppDomain,
		CertbotEmail:       basics.CertbotEmail,
		AppImage:           basics.AppImage,
		PostgresDB:         core.PostgresDB,
		PostgresUser:       core.PostgresUser,
		PostgresPassword:   core.PostgresPassword,
		BackupKeepDays:     core.BackupKeepDays,
		CorsAllowedOrigins: core.CorsAllowedOrigins,
	}

	// ── optional blocks — each is a project choice, not every Spring Boot +
	// Postgres app needs all of these ─────────────────────────────────────
	if err := askIf("Does this app use RS256 JWT keys (mounted from certs/)?", func() error {
		answers.HasJWT = true
		return survey.AskOne(&survey.Input{
			Message: "Path certs/ is mounted at inside the app container:",
			Default: defaultCertsMount(composePath),
			Help:    "must match the app service's ./certs bind mount in the compose file",
		}, &answers.CertsMountPath, survey.WithValidator(survey.Required))
	}); err != nil {
		return err
	}

	if err := askIf("Seed an initial admin account on first boot?", func() error {
		var admin struct {
			AdminEmail    string
			AdminPassword string
			AdminFullName string
		}
		if err := survey.Ask([]*survey.Question{
			{Name: "adminEmail", Prompt: &survey.Input{Message: "Seeded admin email:"}, Validate: survey.Required},
			{Name: "adminPassword", Prompt: &survey.Password{Message: "Seeded admin password:"}, Validate: survey.Required},
			{Name: "adminFullName", Prompt: &survey.Input{Message: "Seeded admin full name:", Default: "Admin"}},
		}, &admin); err != nil {
			return err
		}
		answers.HasAdmin = true
		answers.AdminEmail = admin.AdminEmail
		answers.AdminPassword = admin.AdminPassword
		answers.AdminFullName = admin.AdminFullName
		return nil
	}); err != nil {
		return err
	}

	if err := askIf("Does this app send email via Resend?", func() error {
		var email struct {
			ResendAPIKey string
			ResendFrom   string
		}
		if err := survey.Ask([]*survey.Question{
			{Name: "resendAPIKey", Prompt: &survey.Password{Message: "Resend API key:"}, Validate: survey.Required},
			{Name: "resendFrom", Prompt: &survey.Input{Message: "Resend \"from\" address:", Default: fmt.Sprintf("App <noreply@%s>", basics.AppDomain)}},
		}, &email); err != nil {
			return err
		}
		answers.HasEmail = true
		answers.ResendAPIKey = email.ResendAPIKey
		answers.ResendFrom = email.ResendFrom
		return nil
	}); err != nil {
		return err
	}

	if err := askIf("Does this app use S3-compatible object storage?", func() error {
		var s3 struct {
			S3Endpoint  string
			S3AccessKey string
			S3SecretKey string
			S3Bucket    string
			S3Region    string
		}
		if err := survey.Ask([]*survey.Question{
			{Name: "s3Endpoint", Prompt: &survey.Input{Message: "S3-compatible storage endpoint:", Help: "e.g. https://fsn1.your-objectstorage.com"}, Validate: survey.Required},
			{Name: "s3AccessKey", Prompt: &survey.Password{Message: "S3 access key:"}, Validate: survey.Required},
			{Name: "s3SecretKey", Prompt: &survey.Password{Message: "S3 secret key:"}, Validate: survey.Required},
			{Name: "s3Bucket", Prompt: &survey.Input{Message: "S3 bucket name:"}, Validate: survey.Required},
			{Name: "s3Region", Prompt: &survey.Input{Message: "S3 region:"}, Validate: survey.Required},
		}, &s3); err != nil {
			return err
		}
		answers.HasObjectStorage = true
		answers.S3Endpoint = s3.S3Endpoint
		answers.S3AccessKey = s3.S3AccessKey
		answers.S3SecretKey = s3.S3SecretKey
		answers.S3Bucket = s3.S3Bucket
		answers.S3Region = s3.S3Region
		return nil
	}); err != nil {
		return err
	}
	fmt.Println()

	// ── confirm before touching anything ────────────────────────────────────
	// Secrets are deliberately left out of this summary — it's a sanity check
	// on the settings that are easy to typo (domain, image, db name), not a
	// review of the whole file.
	fmt.Println()
	printRule()
	fmt.Printf("  Domain:        %s\n", answers.AppDomain)
	fmt.Printf("  Image:         %s\n", answers.AppImage)
	fmt.Printf("  Postgres DB:   %s (user %s)\n", answers.PostgresDB, answers.PostgresUser)
	fmt.Printf("  CORS origin:   %s\n", answers.CorsAllowedOrigins)
	if answers.HasAdmin {
		fmt.Printf("  Admin email:   %s\n", answers.AdminEmail)
	}
	if answers.HasObjectStorage {
		fmt.Printf("  S3 bucket:     %s (%s)\n", answers.S3Bucket, answers.S3Region)
	}
	printRule()
	fmt.Println()

	setupAction := "write " + proj.EnvFile
	if answers.HasJWT {
		setupAction += ", generate JWT keys"
	}
	setupAction += fmt.Sprintf(", request a TLS certificate for %s, and start the stack", answers.AppDomain)

	var proceedSetup bool
	if err := survey.AskOne(&survey.Confirm{
		Message: fmt.Sprintf("This will %s. Proceed?", setupAction),
		Default: true,
	}, &proceedSetup); err != nil {
		return err
	}
	if !proceedSetup {
		fmt.Println("Aborted — no changes made.")
		return nil
	}

	content, err := envfile.Generate(answers)
	if err != nil {
		return err
	}

	envPath := filepath.Join(proj.Dir, proj.EnvFile)
	if _, err := os.Stat(envPath); err == nil {
		backup := envPath + ".bak." + time.Now().Format("20060102-150405")
		if err := os.Rename(envPath, backup); err != nil {
			return fmt.Errorf("back up existing %s: %w", proj.EnvFile, err)
		}
		fmt.Printf("Backed up existing %s → %s\n", proj.EnvFile, backup)
	}

	// 0600: this file holds real secrets (DB password, JWT paths, API keys).
	if err := os.WriteFile(envPath, []byte(content), 0600); err != nil {
		return fmt.Errorf("write %s: %w", proj.EnvFile, err)
	}
	fmt.Printf("✓ %s written\n\n", proj.EnvFile)

	// ── JWT keys ─────────────────────────────────────────────────────────────
	if answers.HasJWT {
		certsDir := filepath.Join(proj.Dir, "certs")
		if jwtkeys.Exists(certsDir) {
			fmt.Println("✓ JWT keypair already present in certs/")
		} else {
			fmt.Println("Generating JWT keypair...")
			if err := jwtkeys.Generate(certsDir); err != nil {
				return fmt.Errorf("generate JWT keys: %w", err)
			}
			fmt.Println("✓ JWT keypair generated → certs/")
		}
		// certs/ is 0700 and private.pem 0600, both root-owned. If the image's
		// Dockerfile sets a non-root USER, the app cannot read its own signing
		// key and dies at startup on a file it can see but not open.
		fmt.Println("  Note: keys are root-only. If your Dockerfile sets a non-root USER,")
		fmt.Println("  grant that UID read access, e.g.  sudo chown -R <uid>:<gid> certs/")
		fmt.Println()
	}

	// ── DNS check ────────────────────────────────────────────────────────────
	fmt.Println("Resolving server public IP...")
	serverIP, err := system.GetPublicIP()
	if err != nil {
		return fmt.Errorf("could not determine server public IP: %w", err)
	}
	fmt.Printf("Server IP: %s\n\n", serverIP)

	results := dns.CheckDomains([]string{basics.AppDomain}, serverIP)
	for _, r := range results {
		switch {
		case r.Proxied:
			fmt.Printf("  ✓ %s → Cloudflare edge (proxied — challenge still reaches this server)\n", r.Domain)
		case r.OK:
			fmt.Printf("  ✓ %s → %s\n", r.Domain, r.MatchedIP)
		default:
			fmt.Printf("  ✗ %s → %v (%s)\n", r.Domain, r.ResolvedIPs, r.Err)
		}
	}
	fmt.Println()

	if !dns.AllOK(results) {
		var proceed bool
		if err := survey.AskOne(&survey.Confirm{
			Message: "DNS doesn't point here yet — certificate issuance will fail. Proceed anyway?",
			Default: false,
		}, &proceed); err != nil {
			return err
		}
		if !proceed {
			fmt.Println("Aborted. Point your DNS at this server and re-run setup.")
			return nil
		}
	}

	// ── TLS bootstrap (dummy cert → real cert → reload), then bring up the stack ──
	if err := bootstrapTLS(proj, basics.AppDomain, basics.CertbotEmail); err != nil {
		return fmt.Errorf("TLS bootstrap: %w", err)
	}

	fmt.Println("\nStarting the full stack...")
	if err := proj.Run("up", "-d"); err != nil {
		return fmt.Errorf("compose up: %w", err)
	}

	healthURL := "https://" + basics.AppDomain + "/actuator/health"
	fmt.Printf("\nWaiting for %s to become healthy...\n", healthURL)
	if err := healthcheck.WaitHealthy(healthURL, setupHealthTimeout, setupHealthInterval); err != nil {
		fmt.Printf("✗ %s\n", err)
		fmt.Println("Check logs with: serv-kun springboot logs")
		return nil
	}

	fmt.Printf("\n✓ Live at https://%s\n", basics.AppDomain)
	fmt.Println("Use  serv-kun springboot deploy  for future updates.")
	return nil
}
