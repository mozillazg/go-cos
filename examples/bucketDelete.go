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
	b, _ := cos.ParseBucketFromDomain("testdelete-1253846586.cn-north.myqcloud.com")
	_, err := c.Bucket.Delete(context.Background(), b, startTime, endTime,
		startTime, endTime)
	if err != nil {
		fmt.Println(err)
	}
}
