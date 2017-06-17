package cos

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
)

// ObjectService ...
type ObjectService service

// ObjectGetOptions ...
type ObjectGetOptions struct {
	ResponseContentType        string `url:"response-content-type,omitempty" header:"-"`
	ResponseContentLanguage    string `url:"response-content-language,omitempty" header:"-"`
	ResponseExpires            string `url:"response-expires,omitempty" header:"-"`
	ResponseCacheControl       string `url:"response-cache-control,omitempty" header:"-"`
	ResponseContentDisposition string `url:"response-content-disposition,omitempty" header:"-"`
	ResponseContentEncoding    string `url:"response-content-encoding,omitempty" header:"-"`
	Range                      string `url:"-" header:"Range,omitempty"`
	IfModifiedSince            string `url:"-" header:"If-Modified-Since,omitempty"`
}

// Get Object 请求可以将一个文件（Object）下载至本地。
// 该操作需要对目标 Object 具有读权限或目标 Object 对所有人都开放了读权限（公有读）。
//
// https://www.qcloud.com/document/product/436/7753
func (s *ObjectService) Get(ctx context.Context,
	authTime *AuthTime, name string, w io.Writer, opt *ObjectGetOptions) (*Response, error) {

	u := "/" + encodeURIComponent(name)
	baseURL := s.client.BaseURL.BucketURL
	resp, err := s.client.sendNoBody(ctx, baseURL, u, http.MethodGet, authTime, opt, opt, w)
	return resp, err
}

// ObjectPutHeaderOptions ...
type ObjectPutHeaderOptions struct {
	CacheControl       string `header:"Cache-Control,omitempty" url:"-"`
	ContentDisposition string `header:"Content-Disposition,omitempty" url:"-"`
	ContentEncoding    string `header:"Content-Encoding,omitempty" url:"-"`
	ContentType        string `header:"Content-Type,omitempty" url:"-"`
	Expect             string `header:"Expect,omitempty" url:"-"`
	Expires            string `header:"Expires,omitempty" url:"-"`
	XCosContentSHA1    string `header:"x-cos-content-sha1,omitempty" url:"-"`
	// 自定义的 x-cos-meta-* header
	XCosMetaXXX      *http.Header `header:"x-cos-meta-*,omitempty" url:"-"`
	XCosStorageClass string       `header:"x-cos-storage-class,omitempty" url:"-"`
	// 可选值: Normal, Appendable
	//XCosObjectType string `header:"x-cos-object-type,omitempty" url:"-"`
}

// ObjectPutOptions ...
type ObjectPutOptions struct {
	*ACLHeaderOptions       `header:",omitempty" url:"-" xml:"-"`
	*ObjectPutHeaderOptions `header:",omitempty" url:"-" xml:"-"`
}

// Put Object请求可以将一个文件（Oject）上传至指定Bucket。
//
// https://www.qcloud.com/document/product/436/7749
func (s *ObjectService) Put(ctx context.Context,
	authTime *AuthTime, name string, r io.Reader, opt *ObjectPutOptions) (*Response, error) {

	u := "/" + encodeURIComponent(name)
	baseURL := s.client.BaseURL.BucketURL
	resp, err := s.client.sendWithBody(ctx, baseURL, u, http.MethodPut, authTime, r, nil, opt, nil)
	return resp, err
}

// Delete Object请求可以将一个文件（Object）删除。
//
// https://www.qcloud.com/document/product/436/7743
func (s *ObjectService) Delete(ctx context.Context,
	authTime *AuthTime, name string) (*Response, error) {

	u := "/" + encodeURIComponent(name)
	baseURL := s.client.BaseURL.BucketURL
	resp, err := s.client.sendNoBody(ctx, baseURL, u, http.MethodDelete, authTime, nil, nil, nil)
	return resp, err
}

type ObjectHeadOptions struct {
	IfModifiedSince string `url:"-" header:"If-Modified-Since,omitempty"`
}

// Head Object请求可以取回对应Object的元数据，Head的权限与Get的权限一致
//
// https://www.qcloud.com/document/product/436/7745
func (s *ObjectService) Head(ctx context.Context,
	authTime *AuthTime, name string, opt *ObjectHeadOptions) (*Response, error) {

	u := "/" + encodeURIComponent(name)
	baseURL := s.client.BaseURL.BucketURL
	resp, err := s.client.sendNoBody(ctx, baseURL, u, http.MethodHead, authTime, nil, opt, nil)
	return resp, err
}

type ObjectOptionsOptions struct {
	Origin                      string `url:"-" header:"Origin"`
	AccessControlRequestMethod  string `url:"-" header:"Access-Control-Request-Method"`
	AccessControlRequestHeaders string `url:"-" header:"Access-Control-Request-Headers,omitempty"`
}

// Options Object请求实现跨域访问的预请求。即发出一个 OPTIONS 请求给服务器以确认是否可以进行跨域操作。
//
// 当CORS配置不存在时，请求返回403 Forbidden。
//
// https://www.qcloud.com/document/product/436/8288
func (s *ObjectService) Options(ctx context.Context,
	authTime *AuthTime, name string, opt *ObjectOptionsOptions) (*Response, error) {

	u := "/" + encodeURIComponent(name)
	baseURL := s.client.BaseURL.BucketURL
	resp, err := s.client.sendNoBody(ctx, baseURL, u, http.MethodOptions, authTime, nil, opt, nil)
	return resp, err
}

// Append
//
// Append请求可以将一个文件（Object）以分块追加的方式上传至 Bucket 中。使用Append Upload的文件必须事前被设定为Appendable。
// 当Appendable的文件被执行Put Object的操作以后，文件被覆盖，属性改变为Normal。
//
// 文件属性可以在Head Object操作中被查询到，当您发起Head Object请求时，会返回自定义Header『x-cos-object-type』，该Header只有两个枚举值：Normal或者Appendable。
//
// 追加上传建议文件大小1M - 5G。如果position的值和当前Object的长度不致，COS会返回409错误。
// 如果Append一个Normal的Object，COS会返回409 ObjectNotAppendable。
//
// Appendable的文件不可以被复制，不参与版本管理，不参与生命周期管理，不可跨区域复制。
//
// https://www.qcloud.com/document/product/436/7741
func (s *ObjectService) Append(ctx context.Context,
	authTime *AuthTime, name string, position int,
	r io.Reader, opt *ObjectPutOptions) (*Response, error) {

	u := fmt.Sprintf("/%s?append&position=%d", encodeURIComponent(name), position)
	baseURL := s.client.BaseURL.BucketURL
	resp, err := s.client.sendWithBody(ctx, baseURL, u, http.MethodPost, authTime, r, nil, opt, nil)
	return resp, err
}

type ObjectForDelete struct {
	Key string
}

type ObjectDeleteMultiOptions struct {
	XMLName xml.Name           `xml:"Delete" header:"-"`
	Quiet   bool               `xml:"Quiet" header:"-"`
	Objects []*ObjectForDelete `xml:"Object" header:"-"`
	//XCosSha1 string `xml:"-" header:"x-cos-sha1"`
}

type ObjectDeleteMultiResult struct {
	XMLName        xml.Name `xml:"DeleteResult"`
	DeletedObjects []*struct {
		Key string
	} `xml:"Deleted,omitempty"`
	Errors []*struct {
		Key     string
		Code    string
		Message string
	} `xml:"Error,omitempty"`
}

// DeleteMulti
//
// Delete Multiple Object请求实现批量删除文件，最大支持单次删除1000个文件。
// 对于返回结果，COS提供Verbose和Quiet两种结果模式。Verbose模式将返回每个Object的删除结果；
// Quiet模式只返回报错的Object信息。
//
// 此请求必须携带x-cos-sha1用来校验Body的完整性。
//
// https://www.qcloud.com/document/product/436/8289
func (s *ObjectService) DeleteMulti(ctx context.Context,
	authTime *AuthTime, opt *ObjectDeleteMultiOptions,
) (*ObjectDeleteMultiResult, *Response, error) {
	u := "/?delete"
	baseURL := s.client.BaseURL.BucketURL
	var res ObjectDeleteMultiResult
	resp, err := s.client.sendWithBody(ctx, baseURL, u, http.MethodPost, authTime, opt, nil, nil, &res)
	return &res, resp, err
}
