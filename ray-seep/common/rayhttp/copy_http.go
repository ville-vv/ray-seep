// @File     : header
// @Author   : Ville
// @Time     : 19-9-25 下午4:55
// http
package rayhttp

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"sync"
	"time"
)

const (
	NotAuthorized = `HTTP/1.0 401 Not Authorized
WWW-Authenticate: Basic realm="ngrok"
Content-Length: 23

Authorization required
`

	NotFound = `HTTP/1.0 404 Not Found
Content-Length: %d

Tunnel %s not found
`

	BadRequest = `HTTP/1.0 400 Bad Request
Content-Length: 12

Bad Request
`
	Success = `HTTP/1.0 200 Ok
Content-Length: 13

success hello
`
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
	defer c.Unlock()
	if c.buf == nil {
		return c.Conn.Read(p)
	}
	n, err = c.buf.Read(p)
	if err == io.EOF {
		c.buf = nil
		var n2 int
		n2, err = c.Conn.Read(p[n:])
		n += n2
	}
	return
}

type CopyHttp struct {
	*buildRequest
	request *http.Request
	body    []byte
}

// ToHttp 转为HTTP格式，获取http消息
func ToHttp(c net.Conn) (hp *CopyHttp, err error) {
	rq, rd := newBuildRequest(c)

	hp = &CopyHttp{
		buildRequest: rq,
	}
	// 读取HTTP请求
	if hp.request, err = http.ReadRequest(bufio.NewReader(rd)); err != nil {
		return
	}

	defer hp.request.Body.Close()
	//hp.body = hp.buf.Bytes()
	//bf := hp.buf.Bytes()
	//hp.body = make([]byte, len(bf))
	//for v := range bf {
	//	hp.body[v] = bf[v]
	//}
	return
}

func (c *CopyHttp) GetBody() []byte {
	bf := make([]byte, len(c.body))
	for v := range c.body {
		bf[v] = c.body[v]
	}
	return bf
}
func (c *CopyHttp) SayBackText(status int, body []byte) {
	resp := http.Response{
		StatusCode:    status,
		ProtoMajor:    1,
		ProtoMinor:    1,
		Request:       c.request,
		Header:        http.Header{},
		ContentLength: int64(len(body)),
		Body:          ioutil.NopCloser(bytes.NewBuffer(body)),
	}
	resp.Header.Add("Content-Type", "text/html;charset=utf-8")
	resp.Header.Add("date", time.Now().Format(time.RFC1123))
	resp.Write(c.buildRequest.Conn)
}

func (c *CopyHttp) Host() string {
	return c.request.Host
}

func (c *CopyHttp) RemoteHost() string {
	return c.request.RemoteAddr
}
