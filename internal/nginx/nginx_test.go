package nginx

import (
	"strings"
	"testing"

	"github.com/lahat-m/serv-kun/internal/system"
)

func TestGenerateConfig(t *testing.T) {
	cfg := Config{
		Domains:     []string{"example.com", "www.example.com"},
		Upstream:    "http://localhost:3000",
		MaxBodySize: "50M",
	}

	content, err := GenerateConfig(cfg)
	if err != nil {
		t.Fatalf("GenerateConfig: %v", err)
	}

	for _, want := range []string{
		"server_name example.com www.example.com;",
		"proxy_pass http://localhost:3000;",
		"client_max_body_size 50M;",
	} {
		if !strings.Contains(content, want) {
			t.Errorf("GenerateConfig() output missing %q\n---\n%s", want, content)
		}
	}
}

func TestConfigPathByOSFamily(t *testing.T) {
	debian := &system.Info{Family: system.Debian}
	if got := ConfigPath(debian, "example.com"); got != "/etc/nginx/sites-available/example.com.conf" {
		t.Errorf("ConfigPath(debian) = %q", got)
	}

	rhel := &system.Info{Family: system.RHEL}
	if got := ConfigPath(rhel, "example.com"); got != "/etc/nginx/conf.d/example.com.conf" {
		t.Errorf("ConfigPath(rhel) = %q", got)
	}
}
