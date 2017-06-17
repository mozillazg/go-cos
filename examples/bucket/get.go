package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"net/url"

	"bitbucket.org/mozillazg/go-cos"
)

func main() {
	u, _ := url.Parse("https://test-1253846586.cn-north.myqcloud.com")
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

	opt := &cos.BucketGetOptions{
		Prefix: "test",
	}
	v, _, err := c.Bucket.Get(context.Background(), cos.NewAuthTime(time.Hour), opt)
	if err != nil {
		fmt.Println(err)
	}

	for _, c := range v.Contents {
		fmt.Printf("%s, %d\n", c.Key, c.Size)
	}
}
