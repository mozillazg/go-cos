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
	startTime := time.Now()
	endTime := startTime.Add(time.Hour)
	b, _ := cos.ParseBucketFromDomain("test-1253846586.cn-north.myqcloud.com")
	v, _, err := c.Bucket.ListMultipartUploads(context.Background(), b, startTime, endTime,
		startTime, endTime, nil)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%#v", v)

	opt := &cos.ListMultipartUploadsOptions{
		Prefix: "1",
	}
	v, _, err = c.Bucket.ListMultipartUploads(context.Background(), b, startTime, endTime,
		startTime, endTime, opt)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%#v", v)
}
