package cos

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
)

// HostService 指定获取 Service 信息的域名
var HostService = "service.cos.myqcloud.com"

// ServiceService ...
type ServiceService service

// ServiceBucket ...
type ServiceBucket struct {
	Name       string
	Location   string
	CreateDate string
}

// Service ...
type Service struct {
	Owner   Owner
	Buckets []ServiceBucket `xml:"Buckets>Bucket,omitempty"`
}

// ServiceResult ...
type ServiceResult struct {
	XMLListAllMyBucketsResult xml.Name `xml:"ListAllMyBucketsResult"`
	Service
}

// Get Service 接口实现获取该用户下所有Bucket列表。
//
// 该API接口需要使用Authorization签名认证，
// 且只能获取签名中AccessID所属账户的Bucket列表。
//
// https://www.qcloud.com/document/product/436/8291
func (s *ServiceService) Get(ctx context.Context, authTime AuthTime) (service *Service, resp *Response, err error) {
	var res ServiceResult
	u := "/"
	baseURL := getServiceBaseURL(s.client.Secure)
	resp, err = s.client.sendNoBody(ctx, baseURL, u, http.MethodGet, authTime, nil, nil, &res)
	if err != nil {
		return
	}
	service = &res.Service
	return
}

func getServiceBaseURL(secure bool) *url.URL {
	scheme := "http"
	if secure {
		scheme = "https"
	}
	u, _ := url.Parse(fmt.Sprintf("%s://%s", scheme, HostService))
	return u
}
