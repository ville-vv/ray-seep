// @File     : header
// @Author   : Ville
// @Time     : 19-9-25 下午4:55
// http
package http

import (
	"bufio"
	"bytes"
	"io"
	"net"
	"net/http"
	"sync"
)

type buildRequest struct {
	sync.Mutex
	net.Conn
	buf *bytes.Buffer
}

func newBuildRequest(c net.Conn) (*buildRequest, io.Reader) {
	b := &buildRequest{
		Conn: c,
		buf:  bytes.NewBuffer(make([]byte, 0, 4096)),
	}
	// 如果有触发 io.reader 就会把数据写入到 buf 中
	return b, io.TeeReader(c, b.buf)
}

// Read zai在 Request 中读取了http.Request 后conn 中的数据会写入到 buf 里面，如果外部再次调用Read
// 那么需要吧原来的 bytes 读入，然后再使用 conn 中的read
func (c *buildRequest) Read(p []byte) (n int, err error) {
	c.Lock()
	if c.buf == nil {
		c.Unlock()
		return c.Conn.Read(p)
	}
	n, err = c.buf.Read(p)
	if err == io.EOF {
		c.buf = nil
		var n2 int
		n2, err = c.Conn.Read(p[n:])
		n += n2
	}
	c.Unlock()
	return
}

type CopyHttp struct {
	*buildRequest
	request *http.Request
}

// ToHttp 转为HTTP格式，获取http消息
func NewCopyHttp(c net.Conn) (hp *CopyHttp, err error) {
	rq, rd := newBuildRequest(c)

	hp = &CopyHttp{
		buildRequest: rq,
	}
	// 读取HTTP请求
	if hp.request, err = http.ReadRequest(bufio.NewReader(rd)); err != nil {
		return
	}
	hp.request.Body.Close()
	return
}

func (c *CopyHttp) SayBack(status int, body []byte) {
	resp := http.Response{
		StatusCode: status,
		ProtoMajor: 1,
		ProtoMinor: 1,
		Request:    c.request,
		Header:     c.request.Header,
		//Body:       ioutil.NopCloser(bytes.NewBuffer(body)),
	}
	resp.Write(c.Conn)
}
