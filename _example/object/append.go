package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"time"

	"net/http"

	"github.com/mozillazg/go-cos"
	"github.com/mozillazg/go-cos/debug"
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
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  os.Getenv("COS_SECRETID"),
			SecretKey: os.Getenv("COS_SECRETKEY"),
			Transport: &debug.DebugRequestTransport{
				RequestHeader:  true,
				RequestBody:    false,
				ResponseHeader: true,
				ResponseBody:   true,
			},
		},
	})

	startTime := time.Now()

	name := fmt.Sprintf("test/test_object_append_%s", startTime.Format(time.RFC3339))
	data := genBigData(1024 * 1024 * 1)
	length := len(data)
	r := bytes.NewReader(data)

	ctx := context.Background()

	// 第一次就必须 append
	resp, err := c.Object.Append(ctx, name, 0, r, nil)
	if err != nil {
		panic(err)
		return
	}
	defer resp.Body.Close()
	fmt.Printf("%s\n", resp.Status)

	// head
	resp, err = c.Object.Head(ctx, name, nil)
	if err != nil {
		panic(err)
		return
	}
	defer resp.Body.Close()

	// 再次 append
	data = genBigData(1024 * 1024 * 5)
	r = bytes.NewReader(data)
	resp, err = c.Object.Append(context.Background(), name, length, r, nil)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Printf("%s\n", resp.Status)
}
