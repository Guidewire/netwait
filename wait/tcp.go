package wait

import (
	"context"
	"fmt"
	"net"
)

var TcpWaiter = &RetryWaiter{
	Check: checkTcp,
}

func checkTcp(ctx context.Context, address string) error {
	dialer := net.Dialer{}
	conn, err := dialer.DialContext(ctx, "tcp", address)
	if err != nil {
		return err
	}
	if conn == nil {
		return fmt.Errorf("dial tcp %s: nil connection", address)
	} else {
		defer conn.Close()
		return nil
	}
}
