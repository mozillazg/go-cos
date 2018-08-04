package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"

	"io/ioutil"

	"github.com/mozillazg/go-cos"
)

func upload(c *cos.Client, key string) {
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
	c.Object.Put(context.Background(), key, f, opt)
	return
}

func main() {
	u, _ := url.Parse(os.Getenv("COS_BUCKET_URL"))
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, nil)

	key := "test/anonymous_get.go"
	upload(c, name)

	resp, err := c.Object.Get(context.Background(), key, nil)
	if err != nil {
		panic(err)
		return
	}
	bs, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	fmt.Printf("%s\n", string(bs))
}
