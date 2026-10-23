package cmd

import (
	"io"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// execRoot runs the real command tree the way main() does, so these tests
// exercise the path a user actually hits. Calling a command's Args validator
// directly would pass even when the binary doesn't: cobra checks Runnable()
// first, and for a parent with no Run it returns "print help" long before
// reaching the validator.
func execRoot(t *testing.T, args ...string) error {
	t.Helper()

	root := rootCmd
	root.SetArgs(args)
	root.SetOut(io.Discard)
	root.SetErr(io.Discard)
	defer func() {
		root.SetArgs(nil)
		root.SetOut(nil)
		root.SetErr(nil)
	}()

	return root.Execute()
}

// `serv-kun postgres backupp` must fail, not print help and exit 0 — a typo in
// a script under `set -e` would otherwise pass silently having done nothing.
func TestParentCommandsRejectUnknownSubcommands(t *testing.T) {
	for _, parent := range parentCommands() {
		name := parent.Name()
		err := execRoot(t, name, "bogus-subcommand")
		if err == nil {
			t.Errorf("`serv-kun %s bogus-subcommand` returned no error; "+
				"the binary would exit 0 on a typo", name)
			continue
		}
		if !strings.Contains(err.Error(), "bogus-subcommand") {
			t.Errorf("`serv-kun %s bogus-subcommand` error doesn't name the bad "+
				"argument: %v", name, err)
		}
	}
}

// Bare `serv-kun <parent>` must still print help and succeed.
func TestParentCommandsAcceptNoArguments(t *testing.T) {
	for _, parent := range parentCommands() {
		if err := execRoot(t, parent.Name()); err != nil {
			t.Errorf("`serv-kun %s` returned an error: %v", parent.Name(), err)
		}
	}
}

// Known subcommands must still route correctly — the fix must not break them.
func TestKnownSubcommandsStillRoute(t *testing.T) {
	// --help short-circuits before any RunE, so this checks routing without
	// executing a wizard that would demand root and a live system.
	for _, path := range [][]string{
		{"nginx", "setup"}, {"nginx", "ssl"}, {"ssh", "harden"}, {"ssh", "audit"},
		{"postgres", "setup"}, {"postgres", "backup"}, {"postgres", "restore"},
		{"springboot", "init"}, {"springboot", "deploy"}, {"springboot", "clean"},
	} {
		args := append(append([]string{}, path...), "--help")
		if err := execRoot(t, args...); err != nil {
			t.Errorf("`serv-kun %s --help` failed: %v", strings.Join(path, " "), err)
		}
	}
}

// Both halves of the fix are required: NoArgs alone is unreachable on a
// non-runnable command, and RunE alone would accept stray arguments.
func TestParentCommandsAreGroupingCommands(t *testing.T) {
	for _, parent := range parentCommands() {
		if !parent.HasSubCommands() {
			t.Errorf("%q is listed as a parent but has no subcommands", parent.Name())
		}
		if parent.Args == nil {
			t.Errorf("%q has no Args validator; set cobra.NoArgs", parent.Name())
		}
		if parent.RunE == nil && parent.Run == nil {
			t.Errorf("%q is not runnable, so cobra never reaches its Args validator "+
				"and a mistyped subcommand exits 0", parent.Name())
		}
	}
}

// Guards against a new grouping command being added to the tree without being
// registered in parentCommands, which would silently reintroduce the bug.
func TestEveryGroupingCommandIsListedAsAParent(t *testing.T) {
	listed := make(map[string]bool)
	for _, p := range parentCommands() {
		listed[p.Name()] = true
	}

	for _, c := range rootCmd.Commands() {
		// Cobra adds "help" and "completion" itself.
		if c.Name() == "help" || c.Name() == "completion" {
			continue
		}
		if c.HasSubCommands() && !listed[c.Name()] {
			t.Errorf("%q groups subcommands but isn't in parentCommands(), so nothing enforces NoArgs on it", c.Name())
		}
	}
}

// The subcommands themselves must remain reachable.
func TestSubcommandsAreRegistered(t *testing.T) {
	want := map[string][]string{
		"nginx":      {"setup", "ssl"},
		"ssh":        {"harden", "audit"},
		"postgres":   {"setup", "backup", "restore"},
		"springboot": {"init", "setup", "deploy", "clean", "status", "logs"},
	}

	for _, parent := range parentCommands() {
		names := make(map[string]bool)
		for _, sub := range parent.Commands() {
			names[sub.Name()] = true
		}
		for _, sub := range want[parent.Name()] {
			if !names[sub] {
				t.Errorf("%q is missing subcommand %q", parent.Name(), sub)
			}
		}
	}
}

// The root help text is the tool's front page; keep it in step with the
// commands that actually exist.
func TestRootHelpListsEveryTopLevelCommand(t *testing.T) {
	for _, c := range rootCmd.Commands() {
		if c.Name() == "help" || c.Name() == "completion" {
			continue
		}
		if !strings.Contains(rootCmd.Long, c.Name()) {
			t.Errorf("root help doesn't mention %q", c.Name())
		}
	}
}

func TestRootCommandHasNoRunOfItsOwn(t *testing.T) {
	if rootCmd.Run != nil || rootCmd.RunE != nil {
		t.Error("rootCmd should print help rather than run something")
	}
	var _ *cobra.Command = rootCmd
}
