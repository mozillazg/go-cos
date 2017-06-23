package cos

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"
)

func TestDebugRequestTransport(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("X-Test-Response", "2333")
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte("test response body"))
	})

	w := bytes.NewBufferString("")

	client.client.Transport = &DebugRequestTransport{
		RequestHeader:  true,
		RequestBody:    true,
		ResponseHeader: true,
		ResponseBody:   true,
		Writer:         w,
	}

	body := bytes.NewReader([]byte("test_request body"))
	req, _ := http.NewRequest("GET", client.BaseURL.BucketURL.String(), body)
	req.Header.Add("X-Test-Debug", "123")
	client.doAPI(context.Background(), req, nil, true)

	b := make([]byte, 800)
	w.Read(b)
	info := string(b)
	if !strings.Contains(info, "GET / HTTP/1.1\r\n") ||
		!strings.Contains(info, "X-Test-Debug: 123\r\n") {
		t.Errorf("DebugRequestTransport debug info %#v don't contains request header", info)
	}
	if !strings.Contains(info, "\r\n\r\ntest_request body") {
		t.Errorf("DebugRequestTransport debug info  %#v don't contains request body", info)
	}

	if !strings.Contains(info, "HTTP/1.1 502 Bad Gateway\r\n") ||
		!strings.Contains(info, "X-Test-Response: 2333\r\n") {
		t.Errorf("DebugRequestTransport debug info  %#v don't contains response header", info)
	}

	if !strings.Contains(info, "\r\n\r\ntest response body") {
		t.Errorf("DebugRequestTransport debug info  %#v don't contains response body", info)
	}
}

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
