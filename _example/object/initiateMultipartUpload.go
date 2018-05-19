package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"time"

	"net/http"

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

	name := "test_multipart" + time.Now().Format(time.RFC3339)
	v, _, err := c.Object.InitiateMultipartUpload(context.Background(), name, nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", v.UploadID)
}
