package cos

import (
	"fmt"
	"testing"
)

func Test_calSHA1Digest(t *testing.T) {
	want := "a94a8fe5ccb19ba61c4c0873d391e987982fbbd3"
	got := fmt.Sprintf("%x", calSHA1Digest([]byte("test")))
	if got != want {

		t.Errorf("calSHA1Digest request sha1: %+v, want %+v", got, want)
	}
}

func Test_calMD5Digest(t *testing.T) {
	want := "098f6bcd4621d373cade4e832627b4f6"
	got := fmt.Sprintf("%x", calMD5Digest([]byte("test")))
	if got != want {

		t.Errorf("calMD5Digest request md5: %+v, want %+v", got, want)
	}
}

func Test_encodeURIComponent(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "empty",
			args: args{""},
			want: "",
		},
		{
			name: "no escape",
			args: args{"0123456789abcdefghijkhlmnopqrstuvwxyzABCDEFGHIJKHLMNOPQRSTUVWXYZ-_.!~*'()"},
			want: "0123456789abcdefghijkhlmnopqrstuvwxyzABCDEFGHIJKHLMNOPQRSTUVWXYZ-_.!~*'()",
		},
		{
			name: "escape",
			args: args{"+ $@#/"},
			want: "%2B%20%24%40%23%2F",
		},
		{
			name: "escape+no",
			args: args{"+ $abc@#13/0"},
			want: "%2B%20%24abc%40%2313%2F0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := encodeURIComponent(tt.args.s); got != tt.want {
				t.Errorf("encodeURIComponent() = %v, want %v", got, tt.want)
			}
		})
	}
}
