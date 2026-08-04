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
	Proxied     bool // resolves to Cloudflare's edge rather than the origin
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

	// A proxied Cloudflare record resolves to Cloudflare's edge, never the
	// origin, so an IP mismatch here is expected rather than a misconfigured
	// record — and the HTTP-01 challenge still reaches us, because Cloudflare
	// forwards port 80 to the origin. Treating this as failure is what made
	// the old check cry wolf on the exact setup DEPLOY.md recommends.
	if allCloudflare(ips) {
		r.OK = true
		r.Proxied = true
		return r
	}

	r.Err = fmt.Errorf("resolves to %v, server IP is %s", ips, serverIP)
	return r
}
