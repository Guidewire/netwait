package wait

import (
	"context"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
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

	waiter := HttpWaiter{}
	err := waiter.Wait(ctx, server.URL)
	g.Expect(err).ToNot(HaveOccurred())
}
