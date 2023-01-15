package wait

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	. "github.com/onsi/gomega"
)

func TestTcpWaiter_Wait(t *testing.T) {
	g := NewGomegaWithT(t)

	// mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		// noop
	}))
	defer server.Close()

	u, err := url.Parse(server.URL)
	g.Expect(err).ToNot(HaveOccurred())

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	waiter := TcpWaiter{}
	err = waiter.Wait(ctx, u.Host, nil)
	g.Expect(err).ToNot(HaveOccurred())
}

func TestTcpWaiter_Wait_timeout(t *testing.T) {
	g := NewGomegaWithT(t)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	waiter := TcpWaiter{}
	err := waiter.Wait(ctx, "localhost:65000", nil)
	g.Expect(err).To(HaveOccurred())
}
