package cos

import (
	"context"
	"encoding/xml"
	"net/http"
)

// BucketLifecycleExpiration ...
type BucketLifecycleExpiration struct {
	Date string `xml:"Date,omitempty"`
	Days int    `xml:"Days,omitempty"`
}

// BucketLifecycleTransition ...
type BucketLifecycleTransition struct {
	Date         string `xml:"Date,omitempty"`
	Days         int    `xml:"Days,omitempty"`
	StorageClass string
}

// BucketLifecycleAbortIncompleteMultipartUpload ...
type BucketLifecycleAbortIncompleteMultipartUpload struct {
	DaysAfterInititation string `xml:"DaysAfterInititation,omitempty"`
}

// BucketLifecycleRule ...
type BucketLifecycleRule struct {
	ID                             string `xml:"ID,omitempty"`
	Prefix                         string
	Status                         string
	Transition                     *BucketLifecycleTransition                     `xml:"Transition,omitempty"`
	Expiration                     *BucketLifecycleExpiration                     `xml:"Expiration,omitempty"`
	AbortIncompleteMultipartUpload *BucketLifecycleAbortIncompleteMultipartUpload `xml:"AbortIncompleteMultipartUpload,omitempty"`
}

// BucketGetLifecycleResult ...
type BucketGetLifecycleResult struct {
	XMLName xml.Name               `xml:"LifecycleConfiguration"`
	Rules   []*BucketLifecycleRule `xml:"Rule,omitempty"`
}

// GetLifecycle Get Bucket Lifecycle请求实现读取生命周期管理的配置。当配置不存在时，返回404 Not Found。
//
// （目前只支持华南园区）
//
// https://www.qcloud.com/document/product/436/8278
func (s *BucketService) GetLifecycle(ctx context.Context,
	authTime *AuthTime) (*BucketGetLifecycleResult, *Response, error) {
	u := "/?lifecycle"
	baseURL := s.client.BaseURL.BucketURL
	var res BucketGetLifecycleResult
	resp, err := s.client.sendNoBody(ctx, baseURL, u, http.MethodGet, authTime, nil, nil, &res)
	return &res, resp, err
}

type BucketPutLifecycleOptions struct {
	XMLName xml.Name               `xml:"LifecycleConfiguration"`
	Rules   []*BucketLifecycleRule `xml:"Rule,omitempty"`
}

// PutLifecycle Put Bucket Lifecycle请求实现设置生命周期管理的功能。您可以通过该请求实现数据的生命周期管理配置和定期删除。
//
// 此请求为覆盖操作，上传新的配置文件将覆盖之前的配置文件。生命周期管理对文件和文件夹同时生效。
//
// （目前只支持华南园区）
//
// https://www.qcloud.com/document/product/436/8280
func (s *BucketService) PutLifecycle(ctx context.Context,
	authTime *AuthTime, opt *BucketPutLifecycleOptions) (*Response, error) {
	u := "/?lifecycle"
	baseURL := s.client.BaseURL.BucketURL
	resp, err := s.client.sendWithBody(ctx, baseURL, u, http.MethodPut, authTime, opt, nil, nil, nil)
	return resp, err
}

// DeleteLifecycle Delete Bucket Lifecycle请求实现删除生命周期管理。
//
// （目前只支持华南园区）
//
// https://www.qcloud.com/document/product/436/8284
func (s *BucketService) DeleteLifecycle(ctx context.Context,
	authTime *AuthTime) (*Response, error) {
	u := "/?lifecycle"
	baseURL := s.client.BaseURL.BucketURL
	resp, err := s.client.sendNoBody(ctx, baseURL, u, http.MethodDelete, authTime, nil, nil, nil)
	return resp, err
}
