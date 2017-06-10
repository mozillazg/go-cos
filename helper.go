package cos

import (
	"crypto/md5"
)

func calMD5Digest(msg []byte) []byte {
	m := md5.New()
	m.Write(msg)
	return m.Sum(nil)
}

// contentMD5 computes and sets the HTTP Content-MD5 header for requests that
// require it.
//func contentMD5(req *http.Request) (err error){
//	h := md5.New()
//
//	// hash the body.  seek back to the first position after reading to reset
//	// the body for transmission.  copy errors may be assumed to be from the
//	// body.
//	_, err = io.Copy(h, req.Body)
//	if err != nil {
//		return
//	}
//	io.ReadSeeker
//	_, err = req.Body.Seek(0, 0)
//	if err != nil {
//		return
//	}
//
//	// encode the md5 checksum in base64 and set the request header.
//	sum := h.Sum(nil)
//	sum64 := make([]byte, base64.StdEncoding.EncodedLen(len(sum)))
//	base64.StdEncoding.Encode(sum64, sum)
//	req.Header.Set("Content-MD5", string(sum64))
//}
