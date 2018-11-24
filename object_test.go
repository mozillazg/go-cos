package cos

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"
)

func TestObjectService_Get(t *testing.T) {
	setup()
	defer teardown()
	name := "test/hello.txt"

	mux.HandleFunc("/test/hello.txt", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		vs := values{
			"response-content-type": "text/html",
		}
		testFormValues(t, r, vs)
		testHeader(t, r, "Range", "bytes=0-3")
		fmt.Fprint(w, `hello`)
	})

	opt := &ObjectGetOptions{
		ResponseContentType: "text/html",
		Range:               "bytes=0-3",
	}

	resp, err := client.Object.Get(context.Background(), name, opt)
	if err != nil {
		t.Fatalf("Object.Get returned error: %v", err)
	}

	b, _ := ioutil.ReadAll(resp.Body)
	ref := string(b)
	want := "hello"
	if !reflect.DeepEqual(ref, want) {
		t.Errorf("Object.Get returned %+v, want %+v", ref, want)
	}

}

func TestObjectService_Get_with_PresignedURL(t *testing.T) {
	setup()
	defer teardown()
	name := "test/hello.txt"

	mux.HandleFunc("/233/PresignedURL", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		vs := values{
			"response-content-type": "text/html",
		}
		testFormValues(t, r, vs)
		testHeader(t, r, "Range", "bytes=0-3")
		fmt.Fprint(w, `hello`)
	})

	opt := &ObjectGetOptions{
		ResponseContentType: "text/html",
		Range:               "bytes=0-3",
	}
	PresignedURL, _ := url.Parse(
		fmt.Sprintf("%s/%s", client.BaseURL.BucketURL.String(), "233/PresignedURL"))
	opt.PresignedURL = PresignedURL

	resp, err := client.Object.Get(context.Background(), name, opt)
	if err != nil {
		t.Fatalf("Object.Get returned error: %v", err)
	}

	b, _ := ioutil.ReadAll(resp.Body)
	ref := string(b)
	want := "hello"
	if !reflect.DeepEqual(ref, want) {
		t.Errorf("Object.Get returned %+v, want %+v", ref, want)
	}

}

func TestObjectService_Put(t *testing.T) {
	setup()
	defer teardown()

	opt := &ObjectPutOptions{
		ObjectPutHeaderOptions: &ObjectPutHeaderOptions{
			ContentType: "text/html",
		},
		ACLHeaderOptions: &ACLHeaderOptions{
			XCosACL: "private",
		},
	}
	name := "test/hello.txt"

	mux.HandleFunc("/test/hello.txt", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testHeader(t, r, "x-cos-acl", "private")
		testHeader(t, r, "Content-Type", "text/html")

		b, _ := ioutil.ReadAll(r.Body)
		v := string(b)
		want := "hello"
		if !reflect.DeepEqual(v, want) {
			t.Errorf("Object.Put request body: %#v, want %#v", v, want)
		}
	})

	r := bytes.NewReader([]byte("hello"))
	_, err := client.Object.Put(context.Background(), name, r, opt)
	if err != nil {
		t.Fatalf("Object.Put returned error: %v", err)
	}

}

func TestObjectService_Put_with_PresignedURL(t *testing.T) {
	setup()
	defer teardown()

	opt := &ObjectPutOptions{
		ObjectPutHeaderOptions: &ObjectPutHeaderOptions{
			ContentType: "text/html",
		},
		ACLHeaderOptions: &ACLHeaderOptions{
			XCosACL: "private",
		},
	}
	PresignedURL, _ := url.Parse(
		fmt.Sprintf("%s/%s", client.BaseURL.BucketURL.String(), "233/PresignedURL"))
	opt.PresignedURL = PresignedURL

	name := "test/hello.txt"

	mux.HandleFunc("/233/PresignedURL", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testHeader(t, r, "x-cos-acl", "private")
		testHeader(t, r, "Content-Type", "text/html")

		b, _ := ioutil.ReadAll(r.Body)
		v := string(b)
		want := "hello"
		if !reflect.DeepEqual(v, want) {
			t.Errorf("Object.Put request body: %#v, want %#v", v, want)
		}
	})

	r := bytes.NewReader([]byte("hello"))
	_, err := client.Object.Put(context.Background(), name, r, opt)
	if err != nil {
		t.Fatalf("Object.Put returned error: %v", err)
	}

}

func TestObjectService_Put_not_close(t *testing.T) {
	setup()
	defer teardown()

	opt := &ObjectPutOptions{
		ObjectPutHeaderOptions: &ObjectPutHeaderOptions{
			ContentType:   "text/html",
			ContentLength: 5,
		},
		ACLHeaderOptions: &ACLHeaderOptions{
			XCosACL: "private",
		},
	}
	name := "test/hello.txt"

	mux.HandleFunc("/test/hello.txt", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		testHeader(t, r, "x-cos-acl", "private")
		testHeader(t, r, "Content-Type", "text/html")

		b, _ := ioutil.ReadAll(r.Body)
		v := string(b)
		want := "hello"
		if !reflect.DeepEqual(v, want) {
			t.Errorf("Object.Put request body: %#v, want %#v", v, want)
		}
		if r.ContentLength != 5 {
			t.Error("ContentLength should be 5")
		}
	})

	r := newTraceCloser(bytes.NewReader([]byte("hello")))
	_, err := client.Object.Put(context.Background(), name, r, opt)
	if err != nil {
		t.Fatalf("Object.Put returned error: %v", err)
	}
	if r.Called {
		t.Fatal("Should not close input")
	}

}

func TestObjectService_Delete(t *testing.T) {
	setup()
	defer teardown()
	name := "test/hello.txt"

	mux.HandleFunc("/test/hello.txt", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		w.WriteHeader(http.StatusNoContent)
	})

	_, err := client.Object.Delete(context.Background(), name)
	if err != nil {
		t.Fatalf("Object.Delete returned error: %v", err)
	}
}

func TestObjectService_Head(t *testing.T) {
	setup()
	defer teardown()
	name := "test/hello.txt"

	mux.HandleFunc("/test/hello.txt", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "HEAD")
		testHeader(t, r, "If-Modified-Since", "Mon, 12 Jun 2017 05:36:19 GMT")
	})

	opt := &ObjectHeadOptions{
		IfModifiedSince: "Mon, 12 Jun 2017 05:36:19 GMT",
	}

	_, err := client.Object.Head(context.Background(), name, opt)
	if err != nil {
		t.Fatalf("Object.Head returned error: %v", err)
	}

}

func TestObjectService_Options(t *testing.T) {
	setup()
	defer teardown()
	name := "test/hello.txt"

	mux.HandleFunc("/test/hello.txt", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodOptions)
		testHeader(t, r, "Access-Control-Request-Method", "PUT")
		testHeader(t, r, "Origin", "www.qq.com")
	})

	opt := &ObjectOptionsOptions{
		Origin: "www.qq.com",
		AccessControlRequestMethod: "PUT",
	}

	_, err := client.Object.Options(context.Background(), name, opt)
	if err != nil {
		t.Fatalf("Object.Options returned error: %v", err)
	}

}

func TestObjectService_Append(t *testing.T) {
	setup()
	defer teardown()

	opt := &ObjectPutOptions{
		ObjectPutHeaderOptions: &ObjectPutHeaderOptions{
			ContentType: "text/html",
		},
		ACLHeaderOptions: &ACLHeaderOptions{
			XCosACL: "private",
		},
	}
	name := "test/hello.txt"
	position := 0

	mux.HandleFunc("/test/hello.txt", func(w http.ResponseWriter, r *http.Request) {
		vs := values{
			"append":   "",
			"position": "0",
		}
		testFormValues(t, r, vs)

		testMethod(t, r, http.MethodPost)
		testHeader(t, r, "x-cos-acl", "private")
		testHeader(t, r, "Content-Type", "text/html")

		b, _ := ioutil.ReadAll(r.Body)
		v := string(b)
		want := "hello"
		if !reflect.DeepEqual(v, want) {
			t.Errorf("Object.Append request body: %#v, want %#v", v, want)
		}
	})

	r := bytes.NewReader([]byte("hello"))
	_, err := client.Object.Append(context.Background(), name, position, r, opt)
	if err != nil {
		t.Fatalf("Object.Append returned error: %v", err)
	}
}

func TestObjectService_Append_not_close(t *testing.T) {
	setup()
	defer teardown()

	opt := &ObjectPutOptions{
		ObjectPutHeaderOptions: &ObjectPutHeaderOptions{
			ContentType:   "text/html",
			ContentLength: 5,
		},
		ACLHeaderOptions: &ACLHeaderOptions{
			XCosACL: "private",
		},
	}
	name := "test/hello.txt"
	position := 0

	mux.HandleFunc("/test/hello.txt", func(w http.ResponseWriter, r *http.Request) {
		vs := values{
			"append":   "",
			"position": "0",
		}
		testFormValues(t, r, vs)

		testMethod(t, r, http.MethodPost)
		testHeader(t, r, "x-cos-acl", "private")
		testHeader(t, r, "Content-Type", "text/html")

		b, _ := ioutil.ReadAll(r.Body)
		v := string(b)
		want := "hello"
		if !reflect.DeepEqual(v, want) {
			t.Errorf("Object.Append request body: %#v, want %#v", v, want)
		}
		if r.ContentLength != 5 {
			t.Error("ContentLength should be 5")
		}
	})

	r := newTraceCloser(bytes.NewReader([]byte("hello")))
	_, err := client.Object.Append(context.Background(), name, position, r, opt)
	if err != nil {
		t.Fatalf("Object.Append returned error: %v", err)
	}
	if r.Called {
		t.Fatal("Should not close input")
	}

}

func TestObjectService_DeleteMulti(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		vs := values{
			"delete": "",
		}
		testFormValues(t, r, vs)
		fmt.Fprint(w, `<DeleteResult>
	<Deleted>
		<Key>test1</Key>
	</Deleted>
	<Deleted>
		<Key>test3</Key>
	</Deleted>
	<Deleted>
		<Key>test2</Key>
	</Deleted>
</DeleteResult>`)
	})

	opt := &ObjectDeleteMultiOptions{
		Objects: []Object{
			{
				Key: "test1",
			},
			{
				Key: "test3",
			},
			{
				Key: "test2",
			},
		},
	}

	ref, _, err := client.Object.DeleteMulti(context.Background(), opt)
	if err != nil {
		t.Fatalf("Object.DeleteMulti returned error: %v", err)
	}

	want := &ObjectDeleteMultiResult{
		XMLName: xml.Name{Local: "DeleteResult"},
		DeletedObjects: []Object{
			{
				Key: "test1",
			},
			{
				Key: "test3",
			},
			{
				Key: "test2",
			},
		},
	}

	if !reflect.DeepEqual(ref, want) {
		t.Errorf("Object.DeleteMulti returned %+v, want %+v", ref, want)
	}

}

func TestObjectService_Copy(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/test.go.copy", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		fmt.Fprint(w, `<CopyObjectResult>
		<ETag>"098f6bcd4621d373cade4e832627b4f6"</ETag>
		<LastModified>2017-12-13T14:53:12</LastModified>
	</CopyObjectResult>`)
	})

	sourceURL := "test-1253846586.cn-north.myqcloud.com/test.source"
	ref, _, err := client.Object.Copy(context.Background(), "test.go.copy", sourceURL, nil)
	if err != nil {
		t.Fatalf("Object.Copy returned error: %v", err)
	}

	want := &ObjectCopyResult{
		XMLName:      xml.Name{Local: "CopyObjectResult"},
		ETag:         `"098f6bcd4621d373cade4e832627b4f6"`,
		LastModified: "2017-12-13T14:53:12",
	}

	if !reflect.DeepEqual(ref, want) {
		t.Errorf("Object.Copy returned %+v, want %+v", ref, want)
	}
}

func TestObjectService_PresignedURL(t *testing.T) {
	testTable := map[string]string{
		http.MethodGet: `q-sign-algorithm=sha1&q-ak=QmFzZTY0IGlzIGEgZ2VuZXJp&q-sign-time=1480932292;1481012292&q-key-time=1480932292;1481012292&q-header-list=&q-url-param-list=&q-signature=a5de76b0734f084a7ea24413f7168b4bdbe5676c`,
		http.MethodPut: `q-sign-algorithm=sha1&q-ak=QmFzZTY0IGlzIGEgZ2VuZXJp&q-sign-time=1480932292;1481012292&q-key-time=1480932292;1481012292&q-header-list=&q-url-param-list=&q-signature=b13e488d105301cdc627f33448d3d4237f418256`,
	}

	for method, expectAuthorization := range testTable {
		b, _ := NewBaseURL("https://testbucket-125000000.cos.ap-beijing-1.myqcloud.com")
		c := NewClient(b, nil)
		secretID := "QmFzZTY0IGlzIGEgZ2VuZXJp"
		secretKey := "AKIDZfbOA78asKUYBcXFrJD0a1ICvR98JM"
		startTime := time.Unix(int64(1480932292), 0)
		endTime := time.Unix(int64(1481012292), 0)

		auth := Auth{
			SecretID:  secretID,
			SecretKey: secretKey,
		}
		authTime := &AuthTime{
			SignStartTime: startTime,
			SignEndTime:   endTime,
			KeyStartTime:  startTime,
			KeyEndTime:    endTime,
		}
		opt := &objectPresignedURLTestingOptions{
			authTime: authTime,
		}
		signURL, err := c.Object.PresignedURL(context.Background(), method, "testfile2", auth, opt)
		if err != nil {
			t.Errorf("PresignedURL returned error: %v", err)
		}

		sign := signURL.Query().Get("sign")
		if strings.Compare(sign, expectAuthorization) != 0 {
			t.Errorf("PresignedURL %s contain sign %#v, want %#v", method, sign, expectAuthorization)
		}
	}
}

func TestObjectService_PresignedURL_withoutMockAuthTime(t *testing.T) {
	testTable := map[string]string{
		http.MethodGet: `q-sign-algorithm=sha1&q-ak=QmFzZTY0IGlzIGEgZ2VuZXJp&q-sign-time=1480932292;1481012292&q-key-time=1480932292;1481012292&q-header-list=&q-url-param-list=&q-signature=a5de76b0734f084a7ea24413f7168b4bdbe5676c`,
		http.MethodPut: `q-sign-algorithm=sha1&q-ak=QmFzZTY0IGlzIGEgZ2VuZXJp&q-sign-time=1480932292;1481012292&q-key-time=1480932292;1481012292&q-header-list=&q-url-param-list=&q-signature=b13e488d105301cdc627f33448d3d4237f418256`,
	}

	for method, expectAuthorization := range testTable {
		b, _ := NewBaseURL("https://testbucket-125000000.cos.ap-beijing-1.myqcloud.com")
		c := NewClient(b, nil)
		secretID := "QmFzZTY0IGlzIGEgZ2VuZXJp"
		secretKey := "AKIDZfbOA78asKUYBcXFrJD0a1ICvR98JM"

		auth := Auth{
			SecretID:  secretID,
			SecretKey: secretKey,
		}
		signURL, err := c.Object.PresignedURL(context.Background(), method, "testfile2", auth, nil)
		if err != nil {
			t.Errorf("PresignedURL returned error: %v", err)
		}

		sign := signURL.Query().Get("sign")
		sign = strings.Replace(sign, ";", "-", -1)
		var expectedKeys []string
		var gotKeys []string
		expectAuthorization = strings.Replace(expectAuthorization, ";", "-", -1)
		if v, _ := url.ParseQuery(expectAuthorization); v != nil {
			for k := range v {
				expectedKeys = append(expectedKeys, k)
			}
		}
		if v, _ := url.ParseQuery(sign); v != nil {
			for k := range v {
				gotKeys = append(gotKeys, k)
			}
		}
		sort.Strings(expectedKeys)
		sort.Strings(gotKeys)
		if !reflect.DeepEqual(gotKeys, expectedKeys) {
			t.Errorf("PresignedURL %s contain sign \n%#v, want \n%#v", method, gotKeys, expectedKeys)
		}
	}
}
