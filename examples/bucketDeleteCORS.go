package main

import (
	"bitbucket.org/mozillazg/go-cos"
	"context"
	"fmt"
	"os"
	"time"
)

func main() {
	b, _ := cos.ParseBucketFromDomain("test-1253846586.cn-north.myqcloud.com")
	c := cos.NewClient(os.Getenv("COS_SECRETID"), os.Getenv("COS_SECRETKEY"), b, nil)
	c.Secure = false
	startTime := time.Now()
	endTime := startTime.Add(time.Hour)
	_, err := c.Bucket.DeleteCORS(context.Background(), cos.NewAuthTime(
		startTime, endTime,
		startTime, endTime))
	if err != nil {
		fmt.Println(err)
	}
}
