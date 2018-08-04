package cos

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestNewAuthorization(t *testing.T) {
	expectAuthorization := `q-sign-algorithm=sha1&q-ak=QmFzZTY0IGlzIGEgZ2VuZXJp&q-sign-time=1480932292;1481012292&q-key-time=1480932292;1481012292&q-header-list=host;x-cos-content-sha1;x-cos-stroage-class&q-url-param-list=&q-signature=91f7814df035319aa08d47e5a7a66ea989d57301`
	secretID := "QmFzZTY0IGlzIGEgZ2VuZXJp"
	secretKey := "AKIDZfbOA78asKUYBcXFrJD0a1ICvR98JM"
	host := "testbucket-125000000.cos.ap-beijing-1.myqcloud.com"
	uri := "https://testbucket-125000000.cos.ap-beijing-1.myqcloud.com/testfile2"
	startTime := time.Unix(int64(1480932292), 0)
	endTime := time.Unix(int64(1481012292), 0)

	req, _ := http.NewRequest("PUT", uri, nil)
	req.Header.Add("Host", host)
	req.Header.Add("x-cos-content-sha1", "db8ac1c259eb89d4a131b253bacfca5f319d54f2")
	req.Header.Add("x-cos-stroage-class", "nearline")

	authTime := &AuthTime{
		SignStartTime: startTime,
		SignEndTime:   endTime,
		KeyStartTime:  startTime,
		KeyEndTime:    endTime,
	}
	auth := newAuthorization(Auth{
		SecretID:  secretID,
		SecretKey: secretKey,
	}, req, *authTime)

	if auth != expectAuthorization {
		t.Errorf("NewAuthorization returned \n%#v, want \n%#v", auth, expectAuthorization)
	}
}

func TestAuthorizationTransport(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" {
			t.Error("AuthorizationTransport didn't add Authorization header")
		}
	})

	client.client.Transport = &AuthorizationTransport{}
	req, _ := http.NewRequest("GET", client.BaseURL.BucketURL.String(), nil)
	req.Header.Set("X-Testing", "0")
	client.doAPI(context.Background(), req, nil, true)
}

func TestAuthorizationTransport_skip_PresignedURL(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, exist := r.Header["Authorization"]
		if exist {
			t.Error("AuthorizationTransport add Authorization header when use PresignedURL")
		}
	})

	client.client.Transport = &AuthorizationTransport{}
	sign := "q-sign-algorithm=sha1&q-ak=QmFzZTY0IGlzIGEgZ2VuZXJp&q-sign-time=1480932292;1481012292&q-key-time=1480932292;1481012292&q-header-list=&q-url-param-list=&q-signature=a5de76b0734f084a7ea24413f7168b4bdbe5676c"
	u := fmt.Sprintf("%s?sign=%s", client.BaseURL.BucketURL.String(), sign)
	req, _ := http.NewRequest("GET", u, nil)
	client.doAPI(context.Background(), req, nil, true)
}

func TestAuthorizationTransport_with_another_transport(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" {
			t.Error("AuthorizationTransport didn't add Authorization header")
		}
	})

	tr := &testingTransport{}
	client.client.Transport = &AuthorizationTransport{
		Transport: tr,
	}
	req, _ := http.NewRequest("GET", client.BaseURL.BucketURL.String(), nil)
	req.Header.Set("X-Testing", "0")
	client.doAPI(context.Background(), req, nil, true)
	if tr.called != 1 {
		t.Error("AuthorizationTransport not call another Transport")
	}
}

type testingTransport struct {
	called int
}

func (t *testingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.called++
	return http.DefaultTransport.RoundTrip(req)
}
