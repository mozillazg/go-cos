package cos

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func Test_checkResponse_error(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/test_409", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusConflict)
		fmt.Fprint(w, `<?xml version='1.0' encoding='utf-8' ?>
<Error>
	<Code>BucketAlreadyExists</Code>
	<Message>The requested bucket name is not available.</Message>
	<Resource>testdelete-1253846586.cn-north.myqcloud.com</Resource>
	<RequestId>NTk0NTRjZjZfNTViMjM1XzlkMV9hZTZh</RequestId>
	<TraceId>OGVmYzZiMmQzYjA2OWNhODk0NTRkMTBiOWVmMDAxODc0OWRkZjk0ZDM1NmI1M2E2MTRlY2MzZDhmNmI5MWI1OTBjYzE2MjAxN2M1MzJiOTdkZjMxMDVlYTZjN2FiMmI0NTk3NWFiNjAyMzdlM2RlMmVmOGNiNWIxYjYwNDFhYmQ=</TraceId>
</Error>`)
	})

	_, err := client.send(context.TODO(), &sendOptions{
		baseURL: client.BaseURL.ServiceURL,
		uri:     "/test_409",
		method:  "GET",
	})

	if e, ok := err.(*ErrorResponse); ok {
		if e.Error() == "" {
			t.Errorf("Expected e.Error() not empty, got %+v", e.Error())
		}
		if e.Code != "BucketAlreadyExists" {
			t.Errorf("Expected BucketAlreadyExists error, got %+v", e.Code)
		}
	} else {
		t.Errorf("Expected ErrorResponse error, got %+v", err)
	}
}

func Test_checkResponse_header(t *testing.T) {
	setup()
	defer teardown()
	reqID := "NTk0NTRjZjZfNTViMjM1XzlkMV9hZTZh"
	traceID := "OGVmYzZiMmQzYjA2OWNhODk0NTRkMTBiOWVmMDAxODc0OWRkZjk0ZDM1NmI1M2E2MTRlY2MzZDhmNmI5MWI1OTBjYzE2MjAxN2M1MzJiOTdkZjMxMDVlYTZjN2FiMmI0NTk3NWFiNjAyMzdlM2RlMmVmOGNiNWIxYjYwNDFhYmQ="

	mux.HandleFunc("/test_409", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(xCosRequestID, reqID)
		w.Header().Set(xCosTraceID, traceID)
		w.WriteHeader(http.StatusConflict)
		fmt.Fprint(w, `<?xml version='1.0' encoding='utf-8' ?>
<Error>
	<Code>BucketAlreadyExists</Code>
	<Message>The requested bucket name is not available.</Message>
	<Resource>testdelete-1253846586.cn-north.myqcloud.com</Resource>
</Error>`)
	})

	_, err := client.send(context.TODO(), &sendOptions{
		baseURL: client.BaseURL.ServiceURL,
		uri:     "/test_409",
		method:  "GET",
	})

	if e, ok := err.(*ErrorResponse); ok {
		if e.Error() == "" {
			t.Errorf("Expected e.Error() not empty, got %+v", e.Error())
		}
		if e.Code != "BucketAlreadyExists" {
			t.Errorf("Expected BucketAlreadyExists error, got %+v", e.Code)
		}
		if e.RequestID != reqID {
			t.Errorf("Expected use header field when RequestId is missing, got %+v", e.RequestID)
		}
		if e.TraceID != traceID {
			t.Errorf("Expected use header field when TraceId is missing, got %+v", e.TraceID)
		}
	} else {
		t.Errorf("Expected ErrorResponse error, got %+v", err)
	}
}

func Test_checkResponse_no_error(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/test_200", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `test`)
	})

	_, err := client.send(context.TODO(), &sendOptions{
		baseURL: client.BaseURL.ServiceURL,
		uri:     "/test_200",
		method:  "GET",
	})
	if err != nil {
		t.Errorf("Expected error == nil, got %+v", err)
	}
}

func Test_error_network_error(t *testing.T) {
	setup()
	defer teardown()

	_, err := client.send(context.TODO(), &sendOptions{
		baseURL: &url.URL{Scheme: "http", Host: "127.0.0.1:0"},
		uri:     "/233",
		method:  "GET",
	})
	if !(strings.Contains(err.Error(), "can't assign requested address") ||
		strings.Contains(err.Error(), "connection refused")) {
		t.Errorf(
			`Expected error contains "can't assign requested address" or "connection refused",
			got %+v`, err)
	}
}

func Test_error_cancel_error(t *testing.T) {
	setup()
	defer teardown()

	ctx := context.TODO()
	ctx, cancel := context.WithCancel(ctx)
	cancel()
	_, err := client.send(ctx, &sendOptions{
		baseURL: &url.URL{Scheme: "http", Host: "127.0.0.1:0"},
		uri:     "/233",
		method:  "GET",
	})
	if !strings.Contains(err.Error(), "context canceled") {
		t.Errorf(`Expected error contains "context canceled", got %+v`, err)
	}
}
