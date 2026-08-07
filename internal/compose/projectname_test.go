package compose

import (
	"os"
	"path/filepath"
	"testing"
)

func writeFile(t *testing.T, dir, name, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0644); err != nil {
		t.Fatalf("WriteFile %s: %v", name, err)
	}
}

func TestNormalizeProjectName(t *testing.T) {
	for in, want := range map[string]string{
		"splm-mis":  "splm-mis",
		"SPLM-MIS":  "splm-mis",
		"splm mis":  "splmmis",
		"my.app":    "myapp",
		"_leading":  "leading",
		"-leading":  "leading",
		"trailing_": "trailing_",
		"":          "",
	} {
		if got := NormalizeProjectName(in); got != want {
			t.Errorf("NormalizeProjectName(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestProjectNameOfFallsBackToTheDirectoryName(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "SPLM-MIS")
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	writeFile(t, dir, "docker-compose.prod.yml", "services:\n  app:\n    image: x\n")

	if got := ProjectNameOf(dir, "docker-compose.prod.yml"); got != "splm-mis" {
		t.Errorf("ProjectNameOf() = %q, want %q", got, "splm-mis")
	}
}

func TestProjectNameOfReadsTheTopLevelNameKey(t *testing.T) {
	dir := t.TempDir()
	for _, c := range []struct{ body, want string }{
		{"name: splm-mis-prod\nservices:\n  app:\n    image: x\n", "splm-mis-prod"},
		{"name: \"Quoted-Name\"\nservices: {}\n", "quoted-name"},
		{"name: withcomment # why\nservices: {}\n", "withcomment"},
		// Indented, so not the project name — it's a key inside a service.
		{"services:\n  app:\n    name: notthis\n", filepath.Base(dir)},
	} {
		writeFile(t, dir, "docker-compose.prod.yml", c.body)
		want := NormalizeProjectName(c.want)
		if got := ProjectNameOf(dir, "docker-compose.prod.yml"); got != want {
			t.Errorf("ProjectNameOf() = %q, want %q\n---\n%s", got, want, c.body)
		}
	}
}

// The deployment this check exists for: a repo carrying both a dev
// docker-compose.yml and serv-kun's docker-compose.prod.yml in one directory,
// each declaring a service called "postgres". Compose resolves both to the
// directory name, and a dev `up` recreates the production database container.
func TestCollisionsFindsADevComposeFileSharingTheProjectName(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "splm-mis")
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	writeFile(t, dir, "docker-compose.prod.yml", `services:
  app:
    container_name: mis-app
  postgres:
    container_name: mis-db
  nginx:
    container_name: mis-nginx
`)
	writeFile(t, dir, "docker-compose.yml", `name: splm-mis

services:
  postgres:
    container_name: splm-mis-postgres
  minio:
    image: minio/minio:latest
`)

	got := Collisions(dir, "docker-compose.prod.yml")
	if len(got) != 1 {
		t.Fatalf("Collisions() = %+v, want one entry for docker-compose.yml", got)
	}
	if got[0].File != "docker-compose.yml" {
		t.Errorf("File = %q, want docker-compose.yml", got[0].File)
	}
	if len(got[0].Services) != 1 || got[0].Services[0] != "postgres" {
		t.Errorf("Services = %v, want [postgres] — the service both files declare", got[0].Services)
	}
}

func TestCollisionsIgnoresAFileWithItsOwnName(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "splm-mis")
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	writeFile(t, dir, "docker-compose.prod.yml", "name: splm-mis-prod\nservices:\n  postgres:\n    image: postgres\n")
	writeFile(t, dir, "docker-compose.yml", "services:\n  postgres:\n    image: postgres\n")

	if got := Collisions(dir, "docker-compose.prod.yml"); len(got) != 0 {
		t.Errorf("Collisions() = %+v, want none — the prod file names its own project", got)
	}
}

// An override file is merged into another file's project on purpose, so the
// shared name is the feature, not the bug. Reporting it would train operators
// to ignore this warning.
func TestCollisionsIgnoresOverrideAndNonComposeFiles(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "splm-mis")
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	writeFile(t, dir, "docker-compose.prod.yml", "services:\n  app:\n    image: x\n")
	writeFile(t, dir, "docker-compose.override.yml", "services:\n  app:\n    image: y\n")
	writeFile(t, dir, "compose.override.yaml", "services:\n  app:\n    image: y\n")
	writeFile(t, dir, "application.yml", "services:\n  app:\n    image: y\n")
	writeFile(t, dir, "notes.md", "services:\n")

	if got := Collisions(dir, "docker-compose.prod.yml"); len(got) != 0 {
		t.Errorf("Collisions() = %+v, want none", got)
	}
}

func TestCollisionsReportsAFileWithNoSharedServices(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "splm-mis")
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	writeFile(t, dir, "docker-compose.prod.yml", "services:\n  app:\n    image: x\n")
	writeFile(t, dir, "docker-compose.tools.yml", "services:\n  adminer:\n    image: adminer\n")

	// Still a collision: --remove-orphans in either direction removes the
	// other file's containers whether or not any service name is shared.
	got := Collisions(dir, "docker-compose.prod.yml")
	if len(got) != 1 {
		t.Fatalf("Collisions() = %+v, want one entry", got)
	}
	if len(got[0].Services) != 0 {
		t.Errorf("Services = %v, want none", got[0].Services)
	}
}

func TestServiceNamesSkipsNestedKeys(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "docker-compose.yml", `# a comment
name: x

services:

  app:
    image: x
    environment:
      POSTGRES_DB: y
    depends_on:
      postgres:
        condition: service_healthy

  postgres:
    image: postgres

volumes:
  postgres-data:
`)

	got := serviceNames(filepath.Join(dir, "docker-compose.yml"))
	want := []string{"app", "postgres"}
	if len(got) != len(want) {
		t.Fatalf("serviceNames() = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("serviceNames()[%d] = %q, want %q (full: %v)", i, got[i], want[i], got)
		}
	}
}
