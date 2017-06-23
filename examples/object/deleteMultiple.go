package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"time"

	"bytes"
	"io"

	"math/rand"

	"bitbucket.org/mozillazg/go-cos"
)

func genBigData(blockSize int) []byte {
	b := make([]byte, blockSize)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return b
}

func uploadMulti(c *cos.Client, authTime *cos.AuthTime) []string {
	names := []string{}
	data := genBigData(1024 * 1024 * 1)
	ctx := context.Background()
	var r io.Reader
	var name string
	n := 3

	for n > 0 {
		name = fmt.Sprintf("test/test_multi_delete_%s", time.Now().Format(time.RFC3339))
		r = bytes.NewReader(data)

		c.Object.Put(ctx, authTime, name, r, nil)
		names = append(names, name)
		n--
	}
	return names
}

func main() {
	u, _ := url.Parse("https://test-1253846586.cn-north.myqcloud.com")
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(os.Getenv("COS_SECRETID"), os.Getenv("COS_SECRETKEY"), b, nil)
	transport := &cos.DebugRequestTransport{
		RequestHeader:  true,
		RequestBody:    false,
		ResponseHeader: true,
		ResponseBody:   true,
	}
	c.Client.Transport = transport
	ctx := context.Background()
	authTime := cos.NewAuthTime(time.Hour)

	names := uploadMulti(c, authTime)
	names = append(names, []string{"a", "b", "c", "a+bc/xx&?+# "}...)
	obs := []*cos.ObjectForDelete{}
	for _, v := range names {
		obs = append(obs, &cos.ObjectForDelete{Key: v})
	}
	//sha1 := ""
	opt := &cos.ObjectDeleteMultiOptions{
		Objects: obs,
		//XCosSha1: sha1,
		//Quiet: true,
	}
	transport.RequestBody = true

	v, _, err := c.Object.DeleteMulti(ctx, authTime, opt)
	if err != nil {
		fmt.Println(err)
	}

	for _, x := range v.DeletedObjects {
		fmt.Printf("deleted %s\n", x.Key)
	}
	for _, x := range v.Errors {
		fmt.Printf("error %s, %s, %s\n", x.Key, x.Code, x.Message)
	}
}
