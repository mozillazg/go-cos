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

// NeedSignHeaders 需要校验的 Headers 列表
var NeedSignHeaders = map[string]bool{
	"host":                     true,
	"range":                    true,
	"x-cos-acl":                true,
	"x-cos-grant-read":         true,
	"x-cos-grant-write":        true,
	"x-cos-grant-full-control": true,
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
	formatURI := uri.EscapedPath()

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
			ps.Add(key, strings.ToLower(value))
			signedParameterList = append(signedParameterList, key)
		}
	}
	formatParameters = strings.ToLower(ps.Encode())
	sort.Strings(signedParameterList)
	return
}

// GenFormatHeaders 生成 FormatHeaders 和 SignedHeaderList
func GenFormatHeaders(headers http.Header) (formatHeaders string, signedHeaderList []string) {
	hs := url.Values{}
	for key, values := range headers {
		for _, value := range values {
			key = strings.ToLower(key)
			for k := range NeedSignHeaders {
				key = strings.ToLower(key)
				if key == k {
					hs.Add(key, strings.ToLower(value))
					signedHeaderList = append(signedHeaderList, key)
					break
				}
			}
		}
	}
	formatHeaders = strings.ToLower(hs.Encode())
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
