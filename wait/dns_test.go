package wait

import (
	"context"
	. "github.com/onsi/gomega"
	"testing"
	"time"
)

func TestDnsWaiter_Wait(t *testing.T) {
	g := NewGomegaWithT(t)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	waiter := DnsWaiter{}
	err := waiter.Wait(ctx, "localhost")
	g.Expect(err).ToNot(HaveOccurred())
}
