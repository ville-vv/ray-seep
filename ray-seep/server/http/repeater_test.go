package http

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"ray-seep/ray-seep/common/rayhttp"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestNetRepeater_relay(t *testing.T) {

	l2, err := net.Listen("tcp", ":34981")
	if err != nil {
		t.Error(err)
		return
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		wg.Done()
		repeat := &NetRepeater{}
		for {
			cn2, err := l2.Accept()
			if err != nil {
				t.Error(err)
				os.Exit(1)
			}

			copyHttp, err := rayhttp.ToHttp(cn2)
			fmt.Println("访问的地址：", copyHttp.Host())
			fmt.Println("内容：", string(copyHttp.GetBody()))
			ccn, err := net.Dial("tcp", "www.villeboss.com:10081")
			if err != nil {
				t.Log(err)
				cn2.Close()
				return
			}
			_, _, err = repeat.relay(ccn, copyHttp)
			if err != nil {
				cn2.Close()
				t.Error(err)
				return
			}
		}
	}()
	wg.Wait()
	//resp, err := http.Get("http://127.0.0.1:34981/p/b4102e3e3e96")
	resp, err := http.Post("http://127.0.0.1:34981/api/user/callback?", "application/json", strings.NewReader("4444444"))
	if err != nil {
		t.Error(err)
		return
	}
	bodyResp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		resp.Body.Close()
		t.Error(err)
		return
	}
	resp.Body.Close()
	fmt.Println("返回结果：", string(bodyResp))
	time.Sleep(time.Second * 3)
}
