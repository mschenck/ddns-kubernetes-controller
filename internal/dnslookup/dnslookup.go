package dnslookup

import (
	"context"
	"fmt"
	"net"
)

func DnsLookup(ctx context.Context, record, zone string) (ip string, err error) {
	var resolver net.Resolver
	var ns []*net.NS
	var fqdn string
	var ips []string

	// fetch NS records
	ns, err = resolver.LookupNS(ctx, zone)
	if err != nil {
		return
	}

	// Update Resolver lookup host
	resolver.Dial = func(ctx context.Context, network, address string) (net.Conn, error) {
		var dialer net.Dialer
		return dialer.DialContext(ctx, network, fmt.Sprintf("%s:53", ns[0].Host))
	}

	// Lookup record (against zone NS records)
	fqdn = fmt.Sprintf("%s.%s", record, zone)
	ips, err = resolver.LookupHost(ctx, fqdn)
	if err != nil {
		return
	}
	ip = ips[0]

	return
}
