package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/mozillazg/go-cos"
	"github.com/mozillazg/go-cos/debug"
)

func main() {
	b, _ := cos.NewBaseURL(os.Getenv("COS_BUCKET_URL"))
	auth := cos.Auth{
		SecretID:  os.Getenv("COS_SECRETID"),
		SecretKey: os.Getenv("COS_SECRETKEY"),
		Expire:    time.Hour,
	}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  auth.SecretID,
			SecretKey: auth.SecretKey,
			Expire:    auth.Expire,
			Transport: &debug.DebugRequestTransport{
				RequestHeader:  true,
				RequestBody:    true,
				ResponseHeader: true,
				ResponseBody:   true,
			},
		},
	})

	name := "test/hello.txt"
	ctx := context.Background()

	// 通过生成签名 header 下载文件
	resp, err := c.Object.Get(ctx, name, nil)
	if err != nil {
		panic(err)
	}
	bs, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	fmt.Printf("%s\n", string(bs))

	// 获取预签名授权 URL
	presignedURL, err := c.Object.PresignedURL(ctx, http.MethodGet, name, auth, nil)
	if err != nil {
		panic(err)
	}

	// 通过预签名授权 URL 下载文件
	resp2, err := http.Get(presignedURL.String())
	if err != nil {
		panic(err)
	}
	bs2, _ := ioutil.ReadAll(resp2.Body)
	resp2.Body.Close()
	fmt.Printf("%s\n", string(bs2))

	fmt.Printf("%v\n\n", bytes.Compare(bs2, bs) == 0)

	// c.Object.Get 使用 预签名授权 URL
	c2 := cos.NewClient(b, &http.Client{
		Transport: &debug.DebugRequestTransport{
			RequestHeader:  true,
			RequestBody:    true,
			ResponseHeader: true,
			ResponseBody:   true,
		},
	})
	resp3, err := c2.Object.Get(ctx, name, &cos.ObjectGetOptions{
		PresignedURL: presignedURL,
	})
	if err != nil {
		panic(err)
	}
	resp3.Body.Close()
}
