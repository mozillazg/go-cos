package cos

import (
	"context"
	"net/http"
)

// ObjectPutACLOptions ...
type ObjectRestoreOptions struct {
		contentMD5        string `url:"response-content-type,omitempty" header:"-"`
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
func (s *ObjectService) Restore(ctx context.Context, name string, opt *ObjectRestoreOptions) (*Response, error) {
	sendOpt := sendOptions{
		baseURL:          s.client.BaseURL.BucketURL,
		uri:              "/" + encodeURIComponent(name) + "?restore",
		method:           http.MethodPost,
		optQuery:         opt,
		optHeader:        opt,
		disableCloseBody: true,
	}
	resp, err := s.client.send(ctx, &sendOpt)
	return resp, err
}
