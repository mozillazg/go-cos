package main

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"net/http"

	"github.com/mozillazg/go-cos"
	"github.com/mozillazg/go-cos/debug"
)

func main() {
	u, _ := url.Parse("https://test-1253846586.cos.ap-guangzhou.myqcloud.com")
	b := &cos.BaseURL{
		BucketURL: u,
	}
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

	v, _, err := c.Bucket.GetACL(context.Background())
	if err != nil {
		panic(err)
	}
	for _, a := range v.AccessControlList {
		fmt.Printf("%s, %s, %s\n", a.Grantee.Type, a.Grantee.ID, a.Permission)
	}
}
