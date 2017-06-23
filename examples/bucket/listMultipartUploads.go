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

	opt := &cos.ListMultipartUploadsOptions{
		Prefix: "t",
	}
	v, _, err := c.Bucket.ListMultipartUploads(context.Background(), cos.NewAuthTime(time.Hour), opt)
	if err != nil {
		panic(err)
	}
	for _, p := range v.Uploads {
		fmt.Printf("%s\n", p.Key)
	}
}
