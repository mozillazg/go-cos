package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"strings"
	"time"

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

func uploadPart(c *cos.Client, authTime *cos.AuthTime,
	name string, uploadID string, blockSize, n int) string {

	b := make([]byte, blockSize)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	s := fmt.Sprintf("%X", b)
	f := strings.NewReader(s)

	resp, err := c.Object.UploadPart(
		context.Background(), authTime, name, uploadID, n, f, nil,
	)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%s\n", resp.Status)
	return resp.Header.Get("Etag")
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
	name := "test/test_complete_upload.go"
	up := initUpload(c, authTime, name)
	uploadID := up.UploadID
	blockSize := 1024 * 1024 * 3

	opt := &cos.ObjectCompleteMultipartUploadOption{}
	for i := 1; i < 5; i++ {
		etag := uploadPart(c, authTime, name, uploadID, blockSize, i)
		opt.Parts = append(opt.Parts, &cos.ObjectPart{
			PartNumber: i, ETag: etag},
		)
	}

	c.Client.Transport = &cos.DebugRequestTransport{
		RequestHeader:  true,
		RequestBody:    true,
		ResponseHeader: true,
		ResponseBody:   true,
	}
	v, resp, err := c.Object.CompleteMultipartUpload(
		context.Background(), authTime, name, uploadID, opt,
	)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%s\n", resp.Status)
	fmt.Printf("%#v\n", v)
	fmt.Printf("%s\n", v.Location)
}
