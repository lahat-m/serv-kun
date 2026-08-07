package scaffold

import (
	"strings"
	"testing"
)

func testParams() Params {
	return Params{
		AppName:      "myapp",
		AppPort:      "8080",
		MaxBodySize:  "50M",
		IncludeCerts: true,
		ProjectName:  "myapp-prod",
	}
}

// The name has to be a top-level key — at column 0 — or compose treats it as
// part of whatever block precedes it and silently falls back to the directory
// name, which is the default this exists to escape.
func TestComposeDeclaresProjectNameAtTopLevel(t *testing.T) {
	files, err := Generate(testParams())
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	compose := files[0].Content

	if !strings.Contains(compose, "\nname: myapp-prod\n") {
		t.Errorf("compose file has no top-level name: key\n---\n%s", compose)
	}
	if strings.Index(compose, "\nname: myapp-prod\n") > strings.Index(compose, "\nservices:") {
		t.Errorf("name: must come before services:\n---\n%s", compose)
	}
}

func TestGenerateReturnsAllFiles(t *testing.T) {
	files, err := Generate(testParams())
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}

	want := []string{
		"docker-compose.prod.yml",
		"nginx/templates/app.conf.template",
		"nginx/conf.d/hardening.conf",
		"scripts/backup-postgres.sh",
		".github/workflows/deploy.yml",
		"DEPLOY.md",
	}
	if len(files) != len(want) {
		t.Fatalf("Generate() returned %d files, want %d", len(files), len(want))
	}
	for i, w := range want {
		if files[i].Path != w {
			t.Errorf("files[%d].Path = %q, want %q", i, files[i].Path, w)
		}
		if files[i].Content == "" {
			t.Errorf("files[%d] (%s) has empty content", i, files[i].Path)
		}
	}
}

func TestBackupScriptIsExecutable(t *testing.T) {
	files, err := Generate(testParams())
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	for _, f := range files {
		if f.Path == "scripts/backup-postgres.sh" && !f.Exec {
			t.Error("scripts/backup-postgres.sh should be marked Exec")
		}
	}
}

func TestComposeUsesAppNameAndPort(t *testing.T) {
	files, err := Generate(testParams())
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	compose := fileContent(t, files, "docker-compose.prod.yml")

	for _, want := range []string{
		// The nginx image only runs its /docker-entrypoint.d/ scripts — the
		// envsubst pass that renders templates/ into conf.d/ — when the command
		// starts with "nginx". The cert-reload loop overrides command with a
		// shell, so the entrypoint has to be invoked explicitly; without it
		// nginx serves only its default site and ACME challenges 404.
		"exec /docker-entrypoint.sh nginx -g 'daemon off;'",
		"container_name: myapp-app",
		"container_name: myapp-db",
		"container_name: myapp-nginx",
		"./certs:/etc/myapp/certs:ro",
	} {
		if !strings.Contains(compose, want) {
			t.Errorf("compose output missing %q\n---\n%s", want, compose)
		}
	}
}

func TestComposeOmitsCertsWhenDisabled(t *testing.T) {
	p := testParams()
	p.IncludeCerts = false
	files, err := Generate(p)
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	compose := fileContent(t, files, "docker-compose.prod.yml")

	if strings.Contains(compose, "./certs:") {
		t.Errorf("compose output should not mount ./certs when IncludeCerts is false\n---\n%s", compose)
	}
}

func TestNginxConfUsesPortAndMaxBodySize(t *testing.T) {
	files, err := Generate(testParams())
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	conf := fileContent(t, files, "nginx/templates/app.conf.template")

	for _, want := range []string{
		"proxy_pass http://$upstream_app:8080;",
		"client_max_body_size 50M;",
	} {
		if !strings.Contains(conf, want) {
			t.Errorf("nginx conf missing %q\n---\n%s", want, conf)
		}
	}
}

func fileContent(t *testing.T, files []File, path string) string {
	t.Helper()
	for _, f := range files {
		if f.Path == path {
			return f.Content
		}
	}
	t.Fatalf("no file with path %q", path)
	return ""
}
