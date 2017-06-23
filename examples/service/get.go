package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"bitbucket.org/mozillazg/go-cos"
)

func main() {
	c := cos.NewClient(os.Getenv("COS_SECRETID"), os.Getenv("COS_SECRETKEY"),
		nil, nil,
	)
	c.Client.Transport = &cos.DebugRequestTransport{
		RequestHeader:  true,
		RequestBody:    true,
		ResponseHeader: true,
		ResponseBody:   true,
	}

	s, _, err := c.Service.Get(context.Background(), cos.NewAuthTime(time.Hour))
	if err != nil {
		panic(err)
	}

	for _, b := range s.Buckets {
		fmt.Printf("%#v\n", b)
	}
}
