// @File     : uitl
// @Author   : Ville
// @Time     : 19-9-23 下午6:24
// common
package util

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"io"
)

// 生成32位md5字串
func GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func GetMd5(s string) []byte {
	h := md5.New()
	h.Write([]byte(s))
	src := h.Sum(nil)
	dst := make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(dst, src)
	return dst
}

// 生成Guid字串
func RandToken() string {
	b := make([]byte, 48)
	io.ReadFull(rand.Reader, b)
	return GetMd5String(base64.URLEncoding.EncodeToString(b))
}
