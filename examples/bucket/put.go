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
	u, _ := url.Parse("https://testdelete-1253846586.cn-north.myqcloud.com")
	b := &cos.BaseURL{
		BucketURL: u,
	}
	c := cos.NewClient(os.Getenv("COS_SECRETID"), os.Getenv("COS_SECRETKEY"), b, nil)
	c.Client.Transport = &cos.DebugRequestTransport{
		RequestHeader:  true,
		RequestBody:    true,
		ResponseHeader: true,
		ResponseBody:   true,
	}

	//opt := &cos.BucketPutOptions{
	//	XCosACL: "public-read",
	//}
	_, err := c.Bucket.Put(context.Background(), cos.NewAuthTime(time.Hour), nil)
	if err != nil {
		fmt.Println(err)
	}
}
