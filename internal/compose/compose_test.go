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
