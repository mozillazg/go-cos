package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

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

	name := "test/postFile.go"
	f, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	s, err := f.Stat()
	if err != nil {
		panic(err)
	}
	fmt.Println(s.Size())
	opt := &cos.ObjectPostOptions{
		ContentType: "text/html",
	}
	auth := cos.Auth{
		SecretID:  os.Getenv("COS_SECRETID"),
		SecretKey: os.Getenv("COS_SECRETKEY"),
		Expire:    time.Hour,
	}
	result, _, err := c.Object.Post(context.Background(), name, f, auth, opt)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", result)
}
