package cmd

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	survey "github.com/AlecAivazis/survey/v2"
)

// hostnameRe matches a dotted hostname: labels of alphanumerics and inner
// hyphens, with an alphabetic TLD of at least two characters.
var hostnameRe = regexp.MustCompile(
	`^[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*\.[a-zA-Z]{2,}$`)

// checkDomain validates a single hostname. Beyond catching typos before
// certbot rejects them, this is what keeps whitespace and shell characters
// out of a value that ends up in file paths and container command arguments
// downstream — "example .com" would otherwise word-split into two operands.
func checkDomain(s string) error {
	switch {
	case s == "":
		return errors.New("domain is required")
	case strings.ContainsAny(s, " \t"):
		return errors.New("domain must not contain spaces")
	case strings.Contains(s, "://"):
		return errors.New("enter the hostname only, without http:// or https://")
	case strings.Contains(s, "/"):
		return errors.New("enter the hostname only, without a path")
	case strings.Contains(s, ":"):
		return errors.New("enter the hostname only, without a port")
	case len(s) > 253:
		return errors.New("domain is longer than 253 characters")
	case !hostnameRe.MatchString(s):
		return fmt.Errorf("%q doesn't look like a domain (expected something like example.com)", s)
	}
	return nil
}

// validateDomain is a survey validator for a single hostname.
func validateDomain(ans interface{}) error {
	s, _ := ans.(string)
	return checkDomain(strings.TrimSpace(s))
}

// validateDomainList is a survey validator for a space-separated list of
// hostnames, as the nginx wizards accept.
func validateDomainList(ans interface{}) error {
	s, _ := ans.(string)
	fields := strings.Fields(s)
	if len(fields) == 0 {
		return errors.New("at least one domain is required")
	}
	for _, f := range fields {
		if err := checkDomain(f); err != nil {
			return err
		}
	}
	return nil
}

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
