package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"strings"

	"net/http"

	"github.com/mozillazg/go-cos"
	"github.com/mozillazg/go-cos/debug"
)

func initUpload(c *cos.Client, name string) *cos.InitiateMultipartUploadResult {
	v, _, err := c.Object.InitiateMultipartUpload(context.Background(), name, nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", v)
	return v
}

func uploadPart(c *cos.Client, name string, uploadID string, blockSize, n int) string {

	b := make([]byte, blockSize)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	s := fmt.Sprintf("%X", b)
	f := strings.NewReader(s)

	resp, err := c.Object.UploadPart(
		context.Background(), name, uploadID, n, f, nil,
	)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Printf("%s\n", resp.Status)
	return resp.Header.Get("Etag")
}

func main() {
	u, _ := url.Parse(os.Getenv("COS_BUCKET_URL"))
	b := &cos.BaseURL{BucketURL: u}
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

	name := "test/test_list_parts.go"
	up := initUpload(c, name)
	uploadID := up.UploadID
	ctx := context.Background()
	blockSize := 1024 * 1024 * 3

	for i := 1; i < 5; i++ {
		uploadPart(c, name, uploadID, blockSize, i)
	}

	v, _, err := c.Object.ListParts(ctx, name, uploadID)
	if err != nil {
		panic(err)
		return
	}
	for _, p := range v.Parts {
		fmt.Printf("%d, %s, %d\n", p.PartNumber, p.ETag, p.Size)
	}
	fmt.Printf("%s\n", v.Initiator.ID)
	fmt.Printf("%s\n", v.Owner.ID)
}
