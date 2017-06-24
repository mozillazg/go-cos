package cos

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"os"
)

// 计算 md5 或 sha1 时的分块大小
const calDigestBlockSize = 1024 * 1024 * 10

func calMD5Digest(msg []byte) []byte {
	// TODO: 分块计算,减少内存消耗
	m := md5.New()
	m.Write(msg)
	return m.Sum(nil)
}

func calSHA1Digest(msg []byte) []byte {
	// TODO: 分块计算,减少内存消耗
	m := sha1.New()
	m.Write(msg)
	return m.Sum(nil)
}

// DebugRequestTransport 会打印请求和响应信息, 方便调试.
type DebugRequestTransport struct {
	RequestHeader  bool
	RequestBody    bool // RequestHeader 为 true 时,这个选项才会生效
	ResponseHeader bool
	ResponseBody   bool // ResponseHeader 为 true 时,这个选项才会生效

	// debug 信息输出到 Writer 中, 默认是 os.Stderr
	Writer io.Writer

	Transport http.RoundTripper
}

// RoundTrip implements the RoundTripper interface.
func (t *DebugRequestTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req = cloneRequest(req) // per RoundTrip contract
	w := t.Writer
	if w == nil {
		w = os.Stderr
	}

	if t.RequestHeader {
		a, _ := httputil.DumpRequestOut(req, t.RequestBody)
		fmt.Fprintf(w, "%s\n\n", string(a))
	}

	resp, err := t.transport().RoundTrip(req)
	if err != nil {
		return resp, err
	}

	if t.ResponseHeader {

		b, _ := httputil.DumpResponse(resp, t.ResponseBody)
		fmt.Fprintf(w, "%s\n", string(b))
	}

	return resp, err
}

func (t *DebugRequestTransport) transport() http.RoundTripper {
	if t.Transport != nil {
		return t.Transport
	}
	return http.DefaultTransport
}

// cloneRequest returns a clone of the provided *http.Request. The clone is a
// shallow copy of the struct and its Header map.
func cloneRequest(r *http.Request) *http.Request {
	// shallow copy of the struct
	r2 := new(http.Request)
	*r2 = *r
	// deep copy of the Header
	r2.Header = make(http.Header, len(r.Header))
	for k, s := range r.Header {
		r2.Header[k] = append([]string(nil), s...)
	}
	return r2
}

// encodeURIComponent like same function in javascript
//
// https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/encodeURIComponent
//
// http://www.ecma-international.org/ecma-262/6.0/#sec-uri-syntax-and-semantics
func encodeURIComponent(s string) string {
	var b bytes.Buffer
	written := 0

	for i, n := 0, len(s); i < n; i++ {
		c := s[i]

		switch c {
		case '-', '_', '.', '!', '~', '*', '\'', '(', ')':
			continue
		default:
			// Unreserved according to RFC 3986 sec 2.3
			if 'a' <= c && c <= 'z' {

				continue

			}
			if 'A' <= c && c <= 'Z' {

				continue

			}
			if '0' <= c && c <= '9' {

				continue
			}
		}

		b.WriteString(s[written:i])
		fmt.Fprintf(&b, "%%%02x", c)
		written = i + 1
	}

	if written == 0 {
		return s
	}
	b.WriteString(s[written:])
	return b.String()
}
