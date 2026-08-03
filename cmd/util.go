package cmd

import (
	"fmt"

	survey "github.com/AlecAivazis/survey/v2"
)

func printBanner(title string) {
	fmt.Printf("\n=== %s ===\n\n", title)
}

func printRule() {
	fmt.Println("──────────────────────────────────────────────────────────────")
}

// askIf asks a yes/no question (default yes) and, if confirmed, runs then —
// used to gate optional blocks of prompts (JWT, admin account, email,
// object storage) that not every project needs.
func askIf(message string, then func() error) error {
	var yes bool
	if err := survey.AskOne(&survey.Confirm{Message: message, Default: true}, &yes); err != nil {
		return err
	}
	if !yes {
		return nil
	}
	return then()
}
