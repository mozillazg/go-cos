package cos

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"time"

	"bitbucket.org/mozillazg/go-httpheader"
	"github.com/google/go-querystring/query"
)

const (
	// Version ...
	Version               = "0.2.0"
	userAgent             = "go-cos/" + Version
	contentTypeXML        = "application/xml"
	defaultServiceBaseURL = "https://service.cos.myqcloud.com"
	bucketURLFormat       = "%s://%s-%s.%s.myqcloud.com"
)

// BaseURL 访问各 API 所需的基础 URL
type BaseURL struct {
	// 访问 bucket, object 相关 API 的基础 URL（不包含 path 部分）: http://example.com
	BucketURL *url.URL
	// 访问 service API 的基础 URL（不包含 path 部分）: http://example.com
	ServiceURL *url.URL
}

// NewBucketURL 生成 BaseURL 所需的 BucketURL
//
//   bucketName: bucket 名称
//   AppID: 应用 ID
//   Region: 区域代码: cn-east, cn-south, cn-north
//   secure: 是否使用 https
func NewBucketURL(bucketName, AppID, Region string, secure bool) *url.URL {
	scheme := "https"
	if !secure {
		scheme = "http"
	}
	urlStr := fmt.Sprintf(bucketURLFormat, scheme, bucketName, AppID, Region)
	u, _ := url.Parse(urlStr)
	return u
}

// A Client manages communication with the COS API.
type Client struct {
	Client    *http.Client
	secretID  string
	secretKey string

	UserAgent string
	BaseURL   *BaseURL

	common service

	Service *ServiceService
	Bucket  *BucketService
	Object  *ObjectService
}

type service struct {
	client *Client
}

// NewClient returns a new COS API client.
func NewClient(secretID, secretKey string, uri *BaseURL, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	baseURL := &BaseURL{}
	if uri != nil {
		baseURL.BucketURL = uri.BucketURL
		baseURL.ServiceURL = uri.ServiceURL
	}
	if baseURL.ServiceURL == nil {
		baseURL.ServiceURL, _ = url.Parse(defaultServiceBaseURL)
	}

	c := &Client{
		Client:    httpClient,
		secretID:  secretID,
		secretKey: secretKey,
		UserAgent: userAgent,
		BaseURL:   baseURL,
	}
	c.common.client = c
	c.Service = (*ServiceService)(&c.common)
	c.Bucket = (*BucketService)(&c.common)
	c.Object = (*ObjectService)(&c.common)
	return c
}

func (c *Client) newRequest(ctx context.Context, baseURL *url.URL, uri, method string,
	body interface{}, optQuery interface{}, optHeader interface{}) (req *http.Request, err error) {
	uri, err = addURLOptions(uri, optQuery)
	if err != nil {
		return
	}
	u, _ := url.Parse(uri)
	urlStr := baseURL.ResolveReference(u).String()

	var reader io.Reader
	contentType := ""
	contentMD5 := ""
	contentLength := ""
	xsha1 := ""
	if body != nil {
		// 上传文件
		if r, ok := body.(io.Reader); ok {
			reader = r
		} else {
			b, err := xml.Marshal(body)
			if err != nil {
				return nil, err
			}
			contentType = contentTypeXML
			reader = bytes.NewReader(b)
			contentMD5 = base64.StdEncoding.EncodeToString(calMD5Digest(b))
			//xsha1 = base64.StdEncoding.EncodeToString(calSHA1Digest(b))
			contentLength = fmt.Sprintf("%d", len(b))
		}
	} else {
		contentType = contentTypeXML
	}

	req, err = http.NewRequest(method, urlStr, reader)
	if err != nil {
		return
	}

	req.Header, err = addHeaderOptions(req.Header, optHeader)
	if err != nil {
		return
	}

	if contentLength != "" {
		req.Header.Set("Content-Length", contentLength)
	}
	if contentMD5 != "" {
		req.Header["Content-MD5"] = []string{contentMD5}
	}
	if xsha1 != "" {
		req.Header.Set("x-cos-sha1", xsha1)
	}
	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}
	if req.Header.Get("Content-Type") == "" && contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	return
}

func (c *Client) doAPI(ctx context.Context, req *http.Request, ret interface{},
	authTime *AuthTime) (*Response, error) {
	req = req.WithContext(ctx)

	if authTime != nil {
		AddAuthorization(
			c.secretID, c.secretKey, req,
			authTime.signStartTime, authTime.signEndTime,
			authTime.keyStartTime, authTime.keyEndTime,
		)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		// If we got an error, and the context has been canceled,
		// the context's error is probably more useful.
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		return nil, err
	}

	defer func() {
		// Drain up to 512 bytes and close the body to let the Transport reuse the connection
		io.CopyN(ioutil.Discard, resp.Body, 512)
		resp.Body.Close()
	}()

	response := newResponse(resp)

	err = checkResponse(resp)
	if err != nil {
		// even though there was an error, we still return the response
		// in case the caller wants to inspect it further
		return response, err
	}

	if ret != nil {
		if w, ok := ret.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			err = xml.NewDecoder(resp.Body).Decode(ret)
			if err == io.EOF {
				err = nil // ignore EOF errors caused by empty response body
			}
		}
	}

	return response, err
}

func (c *Client) sendWithBody(ctx context.Context, baseURL *url.URL, uri, method string,
	authTime *AuthTime, body interface{},
	optQuery interface{}, optHeader interface{}, ret interface{}) (resp *Response, err error) {
	req, err := c.newRequest(ctx, baseURL, uri, method, body, optQuery, optHeader)
	if err != nil {
		return
	}

	resp, err = c.doAPI(ctx, req, ret, authTime)
	if err != nil {
		return
	}
	return
}

func (c *Client) sendNoBody(ctx context.Context, baseURL *url.URL, uri, method string,
	authTime *AuthTime,
	optQuery interface{}, optHeader interface{}, ret interface{}) (resp *Response, err error) {
	return c.sendWithBody(ctx, baseURL, uri, method, authTime, nil, optQuery, optHeader, ret)
}

// addURLOptions adds the parameters in opt as URL query parameters to s. opt
// must be a struct whose fields may contain "url" tags.
func addURLOptions(s string, opt interface{}) (string, error) {
	v := reflect.ValueOf(opt)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qs, err := query.Values(opt)
	if err != nil {
		return s, err
	}

	// 保留原有的参数，并且放在前面。因为 cos 的 url 路由是以第一个参数作为路由的
	// e.g. /?uploads
	q := u.RawQuery
	rq := qs.Encode()
	if q != "" {
		if rq != "" {
			u.RawQuery = fmt.Sprintf("%s&%s", q, qs.Encode())
		}
	} else {
		u.RawQuery = rq
	}
	return u.String(), nil
}

// addHeaderOptions adds the parameters in opt as Header fields to req. opt
// must be a struct whose fields may contain "header" tags.
func addHeaderOptions(header http.Header, opt interface{}) (http.Header, error) {
	v := reflect.ValueOf(opt)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return header, nil
	}

	h, err := httpheader.Header(opt)
	if err != nil {
		return nil, err
	}

	for key, values := range h {
		for _, value := range values {
			header.Add(key, value)
		}
	}
	return header, nil
}

// Owner ...
type Owner struct {
	UIN string `xml:"uin"`
}

// AuthTime 用于生成签名所需的 q-sign-time 和 q-key-time 相关参数
type AuthTime struct {
	signStartTime time.Time
	signEndTime   time.Time
	keyStartTime  time.Time
	keyEndTime    time.Time
}

// NewAuthTime ...
//
//   expire: 从现在开始多久过期.
func NewAuthTime(expire time.Duration) *AuthTime {
	signStartTime := time.Now()
	keyStartTime := signStartTime
	signEndTime := signStartTime.Add(expire)
	keyEndTime := signEndTime
	return &AuthTime{
		signStartTime, signEndTime,
		keyStartTime, keyEndTime,
	}
}

// Response API 响应
type Response struct {
	*http.Response
}

func newResponse(resp *http.Response) *Response {
	return &Response{
		Response: resp,
	}
}

// ACLHeaderOptions ...
type ACLHeaderOptions struct {
	XCosACL              string `header:"x-cos-acl,omitempty" url:"-" xml:"-"`
	XCosGrantRead        string `header:"x-cos-grant-read,omitempty" url:"-" xml:"-"`
	XCosGrantWrite       string `header:"x-cos-grant-write,omitempty" url:"-" xml:"-"`
	XCosGrantFullControl string `header:"x-cos-grant-full-control,omitempty" url:"-" xml:"-"`
}
