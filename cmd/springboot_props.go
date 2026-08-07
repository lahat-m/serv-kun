package cmd

import (
	"fmt"
	"sort"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/lahat-m/serv-kun/internal/compose"
	"github.com/lahat-m/serv-kun/internal/envfile"
	"github.com/lahat-m/serv-kun/internal/springprops"
)

// setupProfile is the Spring profile the generated compose file activates for
// the app service, and so the one whose properties decide what .env.prod needs.
const setupProfile = "prod"

// askProjectRequirements asks for the variables this particular project needs
// and the wizard's own template doesn't write.
//
// The template covers one application's variable names. Another project names
// the same things differently — app.storage.* rather than app.s3.*, a signing
// secret of its own — and nothing catches the difference: Spring leaves an
// unresolvable placeholder as literal text, so the value reaching the app is
// "${S3_ENDPOINT}" and the failure surfaces as a crash-looping container
// forty stack frames from the cause.
//
// defined is what the generated file already covers. Returns the extra
// answers in the order the properties declared them.
func askProjectRequirements(proj compose.Project, defined map[string]bool) ([]envfile.KV, error) {
	files, yamlOnly := springprops.Locate(proj.Dir, setupProfile)
	if yamlOnly {
		fmt.Println("Note: this project configures Spring through YAML, which the wizard can't")
		fmt.Println("read yet — check its required variables against the generated .env.prod.")
		fmt.Println()
		return nil, nil
	}
	if len(files) == 0 {
		// A project deployed from a published image has no sources here.
		// Nothing to check against, and nothing worth saying about it.
		return nil, nil
	}

	reqs, err := springprops.Scan(files...)
	if err != nil {
		// Reading the project's own configuration is a convenience, not a
		// precondition — a parse failure must not stop a deployment.
		fmt.Printf("Note: could not read this project's Spring configuration (%v)\n\n", err)
		return nil, nil
	}

	missing := springprops.Unsatisfied(reqs, defined)
	if len(missing) == 0 {
		fmt.Printf("✓ %s covers every variable this project's configuration requires\n\n", proj.EnvFile)
		return nil, nil
	}

	fmt.Printf("This project requires %d variable(s) the wizard doesn't write:\n\n", len(missing))
	for _, req := range missing {
		fmt.Printf("  %-28s → %s\n", req.Var, req.Property)
	}
	fmt.Println("\nLeave one blank only if you know the app can start without it.")
	fmt.Println()

	var extra []envfile.KV
	var blank []springprops.Requirement
	for _, req := range missing {
		var value string
		if err := survey.AskOne(&survey.Input{
			Message: req.Var + ":",
			Help:    requirementHelp(req),
		}, &value); err != nil {
			return nil, err
		}
		if value == "" {
			blank = append(blank, req)
			continue
		}
		extra = append(extra, envfile.KV{Key: req.Var, Value: value})
	}

	if len(blank) > 0 {
		fmt.Println()
		for _, req := range blank {
			fmt.Printf("⚠ %s left blank — %s will resolve to the literal \"${%s}\"\n",
				req.Var, req.Property, req.Var)
		}
	}
	fmt.Println()
	return extra, nil
}

// requirementHelp explains where the requirement came from, falling back to
// the property name when the project's configuration carries no comment.
func requirementHelp(req springprops.Requirement) string {
	if req.Help != "" {
		return req.Help
	}
	return "required by " + req.Property + ", which has no fallback value"
}

// carriedOver returns the keys of an existing env file that the newly
// generated one doesn't write, so overwriting it doesn't quietly discard
// whatever was added by hand.
func carriedOver(envPath string, defined map[string]bool) []envfile.KV {
	existing, err := envfile.Read(envPath)
	if err != nil {
		return nil
	}

	var keys []string
	for key := range existing {
		if !defined[key] {
			keys = append(keys, key)
		}
	}
	sort.Strings(keys)

	carried := make([]envfile.KV, 0, len(keys))
	for _, key := range keys {
		carried = append(carried, envfile.KV{Key: key, Value: existing[key]})
	}
	return carried
}
