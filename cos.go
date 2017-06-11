package cos

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"reflect"
	"time"

	"bitbucket.org/mozillazg/go-httpheader"
	"github.com/google/go-querystring/query"
	"io"
)

const (
	// Version ...
	Version        = "0.2.0"
	userAgent      = "go-cos/" + Version
	contentTypeXML = "application/xml"
)

// A Client manages communication with the COS API.
type Client struct {
	Client    *http.Client
	secretID  string
	secretKey string

	UserAgent   string
	ContentType string
	Secure      bool // 是否使用 https

	common service

	Service *ServiceService
	Bucket  *BucketService
}

type service struct {
	client *Client
	bucket *Bucket
}

func (s *service) SetBucket(b *Bucket) {
	s.bucket = b
}

// NewClient returns a new COS API client.
func NewClient(secretID, secretKey string, b *Bucket, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	c := &Client{
		Client:      httpClient,
		secretID:    secretID,
		secretKey:   secretKey,
		UserAgent:   userAgent,
		ContentType: contentTypeXML,
		Secure:      true,
	}
	c.common.client = c
	c.common.bucket = b
	c.Service = (*ServiceService)(&c.common)
	c.Bucket = (*BucketService)(&c.common)
	return c
}

// SetTimeout 设置超时时间
func (c *Client) SetTimeout(t time.Duration) {
	c.Client.Timeout = t
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
	var bXML []byte
	if body != nil {
		bXML, err = xml.Marshal(body)
		if err != nil {
			return
		}
		reader = bytes.NewReader(bXML)
	}

	req, err = http.NewRequest(method, urlStr, reader)
	if err != nil {
		return
	}

	req.Header, err = addHeaderOptions(req, optHeader)
	if err != nil {
		return
	}
	if body != nil {
		req.Header.Set("Content-Length", len(bXML))
		req.Header.Set("Content-MD5", base64.StdEncoding.EncodeToString(calMD5Digest(bXML)))
	}
	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}
	return
}

func (c *Client) doAPI(ctx context.Context, req *http.Request, ret interface{},
	authTime AuthTime) (resp *http.Response, err error) {
	req = req.WithContext(ctx)

	if authTime != nil {
		AddAuthorization(
			c.secretID, c.secretKey, req,
			authTime.signStartTime, authTime.signEndTime,
			authTime.keyStartTime, authTime.keyEndTime,
		)
	}

	a, _ := httputil.DumpRequest(req, true)
	fmt.Println(string(a))

	resp, err = c.Client.Do(req)
	if err != nil {
		return
	}

	b, _ := httputil.DumpResponse(resp, true)
	fmt.Println(string(b))

	if resp.StatusCode >= http.StatusBadRequest {
		var e ErrorResponse
		err = xml.NewDecoder(resp.Body).Decode(&e)
		if err == nil {
			e.Response = resp
			err = &e
		}
		return
	}

	if ret != nil {
		err = xml.NewDecoder(resp.Body).Decode(&ret)
	}
	return
}

func (c *Client) sendWithBody(ctx context.Context, baseURL *url.URL, uri, method string,
	authTime AuthTime, body interface{},
	optQuery interface{}, optHeader interface{}, ret interface{}) (resp *http.Response, err error) {
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
	authTime AuthTime,
	optQuery interface{}, optHeader interface{}, ret interface{}) (resp *http.Response, err error) {
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
func addHeaderOptions(req http.Request, opt interface{}) (http.Request, error) {
	v := reflect.ValueOf(opt)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return req, nil
	}

	h, err := httpheader.Header(opt)
	if err != nil {
		return nil, err
	}

	for key, values := range h {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}
	return req, nil
}

// Owner ...
type Owner struct {
	UIN string `xml:"uin"`
}

// Initiator ...
type Initiator struct {
	UID string
}

// Opt 定义请求参数
//type Opt struct {
//	query  interface{} // url 参数
//	header interface{} // request header 参数
//}

//// NewOpt ...
//func NewOpt(query, header interface{}) {
//
//}

// AuthTime 用于生成签名所需的 q-sign-time 和 q-key-time 相关参数
type AuthTime struct {
	signStartTime time.Time
	signEndTime   time.Time
	keyStartTime  time.Time
	keyEndTime    time.Time
}

// NewAuthTime ...
func NewAuthTime(signStartTime, signEndTime,
	keyStartTime, keyEndTime time.Time) AuthTime {
	return AuthTime{
		signStartTime, signEndTime,
		keyStartTime, keyEndTime,
	}
}
