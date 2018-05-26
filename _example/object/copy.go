package main

import (
	"context"
	"net/url"
	"os"
	"strings"

	"net/http"

	"fmt"
	"io/ioutil"
	"time"

	"github.com/mozillazg/go-cos"
	"github.com/mozillazg/go-cos/debug"
)

func main() {
	u, _ := url.Parse(os.Getenv("COS_BUCKET_URL"))
	b := &cos.BaseURL{BucketURL: u}
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

	source := "test/objectMove1.go"
	expected := "test"
	f := strings.NewReader(expected)

	_, err := c.Object.Put(context.Background(), source, f, nil)
	if err != nil {
		panic(err)
	}

	soruceURL := fmt.Sprintf("%s/%s", u.Host, source)
	dest := fmt.Sprintf("test/objectMove_%d.go", time.Now().Nanosecond())
	//opt := &cos.ObjectCopyOptions{}
	res, _, err := c.Object.Copy(context.Background(), dest, soruceURL, nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n\n", res)

	resp, err := c.Object.Get(context.Background(), dest, nil)
	if err != nil {
		panic(err)
	}
	bs, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	result := string(bs)
	if result != expected {
		panic(fmt.Sprintf("%s != %s", result, expected))
	}
}
