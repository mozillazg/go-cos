package cos

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
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

	(client.Sender).(*DefaultSender).Transport = &AuthorizationTransport{}
	req, _ := http.NewRequest("GET", client.BaseURL.BucketURL.String(), nil)
	req.Header.Set("X-Testing", "0")
	client.doAPI(context.Background(), Caller{}, req, nil, true)
}

func TestAuthorizationTransportWithSessionToken(t *testing.T) {
	setup()
	defer teardown()

	sessionToken := "CxQQbwSzzX5obZm23yEcyQtpROuDB0Q60d322a47737c8241991d12dc4b8387c7J6NL50eH1BYN6VnFYB_Ml6oPZzUxz5wxDGVvvgxZXr1m-4HvmkvmMH4YB02XdVPapKp7oGnrMous2jsSTALo4iU2fuRclbVw-czYwggSxuNxXAwmqcT1HpD3h3zc3e24sryIhJKqzSOczQZjtGrxSSQ4K23o9Mx8VHgrosliU0aIiI2KFhxJhij03SzDDOQcBAwpFZyM0NvpOdN6b14yJbrt9bAzYGNjX-PeU3MXfi0"

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" {
			t.Error("AuthorizationTransport didn't add Authorization header")
		}
		token := r.Header.Get("x-cos-security-token")
		if token == "" {
			t.Error("AuthorizationTransport didn't add x-cos-security-token header")
		}
		if token != sessionToken {
			t.Errorf("AuthorizationTransport didn't add expected x-cos-security-token header, expected: %s, got: %s", sessionToken, token)
		}
	})

	(client.Sender).(*DefaultSender).Transport = &AuthorizationTransport{
		SecretID:     "233",
		SecretKey:    "666",
		SessionToken: sessionToken,
	}
	req, _ := http.NewRequest("GET", client.BaseURL.BucketURL.String(), nil)
	req.Header.Set("X-Testing", "0")
	client.doAPI(context.Background(), Caller{}, req, nil, true)
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

	(client.Sender).(*DefaultSender).Transport = &AuthorizationTransport{}
	sign := "q-sign-algorithm=sha1&q-ak=QmFzZTY0IGlzIGEgZ2VuZXJp&q-sign-time=1480932292;1481012292&q-key-time=1480932292;1481012292&q-header-list=&q-url-param-list=&q-signature=a5de76b0734f084a7ea24413f7168b4bdbe5676c"
	u := fmt.Sprintf("%s?sign=%s", client.BaseURL.BucketURL.String(), sign)
	req, _ := http.NewRequest("GET", u, nil)
	client.doAPI(context.Background(), Caller{}, req, nil, true)
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
	(client.Sender).(*DefaultSender).Transport = &AuthorizationTransport{
		Transport: tr,
	}
	req, _ := http.NewRequest("GET", client.BaseURL.BucketURL.String(), nil)
	req.Header.Set("X-Testing", "0")
	client.doAPI(context.Background(), Caller{}, req, nil, true)
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

func Test_camSafeURLEncode(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "no replace",
			args: args{"1234 +abc0AB#@"},
			want: "1234%20%2Babc0AB%23%40",
		},
		{
			name: "replace",
			args: args{"1234 +abc0AB#@,!'()*"},
			want: "1234%20%2Babc0AB%23%40%2C%21%27%28%29%2A",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := camSafeURLEncode(tt.args.s); got != tt.want {
				t.Errorf("camSafeURLEncode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_valuesForSign_Encode(t *testing.T) {
	tests := []struct {
		name string
		vs   valuesForSign
		want string
	}{
		{
			name: "test escape",
			vs: valuesForSign{
				"test+233": {"value 666"},
				"test+234": {"value 667"},
			},
			want: "test%2B233=value%20666&test%2B234=value%20667",
		},
		{
			name: "test order",
			vs: valuesForSign{
				"test_233": {"value_666"},
				"233":      {"value_2"},
				"test_666": {"value_123"},
			},
			want: "233=value_2&test_233=value_666&test_666=value_123",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.vs.Encode(); got != tt.want {
				t.Errorf("valuesForSign.Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_valuesForSign_Add(t *testing.T) {
	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name string
		vs   valuesForSign
		args args
		want valuesForSign
	}{
		{
			name: "add new key",
			vs:   valuesForSign{},
			args: args{"test_key", "value_233"},
			want: valuesForSign{"test_key": {"value_233"}},
		},
		{
			name: "extend key",
			vs:   valuesForSign{"test_key": {"value_233"}},
			args: args{"test_key", "value_666"},
			want: valuesForSign{"test_key": {"value_233", "value_666"}},
		},
		{
			name: "key to lower(add)",
			vs:   valuesForSign{},
			args: args{"TEST_KEY", "value_233"},
			want: valuesForSign{"test_key": {"value_233"}},
		},
		{
			name: "key to lower(extend)",
			vs:   valuesForSign{"test_key": {"value_233"}},
			args: args{"TEST_KEY", "value_666"},
			want: valuesForSign{"test_key": {"value_233", "value_666"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.vs.Add(tt.args.key, tt.args.value)
			if !reflect.DeepEqual(tt.vs, tt.want) {
				t.Errorf("%v, want %v", tt.vs, tt.want)
			}
		})
	}
}

func Test_genFormatParameters(t *testing.T) {
	type args struct {
		parameters url.Values
	}
	tests := []struct {
		name                    string
		args                    args
		wantFormatParameters    string
		wantSignedParameterList []string
	}{
		{
			name: "test order",
			args: args{url.Values{
				"test_key_233": {"666"},
				"233":          {"222"},
				"test_key_2":   {"value"},
			}},
			wantFormatParameters:    "233=222&test_key_2=value&test_key_233=666",
			wantSignedParameterList: []string{"233", "test_key_2", "test_key_233"},
		},
		{
			name: "test escape",
			args: args{url.Values{
				"Test+key": {"666 value"},
				"233 666":  {"22+2"},
			}},
			wantFormatParameters:    "233%20666=22%2B2&test%2Bkey=666%20value",
			wantSignedParameterList: []string{"233 666", "test+key"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFormatParameters, gotSignedParameterList := genFormatParameters(tt.args.parameters)
			if gotFormatParameters != tt.wantFormatParameters {
				t.Errorf("genFormatParameters() gotFormatParameters = %v, want %v", gotFormatParameters, tt.wantFormatParameters)
			}
			if !reflect.DeepEqual(gotSignedParameterList, tt.wantSignedParameterList) {
				t.Errorf("genFormatParameters() gotSignedParameterList = %v, want %v", gotSignedParameterList, tt.wantSignedParameterList)
			}
		})
	}
}

func Test_genFormatHeaders(t *testing.T) {
	type args struct {
		headers http.Header
	}
	tests := []struct {
		name                 string
		args                 args
		wantFormatHeaders    string
		wantSignedHeaderList []string
	}{
		{
			name: "test order",
			args: args{http.Header{
				"host":           {"example.com"},
				"content-length": {"22"},
				"content-md5":    {"xxx222"},
			}},
			wantFormatHeaders:    "content-length=22&content-md5=xxx222&host=example.com",
			wantSignedHeaderList: []string{"content-length", "content-md5", "host"},
		},
		{
			name: "test escape",
			args: args{http.Header{
				"host":                {"example.com"},
				"content-length":      {"22"},
				"Content-Disposition": {"attachment; filename=hello - world!(+).go"},
			}},
			wantFormatHeaders:    "content-disposition=attachment%3B%20filename%3Dhello%20-%20world%21%28%2B%29.go&content-length=22&host=example.com",
			wantSignedHeaderList: []string{"content-disposition", "content-length", "host"},
		},
		{
			name: "test skip key",
			args: args{http.Header{
				"Host":           {"example.com"},
				"content-length": {"22"},
				"x-cos-xyz":      {"lala"},
				"Content-Type":   {"text/html"},
			}},
			wantFormatHeaders:    "content-length=22&host=example.com&x-cos-xyz=lala",
			wantSignedHeaderList: []string{"content-length", "host", "x-cos-xyz"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFormatHeaders, gotSignedHeaderList := genFormatHeaders(tt.args.headers)
			if gotFormatHeaders != tt.wantFormatHeaders {
				t.Errorf("genFormatHeaders() gotFormatHeaders = %v, want %v", gotFormatHeaders, tt.wantFormatHeaders)
			}
			if !reflect.DeepEqual(gotSignedHeaderList, tt.wantSignedHeaderList) {
				t.Errorf("genFormatHeaders() gotSignedHeaderList = %v, want %v", gotSignedHeaderList, tt.wantSignedHeaderList)
			}
		})
	}
}
