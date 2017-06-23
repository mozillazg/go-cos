package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"time"

	"bitbucket.org/mozillazg/go-cos"
)

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

	opt := &cos.ObjectPutACLOptions{
		Header: &cos.ACLHeaderOptions{
			XCosACL: "private",
		},
	}
	name := "test/hello.txt"
	_, err := c.Object.PutACL(context.Background(), cos.NewAuthTime(time.Hour), name, opt)
	if err != nil {
		fmt.Println(err)
	}

	// with body
	opt = &cos.ObjectPutACLOptions{
		Body: &cos.BucketGetACLResult{
			Owner: &cos.Owner{
				UIN: "100000760461",
			},
			AccessControlList: []*cos.BucketACLGrant{
				{
					Grantee: &cos.BucketACLGrantee{
						Type: "RootAccount",
						UIN:  "100000760461",
					},

					Permission: "FULL_CONTROL",
				},
			},
		},
	}

	_, err = c.Object.PutACL(context.Background(), cos.NewAuthTime(time.Hour), name, opt)
	if err != nil {
		fmt.Println(err)
	}
}
