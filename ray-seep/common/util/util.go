// @File     : uitl
// @Author   : Ville
// @Time     : 19-9-23 下午6:24
// common
package util

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	mrand "math/rand"
	"os"
	"strconv"
	"time"
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

func HmacSha256String(secretKey string, content string) string {
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(content))
	return hex.EncodeToString(h.Sum(nil))
}
func GenRandID() int {
	ai := fmt.Sprintf("%06v", mrand.New(mrand.NewSource(time.Now().UnixNano())).Int31n(1000000))
	id, _ := strconv.Atoi(ai)
	return id
}

func RandString(len int) string {
	b := make([]byte, len)
	_, _ = io.ReadFull(rand.Reader, b)
	return GetMd5String(base64.URLEncoding.EncodeToString(b))
}

// 生成Guid字串
func RandToken() string {
	return RandString(48)
}

func WritePid(path string) error {
	pid := os.Getpid()
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	_, err = f.WriteString(fmt.Sprintf("%d", pid))
	if err != nil {
		return err
	}
	f.Sync()
	return nil
}
