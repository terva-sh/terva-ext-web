package fetch

import (
	"net"
	"strings"
)

// AllowList holds the local-address escape hatch: targets that resolve to a
// private/reserved address are refused by default UNLESS they match an entry
// here. Entries are classified at parse time into hostnames, IPs, and CIDRs.
type AllowList struct {
	hosts map[string]bool
	ips   []net.IP
	nets  []*net.IPNet
}

// ParseAllowList classifies each entry as a CIDR, an IP, or (failing those) a
// hostname.
func ParseAllowList(entries []string) AllowList {
	a := AllowList{hosts: map[string]bool{}}
	for _, e := range entries {
		e = strings.TrimSpace(e)
		if e == "" {
			continue
		}
		if _, n, err := net.ParseCIDR(e); err == nil {
			a.nets = append(a.nets, n)
			continue
		}
		if ip := net.ParseIP(e); ip != nil {
			a.ips = append(a.ips, ip)
			continue
		}
		a.hosts[strings.ToLower(e)] = true
	}
	return a
}

// HostAllowed reports whether host was explicitly allowlisted by name.
func (a AllowList) HostAllowed(host string) bool {
	return a.hosts[strings.ToLower(strings.TrimSuffix(host, "."))]
}

// IPAllowed reports whether ip matches an allowlisted IP or CIDR.
func (a AllowList) IPAllowed(ip net.IP) bool {
	for _, x := range a.ips {
		if x.Equal(ip) {
			return true
		}
	}
	for _, n := range a.nets {
		if n.Contains(ip) {
			return true
		}
	}
	return false
}

// permitted decides whether a connection to ip (for request host) is allowed:
// any public IP is fine; a private/reserved IP only if the host or the IP is
// on the allowlist.
func (a AllowList) permitted(host string, ip net.IP) bool {
	if !isBlockedIP(ip) {
		return true
	}
	return a.HostAllowed(host) || a.IPAllowed(ip)
}

// isBlockedIP reports whether ip is in a range SSRF protection refuses by
// default: loopback, RFC1918/ULA private, link-local (incl. the cloud metadata
// address 169.254.169.254), multicast, unspecified, and CGNAT 100.64.0.0/10.
func isBlockedIP(ip net.IP) bool {
	if ip == nil {
		return true
	}
	if ip.IsLoopback() || ip.IsPrivate() || ip.IsUnspecified() ||
		ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() || ip.IsMulticast() {
		return true
	}
	if v4 := ip.To4(); v4 != nil && v4[0] == 100 && v4[1] >= 64 && v4[1] <= 127 {
		return true // 100.64.0.0/10 carrier-grade NAT
	}
	return false
}
