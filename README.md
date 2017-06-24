# go-cos

腾讯云对象存储服务 COS(Cloud Object Storage) Go SDK（API 版本：V4 版本的 XML API）。

[![Build Status](https://img.shields.io/travis/mozillazg/go-cos/master.svg)](https://travis-ci.org/mozillazg/go-cos)
[![Coverage Status](https://img.shields.io/coveralls/mozillazg/go-cos/master.svg)](https://coveralls.io/r/mozillazg/go-cos?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/mozillazg/go-cos)](https://goreportcard.com/report/github.com/mozillazg/go-cos)
[![GoDoc](https://godoc.org/github.com/mozillazg/go-cos?status.svg)](https://godoc.org/github.com/mozillazg/go-cos)

## install

`go get -u github.com/mozillazg/go-cos`


## usage

所有的 API 在 [examples](./examples/) 目录下都有对应的使用示例。

## TODO

Service API:

* [x] Get Service

Bucket API:

* [x] Get Bucket
* [x] Get Bucket ACL
* [x] Get Bucket CORS
* [x] Get Bucket Location
* [x] Get Buket Lifecycle
* [x] Get Bucket Tagging
* [x] Put Bucket
* [x] Put Bucket ACL
* [x] Put Bucket CORS
* [x] Put Bucket Lifecycle
* [x] Put Bucket Tagging
* [x] Delete Bucket
* [x] Delete Bucket CORS
* [x] Delete Bucket Lifecycle
* [x] Delete Bucket Tagging
* [x] Head Bucket
* [x] List Multipart Uploads

Object API:

* [x] Append Object
* [x] Get Object
* [x] Get Object ACL
* [x] Put Object
* [x] Put Object ACL
* [x] Delete Object
* [x] Delete Multiple Object
* [x] Head Object
* [x] Options Object
* [x] Initiate Multipart Upload
* [x] Upload Part
* [x] List Parts
* [x] Complete Multipart Upload
* [x] Abort Multipart Upload
