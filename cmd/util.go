package cmd

import (
	"errors"
	"fmt"
	"net"
	"regexp"
	"strconv"
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
// hostnames. Use it where certbot is involved and every entry must be a real
// public domain.
func validateDomainList(ans interface{}) error {
	return validateList(ans, checkDomain, "at least one domain is required")
}

// serverNameRe is checkDomain's pattern without the mandatory TLD, so dotless
// internal names like "myserver" are accepted.
var serverNameRe = regexp.MustCompile(
	`^[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`)

// checkServerName validates an nginx server_name, which may legitimately be an
// IP address or a dotless internal hostname — a reverse proxy doesn't need a
// publicly resolvable name, only certbot does. Still rejects whitespace and
// URL syntax, since those are typos rather than valid configurations.
func checkServerName(s string) error {
	switch {
	case s == "":
		return errors.New("a server name is required")
	case strings.ContainsAny(s, " \t"):
		return errors.New("server name must not contain spaces")
	case strings.Contains(s, "://"):
		return errors.New("enter the host only, without http:// or https://")
	case strings.Contains(s, "/"):
		return errors.New("enter the host only, without a path")
	case len(s) > 253:
		return errors.New("server name is longer than 253 characters")
	}

	if net.ParseIP(s) != nil {
		return nil
	}
	if !serverNameRe.MatchString(s) {
		return fmt.Errorf("%q isn't a valid hostname or IP address", s)
	}
	return nil
}

// validateServerNameList is a survey validator for nginx server_name values.
func validateServerNameList(ans interface{}) error {
	return validateList(ans, checkServerName, "at least one server name is required")
}

// checkPort validates a TCP/UDP port number.
func checkPort(s string) error {
	n, err := strconv.Atoi(s)
	if err != nil {
		return fmt.Errorf("%q isn't a number", s)
	}
	if n < 1 || n > 65535 {
		return fmt.Errorf("port %d is out of range (1-65535)", n)
	}
	return nil
}

// validatePort is a survey validator for a single required port.
func validatePort(ans interface{}) error {
	s, _ := ans.(string)
	return checkPort(strings.TrimSpace(s))
}

// validateOptionalPortList accepts a blank answer, or a space-separated list
// of port numbers. Catching a typo here beats letting ufw reject it after the
// rest of the rules have already been applied.
func validateOptionalPortList(ans interface{}) error {
	s, _ := ans.(string)
	for _, f := range strings.Fields(s) {
		if err := checkPort(f); err != nil {
			return err
		}
	}
	return nil
}

func validateList(ans interface{}, check func(string) error, emptyMsg string) error {
	s, _ := ans.(string)
	fields := strings.Fields(s)
	if len(fields) == 0 {
		return errors.New(emptyMsg)
	}
	for _, f := range fields {
		if err := check(f); err != nil {
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
