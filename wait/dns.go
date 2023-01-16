package wait

import (
	"context"
	"net"
)

var DnsWaiter = &RetryWaiter{
	Check: checkDns,
}

func checkDns(ctx context.Context, hostname string) error {
	resolver := net.Resolver{}
	_, err := resolver.LookupHost(ctx, hostname)
	if err != nil {
		return err
	}
	return nil
}
