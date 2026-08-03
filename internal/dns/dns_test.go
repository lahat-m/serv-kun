package dns

import "testing"

// These exercise the matching logic in checkDomain/AllOK against "localhost",
// which resolves without needing external network access (only /etc/hosts).

func TestCheckDomainsMatch(t *testing.T) {
	results := CheckDomains([]string{"localhost"}, "127.0.0.1")
	if len(results) != 1 {
		t.Fatalf("got %d results, want 1", len(results))
	}
	if !results[0].OK {
		t.Errorf("localhost vs 127.0.0.1: OK = false, err = %v", results[0].Err)
	}
	if !AllOK(results) {
		t.Error("AllOK() = false, want true")
	}
}

func TestCheckDomainsMismatch(t *testing.T) {
	// TEST-NET-3 (RFC 5737) — guaranteed not to be localhost's address.
	results := CheckDomains([]string{"localhost"}, "203.0.113.1")
	if len(results) != 1 {
		t.Fatalf("got %d results, want 1", len(results))
	}
	if results[0].OK {
		t.Error("localhost vs 203.0.113.1: OK = true, want false")
	}
	if results[0].Err == nil {
		t.Error("expected a non-nil Err on mismatch")
	}
	if AllOK(results) {
		t.Error("AllOK() = true, want false")
	}
}

func TestAllOKEmpty(t *testing.T) {
	if !AllOK(nil) {
		t.Error("AllOK(nil) = false, want true (vacuous truth)")
	}
}
