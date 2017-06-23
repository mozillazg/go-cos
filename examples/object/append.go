package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"time"

	"bitbucket.org/mozillazg/go-cos"
)

func genBigData(blockSize int) []byte {
	b := make([]byte, blockSize)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return b
}

func main() {
	// "https://test-1253846586.cn-north.myqcloud.com",
	u, _ := url.Parse("https://huadong-1253846586.cn-east.myqcloud.com")
	b := &cos.BaseURL{
		BucketURL: u,
	}
	c := cos.NewClient(os.Getenv("COS_SECRETID"), os.Getenv("COS_SECRETKEY"), b, nil)
	c.Client.Transport = &cos.DebugRequestTransport{
		RequestHeader:  true,
		RequestBody:    false,
		ResponseHeader: true,
		ResponseBody:   true,
	}

	startTime := time.Now()
	authTime := cos.NewAuthTime(time.Hour)
	name := fmt.Sprintf("test/test_object_append_%s", startTime.Format(time.RFC3339))
	data := genBigData(1024 * 1024 * 1)
	length := len(data)
	r := bytes.NewReader(data)

	ctx := context.Background()

	// 第一次就必须 append
	resp, err := c.Object.Append(ctx, authTime, name, 0, r, nil)
	if err != nil {
		panic(err)
		return
	}
	fmt.Printf("%s\n", resp.Status)

	// head
	if _, err = c.Object.Head(ctx, authTime, name, nil); err != nil {
		panic(err)
		return
	}

	// 再次 append
	data = genBigData(1024 * 1024 * 5)
	r = bytes.NewReader(data)
	resp, err = c.Object.Append(context.Background(), authTime, name, length, r, nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", resp.Status)
}
