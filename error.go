package cos

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"io/ioutil"
)

// ErrorResponse 包含 API 返回的错误信息
// https://www.qcloud.com/document/product/436/7730
type ErrorResponse struct {
	XMLName   xml.Name       `xml:"Error"`
	Response  *http.Response `xml:"-"`
	Code      string
	Message   string
	Resource  string
	RequestID string `xml:"RequestId"`
	TraceID   string `xml:"TraceId"`
}

// Error ...
func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %v(Message: %v, RequestId: %v, TraceId: %v)",
		r.Response.Request.Method, r.Response.Request.URL,
		r.Response.StatusCode, r.Code, r.Message, r.RequestID, r.TraceID)
}

// 检查 response 是否时出错时的返回的 response
func checkResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}
	errorResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && data != nil {
		xml.Unmarshal(data, errorResponse)
	}
	return errorResponse
}
