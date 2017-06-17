package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"net/url"
	"strings"

	"bitbucket.org/mozillazg/go-cos"
)

func initUpload(c *cos.Client, authTime *cos.AuthTime,
	name string,
) *cos.ObjectInitiateMultipartUploadResult {
	v, _, err := c.Object.InitiateMultipartUpload(context.Background(), authTime, name, nil)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%#v\n", v)
	return v
}

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

	authTime := cos.NewAuthTime(time.Hour)
	name := "test/test_multi_upload.go"
	up := initUpload(c, authTime, name)
	uploadID := up.UploadID

	f := strings.NewReader("test heoo")
	_, err := c.Object.UploadPart(
		context.Background(), authTime, name, uploadID, 1, f, nil,
	)
	if err != nil {
		fmt.Println(err)
	}
}
