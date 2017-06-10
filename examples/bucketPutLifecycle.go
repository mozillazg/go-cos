package main

import (
	"bitbucket.org/mozillazg/go-cos"
	"context"
	"fmt"
	"os"
	"time"
)

func main() {
	b, _ := cos.ParseBucketFromDomain("huanan-1253846586.cn-south.myqcloud.com")
	c := cos.NewClient(os.Getenv("COS_SECRETID"), os.Getenv("COS_SECRETKEY"), b, nil)
	c.Secure = false
	startTime := time.Now()
	endTime := startTime.Add(time.Hour)
	lc := &cos.BucketLifecycleResult{
		Rules: []cos.BucketLifecycleRule{
			{
				ID:     "1234",
				Prefix: "test",
				Status: "Enabled",
				Transition: cos.BucketLifecycleTransition{
					Days:         10,
					StorageClass: "Standard",
				},
			},
		},
	}
	_, err := c.Bucket.PutLifecycle(context.Background(), startTime, endTime,
		startTime, endTime, lc)
	if err != nil {
		fmt.Println(err)
	}
}
