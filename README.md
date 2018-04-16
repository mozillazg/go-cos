# go-cos

腾讯云对象存储服务 COS(Cloud Object Storage) Go SDK（API 版本：V5 版本的 XML API）。

[![Build Status](https://img.shields.io/travis/mozillazg/go-cos/master.svg)](https://travis-ci.org/mozillazg/go-cos)
[![Coverage Status](https://img.shields.io/coveralls/mozillazg/go-cos/master.svg)](https://coveralls.io/r/mozillazg/go-cos?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/mozillazg/go-cos)](https://goreportcard.com/report/github.com/mozillazg/go-cos)
[![GoDoc](https://godoc.org/github.com/mozillazg/go-cos?status.svg)](https://godoc.org/github.com/mozillazg/go-cos)

## Install

`go get -u github.com/mozillazg/go-cos`


## Usage

```go
package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/mozillazg/go-cos"
)

func main() {
	//将<bucket>和<region>修改为真实的信息
	//bucket的命名规则为{name}-{appid} ，此处填写的存储桶名称必须为此格式
	u, _ := url.Parse("https://<bucket>.cos.<region>.myqcloud.com")
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		//设置超时时间
		Timeout: 100 * time.Second,
		Transport: &cos.AuthorizationTransport{
			//如实填写账号和密钥，也可以设置为环境变量
			SecretID:  os.Getenv("COS_SECRETID"),
			SecretKey: os.Getenv("COS_SECRETKEY"),
		},
	})

	name := "test/hello.txt"
	resp, err := c.Object.Get(context.Background(), name, nil)
	if err != nil {
		panic(err)
	}
	bs, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	fmt.Printf("%s\n", string(bs))
}
```

所有的 API 在 [_example](./_example/) 目录下都有对应的使用示例。

## TODO

Service API:

* [x] Get Service（使用示例：[service/get.go](./_example/service/get.go)）

Bucket API:

* [x] Get Bucket（使用示例：[bucket/get.go](./_example/bucket/get.go)）
* [x] Get Bucket ACL（使用示例：[bucket/getACL.go](./_example/bucket/getACL.go)）
* [x] Get Bucket CORS（使用示例：[bucket/getCORS.go](./_example/bucket/getCORS.go)）
* [x] Get Bucket Location（使用示例：[bucket/getLocation.go](./_example/bucket/getLocation.go)）
* [x] Get Buket Lifecycle（使用示例：[bucket/getLifecycle.go](./_example/bucket/getLifecycle.go)）
* [x] Get Bucket Tagging（使用示例：[bucket/getTagging.go](./_example/bucket/getTagging.go)）
* [x] Put Bucket（使用示例：[bucket/put.go](./_example/bucket/put.go)）
* [x] Put Bucket ACL（使用示例：[bucket/putACL.go](./_example/bucket/putACL.go)）
* [x] Put Bucket CORS（使用示例：[bucket/putCORS.go](./_example/bucket/putCORS.go)）
* [x] Put Bucket Lifecycle（使用示例：[bucket/putLifecycle.go](./_example/bucket/putLifecycle.go)）
* [x] Put Bucket Tagging（使用示例：[bucket/putTagging.go](./_example/bucket/putTagging.go)）
* [x] Delete Bucket（使用示例：[bucket/delete.go](./_example/bucket/delete.go)）
* [x] Delete Bucket CORS（使用示例：[bucket/deleteCORS.go](./_example/bucket/deleteCORS.go)）
* [x] Delete Bucket Lifecycle（使用示例：[bucket/deleteLifecycle.go](./_example/bucket/deleteLifecycle.go)）
* [x] Delete Bucket Tagging（使用示例：[bucket/deleteTagging.go](./_example/bucket/deleteTagging.go)）
* [x] Head Bucket（使用示例：[bucket/head.go](./_example/bucket/head.go)）
* [x] List Multipart Uploads（使用示例：[bucket/listMultipartUploads.go](./_example/bucket/listMultipartUploads.go)）

Object API:

* [x] Append Object（使用示例：[object/append.go](./_example/object/append.go)）
* [x] Get Object（使用示例：[object/get.go](./_example/object/get.go)）
* [x] Get Object ACL（使用示例：[object/getACL.go](./_example/object/getACL.go)）
* [x] Put Object（使用示例：[object/put.go](./_example/object/put.go)）
* [x] Put Object ACL（使用示例：[object/putACL.go](./_example/object/putACL.go)）
* [x] Put Object Copy（使用示例：[object/copy.go](./_example/object/copy.go)）
* [x] Delete Object（使用示例：[object/delete.go](./_example/object/delete.go)）
* [x] Delete Multiple Object（使用示例：[object/deleteMultiple.go](./_example/object/deleteMultiple.go)）
* [x] Head Object（使用示例：[object/head.go](./_example/object/head.go)）
* [x] Options Object（使用示例：[object/options.go](./_example/object/options.go)）
* [x] Initiate Multipart Upload（使用示例：[object/initiateMultipartUpload.go](./_example/object/initiateMultipartUpload.go)）
* [x] Upload Part（使用示例：[object/uploadPart.go](./_example/object/uploadPart.go)）
* [x] List Parts（使用示例：[object/listParts.go](./_example/object/listParts.go)）
* [x] Complete Multipart Upload（使用示例：[object/completeMultipartUpload.go](./_example/object/completeMultipartUpload.go)）
* [x] Abort Multipart Upload（使用示例：[object/abortMultipartUpload.go](./_example/object/abortMultipartUpload.go)）
* [x] Mutipart Upload（使用示例：[object/MutiUpload.go.go](./_example/object/MutiUpload.go)）
