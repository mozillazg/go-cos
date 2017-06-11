package cos

import (
	"crypto/md5"
)

func calMD5Digest(msg []byte) []byte {
	m := md5.New()
	m.Write(msg)
	return m.Sum(nil)
}

// Bool is a helper routine that allocates a new bool value
// to store v and returns a pointer to it.
func Bool(v bool) *bool { return &v }

// Int is a helper routine that allocates a new int value
// to store v and returns a pointer to it.
func Int(v int) *int { return &v }

// String is a helper routine that allocates a new string value
// to store v and returns a pointer to it.
func String(v string) *string { return &v }
