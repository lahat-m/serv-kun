package envfile

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func testAnswers() Answers {
	return Answers{
		AppDomain:          "example.com",
		CertbotEmail:       "admin@example.com",
		AppImage:           "ghcr.io/org/app:latest",
		PostgresDB:         "app",
		PostgresUser:       "app",
		PostgresPassword:   "s3cr3t",
		BackupKeepDays:     "7",
		CorsAllowedOrigins: "https://example.com",

		HasJWT:         true,
		CertsMountPath: "/etc/app/certs",

		HasAdmin:      true,
		AdminEmail:    "root@example.com",
		AdminPassword: "hunter2",
		AdminFullName: "Admin",

		HasEmail:     true,
		ResendAPIKey: "re_xxx",
		ResendFrom:   "App <noreply@example.com>",

		HasObjectStorage: true,
		S3Endpoint:       "https://fsn1.your-objectstorage.com",
		S3AccessKey:      "AKIA...",
		S3SecretKey:      "shh",
		S3Bucket:         "app-documents",
		S3Region:         "fsn1",
	}
}

func TestGenerateContainsAllFieldsWhenEnabled(t *testing.T) {
	content, err := Generate(testAnswers())
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}

	for _, want := range []string{
		"APP_DOMAIN=example.com",
		"CERTBOT_EMAIL=admin@example.com",
		"APP_IMAGE=ghcr.io/org/app:latest",
		"POSTGRES_DB=app",
		"SPRING_DATASOURCE_URL=jdbc:postgresql://postgres:5432/app",
		"APP_CORS_ALLOWED_ORIGINS=https://example.com",
		"JWT_PUBLIC_KEY=file:/etc/app/certs/public.pem",
		"JWT_PRIVATE_KEY=file:/etc/app/certs/private.pem",
		"APP_ADMIN_EMAIL=root@example.com",
		"RESEND_API_KEY=re_xxx",
		"APP_S3_BUCKET=app-documents",
	} {
		if !strings.Contains(content, want) {
			t.Errorf("Generate() output missing %q\n---\n%s", want, content)
		}
	}
}

func TestGenerateOmitsDisabledSections(t *testing.T) {
	a := testAnswers()
	a.HasJWT = false
	a.HasAdmin = false
	a.HasEmail = false
	a.HasObjectStorage = false

	content, err := Generate(a)
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}

	for _, unwanted := range []string{
		"JWT_PUBLIC_KEY", "JWT_PRIVATE_KEY",
		"APP_ADMIN_EMAIL", "APP_ADMIN_PASSWORD",
		"RESEND_API_KEY", "RESEND_FROM",
		"APP_S3_ENDPOINT", "APP_S3_ACCESS_KEY", "APP_S3_BUCKET",
	} {
		if strings.Contains(content, unwanted) {
			t.Errorf("Generate() output should not contain %q when its section is disabled\n---\n%s", unwanted, content)
		}
	}

	// The always-present fields must still be there.
	for _, want := range []string{"APP_DOMAIN=example.com", "POSTGRES_DB=app", "APP_CORS_ALLOWED_ORIGINS=https://example.com"} {
		if !strings.Contains(content, want) {
			t.Errorf("Generate() output missing always-present field %q\n---\n%s", want, content)
		}
	}
}

func TestReadRoundTrip(t *testing.T) {
	content, err := Generate(testAnswers())
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}

	path := filepath.Join(t.TempDir(), ".env.prod")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	kv, err := Read(path)
	if err != nil {
		t.Fatalf("Read: %v", err)
	}

	if kv["APP_DOMAIN"] != "example.com" {
		t.Errorf("APP_DOMAIN = %q, want %q", kv["APP_DOMAIN"], "example.com")
	}
	if kv["APP_IMAGE"] != "ghcr.io/org/app:latest" {
		t.Errorf("APP_IMAGE = %q, want %q", kv["APP_IMAGE"], "ghcr.io/org/app:latest")
	}
}

func TestReadIgnoresCommentsAndBlankLines(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".env.prod")
	body := "# a comment\n\nKEY=value\n  # indented comment\nOTHER = spaced \n"
	if err := os.WriteFile(path, []byte(body), 0600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	kv, err := Read(path)
	if err != nil {
		t.Fatalf("Read: %v", err)
	}

	if len(kv) != 2 {
		t.Fatalf("Read() = %v, want 2 entries", kv)
	}
	if kv["KEY"] != "value" {
		t.Errorf("KEY = %q, want %q", kv["KEY"], "value")
	}
	if kv["OTHER"] != "spaced" {
		t.Errorf("OTHER = %q, want %q", kv["OTHER"], "spaced")
	}
}

func TestReadMissingFile(t *testing.T) {
	if _, err := Read(filepath.Join(t.TempDir(), "does-not-exist")); err == nil {
		t.Fatal("Read() on a missing file: got nil error, want one")
	}
}
