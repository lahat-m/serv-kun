package system

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type OSFamily string

const (
	Debian OSFamily = "debian"
	RHEL   OSFamily = "rhel"
)

type Info struct {
	Family OSFamily
	ID     string // ubuntu, debian, centos, fedora, rhel, rocky ...
	Pretty string // "Ubuntu 22.04 LTS"
}

func Detect() (*Info, error) {
	f, err := os.Open("/etc/os-release")
	if err != nil {
		return nil, fmt.Errorf("cannot read /etc/os-release: %w", err)
	}
	defer f.Close()

	vars := parseKV(f)

	id := strings.ToLower(vars["ID"])
	pretty := strings.Trim(vars["PRETTY_NAME"], `"`)
	if pretty == "" {
		pretty = id
	}

	family, err := resolveFamily(id, vars["ID_LIKE"])
	if err != nil {
		return nil, fmt.Errorf("unsupported OS %q: %w", pretty, err)
	}

	return &Info{Family: family, ID: id, Pretty: pretty}, nil
}

func resolveFamily(id, idLike string) (OSFamily, error) {
	debianIDs := []string{"ubuntu", "debian", "raspbian", "linuxmint", "pop"}
	rhelIDs := []string{"centos", "rhel", "fedora", "rocky", "almalinux", "ol", "amzn"}

	for _, d := range debianIDs {
		if id == d {
			return Debian, nil
		}
	}
	for _, r := range rhelIDs {
		if id == r {
			return RHEL, nil
		}
	}

	like := strings.ToLower(idLike)
	if strings.Contains(like, "debian") || strings.Contains(like, "ubuntu") {
		return Debian, nil
	}
	if strings.Contains(like, "rhel") || strings.Contains(like, "fedora") || strings.Contains(like, "centos") {
		return RHEL, nil
	}

	return "", fmt.Errorf("could not determine OS family from ID_LIKE=%q", idLike)
}

// GetPublicIP tries several well-known IP-echo services in order.
func GetPublicIP() (string, error) {
	services := []string{
		"https://checkip.amazonaws.com",
		"https://api.ipify.org",
		"https://ifconfig.me/ip",
	}

	client := &http.Client{Timeout: 6 * time.Second}

	for _, svc := range services {
		resp, err := client.Get(svc)
		if err != nil {
			continue
		}
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			continue
		}
		ip := strings.TrimSpace(string(body))
		if ip != "" {
			return ip, nil
		}
	}

	return "", fmt.Errorf("could not determine public IP from any external service")
}

func parseKV(r io.Reader) map[string]string {
	m := make(map[string]string)
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		line := sc.Text()
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		m[strings.TrimSpace(parts[0])] = strings.Trim(strings.TrimSpace(parts[1]), `"`)
	}
	return m
}
