package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/mozillazg/go-cos"
	"github.com/mozillazg/go-cos/debug"
)

func main() {
	b, _ := cos.NewBaseURL(os.Getenv("COS_BUCKET_URL"))
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

	name := "test/objectPut.go"
	f := strings.NewReader("test")

	_, err := c.Object.Put(context.Background(), name, f, nil)
	if err != nil {
		panic(err)
	}

	// 测试上传以及特殊字符
	name = "test/put_ + !'()* option.go"
	contentDisposition := "attachment; filename=Hello - world!(+)'*.go"
	f = strings.NewReader("test xxx")
	opt := &cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			ContentType:        "text/html",
			ContentDisposition: contentDisposition,
		},
		ACLHeaderOptions: &cos.ACLHeaderOptions{
			// XCosACL: "public-read",
			XCosACL: "private",
		},
	}
	resp, err := c.Object.Put(context.Background(), name, f, opt)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()

	// 测试特殊字符
	resp, err = c.Object.Get(context.Background(), name, nil)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	if resp.Header.Get("Content-Disposition") != contentDisposition {
		panic(errors.New("wong Content-Disposition"))
	}
}
