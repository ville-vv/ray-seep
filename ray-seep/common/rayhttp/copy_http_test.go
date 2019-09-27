// @File     : header_test
// @Author   : Ville
// @Time     : 19-9-26 上午9:39
// http
package rayhttp

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
)

func dummyReq11(method string) *http.Request {
	return &http.Request{Method: method, Proto: "HTTP/1.1", ProtoMajor: 1, ProtoMinor: 1}
}

func TestHttpResponse(t *testing.T) {
	var braw bytes.Buffer

	resp := http.Response{
		StatusCode: 200,
		ProtoMajor: 1,
		ProtoMinor: 1,
		Request:    dummyReq11("GET"),
		Header:     http.Header{},
		Close:      true,
		Body:       ioutil.NopCloser(bytes.NewBuffer([]byte("123lkjoijlgjlj"))),
	}
	resp.Write(&braw)
	t.Log(braw.String())
}
