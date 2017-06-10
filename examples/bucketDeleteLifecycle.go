package main

import (
	"bitbucket.org/mozillazg/go-cos"
	"context"
	"fmt"
	"os"
	"time"
)

func main() {
	b, _ := cos.ParseBucketFromDomain("testhuanan-1253846586.cn-south.myqcloud.com")
	c := cos.NewClient(os.Getenv("COS_SECRETID"), os.Getenv("COS_SECRETKEY"), b, nil)
	c.Secure = false
	startTime := time.Now()
	endTime := startTime.Add(time.Hour)
	_, err := c.Bucket.DeleteLifecycle(context.Background(), startTime, endTime,
		startTime, endTime)
	if err != nil {
		fmt.Println(err)
	}
}
