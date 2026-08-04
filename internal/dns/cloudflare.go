package dns

import "net"

// cloudflareCIDRs is Cloudflare's published edge range list, from
// https://www.cloudflare.com/ips/ — the same list nginx/conf.d/hardening.conf
// trusts for set_real_ip_from. Cloudflare adds ranges occasionally, so this
// needs an eyeball every so often; a missed range only costs a false
// "doesn't point here" warning, never a wrong success.
var cloudflareCIDRs = []string{
	"173.245.48.0/20",
	"103.21.244.0/22",
	"103.22.200.0/22",
	"103.31.4.0/22",
	"141.101.64.0/18",
	"108.162.192.0/18",
	"190.93.240.0/20",
	"188.114.96.0/20",
	"197.234.240.0/22",
	"198.41.128.0/17",
	"162.158.0.0/15",
	"104.16.0.0/13",
	"104.24.0.0/14",
	"172.64.0.0/13",
	"131.0.72.0/22",
	"2400:cb00::/32",
	"2606:4700::/32",
	"2803:f800::/32",
	"2405:b500::/32",
	"2405:8100::/32",
	"2a06:98c0::/29",
	"2c0f:f248::/32",
}

var cloudflareNets = parseCIDRs(cloudflareCIDRs)

func parseCIDRs(cidrs []string) []*net.IPNet {
	nets := make([]*net.IPNet, 0, len(cidrs))
	for _, c := range cidrs {
		if _, n, err := net.ParseCIDR(c); err == nil {
			nets = append(nets, n)
		}
	}
	return nets
}

// isCloudflare reports whether ip belongs to a Cloudflare edge range.
func isCloudflare(ip string) bool {
	parsed := net.ParseIP(ip)
	if parsed == nil {
		return false
	}
	for _, n := range cloudflareNets {
		if n.Contains(parsed) {
			return true
		}
	}
	return false
}

// allCloudflare reports whether every resolved address is a Cloudflare edge
// IP — the signature of a proxied ("orange cloud") DNS record.
func allCloudflare(ips []string) bool {
	if len(ips) == 0 {
		return false
	}
	for _, ip := range ips {
		if !isCloudflare(ip) {
			return false
		}
	}
	return true
}
