package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"io/ioutil"

	"bitbucket.org/mozillazg/go-cos"
)

func upload(c *cos.Client, name string) {
	f := strings.NewReader("test")
	f = strings.NewReader("test xxx")
	opt := &cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			ContentType: "text/html",
		},
		ACLHeaderOptions: &cos.ACLHeaderOptions{
			XCosACL: "public-read",
		},
	}
	c.Object.Put(context.Background(), cos.NewAuthTime(time.Hour), name, f, opt)
	return
}

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

	name := "test/anonymous_get.go"
	upload(c, name)

	resp, err := c.Object.Get(context.Background(), nil, name, nil)
	if err != nil {
		panic(err)
		return
	}
	bs, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	fmt.Printf("%s\n", string(bs))
}
