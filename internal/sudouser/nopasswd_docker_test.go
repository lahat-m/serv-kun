//go:build docker

// These tests need a working docker daemon and are excluded from the ordinary
// `go test ./...`. Run them with:
//
//	go test -tags docker ./internal/sudouser/
//
// They exist because the only opinion that counts about a sudoers file is
// sudo's. Rendering the right string is easy to unit-test and proves nothing:
// the rule has to parse, apply to the named user, and — the case a passing
// unit test would hide — *not* apply to anything else.
package sudouser

import (
	"os/exec"
	"strings"
	"testing"
)

// checkRuleScript installs sudo, feeds it the rule on stdin, then exercises it
// as the granted user. The two fake binaries stand in for serv-kun and docker
// so the test turns on sudoers matching rather than on either program.
const checkRuleScript = `
set -e
apt-get update -qq >/dev/null 2>&1
apt-get install -y -qq sudo >/dev/null 2>&1

cat > /etc/sudoers.d/serv-kun
chmod 0440 /etc/sudoers.d/serv-kun
visudo -c -f /etc/sudoers.d/serv-kun

useradd -m deployer
printf '#!/bin/sh\necho SERVKUN-RAN\n' > /usr/local/bin/serv-kun
printf '#!/bin/sh\necho DOCKER-RAN\n' > /usr/bin/docker
chmod +x /usr/local/bin/serv-kun /usr/bin/docker

# -n is what makes this a real test of the CI path: no terminal, no prompt.
echo "SERVKUN: $(su deployer -c 'sudo -n serv-kun springboot deploy' 2>&1 || true)"
echo "PRUNE: $(su deployer -c 'sudo -n docker image prune -f' 2>&1 || true)"
echo "UNSCOPED: $(su deployer -c 'sudo -n docker run -v /:/host debian' 2>&1 || true)"
echo "OTHER: $(su deployer -c 'sudo -n cat /etc/shadow' 2>&1 || true)"
`

func TestGeneratedRuleAppliesToTheGrantedCommandsOnly(t *testing.T) {
	rule := NoPasswdRule("deployer", []Command{
		{Path: "/usr/local/bin/serv-kun"},
		{Path: "/usr/bin/docker", Args: []string{"image", "prune", "-f"}},
	})

	cmd := exec.Command("docker", "run", "--rm", "-i", "debian:bookworm-slim", "sh", "-c", checkRuleScript)
	cmd.Stdin = strings.NewReader(rule)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("docker run: %v\n%s", err, out)
	}
	got := string(out)

	if !strings.Contains(got, "parsed OK") {
		t.Fatalf("visudo did not accept the generated rule\n%s", got)
	}

	// Granted, with no password, over a session that has no terminal — the
	// whole reason the file exists.
	if !strings.Contains(got, "SERVKUN: SERVKUN-RAN") {
		t.Errorf("`sudo -n serv-kun springboot deploy` was not permitted\n%s", got)
	}
	if !strings.Contains(got, "PRUNE: DOCKER-RAN") {
		t.Errorf("`sudo -n docker image prune -f` was not permitted\n%s", got)
	}

	// And the scoping half. Pinning docker's arguments is the difference
	// between granting one cleanup command and granting a root shell, since
	// `docker run -v /:/host` mounts the host filesystem.
	for _, label := range []string{"UNSCOPED", "OTHER"} {
		line := lineWithPrefix(got, label+": ")
		if strings.Contains(line, "-RAN") || strings.Contains(line, "root:") {
			t.Errorf("%s should not have been permitted by the rule: %q\n%s", label, line, got)
		}
	}
}

func lineWithPrefix(out, prefix string) string {
	for _, line := range strings.Split(out, "\n") {
		if strings.HasPrefix(line, prefix) {
			return line
		}
	}
	return ""
}
