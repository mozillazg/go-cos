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
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(os.Getenv("COS_SECRETID"), os.Getenv("COS_SECRETKEY"), b, nil)
	c.Client.Transport = &cos.DebugRequestTransport{
		RequestHeader:  true,
		RequestBody:    false,
		ResponseHeader: true,
		ResponseBody:   true,
	}

	name := "test_multipart.txt"
	v, _, err := c.Object.InitiateMultipartUpload(context.Background(), cos.NewAuthTime(time.Hour), name, nil)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%s\n", v.UploadID)
}
