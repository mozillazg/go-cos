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
	startTime := time.Now()

	tg := &cos.BucketPutTaggingOptions{
		TagSet: []*cos.BucketTaggingTag{
			{
				Key:   "test_k2",
				Value: "test_v2",
			},
			{
				Key:   "test_k3",
				Value: "test_v3",
			},
			{
				Key:   startTime.Format("02_Jan_06_15_04_MST"),
				Value: "test_time",
			},
		},
	}
	_, err := c.Bucket.PutTagging(context.Background(), cos.NewAuthTime(time.Hour), tg)
	if err != nil {
		fmt.Println(err)
	}
}
