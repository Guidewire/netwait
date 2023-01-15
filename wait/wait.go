package wait

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/url"
	"strings"
	"time"

	"github.com/avast/retry-go/v4"
	"golang.org/x/sync/errgroup"
)

type NetWaiter interface {
	Wait(ctx context.Context, resource string, retryOptions []retry.Option) error
}

type CompositeMultiWaiter struct{}

var _ NetWaiter = CompositeMultiWaiter{}

func (c CompositeMultiWaiter) Wait(ctx context.Context, resource string, retryOptions []retry.Option) error {
	// look up waiter for resource
	delegate, err := getWaiterForResource(resource)
	if err != nil {
		return err
	}

	delegate = LogWaiterDecorator{delegate: delegate}

	// run wait on delegate
	return delegate.Wait(ctx, resource, retryOptions)
}

func (c CompositeMultiWaiter) WaitMulti(resources []string, timeout time.Duration, retryOptions []retry.Option) error {

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	errs, ctx := errgroup.WithContext(ctx)
	for _, resource := range resources {
		r := resource
		errs.Go(func() error {
			return c.Wait(ctx, r, retryOptions)
		})
	}
	return errs.Wait()
}

// getWaiterForResource attempts to resolve waiter based on format of resource
func getWaiterForResource(resource string) (NetWaiter, error) {
	// Check if resource is a valid URL
	// See: https://stackoverflow.com/questions/31480710/validate-url-with-standard-package-in-go
	u, err := url.ParseRequestURI(resource)
	if err == nil && u.Scheme != "" && u.Host != "" {
		if u.Scheme == "http" || u.Scheme == "https" {
			return HttpWaiter{}, nil
		} else {
			return nil, fmt.Errorf("invalid format: URL scheme must be http(s): %s", resource)
		}
	}

	// non-URL resource must not contain /
	if strings.Contains(resource, "/") {
		return nil, fmt.Errorf("invalid format: non-URL cannot contain '/': %s", resource)
	}

	host, port, err := net.SplitHostPort(resource)
	if err != nil {
		// If parse error "missing port in address" returned, assume DNS
		var addrError *net.AddrError
		if errors.As(err, &addrError) && addrError.Err == "missing port in address" {
			return DnsWaiter{}, nil
		}
	} else if host != "" && port != "" {
		// if parser returned host and port, assume TCP
		return TcpWaiter{}, nil
	}

	return nil, fmt.Errorf("invalid format: %s", resource)
}

// LogWaiterDecorator wraps a NetWaiter and adds logging around Wait()
type LogWaiterDecorator struct {
	delegate NetWaiter
}

var _ NetWaiter = LogWaiterDecorator{}

func (d LogWaiterDecorator) Wait(ctx context.Context, resource string, retryOptions []retry.Option) error {
	err := d.delegate.Wait(ctx, resource, retryOptions)
	if err == nil {
		Println("available:", resource)
	} else {
		Println("unavailable:", resource)
	}
	return err
}

// retryCheck retries a check until the context deadline expires
// func retryCheck(ctx context.Context, check func() error) error {
func retryCheck(ctx context.Context, check func() error, retryOptions []retry.Option) error {
	retryOptions = append(retryOptions, retry.Context(ctx))
	retryOptions = append(retryOptions, retry.Delay(2*time.Second))

	return retry.Do(func() error {
		return check()
	}, retryOptions...)
}
