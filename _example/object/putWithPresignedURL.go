package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
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

	name := "test/objectPut.go"
	ctx := context.Background()
	f := strings.NewReader("test")

	// 通过生成签名 header 上传文件
	_, err := c.Object.Put(ctx, name, f, nil)
	if err != nil {
		panic(err)
	}

	// 获取预签名授权 URL
	opt := &cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			// 指定上传内容的 Content-Type，用于指定下载时 response 的 Content-Type
			ContentType: "text/html",
		},
	}
	presignedURL, err := c.Object.PresignedURL(ctx, http.MethodPut, name, auth, opt)
	if err != nil {
		panic(err)
	}

	// 通过预签名授权 URL 上传
	data := "test upload with presignedURL"
	f = strings.NewReader(data)
	req, err := http.NewRequest(http.MethodPut, presignedURL.String(), f)
	if err != nil {
		panic(err)
	}
	// 指定上传内容的 Content-Type，用于指定下载时 response 的 Content-Type
	req.Header.Set("Content-Type", "text/html")
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	// 验证上传的内容
	resp, err := c.Object.Get(ctx, name, nil)
	if err != nil {
		panic(err)
	}
	bs, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	fmt.Printf("%s\n", string(bs))
	fmt.Printf("Content-Type: %s\n\n", resp.Header.Get("Content-Type"))
	fmt.Printf("%v\n\n", strings.Compare(data, string(bs)) == 0)

	// c.Object.Put 使用 预签名授权 URL
	c2 := cos.NewClient(b, &http.Client{
		Transport: &debug.DebugRequestTransport{
			RequestHeader:  true,
			RequestBody:    true,
			ResponseHeader: true,
			ResponseBody:   true,
		},
	})
	f = strings.NewReader("test c.Object.Put with presignedURL")
	opt.PresignedURL = presignedURL
	resp2, err := c2.Object.Put(ctx, name, f, opt)
	if err != nil {
		panic(err)
	}
	resp2.Body.Close()
}
