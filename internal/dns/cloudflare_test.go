package dns

import "testing"

func TestIsCloudflare(t *testing.T) {
	cloudflare := []string{
		"104.16.0.1",      // 104.16.0.0/13
		"172.64.1.2",      // 172.64.0.0/13
		"131.0.72.5",      // 131.0.72.0/22
		"2606:4700::1111", // 2606:4700::/32
	}
	for _, ip := range cloudflare {
		if !isCloudflare(ip) {
			t.Errorf("isCloudflare(%q) = false, want true", ip)
		}
	}

	origin := []string{
		"203.0.113.10",
		"1.1.1.1", // Cloudflare's resolver, deliberately not in the edge ranges
		"192.0.2.1",
		"2001:db8::1",
		"not-an-ip",
		"",
	}
	for _, ip := range origin {
		if isCloudflare(ip) {
			t.Errorf("isCloudflare(%q) = true, want false", ip)
		}
	}
}

func TestAllCloudflare(t *testing.T) {
	if !allCloudflare([]string{"104.16.0.1", "172.64.1.2"}) {
		t.Error("all-Cloudflare set should report true")
	}
	if allCloudflare([]string{"104.16.0.1", "203.0.113.10"}) {
		t.Error("mixed set should report false")
	}
	if allCloudflare(nil) {
		t.Error("empty set should report false")
	}
}

// A proxied record resolves to Cloudflare's edge and never to the origin, so
// treating the mismatch as failure made the check cry wolf on exactly the
// setup DEPLOY.md recommends.
func TestProxiedRecordCountsAsOK(t *testing.T) {
	r := Result{Domain: "example.com", ResolvedIPs: []string{"104.16.0.1"}}
	if allCloudflare(r.ResolvedIPs) {
		r.OK = true
		r.Proxied = true
	}

	if !r.OK || !r.Proxied {
		t.Errorf("proxied result = %+v, want OK and Proxied", r)
	}
	if !AllOK([]Result{r}) {
		t.Error("AllOK should accept a proxied result")
	}
}

func TestAllOKRejectsFailure(t *testing.T) {
	results := []Result{
		{Domain: "ok.example.com", OK: true},
		{Domain: "bad.example.com", OK: false},
	}
	if AllOK(results) {
		t.Error("AllOK should be false when any result failed")
	}
}
