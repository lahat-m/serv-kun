package springprops

import (
	"os"
	"path/filepath"
	"testing"
)

func writeProps(t *testing.T, name, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	return path
}

func TestScanTakesOnlyPlaceholdersWithoutFallbacks(t *testing.T) {
	path := writeProps(t, "application-prod.properties", `
server.port=${SERVER_PORT:8080}
app.storage.endpoint=${S3_ENDPOINT}
app.email.resend-api-key=${RESEND_API_KEY:}
app.name=mis
#app.commented.out=${NOT_REQUIRED}
`)

	reqs, err := Scan(path)
	if err != nil {
		t.Fatalf("Scan: %v", err)
	}
	if len(reqs) != 1 {
		t.Fatalf("Scan() = %+v, want only app.storage.endpoint", reqs)
	}
	if reqs[0].Var != "S3_ENDPOINT" || reqs[0].Property != "app.storage.endpoint" {
		t.Errorf("Scan()[0] = %+v", reqs[0])
	}
	// An empty fallback (${VAR:}) is still a fallback: the app starts without it.
	if reqs[0].Relaxed != "APP_STORAGE_ENDPOINT" {
		t.Errorf("Relaxed = %q, want APP_STORAGE_ENDPOINT", reqs[0].Relaxed)
	}
}

// The comment above a property is the only explanation of it the operator gets
// at the prompt, so it has to survive the scan.
func TestScanKeepsTheCommentAboveTheProperty(t *testing.T) {
	path := writeProps(t, "application-prod.properties", `
## Membership card QR signing — must be a real, per-environment
## 32-byte Base64 value.
app.member.card-hmac-key=${CARD_HMAC_KEY}

app.site.base-url=${SITE_BASE_URL}
`)

	reqs, err := Scan(path)
	if err != nil {
		t.Fatalf("Scan: %v", err)
	}
	if len(reqs) != 2 {
		t.Fatalf("Scan() returned %d requirements, want 2", len(reqs))
	}
	want := "Membership card QR signing — must be a real, per-environment 32-byte Base64 value."
	if reqs[0].Help != want {
		t.Errorf("Help = %q, want %q", reqs[0].Help, want)
	}
	// A blank line ends the comment block — the next property must not inherit it.
	if reqs[1].Help != "" {
		t.Errorf("app.site.base-url inherited a comment: %q", reqs[1].Help)
	}
}

// Only a value that is *exactly* one placeholder can be satisfied by setting
// the property itself, because an environment variable replaces the whole
// value — including any text around the placeholder.
func TestScanEmbeddedPlaceholderHasNoRelaxedAlternative(t *testing.T) {
	path := writeProps(t, "application-prod.properties",
		"spring.datasource.url=jdbc:postgresql://${DB_HOST}:5432/${DB_NAME}\n")

	reqs, err := Scan(path)
	if err != nil {
		t.Fatalf("Scan: %v", err)
	}
	if len(reqs) != 2 {
		t.Fatalf("Scan() = %+v, want both DB_HOST and DB_NAME", reqs)
	}
	for _, req := range reqs {
		if req.Relaxed != "" {
			t.Errorf("%s offered relaxed name %q for an embedded placeholder", req.Var, req.Relaxed)
		}
	}
}

// A profile that supplies a fallback for a property the base file left bare
// makes it optional — the app starts without it, so the wizard must not ask.
func TestScanProfileFallbackCancelsABaseRequirement(t *testing.T) {
	dir := t.TempDir()
	base := filepath.Join(dir, "application.properties")
	prod := filepath.Join(dir, "application-prod.properties")
	if err := os.WriteFile(base, []byte("app.email.provider=${EMAIL_PROVIDER}\n"), 0644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	if err := os.WriteFile(prod, []byte("app.email.provider=${EMAIL_PROVIDER:resend}\n"), 0644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	reqs, err := Scan(base, prod)
	if err != nil {
		t.Fatalf("Scan: %v", err)
	}
	if len(reqs) != 0 {
		t.Errorf("Scan() = %+v, want nothing — the prod profile supplies a fallback", reqs)
	}
}

func TestRelaxedName(t *testing.T) {
	for property, want := range map[string]string{
		"app.storage.endpoint":     "APP_STORAGE_ENDPOINT",
		"app.member.card-hmac-key": "APP_MEMBER_CARD_HMAC_KEY",
		"jwt.public-key":           "JWT_PUBLIC_KEY",
		"spring.datasource.url":    "SPRING_DATASOURCE_URL",
	} {
		if got := RelaxedName(property); got != want {
			t.Errorf("RelaxedName(%q) = %q, want %q", property, got, want)
		}
	}
}

// The case this package exists for, taken from a real deployment: a project
// whose variable names differ from the wizard's throughout. Eight of these
// properties are satisfied anyway — Spring accepts the property's own relaxed
// name, and the wizard happens to write exactly that — and flagging them would
// bury the seven that genuinely are missing.
func TestUnsatisfiedIgnoresWhatRelaxedBindingAlreadyCovers(t *testing.T) {
	path := writeProps(t, "application-prod.properties", `
jwt.public-key=${JWT_PUBLIC_KEY_PATH}
jwt.private-key=${JWT_PRIVATE_KEY_PATH}
app.member.card-hmac-key=${CARD_HMAC_KEY}
app.admin.password=${ADMIN_SEED_PASSWORD}
app.cors.allowed-origins=${CORS_ALLOWED_ORIGINS}
spring.datasource.url=${DB_URL}
spring.datasource.username=${DB_USERNAME}
spring.datasource.password=${DB_PASSWORD}
app.storage.endpoint=${S3_ENDPOINT}
app.storage.region=${S3_REGION}
app.storage.access-key=${S3_ACCESS_KEY}
app.storage.secret-key=${S3_SECRET_KEY}
app.storage.bucket=${S3_BUCKET}
app.site.base-url=${SITE_BASE_URL}
app.email.resend-api-key=${RESEND_API_KEY}
`)

	// Exactly what the wizard's template writes.
	defined := map[string]bool{}
	for _, key := range []string{
		"APP_DOMAIN", "CERTBOT_EMAIL", "APP_IMAGE",
		"POSTGRES_DB", "POSTGRES_USER", "POSTGRES_PASSWORD",
		"SPRING_DATASOURCE_URL", "SPRING_DATASOURCE_USERNAME", "SPRING_DATASOURCE_PASSWORD",
		"BACKUP_KEEP_DAYS", "APP_CORS_ALLOWED_ORIGINS",
		"JWT_PUBLIC_KEY", "JWT_PRIVATE_KEY",
		"APP_ADMIN_EMAIL", "APP_ADMIN_PASSWORD", "APP_ADMIN_FULL_NAME",
		"RESEND_API_KEY", "RESEND_FROM",
		"APP_S3_ENDPOINT", "APP_S3_ACCESS_KEY", "APP_S3_SECRET_KEY", "APP_S3_BUCKET", "APP_S3_REGION",
	} {
		defined[key] = true
	}

	reqs, err := Scan(path)
	if err != nil {
		t.Fatalf("Scan: %v", err)
	}

	var got []string
	for _, req := range Unsatisfied(reqs, defined) {
		got = append(got, req.Var)
	}
	want := []string{
		"CARD_HMAC_KEY",
		"S3_ENDPOINT", "S3_REGION", "S3_ACCESS_KEY", "S3_SECRET_KEY", "S3_BUCKET",
		"SITE_BASE_URL",
	}

	if len(got) != len(want) {
		t.Fatalf("Unsatisfied() = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("Unsatisfied()[%d] = %q, want %q (full: %v)", i, got[i], want[i], got)
		}
	}
}

func TestLocateReportsYAMLItCannotRead(t *testing.T) {
	dir := t.TempDir()
	resources := filepath.Join(dir, "src", "main", "resources")
	if err := os.MkdirAll(resources, 0755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	if err := os.WriteFile(filepath.Join(resources, "application.yml"), []byte("server:\n  port: 8080\n"), 0644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	files, yamlOnly := Locate(dir, "prod")
	if len(files) != 0 {
		t.Errorf("Locate() returned %v, want none", files)
	}
	if !yamlOnly {
		t.Error("Locate() yamlOnly = false — a YAML-configured project would be reported as fully covered")
	}
}

// No sources on the box is the normal shape of a pull-only deployment, not an
// error, and it must not be confused with a project that has YAML.
func TestLocateOnAProjectWithoutSources(t *testing.T) {
	files, yamlOnly := Locate(t.TempDir(), "prod")
	if len(files) != 0 || yamlOnly {
		t.Errorf("Locate() = %v, yamlOnly=%v; want none, false", files, yamlOnly)
	}
}
