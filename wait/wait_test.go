package wait

import (
	"reflect"
	"testing"
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
			want: HttpWaiter{},
		},
		{
			name: "HTTP resource, short URL",
			args: args{resource: "https://service.fake"},
			want: HttpWaiter{},
		},
		{
			name: "TCP resource",
			args: args{resource: "service.fake:123"},
			want: TcpWaiter{},
		},
		{
			name: "TCP resource",
			args: args{resource: "127.0.0.1:123"},
			want: TcpWaiter{},
		},
		{
			name: "DNS resource",
			args: args{resource: "service.fake"},
			want: DnsWaiter{},
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
