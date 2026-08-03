package cmd

import (
	"errors"
	"fmt"
	"strings"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/lahat-m/serv-kun/internal/sudouser"
	"github.com/lahat-m/serv-kun/internal/system"
	"github.com/spf13/cobra"
)

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Create a Linux user with sudo privileges",
	RunE:  runUser,
}

func runUser(cmd *cobra.Command, _ []string) error {
	requireRoot(cmd)

	printBanner("serv-kun user")

	osInfo, err := system.Detect()
	if err != nil {
		return err
	}

	var username string
	if err := survey.AskOne(&survey.Input{
		Message: "Username for the new sudo user:",
	}, &username, survey.WithValidator(survey.Required)); err != nil {
		return err
	}

	if sudouser.Exists(username) {
		return fmt.Errorf("user %q already exists", username)
	}

	fmt.Printf("\nCreating user %q...\n", username)
	if err := sudouser.Create(username); err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	if err := sudouser.GrantSudo(username, osInfo); err != nil {
		return fmt.Errorf("grant sudo: %w", err)
	}
	fmt.Printf("✓ %q created and added to the %q group\n\n", username, sudouser.SudoGroup(osInfo))

	var authMethods []string
	if err := survey.AskOne(&survey.MultiSelect{
		Message: "How should this user authenticate over SSH?",
		Options: []string{"SSH public key", "Password"},
		Default: []string{"SSH public key"},
	}, &authMethods, survey.WithValidator(survey.Required)); err != nil {
		return err
	}

	for _, method := range authMethods {
		switch method {
		case "SSH public key":
			var pubKey string
			if err := survey.AskOne(&survey.Multiline{
				Message: "Paste the SSH public key:",
			}, &pubKey, survey.WithValidator(validateSSHKey)); err != nil {
				return err
			}
			if err := sudouser.AddSSHKey(username, pubKey); err != nil {
				return fmt.Errorf("add SSH key: %w", err)
			}
			fmt.Println("✓ SSH key installed")

		case "Password":
			var password string
			if err := survey.AskOne(&survey.Password{
				Message: "Set a password for this user:",
			}, &password, survey.WithValidator(survey.Required)); err != nil {
				return err
			}
			if err := sudouser.SetPassword(username, password); err != nil {
				return fmt.Errorf("set password: %w", err)
			}
			fmt.Println("✓ Password set")
		}
	}

	fmt.Printf("\n✓ Sudo user %q is ready\n", username)
	return nil
}

func validateSSHKey(ans interface{}) error {
	s, _ := ans.(string)
	s = strings.TrimSpace(s)
	if s == "" {
		return errors.New("public key is required")
	}
	prefixes := []string{"ssh-rsa", "ssh-ed25519", "ssh-dss", "ecdsa-sha2-"}
	for _, p := range prefixes {
		if strings.HasPrefix(s, p) {
			return nil
		}
	}
	return errors.New("doesn't look like an SSH public key (should start with ssh-rsa, ssh-ed25519, etc.)")
}
