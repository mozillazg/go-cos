package main

import (
	"bitbucket.org/mozillazg/go-cos"
	"context"
	"fmt"
	"os"
	"time"
)

func main() {
	c := cos.NewClient(os.Getenv("COS_SECRETID"), os.Getenv("COS_SECRETKEY"), nil)
	startTime := time.Now()
	endTime := startTime.Add(time.Hour)
	s, _, err := c.Service.Get(context.Background(), startTime, endTime,
		startTime, endTime)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%#v", s)
}
