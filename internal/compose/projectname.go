package compose

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// keepInProjectName mirrors compose-go's NormalizeProjectName: characters
// outside [a-z0-9_-] are *dropped* rather than replaced, and a leading "_" or
// "-" is trimmed. Matching it exactly is the point — the job here is to
// predict the name docker itself will choose, and a near-miss would report a
// collision that doesn't exist or, worse, miss one that does.
var keepInProjectName = regexp.MustCompile(`[a-z0-9_-]`)

// NormalizeProjectName reduces a string to the form docker compose accepts as
// a project name.
func NormalizeProjectName(s string) string {
	s = strings.ToLower(s)
	s = strings.Join(keepInProjectName.FindAllString(s, -1), "")
	return strings.TrimLeft(s, "_-")
}

// topLevelNameRe matches a compose file's top-level "name:" key. Column 0 is
// what makes it top-level: every nested key is indented, and the only other
// place "name:" can appear unindented is inside a block scalar, which is
// indented too.
var topLevelNameRe = regexp.MustCompile(`(?m)^name:[ \t]*(.*)$`)

// ProjectNameOf reports the compose project name a compose file resolves to
// on its own: its top-level "name:" key if it has one, otherwise — the case
// that causes trouble — the normalized name of the directory holding it.
//
// COMPOSE_PROJECT_NAME outranks both and is not consulted: an operator who
// exported it has already chosen a namespace deliberately, and guessing at
// which shell a future command will run in would make this less reliable,
// not more.
func ProjectNameOf(dir, file string) string {
	if data, err := os.ReadFile(filepath.Join(dir, file)); err == nil {
		if m := topLevelNameRe.FindSubmatch(data); m != nil {
			if name := yamlScalar(string(m[1])); name != "" {
				return NormalizeProjectName(name)
			}
		}
	}
	return DefaultProjectName(dir)
}

// DefaultProjectName is the project name compose derives from a directory
// when no compose file in it declares one — which is to say, the name two
// stacks in the same directory would silently end up sharing.
func DefaultProjectName(dir string) string {
	return NormalizeProjectName(dirBase(dir))
}

// Collision is another compose file in the same directory that resolves to
// the same project name — and so shares this one's container namespace.
type Collision struct {
	File     string   // file name, relative to the project directory
	Services []string // service names both files declare, sorted
}

// Collisions reports the other compose files in dir that resolve to the same
// project name as file.
//
// Files named "*.override.y[a]ml" are skipped: an override is merged into
// another file's project by design, so sharing a name is what it is for.
// Everything else that looks like a compose file is fair game, because the
// hazard doesn't depend on how the file is named — only on docker resolving
// two files to one project.
func Collisions(dir, file string) []Collision {
	name := ProjectNameOf(dir, file)
	if name == "" {
		return nil
	}
	ours := serviceNames(filepath.Join(dir, file))

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}

	var found []Collision
	for _, entry := range entries {
		other := entry.Name()
		if entry.IsDir() || other == file || !isComposeFile(other) || isOverrideFile(other) {
			continue
		}
		if ProjectNameOf(dir, other) != name {
			continue
		}
		found = append(found, Collision{
			File:     other,
			Services: sharedServices(ours, serviceNames(filepath.Join(dir, other))),
		})
	}
	sort.Slice(found, func(i, j int) bool { return found[i].File < found[j].File })
	return found
}

func isComposeFile(name string) bool {
	ext := strings.ToLower(filepath.Ext(name))
	if ext != ".yml" && ext != ".yaml" {
		return false
	}
	return strings.Contains(strings.ToLower(name), "compose")
}

func isOverrideFile(name string) bool {
	base := strings.ToLower(strings.TrimSuffix(name, filepath.Ext(name)))
	return strings.HasSuffix(base, ".override")
}

// serviceNames lists the keys under a compose file's top-level "services:"
// block. Only the first indentation level counts: everything deeper belongs
// to a service, not to the block.
func serviceNames(path string) []string {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}

	var (
		names      []string
		indent     string
		inServices bool
	)
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimRight(line, "\r")
		if strings.TrimSpace(line) == "" || strings.HasPrefix(strings.TrimSpace(line), "#") {
			continue
		}
		if !inServices {
			inServices = strings.HasPrefix(line, "services:")
			continue
		}
		// Back at column 0: the services block has ended.
		if line == strings.TrimLeft(line, " \t") {
			break
		}
		lead := line[:len(line)-len(strings.TrimLeft(line, " \t"))]
		if indent == "" {
			indent = lead
		}
		if lead != indent {
			continue
		}
		if key, _, ok := strings.Cut(strings.TrimSpace(line), ":"); ok {
			names = append(names, key)
		}
	}
	return names
}

func sharedServices(a, b []string) []string {
	inB := make(map[string]bool, len(b))
	for _, name := range b {
		inB[name] = true
	}
	var shared []string
	for _, name := range a {
		if inB[name] {
			shared = append(shared, name)
		}
	}
	sort.Strings(shared)
	return shared
}

// dirBase resolves dir to the directory name docker would derive a project
// name from. It has to be absolute first: the usual value is ".", which
// names the working directory, not a directory called "." .
func dirBase(dir string) string {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return filepath.Base(dir)
	}
	return filepath.Base(abs)
}

// yamlScalar strips what a plain YAML scalar may carry around its value: a
// trailing comment, then surrounding quotes.
func yamlScalar(s string) string {
	s = strings.TrimSpace(s)
	// "#" only opens a comment when it follows whitespace or starts the value.
	if i := strings.Index(s, " #"); i >= 0 {
		s = s[:i]
	}
	s = strings.TrimSpace(s)
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'') {
			s = s[1 : len(s)-1]
		}
	}
	return strings.TrimSpace(s)
}
