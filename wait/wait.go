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

type Config struct {
	Timeout time.Duration
	// Limits number of attempts. 0 means unlimited.
	Attempts      uint
	RetryDelay    *time.Duration
	RetryMaxDelay *time.Duration
}

func DefaultConfig() Config {
	return Config{
		Timeout:    1 * time.Minute,
		Attempts:   0,
		RetryDelay: pDuration(2 * time.Second),
	}
}

func pDuration(d time.Duration) *time.Duration {
	return &d
}

type NetWaiter interface {
	Wait(ctx context.Context, resource string, config Config) error
}

type CompositeMultiWaiter struct{}

var _ NetWaiter = CompositeMultiWaiter{}

func (c CompositeMultiWaiter) Wait(ctx context.Context, resource string, config Config) error {
	// look up waiter for resource
	delegate, err := getWaiterForResource(resource)
	if err != nil {
		return err
	}

	delegate = LogWaiterDecorator{delegate: delegate}

	// run wait on delegate
	return delegate.Wait(ctx, resource, config)
}

func (c CompositeMultiWaiter) WaitMulti(resources []string, config Config) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	errs, ctx := errgroup.WithContext(ctx)
	for _, resource := range resources {
		r := resource
		errs.Go(func() error {
			return c.Wait(ctx, r, config)
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
			return HttpWaiter, nil
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
			return DnsWaiter, nil
		}
	} else if host != "" && port != "" {
		// if parser returned host and port, assume TCP
		return TcpWaiter, nil
	}

	return nil, fmt.Errorf("invalid format: %s", resource)
}

// LogWaiterDecorator wraps a NetWaiter and adds logging around Wait()
type LogWaiterDecorator struct {
	delegate NetWaiter
}

var _ NetWaiter = LogWaiterDecorator{}

func (d LogWaiterDecorator) Wait(ctx context.Context, resource string, config Config) error {
	err := d.delegate.Wait(ctx, resource, config)
	if err == nil {
		Println("available:", resource)
	} else {
		Println("unavailable:", resource)
	}
	return err
}

// RetryWaiter is a generic NetWaiter that retries a check on error until the context deadline expires
type RetryWaiter struct {
	Check func(ctx context.Context, resource string) error
}

var _ NetWaiter = RetryWaiter{}

func (w RetryWaiter) Wait(ctx context.Context, resource string, config Config) error {
	retryOptions := []retry.Option{}
	retryOptions = append(retryOptions, retry.Context(ctx))

	if config.RetryDelay != nil {
		retryOptions = append(retryOptions, retry.Delay(*config.RetryDelay))
	}
	if config.RetryMaxDelay != nil {
		retryOptions = append(retryOptions, retry.MaxDelay(*config.RetryMaxDelay))
	}

	attempts := config.Attempts
	if attempts == 0 {
		// Due to a bug in retry, setting attempts to 0 causes timeouts to not return error
		// See: https://github.com/avast/retry-go/issues/83
		attempts = 99999999
	}
	retryOptions = append(retryOptions, retry.Attempts(attempts))

	retryOptions = append(retryOptions, retry.OnRetry(func(n uint, err error) {
		fmt.Println("retrying")
	}))

	return retry.Do(func() error {
		return w.Check(ctx, resource)
	}, retryOptions...)
}
