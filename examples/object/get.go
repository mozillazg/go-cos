package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"time"

	"io/ioutil"

	"bitbucket.org/mozillazg/go-cos"
)

func main() {
	u, _ := url.Parse("https://test-1253846586.cn-north.myqcloud.com")
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(os.Getenv("COS_SECRETID"), os.Getenv("COS_SECRETKEY"), b, nil)
	c.Client.Transport = &cos.DebugRequestTransport{
		RequestHeader:  true,
		RequestBody:    true,
		ResponseHeader: true,
		ResponseBody:   true,
	}

	name := "test/hello.txt"
	resp, err := c.Object.Get(context.Background(), cos.NewAuthTime(time.Hour), name, nil)
	if err != nil {
		fmt.Println(err)
	}
	bs, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	fmt.Printf("%s\n", string(bs))

	// range
	opt := &cos.ObjectGetOptions{
		ResponseContentType: "text/html",
		Range:               "bytes=0-3",
	}
	resp, err = c.Object.Get(context.Background(), cos.NewAuthTime(time.Hour), name, opt)
	if err != nil {
		fmt.Println(err)
	}
	bs, _ = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	fmt.Printf("%s\n", string(bs))
}
