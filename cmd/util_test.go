package cmd

import "testing"

func TestCheckDomainAccepts(t *testing.T) {
	valid := []string{
		"example.com",
		"kuniel.com",
		"api.kuniel.com",
		"a.b.c.example.co.uk",
		"my-app.example.com",
		"xn--80ak6aa92e.com", // punycode
		"123app.example.com",
	}
	for _, d := range valid {
		if err := checkDomain(d); err != nil {
			t.Errorf("checkDomain(%q) = %v, want nil", d, err)
		}
	}
}

func TestCheckDomainRejects(t *testing.T) {
	invalid := []string{
		"",
		"localhost",           // no dot, can't get a public cert
		"example .com",        // the case that word-split an rm -rf
		"example.com ",        // trailing space
		"https://example.com", // scheme
		"example.com/path",    // path
		"example.com:8080",    // port
		"-example.com",        // label starts with a hyphen
		"example-.com",        // label ends with a hyphen
		"example.c",           // one-character TLD
		"example.123",         // numeric TLD
		"exa mple.com",
	}
	for _, d := range invalid {
		if err := checkDomain(d); err == nil {
			t.Errorf("checkDomain(%q) = nil, want an error", d)
		}
	}
}

func TestValidateDomainTolerantOfSurroundingSpace(t *testing.T) {
	// The validator trims before checking; call sites trim the stored answer.
	if err := validateDomain(" example.com "); err != nil {
		t.Errorf("validateDomain with padding = %v, want nil", err)
	}
	if err := validateDomain("not a domain"); err == nil {
		t.Error("validateDomain(\"not a domain\") = nil, want an error")
	}
}

// server_name may be an IP or a dotless internal host — only certbot needs a
// public domain. Rejecting these broke a legitimate reverse-proxy setup.
func TestCheckServerNameAcceptsWhatCheckDomainRejects(t *testing.T) {
	valid := []string{
		"example.com",
		"myserver",     // dotless internal host
		"app.local",    // non-public TLD
		"203.0.113.10", // IPv4
		"2001:db8::1",  // IPv6
		"localhost",
		"box-01.lan",
	}
	for _, s := range valid {
		if err := checkServerName(s); err != nil {
			t.Errorf("checkServerName(%q) = %v, want nil", s, err)
		}
	}
}

func TestCheckServerNameStillRejectsTypos(t *testing.T) {
	invalid := []string{
		"",
		"my server",
		"https://example.com",
		"example.com/path",
		"-bad.example.com",
	}
	for _, s := range invalid {
		if err := checkServerName(s); err == nil {
			t.Errorf("checkServerName(%q) = nil, want an error", s)
		}
	}
}

func TestCheckPort(t *testing.T) {
	for _, s := range []string{"1", "22", "8080", "65535"} {
		if err := checkPort(s); err != nil {
			t.Errorf("checkPort(%q) = %v, want nil", s, err)
		}
	}
	for _, s := range []string{"", "0", "65536", "-1", "abc", "80a"} {
		if err := checkPort(s); err == nil {
			t.Errorf("checkPort(%q) = nil, want an error", s)
		}
	}
}

func TestValidateOptionalPortList(t *testing.T) {
	if err := validateOptionalPortList(""); err != nil {
		t.Errorf("blank list = %v, want nil", err)
	}
	if err := validateOptionalPortList("8080 3000"); err != nil {
		t.Errorf("valid list = %v, want nil", err)
	}
	if err := validateOptionalPortList("8080 99999"); err == nil {
		t.Error("out-of-range port should be rejected")
	}
}

func TestValidateDomainList(t *testing.T) {
	if err := validateDomainList("example.com www.example.com"); err != nil {
		t.Errorf("valid list = %v, want nil", err)
	}
	if err := validateDomainList(""); err == nil {
		t.Error("empty list should be rejected")
	}
	if err := validateDomainList("   "); err == nil {
		t.Error("whitespace-only list should be rejected")
	}
	if err := validateDomainList("example.com not-a-domain"); err == nil {
		t.Error("list containing an invalid entry should be rejected")
	}
}
