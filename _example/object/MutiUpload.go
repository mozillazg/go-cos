package main

import (
	"context"
	"net/url"
	"os"

	"net/http"

	"fmt"

	"github.com/mozillazg/go-cos"
	"github.com/mozillazg/go-cos/debug"
)

func main() {
	u, _ := url.Parse("http://lewzylu02-1252448703.cos.ap-guangzhou.myqcloud.com")
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  os.Getenv("COS_Key"),
			SecretKey: os.Getenv("COS_Secret"),
			Transport: &debug.DebugRequestTransport{
				RequestHeader:  false,
				RequestBody:    false,
				ResponseHeader: false,
				ResponseBody:   false,
			},
		},
	})
	f,err:=os.Open("E:/cos-php-sdk.zip")
	if err!=nil {panic(err)}
	v, _, err := c.Object.MultiUpload(
		context.Background(), "test", f, nil,
	)
	if err!=nil {panic(err)}
	fmt.Println(v)
}