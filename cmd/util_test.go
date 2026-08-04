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
