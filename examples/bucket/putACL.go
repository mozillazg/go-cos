package main

import (
	"context"
	"net/url"
	"os"

	"net/http"

	"bitbucket.org/mozillazg/go-cos"
)

func main() {
	u, _ := url.Parse("https://test-1253846586.cn-north.myqcloud.com")
	b := &cos.BaseURL{
		BucketURL: u,
	}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  os.Getenv("COS_SECRETID"),
			SecretKey: os.Getenv("COS_SECRETKEY"),
			Transport: &cos.DebugRequestTransport{
				RequestHeader:  true,
				RequestBody:    true,
				ResponseHeader: true,
				ResponseBody:   true,
			},
		},
	})

	// with header
	opt := &cos.BucketPutACLOptions{
		Header: &cos.ACLHeaderOptions{
			XCosACL: "private",
		},
	}
	_, err := c.Bucket.PutACL(context.Background(), opt)
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
	_, err = c.Bucket.PutACL(context.Background(), opt)
	if err != nil {
		panic(err)
	}
}
