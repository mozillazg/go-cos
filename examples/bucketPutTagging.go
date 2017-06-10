package main

import (
	"bitbucket.org/mozillazg/go-cos"
	"context"
	"fmt"
	"os"
	"time"
)

func main() {
	c := cos.NewClient(os.Getenv("COS_SECRETID"), os.Getenv("COS_SECRETKEY"), nil)
	c.Secure = false
	startTime := time.Now()
	endTime := startTime.Add(time.Hour)
	b, _ := cos.ParseBucketFromDomain("test-1253846586.cn-north.myqcloud.com")
	tg := &cos.BucketTaggingResult{
		TagSet: []cos.BucketTaggingTag{
			{
				Key:   "test_k2",
				Value: "test_v2",
			},
			{
				Key:   "test_k3",
				Value: "test_v3",
			},
		},
	}
	_, err := c.Bucket.PutTagging(context.Background(), b, startTime, endTime,
		startTime, endTime, tg)
	if err != nil {
		fmt.Println(err)
	}
}
