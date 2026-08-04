package sshd

import (
	"strings"
	"testing"
)

func TestRenderIncludesSelectedDirectives(t *testing.T) {
	out := Render(Options{DisableRootLogin: true, DisablePasswordAuth: true})

	for _, want := range []string{
		"PermitRootLogin no",
		"PasswordAuthentication no",
		"KbdInteractiveAuthentication no",
		"PubkeyAuthentication yes",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("Render() missing %q\n---\n%s", want, out)
		}
	}
}

func TestRenderOmitsUnselectedDirectives(t *testing.T) {
	out := Render(Options{DisableRootLogin: true})

	if !strings.Contains(out, "PermitRootLogin no") {
		t.Errorf("Render() missing PermitRootLogin\n---\n%s", out)
	}
	if strings.Contains(out, "PasswordAuthentication") {
		t.Errorf("Render() set PasswordAuthentication when it wasn't requested\n---\n%s", out)
	}
}

func TestOptionsAny(t *testing.T) {
	if (Options{}).Any() {
		t.Error("empty Options should report no changes")
	}
	if !(Options{DisableRootLogin: true}).Any() {
		t.Error("Options with a change should report Any() == true")
	}
}

func TestRewriteConfigCommentsOutExistingSettings(t *testing.T) {
	const original = `# sample sshd_config
PermitRootLogin yes
PasswordAuthentication yes
X11Forwarding no
`
	out := rewriteConfig(original, Options{DisableRootLogin: true, DisablePasswordAuth: true}.directives())

	if !strings.Contains(out, "# serv-kun: superseded above — PermitRootLogin yes") {
		t.Errorf("existing PermitRootLogin was not commented out\n---\n%s", out)
	}
	if !strings.Contains(out, "# serv-kun: superseded above — PasswordAuthentication yes") {
		t.Errorf("existing PasswordAuthentication was not commented out\n---\n%s", out)
	}
	// Unmanaged directives must survive untouched.
	if !strings.Contains(out, "\nX11Forwarding no") {
		t.Errorf("unmanaged directive was modified\n---\n%s", out)
	}
}

// sshd takes the first value it obtains for a keyword, so our block only
// works if it lands ahead of whatever was already in the file.
func TestRewriteConfigPutsOurBlockFirst(t *testing.T) {
	const original = "PermitRootLogin yes\n"
	out := rewriteConfig(original, Options{DisableRootLogin: true}.directives())

	ours := strings.Index(out, "PermitRootLogin no")
	theirs := strings.Index(out, "# serv-kun: superseded above")
	if ours == -1 || theirs == -1 {
		t.Fatalf("expected both our directive and the commented original\n---\n%s", out)
	}
	if ours > theirs {
		t.Errorf("our block must come first, got ours=%d theirs=%d\n---\n%s", ours, theirs, out)
	}
}

// Directives inside a Match block are scoped overrides — rewriting them
// would change behavior for connections we were never asked about.
func TestRewriteConfigLeavesMatchBlocksAlone(t *testing.T) {
	const original = `PermitRootLogin yes

Match User backup
    PasswordAuthentication yes
`
	out := rewriteConfig(original, Options{DisableRootLogin: true, DisablePasswordAuth: true}.directives())

	if !strings.Contains(out, "    PasswordAuthentication yes") {
		t.Errorf("Match-block directive should be untouched\n---\n%s", out)
	}
	if !strings.Contains(out, "# serv-kun: superseded above — PermitRootLogin yes") {
		t.Errorf("top-level directive should still be commented out\n---\n%s", out)
	}
}

func TestRewriteConfigPreservesTrailingContent(t *testing.T) {
	const original = "Port 22\nX11Forwarding no\n"
	out := rewriteConfig(original, Options{DisableRootLogin: true}.directives())

	if !strings.Contains(out, "Port 22") {
		t.Errorf("Port directive lost\n---\n%s", out)
	}
	if strings.Contains(out, "\n\n\n") {
		t.Errorf("rewrite introduced blank-line runs\n---\n%s", out)
	}
}

func TestIsKeyLine(t *testing.T) {
	valid := []string{
		"ssh-ed25519 AAAAC3Nz lahat",
		"ssh-rsa AAAAB3Nz user@host",
		"ecdsa-sha2-nistp256 AAAAE2Vj user@host",
		// Options may precede the key type.
		`no-port-forwarding,command="/bin/echo" ssh-ed25519 AAAAC3Nz user`,
	}
	for _, line := range valid {
		if !isKeyLine(line) {
			t.Errorf("isKeyLine(%q) = false, want true", line)
		}
	}

	invalid := []string{
		"",
		"not a key at all",
	}
	for _, line := range invalid {
		if isKeyLine(line) {
			t.Errorf("isKeyLine(%q) = true, want false", line)
		}
	}
}
