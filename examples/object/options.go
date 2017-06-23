package main

import (
	"context"
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
	opt := &cos.ObjectOptionsOptions{
		Origin: "http://www.qq.com",
		AccessControlRequestMethod: "PUT",
	}
	_, err := c.Object.Options(context.Background(), cos.NewAuthTime(time.Hour), name, opt)
	if err != nil {
		panic(err)
	}
}
