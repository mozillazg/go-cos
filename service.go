package cos

import (
	"context"
	"encoding/xml"
	"net/http"
)

// ServiceService ...
//
// Service 相关 API
type ServiceService service

// ServiceGetResult ...
type ServiceGetResult struct {
	XMLName xml.Name `xml:"ListAllMyBucketsResult"`
	Owner   *struct {
		UIN string `xml:"uin"`
	} `xml:"Owner"`
	Buckets []*struct {
		Name       string
		Location   string
		CreateDate string
	} `xml:"Buckets>Bucket,omitempty"`
}

// Get Service 接口实现获取该用户下所有Bucket列表。
//
// 该API接口需要使用Authorization签名认证，
// 且只能获取签名中AccessID所属账户的Bucket列表。
//
// https://www.qcloud.com/document/product/436/8291
func (s *ServiceService) Get(ctx context.Context, authTime *AuthTime) (*ServiceGetResult, *Response, error) {
	u := "/"
	baseURL := s.client.BaseURL.ServiceURL
	var res ServiceGetResult
	resp, err := s.client.sendNoBody(ctx, baseURL, u, http.MethodGet, authTime, nil, nil, &res)
	return &res, resp, err
}
