package main

import (
	"bitbucket.org/mozillazg/go-cos"
	"context"
	"fmt"
	"os"
	"time"
)

func main() {
	b, _ := cos.ParseBucketFromDomain("testdelete-1253846586.cn-north.myqcloud.com")
	c := cos.NewClient(os.Getenv("COS_SECRETID"), os.Getenv("COS_SECRETKEY"), b, nil)
	startTime := time.Now()
	endTime := startTime.Add(time.Hour)
	//opt := &cos.BucketPutOptions{
	//	XCosAcl: "public-read",
	//}
	_, err := c.Bucket.Put(context.Background(), startTime, endTime,
		startTime, endTime, nil)
	if err != nil {
		fmt.Println(err)
	}

	opt := &cos.BucketGetOptions{
		MaxKeys: 1,
	}
	v, _, err := c.Bucket.Get(context.Background(), startTime, endTime,
		startTime, endTime, opt)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%#v", v)
}
