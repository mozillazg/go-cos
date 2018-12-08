package cos

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"net/url"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"
)

var (
	// mux is the HTTP request multiplexer used with the test server.
	mux *http.ServeMux

	// client is the COS client being tested.
	client *Client

	// server is a test HTTP server used to provide mock API responses.
	server *httptest.Server
)

// setup sets up a test HTTP server along with a cos.Client that is
// configured to talk to that test server. Tests should register handlers on
// mux which provide mock responses for the API method being tested.
func setup() {
	// test server
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	u, _ := url.Parse(server.URL)
	client = NewClient(&BaseURL{u, u}, nil)
}

// teardown closes the test HTTP server.
func teardown() {
	server.Close()
}

type values map[string]string

func testFormValues(t *testing.T, r *http.Request, values values) {
	want := url.Values{}
	for k, v := range values {
		want.Set(k, v)
	}

	r.ParseForm()
	if got := r.Form; !reflect.DeepEqual(got, want) {
		t.Errorf("Request parameters: %v, want %v", got, want)
	}
}

func testMethod(t *testing.T, r *http.Request, want string) {
	if got := r.Method; got != want {
		t.Errorf("Request method: %v, want %v", got, want)
	}
}

func testHeader(t *testing.T, r *http.Request, header string, want string) {
	if got := r.Header.Get(header); got != want {
		t.Errorf("Header.Get(%q) returned %q, want %q", header, got, want)
	}
}

func testURLParseError(t *testing.T, err error) {
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
	if err, ok := err.(*url.Error); !ok || err.Op != "parse" {
		t.Errorf("Expected URL parse error, got %+v", err)
	}
}

func testBody(t *testing.T, r *http.Request, want string) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Errorf("Error reading request body: %v", err)
	}
	if got := string(b); got != want {
		t.Errorf("request Body is %s, want %s", got, want)
	}
}

// Helper function to test that a value is marshalled to XML as expected.
func testXMLMarshal(t *testing.T, v interface{}, want string) {
	j, err := xml.Marshal(v)
	if err != nil {
		t.Errorf("Unable to marshal JSON for %v", v)
	}

	w := new(bytes.Buffer)
	err = xml.NewEncoder(w).Encode([]byte(want))
	if err != nil {
		t.Errorf("String is not valid json: %s", want)
	}

	if w.String() != string(j) {
		t.Errorf("xml.Marshal(%q) returned %s, want %s", v, j, w)
	}

	// now go the other direction and make sure things unmarshal as expected
	u := reflect.ValueOf(v).Interface()
	if err := xml.Unmarshal([]byte(want), u); err != nil {
		t.Errorf("Unable to unmarshal XML for %v", want)
	}

	if !reflect.DeepEqual(v, u) {
		t.Errorf("xml.Unmarshal(%q) returned %s, want %s", want, u, v)
	}
}

func TestNewClient(t *testing.T) {
	c := NewClient(nil, nil)

	if got, want := c.BaseURL.ServiceURL.String(), defaultServiceBaseURL; got != want {
		t.Errorf("NewClient BaseURL is %v, want %v", got, want)
	}
	if got, want := c.UserAgent, userAgent; got != want {
		t.Errorf("NewClient UserAgent is %v, want %v", got, want)
	}
}

func TestNewBucketURL_secure_false(t *testing.T) {
	got := NewBucketURL("bname", "idx", "ap-beijing", false).String()
	want := "http://bname-idx.cos.ap-beijing.myqcloud.com"
	if got != want {
		t.Errorf("NewBucketURL is %v, want %v", got, want)
	}
}

func TestNewBucketURL_secure_true(t *testing.T) {
	got := NewBucketURL("bname", "idx", "ap-beijing", true).String()
	want := "https://bname-idx.cos.ap-beijing.myqcloud.com"
	if got != want {
		t.Errorf("NewBucketURL is %v, want %v", got, want)
	}
}

func TestNewBaseURL(t *testing.T) {
	bu := "https://test-1253846586.cos.ap-beijing.myqcloud.com"
	got, _ := NewBaseURL(bu)
	if got.BucketURL.String() != bu {
		t.Errorf("bucketURL want %s, but got %s", bu, got.BucketURL.String())
	}
	if got.ServiceURL.String() != defaultServiceBaseURL {
		t.Errorf("serviceURL want %s, but got %s", defaultServiceBaseURL, got.ServiceURL.String())
	}
}

func TestClient_doAPI(t *testing.T) {
	setup()
	defer teardown()

}

func TestNewAuthTime(t *testing.T) {
	a := NewAuthTime(time.Hour)
	if a.SignStartTime != a.KeyStartTime ||
		a.SignEndTime != a.SignEndTime ||
		a.SignStartTime.Add(time.Hour) != a.SignEndTime {
		t.Errorf("NewAuthTime request got %+v is not valid", a)
	}
}

type traceCloser struct {
	io.Reader
	Called bool
}

func (t traceCloser) Close() error {
	t.Called = true
	return nil
}

func newTraceCloser(r io.Reader) traceCloser {
	return traceCloser{r, false}
}

func Test_doAPI_copy_body(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/test_down", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `test`)
	})

	w := bytes.NewBuffer([]byte{})
	resp, err := client.send(context.TODO(), &sendOptions{
		baseURL: client.BaseURL.ServiceURL,
		uri:     "/test_down",
		method:  "GET",
		result:  w,
	})

	if err != nil {
		t.Errorf("Expected error == nil, got %+v", err)
	}
	b, _ := ioutil.ReadAll(resp.Body)
	if len(b) != 0 || string(w.Bytes()) != "test" {
		t.Errorf(
			"Expected body was copy and close, got %+v, %+v",
			string(b), string(w.Bytes()))
	}
}

func Test_Response_header_method(t *testing.T) {
	setup()
	defer teardown()
	reqID := "NTk0NTRjZjZfNTViMjM1XzlkMV9hZTZh"
	traceID := "OGVmYzZiMmQzYjA2OWNhODk0NTRkMTBiOWVmMDAxODc0OWRkZjk0ZDM1NmI1M2E2MTRlY2MzZDhmNmI5MWI1OTBjYzE2MjAxN2M1MzJiOTdkZjMxMDVlYTZjN2FiMmI0NTk3NWFiNjAyMzdlM2RlMmVmOGNiNWIxYjYwNDFhYmQ="
	objType := "normal"
	storageCls := "STANDARD"
	versionID := "xxx-v1" // ?
	encryption := "AES256"

	mux.HandleFunc("/test_down", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(xCosRequestID, reqID)
		w.Header().Set(xCosTraceID, traceID)
		w.Header().Set(xCosObjectType, objType)
		w.Header().Set(xCosStorageClass, storageCls)
		w.Header().Set(xCosVersionID, versionID)
		w.Header().Set(xCosServerSideEncryption, encryption)
		w.Header().Add("x-cos-meta-1", "1")
		w.Header().Add("x-cos-meta-1", "11")
		w.Header().Add("x-cos-meta-2", "2")
		w.Header().Add("x-cos-meta-2", "22")
		w.Header().Add("x-cos-meta-3", "33")
		fmt.Fprint(w, `test`)
	})

	w := bytes.NewBuffer([]byte{})
	resp, err := client.send(context.TODO(), &sendOptions{
		baseURL: client.BaseURL.ServiceURL,
		uri:     "/test_down",
		method:  "GET",
		result:  w,
	})

	if err != nil {
		t.Errorf("Expected error == nil, got %+v", err)
	}
	b, _ := ioutil.ReadAll(resp.Body)
	if len(b) != 0 || string(w.Bytes()) != "test" {
		t.Errorf(
			"Expected body was copy and close, got %+v, %+v",
			string(b), string(w.Bytes()))
	}
	h := resp.MetaHeaders()
	keys := []string{}
	for k := range h {
		keys = append(keys, strings.ToLower(k))
	}
	sort.Strings(keys)
	if resp.RequestID() != reqID ||
		resp.TraceID() != traceID ||
		resp.ObjectType() != objType ||
		resp.StorageClass() != storageCls ||
		resp.VersionID() != versionID ||
		resp.ServerSideEncryption() != encryption ||
		!reflect.DeepEqual(keys,
			[]string{"x-cos-meta-1", "x-cos-meta-2", "x-cos-meta-3"}) {
		t.Errorf("result of response header method is not expected")
	}
	v1 := h[textproto.CanonicalMIMEHeaderKey("x-cos-meta-1")]
	sort.Strings(v1)
	v2 := h[textproto.CanonicalMIMEHeaderKey("x-cos-meta-2")]
	sort.Strings(v2)
	v3 := h[textproto.CanonicalMIMEHeaderKey("x-cos-meta-3")]
	sort.Strings(v3)
	if !reflect.DeepEqual(v1,
		[]string{"1", "11"}) ||
		!reflect.DeepEqual(v2,
			[]string{"2", "22"}) ||
		!reflect.DeepEqual(v3,
			[]string{"33"}) {
		t.Errorf("result of response meta headers is not expected, %s, %s, %s",
			v1, v2, v3)
	}
}
