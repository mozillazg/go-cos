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
	var res GetLocationResult
	sendOpt := sendOptions{
		baseURL:  s.client.BaseURL.BucketURL,
		uri:      "/?location",
		method:   http.MethodGet,
		authTime: authTime,
		result:   &res,
	}
	resp, err := s.client.send(ctx, &sendOpt)
	return &res, resp, err
}
