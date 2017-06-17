package cos

import (
	"context"
	"encoding/xml"
	"net/http"
)

// BucketTaggingTag ...
type BucketTaggingTag struct {
	Key   string
	Value string
}

// BucketGetTaggingResult ...
type BucketGetTaggingResult struct {
	XMLName xml.Name            `xml:"Tagging"`
	TagSet  []*BucketTaggingTag `xml:"TagSet>Tag,omitempty"`
}

// GetTagging ...
//
// Get Bucket Tagging接口实现获取指定Bucket的标签。
//
// https://www.qcloud.com/document/product/436/8277
func (s *BucketService) GetTagging(ctx context.Context,
	authTime *AuthTime) (*BucketGetTaggingResult, *Response, error) {
	u := "/?tagging"
	baseURL := s.client.BaseURL.BucketURL
	var res BucketGetTaggingResult
	resp, err := s.client.sendNoBody(ctx, baseURL, u, http.MethodGet, authTime, nil, nil, &res)
	return &res, resp, err
}

// BucketPutTaggingOptions ...
type BucketPutTaggingOptions struct {
	XMLName xml.Name            `xml:"Tagging"`
	TagSet  []*BucketTaggingTag `xml:"TagSet>Tag,omitempty"`
}

// PutTagging ...
//
// Put Bucket Tagging接口实现给用指定Bucket打标签。用来组织和管理相关Bucket。
//
// 当该请求设置相同Key名称，不同Value时，会返回400。请求成功，则返回204。
//
// https://www.qcloud.com/document/product/436/8281
func (s *BucketService) PutTagging(ctx context.Context,
	authTime *AuthTime, opt *BucketPutTaggingOptions) (*Response, error) {
	u := "/?tagging"
	baseURL := s.client.BaseURL.BucketURL
	resp, err := s.client.sendWithBody(ctx, baseURL, u, http.MethodPut, authTime, opt, nil, nil, nil)
	return resp, err
}

// DeleteTagging ...
//
// Delete Bucket Tagging接口实现删除指定Bucket的标签。
//
// https://www.qcloud.com/document/product/436/8286
func (s *BucketService) DeleteTagging(ctx context.Context,
	authTime *AuthTime) (*Response, error) {
	u := "/?tagging"
	baseURL := s.client.BaseURL.BucketURL
	resp, err := s.client.sendNoBody(ctx, baseURL, u, http.MethodDelete, authTime, nil, nil, nil)
	return resp, err
}
