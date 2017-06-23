package main

import (
	"context"
	//"net/url"
	"os"
	"time"

	"bitbucket.org/mozillazg/go-cos"
)

func main() {
	//u, _ := url.Parse("https://test-1253846586.cn-north.myqcloud.com")
	u := cos.NewBucketURL("test", "1253846586", "cn-north", true)
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(os.Getenv("COS_SECRETID"), os.Getenv("COS_SECRETKEY"), b, nil)
	c.Client.Transport = &cos.DebugRequestTransport{
		RequestHeader:  true,
		RequestBody:    true,
		ResponseHeader: true,
		ResponseBody:   true,
	}

	name := "test/hello.txt"
	_, err := c.Object.Head(context.Background(), cos.NewAuthTime(time.Hour), name, nil)
	if err != nil {
		panic(err)
	}
}
