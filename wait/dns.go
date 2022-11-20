package wait

import (
	"context"
	"net"
)

type DnsWaiter struct{}

var _ NetWaiter = DnsWaiter{}

func (h DnsWaiter) Wait(ctx context.Context, hostname string) error {
	return retryCheck(ctx, func() error {
		return checkDns(ctx, hostname)
	})
}

func checkDns(ctx context.Context, hostname string) error {
	resolver := net.Resolver{}
	_, err := resolver.LookupHost(ctx, hostname)
	if err != nil {
		return err
	}
	return nil
}
