package main

import (
	"context"
	"fmt"
	"os"

	"net/http"

	"github.com/mozillazg/go-cos"
	"github.com/mozillazg/go-cos/debug"
)

func main() {
	b, _ := cos.NewBaseURL(os.Getenv("COS_BUCKET_URL"))
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  os.Getenv("COS_SECRETID"),
			SecretKey: os.Getenv("COS_SECRETKEY"),
			Transport: &debug.DebugRequestTransport{
				RequestHeader:  true,
				RequestBody:    true,
				ResponseHeader: true,
				ResponseBody:   true,
			},
		},
	})

	opt := &cos.BucketGetOptions{
		Prefix:  "test",
		MaxKeys: 3,
	}
	v, resp, err := c.Bucket.Get(context.Background(), opt)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()

	for _, c := range v.Contents {
		fmt.Printf("%s, %d\n", c.Key, c.Size)
	}

	// 测试特殊字符
	opt.Prefix = "test/put_ + !'()* option"
	_, resp, err = c.Bucket.Get(context.Background(), opt)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
}
