package wait

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	. "github.com/onsi/gomega"
)

func TestHttpWaiter_Wait(t *testing.T) {
	g := NewGomegaWithT(t)

	// mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		// noop
	}))
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	waiter := HttpWaiter
	err := waiter.Wait(ctx, server.URL, nil)
	g.Expect(err).ToNot(HaveOccurred())
}
