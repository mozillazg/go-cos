package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

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
				RequestBody:    false,
				ResponseHeader: true,
				ResponseBody:   true,
			},
		},
	})

	name := "test/uploadFile.go"
	f, err := os.Open(os.Args[0])
	if err != nil {
		panic(err)
	}
	defer f.Close()
	s, err := f.Stat()
	if err != nil {
		panic(err)
	}
	fmt.Println(s.Size())

	resp, err = c.Object.Put(context.Background(), name, f, nil)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}
