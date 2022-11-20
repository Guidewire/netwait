package wait

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
)

type HttpWaiter struct{}

var _ NetWaiter = HttpWaiter{}

func (h HttpWaiter) Wait(ctx context.Context, url string) error {
	return retryCheck(ctx, func() error {
		return checkHttp(ctx, url)
	})
}

func checkHttp(ctx context.Context, url string) error {
	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	} else {
		return fmt.Errorf("GET '%s': returned status code %d", url, resp.StatusCode)
	}
}
