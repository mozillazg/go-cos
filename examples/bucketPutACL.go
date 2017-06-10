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
	startTime := time.Now()
	endTime := startTime.Add(time.Hour)
	opt := &cos.BucketPutACLOptions{
		XCosACL: "private",
	}
	_, err := c.Bucket.PutACL(context.Background(), startTime, endTime,
		startTime, endTime, opt, nil)
	if err != nil {
		fmt.Println(err)
	}

}
