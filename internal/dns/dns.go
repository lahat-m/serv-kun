package dns

import (
	"fmt"
	"net"
)

// Result is the outcome of checking whether a domain's A record points at serverIP.
type Result struct {
	Domain      string
	ResolvedIPs []string
	MatchedIP   string
	OK          bool
	Err         error
}

// CheckDomains resolves each domain and reports whether any A record matches serverIP.
func CheckDomains(domains []string, serverIP string) []Result {
	results := make([]Result, len(domains))
	for i, d := range domains {
		results[i] = checkDomain(d, serverIP)
	}
	return results
}

// AllOK returns true only if every result matched.
func AllOK(results []Result) bool {
	for _, r := range results {
		if !r.OK {
			return false
		}
	}
	return true
}

func checkDomain(domain, serverIP string) Result {
	r := Result{Domain: domain}

	ips, err := net.LookupHost(domain)
	if err != nil {
		r.Err = fmt.Errorf("DNS lookup failed: %w", err)
		return r
	}

	r.ResolvedIPs = ips

	for _, ip := range ips {
		if ip == serverIP {
			r.OK = true
			r.MatchedIP = ip
			return r
		}
	}

	r.Err = fmt.Errorf("resolves to %v, server IP is %s", ips, serverIP)
	return r
}
