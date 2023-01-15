package wait

import (
	"context"
	"net"

	"github.com/avast/retry-go/v4"
)

type DnsWaiter struct{}

var _ NetWaiter = DnsWaiter{}

func (h DnsWaiter) Wait(ctx context.Context, hostname string, retryOptions []retry.Option) error {
	return retryCheck(ctx, func() error {
		return checkDns(ctx, hostname)
	}, retryOptions)
}

func checkDns(ctx context.Context, hostname string) error {
	resolver := net.Resolver{}
	_, err := resolver.LookupHost(ctx, hostname)
	if err != nil {
		return err
	}
	return nil
}
