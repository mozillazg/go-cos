package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"time"

	"bitbucket.org/mozillazg/go-cos"
)

func main() {
	u, _ := url.Parse("https://test-1253846586.cn-north.myqcloud.com")
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(os.Getenv("COS_SECRETID"), os.Getenv("COS_SECRETKEY"), b, nil)
	c.Client.Transport = &cos.DebugRequestTransport{
		RequestHeader:  true,
		RequestBody:    true,
		ResponseHeader: true,
		ResponseBody:   true,
	}

	name := "test/hello.txt"
	v, _, err := c.Object.GetACL(context.Background(), cos.NewAuthTime(time.Hour), name)
	if err != nil {
		fmt.Println(err)
	}
	for _, a := range v.AccessControlList {
		fmt.Printf("%s, %s, %s\n", a.Grantee.Type, a.Grantee.UIN, a.Permission)
	}

}
