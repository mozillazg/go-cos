package cos

import (
	"context"
	"encoding/xml"
	"net/http"
)

// ListMultipartUploadsResult ...
type ListMultipartUploadsResult struct {
	XMLName            xml.Name `xml:"ListMultipartUploadsResult"`
	Bucket             string   `xml:"Bucket"`
	EncodingType       string   `xml:"Encoding-Type"`
	KeyMarker          string
	UploadIDMarker     string `xml:"UploadIdMarker"`
	NextKeyMarker      string
	NextUploadIDMarker string `xml:"NextUploadIdMarker"`
	MaxUploads         int
	IsTruncated        bool
	Uploads            []*struct {
		Key          string
		UploadID     string
		StorageClass string
		Initiator    *struct {
			UID string
		}
		Owner     *Owner
		Initiated string
	} `xml:"Upload"`
	Prefix        string
	Delimiter     string   `xml:"delimiter,omitempty"`
	CommonPrefixs []string `xml:"CommonPrefixs>Prefix,omitempty"`
}

// ListMultipartUploadsOptions ...
type ListMultipartUploadsOptions struct {
	Delimiter      string `url:"delimiter,omitempty"`
	EncodingType   string `url:"encoding-type,omitempty"`
	Prefix         string `url:",omitempty"`
	MaxUploads     int    `url:"max-uploads,omitempty"`
	KeyMarker      string `url:"key-marker,omitempty"`
	UploadIDMarker string `url:"upload-id-marker,omitempty"`
}

// ListMultipartUploads ...
//
// List Multipart Uploads用来查询正在进行中的分块上传。单次最多列出1000个正在进行中的分块上传。
//
// https://www.qcloud.com/document/product/436/7736
func (s *BucketService) ListMultipartUploads(ctx context.Context,
	authTime *AuthTime,
	opt *ListMultipartUploadsOptions) (*ListMultipartUploadsResult, *Response, error) {
	u := "/?uploads"
	baseURL := s.client.BaseURL.BucketURL
	var res ListMultipartUploadsResult
	resp, err := s.client.sendNoBody(ctx, baseURL, u, http.MethodGet, authTime, opt, nil, &res)
	return &res, resp, err
}
