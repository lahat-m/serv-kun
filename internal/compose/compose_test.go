package compose

import "testing"

func TestProjectArgs(t *testing.T) {
	p := Project{Dir: "/opt/app", File: "docker-compose.prod.yml", EnvFile: ".env.prod"}

	got := p.Args("up", "-d", "app")
	want := []string{"compose", "-f", "docker-compose.prod.yml", "--env-file", ".env.prod", "up", "-d", "app"}

	if len(got) != len(want) {
		t.Fatalf("Args() = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("Args()[%d] = %q, want %q (full: %v)", i, got[i], want[i], got)
		}
	}
}

// An image that runs as root needs no chown at all, so every spelling of root
// has to be recognised — including the empty string Docker reports for an
// image whose Dockerfile never set USER.
func TestIsRootSpec(t *testing.T) {
	for _, spec := range []string{"", "0", "root", "0:0", "root:root"} {
		if !isRootSpec(spec) {
			t.Errorf("isRootSpec(%q) = false, want true", spec)
		}
	}
	for _, spec := range []string{"1001", "appuser", "appuser:appgroup", "1001:1001"} {
		if isRootSpec(spec) {
			t.Errorf("isRootSpec(%q) = true, want false", spec)
		}
	}
}

func TestParseNumericUserSpec(t *testing.T) {
	tests := []struct {
		spec     string
		uid, gid int
		ok       bool
	}{
		{spec: "1001", uid: 1001, gid: 1001, ok: true},
		{spec: "1001:1002", uid: 1001, gid: 1002, ok: true},
		// The group name is unresolvable without the image's passwd file, but
		// the uid is what a chown actually needs — don't throw it away.
		{spec: "1001:appgroup", uid: 1001, gid: 1001, ok: true},
		// A named user carries no number at all.
		{spec: "appuser", ok: false},
		{spec: "appuser:appgroup", ok: false},
	}

	for _, tt := range tests {
		uid, gid, ok := parseNumericUserSpec(tt.spec)
		if ok != tt.ok {
			t.Errorf("parseNumericUserSpec(%q) ok = %v, want %v", tt.spec, ok, tt.ok)
			continue
		}
		if !tt.ok {
			continue
		}
		if uid != tt.uid || gid != tt.gid {
			t.Errorf("parseNumericUserSpec(%q) = %d:%d, want %d:%d", tt.spec, uid, gid, tt.uid, tt.gid)
		}
	}
}

func TestProjectArgsNoExtra(t *testing.T) {
	p := Project{Dir: ".", File: "docker-compose.prod.yml", EnvFile: ".env.prod"}

	got := p.Args()
	want := []string{"compose", "-f", "docker-compose.prod.yml", "--env-file", ".env.prod"}

	if len(got) != len(want) {
		t.Fatalf("Args() = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("Args()[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}
