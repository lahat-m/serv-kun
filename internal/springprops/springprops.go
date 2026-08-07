// Package springprops reads a Spring Boot application's own configuration to
// find the environment variables it cannot start without.
//
// The wizard writes a fixed set of variable names. Any other project names
// things differently — its own S3 keys, a signing secret, a base URL — and the
// mismatch is invisible until the container crash-loops on a value Spring left
// as the literal text "${S3_ENDPOINT}". Reading what the application actually
// declares turns that into questions asked before anything is written.
package springprops

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Requirement is one property the application cannot start without: the
// property itself, the variable its author named, and — when the whole value
// is that one placeholder — the variable name Spring will equally accept.
type Requirement struct {
	Property string // app.storage.endpoint
	Var      string // S3_ENDPOINT
	Relaxed  string // APP_STORAGE_ENDPOINT, or "" when the placeholder is embedded
	Help     string // the comment above the line, for the prompt
}

// propertyLine matches "key=value" for the subset of keys Spring uses. Keys
// are left-hand side only: a value may contain anything, including further '='.
var propertyLine = regexp.MustCompile(`^([A-Za-z0-9][A-Za-z0-9._-]*)\s*=\s*(.*)$`)

// placeholder finds ${...} spans. The body is captured whole so a default
// (${VAR:fallback}) can be told apart from a bare requirement.
var placeholder = regexp.MustCompile(`\$\{([^}]*)\}`)

// envVarName is what Spring's relaxed binding and a .env file agree on.
var envVarName = regexp.MustCompile(`^[A-Z_][A-Z0-9_]*$`)

// Locate returns the properties files worth scanning, base profile first so a
// profile-specific file can override it. yamlOnly reports a project that
// configures Spring through YAML, which this package does not parse — callers
// should say so rather than report a clean bill of health.
func Locate(projectDir, profile string) (files []string, yamlOnly bool) {
	resources := filepath.Join(projectDir, "src", "main", "resources")

	for _, name := range []string{"application.properties", "application-" + profile + ".properties"} {
		path := filepath.Join(resources, name)
		if _, err := os.Stat(path); err == nil {
			files = append(files, path)
		}
	}
	if len(files) > 0 {
		return files, false
	}

	for _, name := range []string{"application.yml", "application.yaml",
		"application-" + profile + ".yml", "application-" + profile + ".yaml"} {
		if _, err := os.Stat(filepath.Join(resources, name)); err == nil {
			return nil, true
		}
	}
	return nil, false
}

// Scan reads the given files in order and returns every property that resolves
// from an environment variable with no fallback. A later file wins: a profile
// that supplies a default for a property the base file left bare makes it
// optional, and it drops out of the result.
func Scan(paths ...string) ([]Requirement, error) {
	// Keyed per property *and* variable: one value can interpolate more than
	// one of them, and dropping either would let a required variable through.
	var order []string
	found := make(map[string]Requirement)

	for _, path := range paths {
		reqs, optional, err := scanFile(path)
		if err != nil {
			return nil, err
		}
		for _, req := range reqs {
			key := req.Property + "\x00" + req.Var
			if _, seen := found[key]; !seen {
				order = append(order, key)
			}
			found[key] = req
		}
		// A fallback supplied later for a property an earlier file left bare
		// makes every variable behind that property optional.
		for _, property := range optional {
			for key := range found {
				if strings.HasPrefix(key, property+"\x00") {
					delete(found, key)
				}
			}
		}
	}

	var out []Requirement
	for _, key := range order {
		if req, ok := found[key]; ok {
			out = append(out, req)
		}
	}
	return out, nil
}

// scanFile returns the requirements in one file, plus the properties it
// resolves with a fallback — the latter matter only to undo a requirement an
// earlier file recorded.
func scanFile(path string) (reqs []Requirement, optional []string, err error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, fmt.Errorf("open %s: %w", path, err)
	}
	defer f.Close()

	var comment []string
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())

		switch {
		case line == "":
			comment = nil
			continue
		case strings.HasPrefix(line, "#"), strings.HasPrefix(line, "!"):
			comment = append(comment, strings.TrimSpace(strings.TrimLeft(line, "#! ")))
			continue
		}

		m := propertyLine.FindStringSubmatch(line)
		if m == nil {
			comment = nil
			continue
		}
		property, value := m[1], strings.TrimSpace(m[2])

		required, hasFallback := parseValue(property, value)
		for i := range required {
			required[i].Help = joinComment(comment)
		}
		reqs = append(reqs, required...)
		if hasFallback {
			optional = append(optional, property)
		}
		comment = nil
	}
	if err := sc.Err(); err != nil {
		return nil, nil, fmt.Errorf("read %s: %w", path, err)
	}
	return reqs, optional, nil
}

// parseValue pulls the requirements out of one property's value. A value that
// is exactly one placeholder can also be satisfied by the property's own
// relaxed name, because an environment variable outranks the properties file
// and replaces the value outright. An embedded placeholder gets no such
// alternative: setting the property would replace the surrounding text too.
func parseValue(property, value string) (reqs []Requirement, hasFallback bool) {
	for _, m := range placeholder.FindAllStringSubmatch(value, -1) {
		body := m[1]
		if strings.Contains(body, ":") {
			hasFallback = true
			continue
		}
		if !envVarName.MatchString(body) {
			// A nested or lowercase placeholder refers to another property,
			// not to the environment.
			continue
		}
		req := Requirement{Property: property, Var: body}
		if value == m[0] {
			req.Relaxed = RelaxedName(property)
		}
		reqs = append(reqs, req)
	}
	return reqs, hasFallback
}

// RelaxedName is the environment variable Spring maps onto a property name:
// upper case, with dots and dashes as underscores.
func RelaxedName(property string) string {
	name := strings.ToUpper(property)
	name = strings.ReplaceAll(name, ".", "_")
	return strings.ReplaceAll(name, "-", "_")
}

// Unsatisfied returns the requirements that defined does not cover, under
// either name Spring would accept.
func Unsatisfied(reqs []Requirement, defined map[string]bool) []Requirement {
	var out []Requirement
	for _, req := range reqs {
		if defined[req.Var] {
			continue
		}
		if req.Relaxed != "" && defined[req.Relaxed] {
			continue
		}
		out = append(out, req)
	}
	return out
}

// joinComment flattens the comment block above a property into one line of
// prompt help, short enough to read in a terminal.
func joinComment(lines []string) string {
	const maxHelp = 240

	var kept []string
	for _, line := range lines {
		if line != "" {
			kept = append(kept, line)
		}
	}
	help := strings.Join(kept, " ")
	if len(help) > maxHelp {
		help = help[:maxHelp] + "…"
	}
	return help
}
