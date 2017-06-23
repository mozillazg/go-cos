package main

import (
	"context"
	"net/url"
	"os"
	"time"

	"bitbucket.org/mozillazg/go-cos"
)

func main() {
	u, _ := url.Parse("https://test-1253846586.cn-north.myqcloud.com")
	b := &cos.BaseURL{
		BucketURL: u,
	}
	c := cos.NewClient(os.Getenv("COS_SECRETID"), os.Getenv("COS_SECRETKEY"), b, nil)
	c.Client.Transport = &cos.DebugRequestTransport{
		RequestHeader:  true,
		RequestBody:    true,
		ResponseHeader: true,
		ResponseBody:   true,
	}

	authTime := cos.NewAuthTime(time.Hour)

	// with header
	opt := &cos.BucketPutACLOptions{
		Header: &cos.ACLHeaderOptions{
			XCosACL: "private",
		},
	}
	_, err := c.Bucket.PutACL(context.Background(), authTime, opt)
	if err != nil {
		panic(err)
	}

	// with body
	opt = &cos.BucketPutACLOptions{
		Body: &cos.ACLXml{
			Owner: &cos.Owner{
				UIN: "100000760461",
			},
			AccessControlList: []cos.ACLGrant{
				{
					Grantee: &cos.ACLGrantee{
						Type: "RootAccount",
						UIN:  "100000760461",
					},

					Permission: "FULL_CONTROL",
				},
			},
		},
	}
	_, err = c.Bucket.PutACL(context.Background(), authTime, opt)
	if err != nil {
		panic(err)
	}
}
