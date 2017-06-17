package cos

import (
	"context"
	"encoding/xml"
	"net/http"
)

// ObjectGetACLResult ...
type ObjectGetACLResult struct {
	XMLName           xml.Name `xml:"AccessControlPolicy"`
	Owner             *Owner
	AccessControlList []*BucketACLGrant `xml:"AccessControlList>Grant,omitempty"`
}

// GetACL Get Object ACL接口实现使用API读取Object的ACL表，只有所有者有权操作。
//
// https://www.qcloud.com/document/product/436/7744
func (s *ObjectService) GetACL(ctx context.Context,
	authTime *AuthTime, name string) (*ObjectGetACLResult, *Response, error) {

	u := "/" + encodeURIComponent(name) + "?acl"
	baseURL := s.client.BaseURL.BucketURL
	var res ObjectGetACLResult
	resp, err := s.client.sendNoBody(ctx, baseURL, u, http.MethodGet, authTime, nil, nil, &res)
	return &res, resp, err
}

type ObjectPutACLOptions struct {
	Header *ACLHeaderOptions   `url:"-" xml:"-"`
	Body   *BucketGetACLResult `url:"-" header:"-"`
}

// PutACL 使用API写入Object的ACL表，您可以通过Header："x-cos-acl", "x-cos-grant-read" ,
// "x-cos-grant-write" ,"x-cos-grant-full-control"传入ACL信息，
// 也可以通过body以XML格式传入ACL信息，但是只能选择Header和Body其中一种，否则，返回冲突。
//
// Put Object ACL是一个覆盖操作，传入新的ACL将覆盖原有ACL。只有所有者有权操作。
//
// "x-cos-acl"：枚举值为public-read，private；public-read意味这个Object有公有读私有写的权限，
// private意味这个Object有私有读写的权限。
//
// "x-cos-grant-read"：意味被赋予权限的用户拥有该Object的读权限
//
// "x-cos-grant-write"：意味被赋予权限的用户拥有该Object的写权限
//
// "x-cos-grant-full-control"：意味被赋予权限的用户拥有该Object的读写权限
//
// https://www.qcloud.com/document/product/436/7748
func (s *ObjectService) PutACL(ctx context.Context,
	authTime *AuthTime, name string,
	opt *ObjectPutACLOptions) (*Response, error) {

	u := "/" + encodeURIComponent(name) + "?acl"
	baseURL := s.client.BaseURL.BucketURL
	header := opt.Header
	body := opt.Body
	if body != nil {
		header = nil
	}
	resp, err := s.client.sendWithBody(ctx, baseURL, u, http.MethodPut, authTime, body, nil, header, nil)
	return resp, err
}
