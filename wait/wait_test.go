package wait

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	. "github.com/onsi/gomega"
)

func Test_getWaiterForResource(t *testing.T) {
	type args struct {
		resource string
	}
	tests := []struct {
		name    string
		args    args
		want    NetWaiter
		wantErr bool
	}{

		{
			name: "HTTP resource, full URL",
			args: args{resource: "https://service.fake:443/a/b/c?d=e"},
			want: HttpWaiter,
		},
		{
			name: "HTTP resource, short URL",
			args: args{resource: "http://service.fake"},
			want: HttpWaiter,
		},
		{
			name: "TCP resource",
			args: args{resource: "service.fake:123"},
			want: TcpWaiter,
		},
		{
			name: "TCP resource",
			args: args{resource: "127.0.0.1:123"},
			want: TcpWaiter,
		},
		{
			name: "DNS resource",
			args: args{resource: "service.fake"},
			want: DnsWaiter,
		},
		{
			name:    "invalid resource format 1",
			args:    args{resource: "foo/bar"},
			wantErr: true,
		},
		{
			name:    "invalid resource format 2",
			args:    args{resource: "foo:"},
			wantErr: true,
		},
		{
			name:    "invalid resource format 3",
			args:    args{resource: ":123"},
			wantErr: true,
		},
		{
			name:    "invalid resource format 4",
			args:    args{resource: "ssh://service.fake"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getWaiterForResource(tt.args.resource)
			if (err != nil) != tt.wantErr {
				t.Errorf("getWaiterForResource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getWaiterForResource() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRetryWaiter_Wait_success(t *testing.T) {
	g := NewGomegaWithT(t)

	cfg := Config{}
	rw := RetryWaiter{
		Check: func(ctx context.Context, resource string) error {
			return nil
		},
	}

	err := rw.Wait(context.Background(), "ignore", cfg)
	g.Expect(err).NotTo(HaveOccurred())
}

func TestRetryWaiter_Wait_errorThenSuccess(t *testing.T) {
	g := NewGomegaWithT(t)

	cfg := Config{
		RetryDelay: pDuration(0),
	}
	attempt := 1
	errCheckFailed := errors.New("check failed")
	rw := RetryWaiter{
		Check: func(ctx context.Context, resource string) error {
			if attempt < 1 {
				// fail, try again
				attempt++
				return errCheckFailed
			} else {
				// success
				return nil
			}
		},
	}

	err := rw.Wait(context.Background(), "ignore", cfg)
	g.Expect(err).ToNot(HaveOccurred())
}

func TestRetryWaiter_Wait_errorAttemptsExceeded(t *testing.T) {
	g := NewGomegaWithT(t)

	cfg := Config{
		Timeout:    5 * time.Second,
		Attempts:   10,
		RetryDelay: pDuration(0),
	}
	errCheckFailed := errors.New("check failed")
	rw := RetryWaiter{
		Check: func(ctx context.Context, resource string) error {
			return errCheckFailed
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	err := rw.Wait(ctx, "ignore", cfg)
	g.Expect(err).To(HaveOccurred())
	g.Expect(err.Error()).To(ContainSubstring("All attempts fail"))
}

func TestRetryWaiter_Wait_errorTimeout(t *testing.T) {
	g := NewGomegaWithT(t)

	cfg := Config{
		Timeout: 5 * time.Second,
	}
	errCheckFailed := errors.New("check failed")
	rw := RetryWaiter{
		Check: func(ctx context.Context, resource string) error {
			return errCheckFailed
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	err := rw.Wait(ctx, "ignore", cfg)
	g.Expect(err).To(HaveOccurred())
	g.Expect(err.Error()).To(ContainSubstring(errCheckFailed.Error()))
}
