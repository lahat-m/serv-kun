package sudouser

import (
	"path/filepath"
	"strings"
	"testing"
)

// sudo's #includedir skips any file in /etc/sudoers.d whose name contains a
// "." or ends in "~". A rule written to such a name is valid, root-owned, and
// completely inert — the exact failure this constant exists to avoid, and one
// that reproduces as "I wrote the rule and sudo still asks for a password".
func TestSudoersFileNameIsOneSudoWillRead(t *testing.T) {
	base := filepath.Base(SudoersFile)
	if strings.Contains(base, ".") {
		t.Errorf("SudoersFile %q has a dot in its base name — sudo will ignore it", SudoersFile)
	}
	if strings.HasSuffix(base, "~") {
		t.Errorf("SudoersFile %q ends in ~ — sudo will ignore it", SudoersFile)
	}
	if !filepath.IsAbs(SudoersFile) {
		t.Errorf("SudoersFile %q must be absolute", SudoersFile)
	}
}

func TestCommandStringPinsArgumentsOnlyWhenGiven(t *testing.T) {
	anyArgs := Command{Path: "/usr/local/bin/serv-kun"}
	if got := anyArgs.String(); got != "/usr/local/bin/serv-kun" {
		t.Errorf("String() = %q, want the bare path", got)
	}

	pinned := Command{Path: "/usr/bin/docker", Args: []string{"image", "prune", "-f"}}
	if got, want := pinned.String(), "/usr/bin/docker image prune -f"; got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}

func TestNoPasswdRuleGrantsTheListedCommandsToTheUser(t *testing.T) {
	rule := NoPasswdRule("deployer", []Command{
		{Path: "/usr/local/bin/serv-kun"},
		{Path: "/usr/bin/docker", Args: []string{"image", "prune", "-f"}},
	})

	want := "deployer ALL=(ALL) NOPASSWD: /usr/local/bin/serv-kun, /usr/bin/docker image prune -f\n"
	if !strings.Contains(rule, want) {
		t.Errorf("rule does not contain\n%s\n--- got ---\n%s", want, rule)
	}

	// The grant has to be the last word. A comment after it reads, to anyone
	// skimming, like part of the specification.
	lines := strings.Split(strings.TrimRight(rule, "\n"), "\n")
	if last := lines[len(lines)-1]; !strings.HasPrefix(last, "deployer ") {
		t.Errorf("last line is %q, want the grant itself", last)
	}
	for _, line := range lines[:len(lines)-1] {
		if !strings.HasPrefix(line, "#") {
			t.Errorf("unexpected non-comment line before the grant: %q", line)
		}
	}
}

// Refusing early matters because the alternative is a syntactically valid file
// that grants nothing, or grants it to nobody — installed, reported as
// success, and useless.
func TestInstallNoPasswdRefusesIncompleteInput(t *testing.T) {
	commands := []Command{{Path: "/usr/local/bin/serv-kun"}}

	if err := InstallNoPasswd("", commands); err == nil {
		t.Error("InstallNoPasswd with no username returned nil, want an error")
	}
	if err := InstallNoPasswd("deployer", nil); err == nil {
		t.Error("InstallNoPasswd with no commands returned nil, want an error")
	}
}
