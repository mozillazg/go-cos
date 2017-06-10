package main

import (
	"bitbucket.org/mozillazg/go-cos"
	"context"
	"fmt"
	"net/http"
	"os"
	"time"
)

func main() {
	c := cos.NewClient(os.Getenv("COS_SECRETID"), os.Getenv("COS_SECRETKEY"), nil)
	c.Secure = false
	startTime := time.Now()
	endTime := startTime.Add(time.Hour)
	b, _ := cos.ParseBucketFromDomain("test-1253846586.cn-north.myqcloud.com")
	cors := &cos.BucketCORSResult{
		Rules: []cos.BucketCORSRule{
			{
				ID:            "1234",
				AllowedOrigin: "http://www.qq.com",
				AllowedMethod: http.MethodPut,
				AllowedHeader: "x-cos-meta-test",
				MaxAgeSeconds: 500,
				ExposeHeader:  "x-cos-meta-test1",
			},
		},
	}
	_, err := c.Bucket.PutCORS(context.Background(), b, startTime, endTime,
		startTime, endTime, cors)
	if err != nil {
		fmt.Println(err)
	}
}
