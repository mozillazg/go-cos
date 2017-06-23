package cos

import (
	"crypto/hmac"
	"crypto/sha1"
	"fmt"
	"hash"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

const sha1SignAlgorithm = "sha1"
const privateHeaderPrefix = "x-cos-"

// 需要校验的 Headers 列表
var needSignHeaders = map[string]bool{
	"host":                           true,
	"range":                          true,
	"x-cos-acl":                      true,
	"x-cos-grant-read":               true,
	"x-cos-grant-write":              true,
	"x-cos-grant-full-control":       true,
	"response-content-type":          true,
	"response-content-language":      true,
	"response-expires":               true,
	"response-cache-control":         true,
	"response-content-disposition":   true,
	"response-content-encoding":      true,
	"cache-control":                  true,
	"content-disposition":            true,
	"content-encoding":               true,
	"content-type":                   true,
	"content-length":                 true,
	"content-md5":                    true,
	"expect":                         true,
	"expires":                        true,
	"x-cos-content-sha1":             true,
	"x-cos-storage-class":            true,
	"if-modified-since":              true,
	"origin":                         true,
	"access-control-request-method":  true,
	"access-control-request-headers": true,
	"x-cos-object-type":              true,
}

// NewAuthorization 通过一系列步骤生成最终需要的 Authorization 字符串
func NewAuthorization(secretID, secretKey string, req *http.Request,
	signStartTime, signEndTime, keyStartTime, keyEndTime time.Time,
) string {
	signTime := GenSignTime(signStartTime, signEndTime)
	keyTime := GenSignTime(keyStartTime, keyEndTime)
	signKey := CalSignKey(secretKey, keyTime)

	formatHeaders, signedHeaderList := GenFormatHeaders(req.Header)
	formatParameters, signedParameterList := GenFormatParameters(req.URL.Query())
	formatString := GenFormatString(req.Method, *req.URL, formatParameters, formatHeaders)

	stringToSign := CalStringToSign(sha1SignAlgorithm, keyTime, formatString)
	signature := CalSignature(signKey, stringToSign)

	return GenAuthorization(
		secretID, signTime, keyTime, signature, signedHeaderList,
		signedParameterList,
	)
}

// AddAuthorization 给 req 增加签名信息
func AddAuthorization(secretID, secretKey string, req *http.Request,
	signStartTime, signEndTime, keyStartTime, keyEndTime time.Time,
) {
	auth := NewAuthorization(secretID, secretKey, req,
		signStartTime, signEndTime, keyStartTime, keyEndTime,
	)
	req.Header.Set("Authorization", auth)
}

// CalSignKey 计算 SignKey
func CalSignKey(secretKey, keyTime string) string {
	digest := calHMACDigest(secretKey, keyTime, sha1SignAlgorithm)
	return fmt.Sprintf("%x", digest)
}

// CalStringToSign 计算 StringToSign
func CalStringToSign(signAlgorithm, signTime, formatString string) string {
	h := sha1.New()
	h.Write([]byte(formatString))
	return fmt.Sprintf("%s\n%s\n%x\n", signAlgorithm, signTime, h.Sum(nil))
}

// CalSignature 计算 Signature
func CalSignature(signKey, stringToSign string) string {
	digest := calHMACDigest(signKey, stringToSign, sha1SignAlgorithm)
	return fmt.Sprintf("%x", digest)
}

// GenAuthorization 生成 Authorization
func GenAuthorization(secretID, signTime, keyTime, signature string,
	signedHeaderList, signedParameterList []string) string {
	return strings.Join([]string{
		"q-sign-algorithm=" + sha1SignAlgorithm,
		"q-ak=" + secretID,
		"q-sign-time=" + signTime,
		"q-key-time=" + keyTime,
		"q-header-list=" + strings.Join(signedHeaderList, ";"),
		"q-url-param-list=" + strings.Join(signedParameterList, ";"),
		"q-signature=" + signature,
	}, "&")
}

// GenSignTime 生成 SignTime
func GenSignTime(startTime, endTime time.Time) string {
	return fmt.Sprintf("%d;%d", startTime.Unix(), endTime.Unix())
}

// GenFormatString 生成 FormatString
func GenFormatString(method string, uri url.URL, formatParameters, formatHeaders string) string {
	formatMethod := strings.ToLower(method)
	formatURI := uri.Path

	return fmt.Sprintf("%s\n%s\n%s\n%s\n", formatMethod, formatURI,
		formatParameters, formatHeaders,
	)
}

// GenFormatParameters 生成 FormatParameters 和 SignedParameterList
func GenFormatParameters(parameters url.Values) (formatParameters string, signedParameterList []string) {
	ps := url.Values{}
	for key, values := range parameters {
		for _, value := range values {
			key = strings.ToLower(key)
			ps.Add(key, value)
			signedParameterList = append(signedParameterList, key)
		}
	}
	//formatParameters = strings.ToLower(ps.Encode())
	formatParameters = ps.Encode()
	sort.Strings(signedParameterList)
	return
}

// GenFormatHeaders 生成 FormatHeaders 和 SignedHeaderList
func GenFormatHeaders(headers http.Header) (formatHeaders string, signedHeaderList []string) {
	hs := url.Values{}
	for key, values := range headers {
		for _, value := range values {
			key = strings.ToLower(key)
			if isSignHeader(key) {
				hs.Add(key, value)
				signedHeaderList = append(signedHeaderList, key)
			}
		}
	}
	formatHeaders = hs.Encode()
	sort.Strings(signedHeaderList)
	return
}

// HMAC 签名
func calHMACDigest(key, msg, signMethod string) []byte {
	var hashFunc func() hash.Hash
	switch signMethod {
	case "sha1":
		hashFunc = sha1.New
	default:
		hashFunc = sha1.New
	}
	h := hmac.New(hashFunc, []byte(key))
	h.Write([]byte(msg))
	return h.Sum(nil)
}

func isSignHeader(key string) bool {
	for k, v := range needSignHeaders {
		if key == k && v {
			return true
		}
	}
	return strings.HasPrefix(key, privateHeaderPrefix)
}
