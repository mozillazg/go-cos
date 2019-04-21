package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"

	"github.com/mozillazg/go-cos"
	"github.com/mozillazg/go-cos/debug"
)

type MockSender struct{}

func (s *MockSender) Send(ctx context.Context, caller cos.Caller, req *http.Request) (*http.Response, error) {
	// 如果用不到 response 的话，也可以直接 return &http.Response{}, nil
	resp, _ := http.ReadResponse(bufio.NewReader(bytes.NewReader([]byte(`HTTP/1.1 200 OK
Content-Length: 6
Accept-Ranges: bytes
Connection: keep-alive
Content-Type: text/plain; charset=utf-8
Date: Sat, 19 Jan 2019 08:25:27 GMT
Etag: "f572d396fae9206628714fb2ce00f72e94f2258f"
Last-Modified: Mon, 12 Jun 2017 13:36:19 GMT
Server: tencent-cos
X-Cos-Request-Id: NWM0MmRlZjdfMmJhZDM1MGFfNDFkM19hZGI3MQ==

hello
`))), nil)
	return resp, nil
}

type MockerResponseParser struct {
	result *cos.ObjectGetACLResult
}

func (p *MockerResponseParser) ParseResponse(ctx context.Context, caller cos.Caller, resp *http.Response, result interface{}) (*cos.Response, error) {
	b, _ := ioutil.ReadAll(resp.Body)
	if string(b) != "hello\n" {
		panic(string(b))
	}

	// 插入预设的结果
	switch caller.Method {
	case cos.MethodObjectGetACL:
		v := result.(*cos.ObjectGetACLResult)
		*v = *p.result
	}

	return &cos.Response{Response: resp}, nil
}

func main() {
	b, _ := cos.NewBaseURL("http://cos.example.com")
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  os.Getenv("COS_SECRETID"),
			SecretKey: os.Getenv("COS_SECRETKEY"),
			Transport: &debug.DebugRequestTransport{
				RequestHeader:  true,
				RequestBody:    true,
				ResponseHeader: true,
				ResponseBody:   true,
			},
		},
	})
	c.Sender = &MockSender{}
	acl := &cos.ObjectGetACLResult{
		Owner: &cos.Owner{
			ID: "test",
		},
		AccessControlList: []cos.ACLGrant{
			{
				Permission: "READ",
			},
		},
	}
	c.ResponseParser = &MockerResponseParser{acl}

	result, resp, err := c.Object.GetACL(context.Background(), "test/mock.go")
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	fmt.Printf("%#v\n", result)
	if !reflect.DeepEqual(*result, *acl) {
		panic(*result)
	}
}
