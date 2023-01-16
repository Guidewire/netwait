package wait

import (
	"context"
	"testing"
	"time"

	. "github.com/onsi/gomega"
)

func TestDnsWaiter_Wait(t *testing.T) {
	g := NewGomegaWithT(t)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	waiter := DnsWaiter
	err := waiter.Wait(ctx, "localhost")
	g.Expect(err).ToNot(HaveOccurred())
}
