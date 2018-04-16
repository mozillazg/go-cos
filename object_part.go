package cos

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
)

// InitiateMultipartUploadOptions ...
type InitiateMultipartUploadOptions struct {
	*ACLHeaderOptions
	*ObjectPutHeaderOptions
}

// InitiateMultipartUploadResult ...
type InitiateMultipartUploadResult struct {
	XMLName  xml.Name `xml:"InitiateMultipartUploadResult"`
	Bucket   string
	Key      string
	UploadID string `xml:"UploadId"`
}

// InitiateMultipartUpload ...
//
// Initiate Multipart Upload请求实现初始化分片上传，成功执行此请求以后会返回Upload ID用于后续的Upload Part请求。
//
// https://www.qcloud.com/document/product/436/7746
func (s *ObjectService) InitiateMultipartUpload(ctx context.Context, name string, opt *InitiateMultipartUploadOptions) (*InitiateMultipartUploadResult, *Response, error) {
	var res InitiateMultipartUploadResult
	sendOpt := sendOptions{
		baseURL:   s.client.BaseURL.BucketURL,
		uri:       "/" + encodeURIComponent(name) + "?uploads",
		method:    http.MethodPost,
		optHeader: opt,
		result:    &res,
	}
	resp, err := s.client.send(ctx, &sendOpt)
	return &res, resp, err
}

// ObjectUploadPartOptions ...
type ObjectUploadPartOptions struct {
	Expect          string `header:"Expect,omitempty" url:"-"`
	XCosContentMD5 string `header:"x-cos-content-md5" url:"-"`
	ContentLength   int    `header:"Content-Length,omitempty" url:"-"`
}


// UploadPart ...
//
// Upload Part请求实现在初始化以后的分块上传，支持的块的数量为1到10000，块的大小为1 MB 到5 GB。
// 在每次请求Upload Part时候，需要携带partNumber和uploadID，partNumber为块的编号，支持乱序上传。
//
// 当传入uploadID和partNumber都相同的时候，后传入的块将覆盖之前传入的块。当uploadID不存在时会返回404错误，NoSuchUpload.
//
// 当 r 不是 bytes.Buffer/bytes.Reader/strings.Reader 时，必须指定 opt.ContentLength
//
// https://www.qcloud.com/document/product/436/7750
func (s *ObjectService) UploadPart(ctx context.Context, name, uploadID string, partNumber int, r io.Reader, opt *ObjectUploadPartOptions) (*Response, error) {
	u := fmt.Sprintf("/%s?partNumber=%d&uploadId=%s", encodeURIComponent(name), partNumber, uploadID)
	sendOpt := sendOptions{
		baseURL:   s.client.BaseURL.BucketURL,
		uri:       u,
		method:    http.MethodPut,
		optHeader: opt,
		body:      r,
	}
	resp, err := s.client.send(ctx, &sendOpt)
	return resp, err
}

// ObjectListPartsOptions ...
type ObjectListPartsOptions struct {
	EncodingType     string `url:"Encoding-type,omitempty"`
	MaxParts         int    `url:"max-parts,omitempty"`
	PartNumberMarker int    `url:"part-number-marker,omitempty"`
}

// ObjectListPartsResult ...
type ObjectListPartsResult struct {
	XMLName              xml.Name `xml:"ListPartsResult"`
	Bucket               string
	EncodingType         string `xml:"Encoding-type,omitempty"`
	Key                  string
	UploadID             string     `xml:"UploadId"`
	Initiator            *Initiator `xml:"Initiator,omitempty"`
	Owner                *Owner     `xml:"Owner,omitempty"`
	StorageClass         string
	PartNumberMarker     int
	NextPartNumberMarker int `xml:"NextPartNumberMarker,omitempty"`
	MaxParts             int
	IsTruncated          bool
	Parts                []Object `xml:"Part,omitempty"`
}

// ListParts ...
//
// List Parts用来查询特定分块上传中的已上传的块。
//
// https://www.qcloud.com/document/product/436/7747
func (s *ObjectService) ListParts(ctx context.Context, name, uploadID string) (*ObjectListPartsResult, *Response, error) {
	u := fmt.Sprintf("/%s?uploadId=%s", encodeURIComponent(name), uploadID)
	var res ObjectListPartsResult
	sendOpt := sendOptions{
		baseURL: s.client.BaseURL.BucketURL,
		uri:     u,
		method:  http.MethodGet,
		result:  &res,
	}
	resp, err := s.client.send(ctx, &sendOpt)
	return &res, resp, err
}

// CompleteMultipartUploadOptions ...
type CompleteMultipartUploadOptions struct {
	XMLName xml.Name `xml:"CompleteMultipartUpload"`
	Parts   []Object `xml:"Part"`
}

// CompleteMultipartUploadResult ...
type CompleteMultipartUploadResult struct {
	XMLName  xml.Name `xml:"CompleteMultipartUploadResult"`
	Location string
	Bucket   string
	Key      string
	ETag     string
}

// CompleteMultipartUpload ...
//
// Complete Multipart Upload用来实现完成整个分块上传。当您已经使用Upload Parts上传所有块以后，你可以用该API完成上传。
// 在使用该API时，您必须在Body中给出每一个块的PartNumber和ETag，用来校验块的准确性。
//
// 由于分块上传的合并需要数分钟时间，因而当合并分块开始的时候，COS就立即返回200的状态码，在合并的过程中，
// COS会周期性的返回空格信息来保持连接活跃，直到合并完成，COS会在Body中返回合并后块的内容。
//
// 当上传块小于1 MB的时候，在调用该请求时，会返回400 EntityTooSmall；
// 当上传块编号不连续的时候，在调用该请求时，会返回400 InvalidPart；
// 当请求Body中的块信息没有按序号从小到大排列的时候，在调用该请求时，会返回400 InvalidPartOrder；
// 当UploadId不存在的时候，在调用该请求时，会返回404 NoSuchUpload。
//
// 建议您及时完成分块上传或者舍弃分块上传，因为已上传但是未终止的块会占用存储空间进而产生存储费用。
//
// https://www.qcloud.com/document/product/436/7742
func (s *ObjectService) CompleteMultipartUpload(ctx context.Context, name, uploadID string, opt *CompleteMultipartUploadOptions) (*CompleteMultipartUploadResult, *Response, error) {
	u := fmt.Sprintf("/%s?uploadId=%s", encodeURIComponent(name), uploadID)
	var res CompleteMultipartUploadResult
	sendOpt := sendOptions{
		baseURL: s.client.BaseURL.BucketURL,
		uri:     u,
		method:  http.MethodPost,
		body:    opt,
		result:  &res,
	}
	resp, err := s.client.send(ctx, &sendOpt)
	return &res, resp, err
}

// AbortMultipartUpload ...
//
// Abort Multipart Upload用来实现舍弃一个分块上传并删除已上传的块。当您调用Abort Multipart Upload时，
// 如果有正在使用这个Upload Parts上传块的请求，则Upload Parts会返回失败。当该UploadID不存在时，会返回404 NoSuchUpload。
//
// 建议您及时完成分块上传或者舍弃分块上传，因为已上传但是未终止的块会占用存储空间进而产生存储费用。
//
// https://www.qcloud.com/document/product/436/7740
func (s *ObjectService) AbortMultipartUpload(ctx context.Context, name, uploadID string) (*Response, error) {
	u := fmt.Sprintf("/%s?uploadId=%s", encodeURIComponent(name), uploadID)
	sendOpt := sendOptions{
		baseURL: s.client.BaseURL.BucketURL,
		uri:     u,
		method:  http.MethodDelete,
	}
	resp, err := s.client.send(ctx, &sendOpt)
	return resp, err
}
