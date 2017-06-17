package cos

import (
	"context"
	"encoding/xml"
	"net/http"
)

// GetLocationResult ...
type GetLocationResult struct {
	XMLName  xml.Name `xml:"LocationConstraint"`
	Location string   `xml:",chardata"`
}

// GetLocation ...
//
// Get Bucket Location接口获取Bucket所在地域信息，只有Bucket所有者有权限读取信息。
//
// https://www.qcloud.com/document/product/436/8275
func (s *BucketService) GetLocation(ctx context.Context,
	authTime *AuthTime) (*GetLocationResult, *Response, error) {
	u := "/?location"
	baseURL := s.client.BaseURL.BucketURL
	var res GetLocationResult
	resp, err := s.client.sendNoBody(ctx, baseURL, u, http.MethodGet, authTime, nil, nil, &res)
	return &res, resp, err
}
