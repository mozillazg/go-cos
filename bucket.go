package cos

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"time"
)

// HostBucket 指定默认的域名结构
var HostBucket = "%s-%s.%s.myqcloud.com"

// ReBucketDomain 匹配默认的域名结构
var ReBucketDomain = regexp.MustCompile("^(.+)-(\\d+)\\.([\\w-]+)\\.myqcloud\\.com$")

// BucketService ...
type BucketService service

// Bucket ...
type Bucket struct {
	domain string
	Name   string
	AppID  string
	Region string
}

// NewBucket ...
func NewBucket(name, appID, region string) *Bucket {
	return &Bucket{
		domain: fmt.Sprintf(HostBucket, name, appID, region),
		Name:   name,
		AppID:  appID,
		Region: region,
	}
}

// ParseBucketFromDomain 从域名中解析信息然后生成一个 Bucket
func ParseBucketFromDomain(domain string) (b *Bucket, err error) {
	matched := ReBucketDomain.FindStringSubmatch(domain)
	if len(matched) != 4 {
		err = errors.New("invalid bucket domain")
		return
	}
	b = &Bucket{
		domain: domain,
		Name:   matched[1],
		AppID:  matched[2],
		Region: matched[3],
	}
	return
}

// GetBaseURL 获取 Bucket 的基础请求 URL
func (b *Bucket) GetBaseURL(secure bool) string {
	scheme := "http"
	if secure {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s", scheme, b.domain)
}

func (s *BucketService) SetBucket(b *Bucket) {
	s.bucket = b
}

// BucketOwner ...
type BucketOwner struct {
	ID string
}

// ObjectMeta ...
type ObjectMeta struct {
	Key          string
	LastModified string
	ETag         string
	Size         string
	Owner        BucketOwner
	StorageClass string
}

// ListBucket ...
type ListBucket struct {
	Name           string
	Prefix         string `xml:"Prefix,omitempty"`
	Marker         string `xml:"Marker,omitempty"`
	NextMarker     string `xml:"NextMarker,omitempty"`
	Delimiter      string `xml:"Delimiter,omitempty"`
	MaxKeys        int64
	IsTruncated    bool
	Contents       []ObjectMeta `xml:"Contents,omitempty"`
	CommonPrefixes []string     `xml:"CommonPrefixes>Prefix,omitempty"`
	EncodingType   string       `xml:"Encoding-Type,omitempty"`
}

// ListBucketResult ...
type ListBucketResult struct {
	XMLName xml.Name `xml:"ListBucketResult"`
	ListBucket
}

// BucketGetOptions ...
type BucketGetOptions struct {
	Prefix       string `url:"prefix,omitempty"`
	Delimiter    string `url:"delimiter,omitempty"`
	EncodingType string `url:"encoding-type,omitempty"`
	Marker       string `url:"marker,omitempty"`
	MaxKeys      int64  `url:"max-keys,omitempty"`
}

// Get Bucket请求等同于 List Object请求，可以列出该Bucket下部分或者所有Object，发起该请求需要拥有Read权限。
// https://www.qcloud.com/document/product/436/7734
func (s *BucketService) Get(ctx context.Context,
	signStartTime, signEndTime,
	keyStartTime, keyEndTime time.Time, opt *BucketGetOptions) (listBucket *ListBucket,
	resp *http.Response, err error) {
	u := "/"
	baseURL := s.bucket.GetBaseURL(s.client.Secure)
	var res ListBucketResult
	resp, err = s.client.sendNoBody(ctx, u, http.MethodGet, baseURL, signStartTime, signEndTime,
		keyStartTime, keyEndTime, opt, nil, &res)
	if err != nil {
		return
	}
	listBucket = &res.ListBucket
	return
}

// BucketACLGrantee ...
type BucketACLGrantee struct {
	Type       string `xml:"type,attr"`
	UIN        string `xml:"uin"`
	Subaccount string `xml:"Subaccount,omitempty"`
}

// BucketACLGrant ...
type BucketACLGrant struct {
	Grantee    BucketACLGrantee
	Permission string
}

// BucketACL ...
type BucketACL struct {
	Owner             Owner
	AccessControlList []BucketACLGrant `xml:"AccessControlList>Grant,omitempty"`
}

// BucketACLResult ...
type BucketACLResult struct {
	XMLName xml.Name `xml:"AccessControlPolicy"`
	BucketACL
}

// GetACL 使用API读取Bucket的ACL表，只有所有者有权操作。
// https://www.qcloud.com/document/product/436/7733
func (s *BucketService) GetACL(ctx context.Context,
	signStartTime, signEndTime,
	keyStartTime, keyEndTime time.Time) (acl *BucketACL, resp *http.Response, err error) {
	u := "/?acl"
	baseURL := s.bucket.GetBaseURL(s.client.Secure)
	var res BucketACLResult
	resp, err = s.client.sendNoBody(ctx, u, http.MethodGet, baseURL, signStartTime, signEndTime,
		keyStartTime, keyEndTime, nil, nil, &res)
	if err != nil {
		return
	}
	acl = &res.BucketACL
	return
}

// BucketCORSRule ...
type BucketCORSRule struct {
	ID            string `xml:"ID,omitempty"`
	AllowedMethod string
	AllowedOrigin string
	AllowedHeader string `xml:"AllowedHeader,omitempty"`
	MaxAgeSeconds int64  `xml:"MaxAgeSeconds,omitempty"`
	ExposeHeader  string `xml:"ExposeHeader,omitempty"`
}

// BucketCORSResult ...
type BucketCORSResult struct {
	XMLName xml.Name         `xml:"CORSConfiguration"`
	Rules   []BucketCORSRule `xml:"CORSRule,omitempty"`
}

// GetCORS Get Bucket CORS实现跨域访问配置读取。
// https://www.qcloud.com/document/product/436/8274
func (s *BucketService) GetCORS(ctx context.Context,
	signStartTime, signEndTime,
	keyStartTime, keyEndTime time.Time) (cors *[]BucketCORSRule, resp *http.Response, err error) {
	u := "/?cors"
	baseURL := s.bucket.GetBaseURL(s.client.Secure)
	var res BucketCORSResult
	resp, err = s.client.sendNoBody(ctx, u, http.MethodGet, baseURL, signStartTime, signEndTime,
		keyStartTime, keyEndTime, nil, nil, &res)
	if err != nil {
		return
	}
	cors = &res.Rules
	return
}

// BucketLocation ...
type BucketLocation struct {
	Location string `xml:",chardata"`
}

// BucketLocationResult ...
type BucketLocationResult struct {
	XMLName xml.Name `xml:"LocationConstraint"`
	BucketLocation
}

// GetLocation Get Bucket Location接口获取Bucket所在地域信息，只有Bucket所有者有权限读取信息。
// https://www.qcloud.com/document/product/436/8275
func (s *BucketService) GetLocation(ctx context.Context,
	signStartTime, signEndTime,
	keyStartTime, keyEndTime time.Time) (location *BucketLocation, resp *http.Response, err error) {
	u := "/?location"
	baseURL := s.bucket.GetBaseURL(s.client.Secure)
	var res BucketLocationResult
	resp, err = s.client.sendNoBody(ctx, u, http.MethodGet, baseURL, signStartTime, signEndTime,
		keyStartTime, keyEndTime, nil, nil, &res)
	if err != nil {
		return
	}
	location = &res.BucketLocation
	return
}

// BucketLifecycleExpiration ...
type BucketLifecycleExpiration struct {
	Date string `xml:"Date,omitempty"`
	Days int64  `xml:"Days,omitempty"`
}

// BucketLifecycleTransition ...
type BucketLifecycleTransition struct {
	Date         string `xml:"Date,omitempty"`
	Days         int64  `xml:"Days,omitempty"`
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
	Transition                     BucketLifecycleTransition                     `xml:"Transition,omitempty"`
	Expiration                     BucketLifecycleExpiration                     `xml:"Expiration,omitempty"`
	AbortIncompleteMultipartUpload BucketLifecycleAbortIncompleteMultipartUpload `xml:"AbortIncompleteMultipartUpload,omitempty"`
}

// BucketLifecycleResult ...
type BucketLifecycleResult struct {
	XMLName xml.Name              `xml:"LifecycleConfiguration"`
	Rules   []BucketLifecycleRule `xml:"Rule,omitempty"`
}

// GetLifecycle Get Bucket Lifecycle请求实现读取生命周期管理的配置。当配置不存在时，返回404 Not Found。
// （目前只支持华南园区）
// https://www.qcloud.com/document/product/436/8278
func (s *BucketService) GetLifecycle(ctx context.Context,
	signStartTime, signEndTime,
	keyStartTime, keyEndTime time.Time) (lc *[]BucketLifecycleRule, resp *http.Response, err error) {
	u := "/?lifecycle"
	baseURL := s.bucket.GetBaseURL(s.client.Secure)
	var res BucketLifecycleResult
	resp, err = s.client.sendNoBody(ctx, u, http.MethodGet, baseURL, signStartTime, signEndTime,
		keyStartTime, keyEndTime, nil, nil, &res)
	if err != nil {
		return
	}
	lc = &res.Rules
	return
}

// BucketTaggingTag ...
type BucketTaggingTag struct {
	Key   string
	Value string
}

// BucketTaggingResult ...
type BucketTaggingResult struct {
	XMLName xml.Name           `xml:"Tagging"`
	TagSet  []BucketTaggingTag `xml:"TagSet>Tag,omitempty"`
}

// GetTagging Get Bucket Tagging接口实现获取指定Bucket的标签。
// https://www.qcloud.com/document/product/436/8277
func (s *BucketService) GetTagging(ctx context.Context,
	signStartTime, signEndTime,
	keyStartTime, keyEndTime time.Time) (tgs *[]BucketTaggingTag, resp *http.Response, err error) {
	u := "/?tagging"
	baseURL := s.bucket.GetBaseURL(s.client.Secure)
	var res BucketTaggingResult
	resp, err = s.client.sendNoBody(ctx, u, http.MethodGet, baseURL, signStartTime, signEndTime,
		keyStartTime, keyEndTime, nil, nil, &res)
	if err != nil {
		return
	}
	tgs = &res.TagSet
	return
}

// BucketPutOptions ...
type BucketPutOptions struct {
	XCosACL              string `header:"x-cos-acl,omitempty"`
	XCosGrantRead        string `header:"x-cos-grant-read,omitempty"`
	XCosGrantWrite       string `header:"x-cos-grant-write,omitempty"`
	XCosGrantFullControl string `header:"x-cos-grant-full-control,omitempty"`
}

// Put Bucket请求可以在指定账号下创建一个Bucket。
// https://www.qcloud.com/document/product/436/7738
func (s *BucketService) Put(ctx context.Context,
	signStartTime, signEndTime,
	keyStartTime, keyEndTime time.Time, opt *BucketPutOptions) (resp *http.Response, err error) {
	u := "/"
	baseURL := s.bucket.GetBaseURL(s.client.Secure)
	resp, err = s.client.sendWithBody(ctx, u, http.MethodPut, baseURL, signStartTime, signEndTime,
		keyStartTime, keyEndTime, nil, opt, nil, nil)
	return
}

// BucketPutACLOptions ...
type BucketPutACLOptions BucketPutOptions

// PutACL 使用API写入Bucket的ACL表，您可以通过Header："x-cos-acl","x-cos-grant-read",
// "x-cos-grant-write","x-cos-grant-full-control"传入ACL信息，也可以通过body以XML格式传入ACL信息，
// 但是只能选择Header和Body其中一种，否则返回冲突。
// Put Bucket ACL是一个覆盖操作，传入新的ACL将覆盖原有ACL。只有所有者有权操作。
// "x-cos-acl"：枚举值为public-read，private；public-read意味这个Bucket有公有读私有写的权限，
//  private意味这个Bucket有私有读写的权限。
// "x-cos-grant-read"：意味被赋予权限的用户拥有该Bucket的读权限
// "x-cos-grant-write"：意味被赋予权限的用户拥有该Bucket的写权限
// "x-cos-grant-full-control"：意味被赋予权限的用户拥有该Bucket的读写权限
// https://www.qcloud.com/document/product/436/7737
func (s *BucketService) PutACL(ctx context.Context,
	signStartTime, signEndTime,
	keyStartTime, keyEndTime time.Time,
	opt *BucketPutACLOptions, acl *BucketACLResult) (resp *http.Response, err error) {
	u := "/?acl"
	baseURL := s.bucket.GetBaseURL(s.client.Secure)
	resp, err = s.client.sendWithBody(ctx, u, http.MethodPut, baseURL, signStartTime, signEndTime,
		keyStartTime, keyEndTime, acl, nil, opt, nil)
	return
}

// PutCORS Put Bucket CORS实现跨域访问设置，您可以通过传入XML格式的配置文件实现配置，文件大小限制为64 KB。
// https://www.qcloud.com/document/product/436/8279
func (s *BucketService) PutCORS(ctx context.Context,
	signStartTime, signEndTime,
	keyStartTime, keyEndTime time.Time, cos *BucketCORSResult) (resp *http.Response, err error) {
	u := "/?cors"
	baseURL := s.bucket.GetBaseURL(s.client.Secure)
	resp, err = s.client.sendWithBody(ctx, u, http.MethodPut, baseURL, signStartTime, signEndTime,
		keyStartTime, keyEndTime, cos, nil, nil, nil)
	return
}

// PutLifecycle Put Bucket Lifecycle请求实现设置生命周期管理的功能。您可以通过该请求实现数据的生命周期管理配置和定期删除。
// 此请求为覆盖操作，上传新的配置文件将覆盖之前的配置文件。生命周期管理对文件和文件夹同时生效。
// （目前只支持华南园区）
// https://www.qcloud.com/document/product/436/8280
// TODO: fix doesn't work
func (s *BucketService) PutLifecycle(ctx context.Context,
	signStartTime, signEndTime,
	keyStartTime, keyEndTime time.Time, lc *BucketLifecycleResult) (resp *http.Response, err error) {
	u := "/?lifecycle"
	baseURL := s.bucket.GetBaseURL(s.client.Secure)
	resp, err = s.client.sendWithBody(ctx, u, http.MethodPut, baseURL, signStartTime, signEndTime,
		keyStartTime, keyEndTime, lc, nil, nil, nil)
	return
}

// PutTagging Put Bucket Tagging接口实现给用指定Bucket打标签。用来组织和管理相关Bucket。
// 当该请求设置相同Key名称，不同Value时，会返回400。请求成功，则返回204。
// https://www.qcloud.com/document/product/436/8281
func (s *BucketService) PutTagging(ctx context.Context,
	signStartTime, signEndTime,
	keyStartTime, keyEndTime time.Time, tg *BucketTaggingResult) (resp *http.Response, err error) {
	u := "/?tagging"
	baseURL := s.bucket.GetBaseURL(s.client.Secure)
	resp, err = s.client.sendWithBody(ctx, u, http.MethodPut, baseURL, signStartTime, signEndTime,
		keyStartTime, keyEndTime, tg, nil, nil, nil)
	return
}

// Delete Bucket请求可以在指定账号下删除Bucket，删除之前要求Bucket为空。
// https://www.qcloud.com/document/product/436/7732
func (s *BucketService) Delete(ctx context.Context,
	signStartTime, signEndTime,
	keyStartTime, keyEndTime time.Time) (resp *http.Response, err error) {
	u := "/"
	baseURL := s.bucket.GetBaseURL(s.client.Secure)
	resp, err = s.client.sendNoBody(ctx, u, http.MethodDelete, baseURL, signStartTime, signEndTime,
		keyStartTime, keyEndTime, nil, nil, nil)
	return
}

// DeleteCORS Delete Bucket CORS实现跨域访问配置删除。
// https://www.qcloud.com/document/product/436/8283
func (s *BucketService) DeleteCORS(ctx context.Context,
	signStartTime, signEndTime,
	keyStartTime, keyEndTime time.Time) (resp *http.Response, err error) {
	u := "/?cors"
	baseURL := s.bucket.GetBaseURL(s.client.Secure)
	resp, err = s.client.sendNoBody(ctx, u, http.MethodDelete, baseURL, signStartTime, signEndTime,
		keyStartTime, keyEndTime, nil, nil, nil)
	return
}

// DeleteLifecycle Delete Bucket Lifecycle请求实现删除生命周期管理。
// （目前只支持华南园区）
// https://www.qcloud.com/document/product/436/8284
func (s *BucketService) DeleteLifecycle(ctx context.Context,
	signStartTime, signEndTime,
	keyStartTime, keyEndTime time.Time) (resp *http.Response, err error) {
	u := "/?lifecycle"
	baseURL := s.bucket.GetBaseURL(s.client.Secure)
	resp, err = s.client.sendNoBody(ctx, u, http.MethodDelete, baseURL, signStartTime, signEndTime,
		keyStartTime, keyEndTime, nil, nil, nil)
	return
}

// DeleteTagging Delete Bucket Tagging接口实现删除指定Bucket的标签。
// https://www.qcloud.com/document/product/436/8286
func (s *BucketService) DeleteTagging(ctx context.Context,
	signStartTime, signEndTime,
	keyStartTime, keyEndTime time.Time) (resp *http.Response, err error) {
	u := "/?tagging"
	baseURL := s.bucket.GetBaseURL(s.client.Secure)
	resp, err = s.client.sendNoBody(ctx, u, http.MethodDelete, baseURL, signStartTime, signEndTime,
		keyStartTime, keyEndTime, nil, nil, nil)
	return
}

// Head Bucket请求可以确认是否存在该Bucket，是否有权限访问，Head的权限与Read一致。
// 当其存在时，返回 HTTP 状态码200；当无权限时，返回 HTTP 状态码403；
// 当不存在时，返回 HTTP 状态码404。
// https://www.qcloud.com/document/product/436/7735
func (s *BucketService) Head(ctx context.Context,
	signStartTime, signEndTime,
	keyStartTime, keyEndTime time.Time) (resp *http.Response, err error) {
	u := "/"
	baseURL := s.bucket.GetBaseURL(s.client.Secure)
	resp, err = s.client.sendNoBody(ctx, u, http.MethodHead, baseURL, signStartTime, signEndTime,
		keyStartTime, keyEndTime, nil, nil, nil)
	return
}

// MultipartUploadMeta ...
type MultipartUploadMeta struct {
	Key          string
	UploadID     string
	StorageClass string
	Initiator    Initiator
	Owner        Owner
	Initiated    string
}

// MultipartUploads ...
type MultipartUploads struct {
	Bucket             string `xml:"Bucket"`
	EncodingType       string `xml:"Encoding-Type"`
	KeyMarker          string
	UploadIDMarker     string `xml:"UploadIdMarker"`
	NextKeyMarker      string
	NextUploadIDMarker string `xml:"NextUploadIdMarker"`
	MaxUploads         int64
	IsTruncated        bool
	Uploads            []MultipartUploadMeta `xml:"Upload"`
	Prefix             string
	Delimiter          string   `xml:"delimiter,omitempty"`
	CommonPrefixs      []string `xml:"CommonPrefixs>Prefix,omitempty"`
}

// ListMultipartUploadsResult ...
type ListMultipartUploadsResult struct {
	XMLName xml.Name `xml:"ListMultipartUploadsResult"`
	MultipartUploads
}

// ListMultipartUploadsOptions ...
type ListMultipartUploadsOptions struct {
	Delimiter      string `url:"delimiter,omitempty"`
	EncodingType   string `url:"encoding-type,omitempty"`
	Prefix         string `url:",omitempty"`
	MaxUploads     int64  `url:"max-uploads,omitempty"`
	KeyMarker      string `url:"key-marker,omitempty"`
	UploadIDMarker string `url:"upload-id-marker,omitempty"`
}

// ListMultipartUploads List Multipart Uploads用来查询正在进行中的分块上传。单次最多列出1000个正在进行中的分块上传。
// https://www.qcloud.com/document/product/436/7736
func (s *BucketService) ListMultipartUploads(ctx context.Context,
	signStartTime, signEndTime,
	keyStartTime, keyEndTime time.Time,
	opt *ListMultipartUploadsOptions) (uploads *MultipartUploads, resp *http.Response, err error) {
	u := "/?uploads"
	baseURL := s.bucket.GetBaseURL(s.client.Secure)
	var res ListMultipartUploadsResult
	resp, err = s.client.sendNoBody(ctx, u, http.MethodGet, baseURL, signStartTime, signEndTime,
		keyStartTime, keyEndTime, opt, nil, &res)
	if err != nil {
		return
	}
	uploads = &res.MultipartUploads
	return
}
