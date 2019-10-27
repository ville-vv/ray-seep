package http

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"ray-seep/ray-seep/common/rayhttp"
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
			defer cn2.Close()

			copyHttp, err := rayhttp.ToHttp(cn2)
			fmt.Println("访问的地址：", copyHttp.Host())
			//fmt.Println("内容：", string(copyHttp.GetBody()))
			ccn, err := net.Dial("tcp", "www.jianshu.com:80")
			if err != nil {
				cn2.Write([]byte(err.Error()))
				t.Log(err)
				return
			}
			_, _, err = repeat.relay(ccn, copyHttp)
			if err != nil {
				t.Error(err)
				return
			}

		}
	}()
	wg.Wait()
	resp, err := http.Get("http://127.0.0.1:34981/p/b4102e3e3e96")
	//resp, err := http.Post("http://127.0.0.1:34981/api/user/callback?", "application/json", strings.NewReader("4444444"))
	if err != nil {
		t.Error(err)
		return
	}
	defer resp.Body.Close()
	bodyResp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println("返回结果：", string(bodyResp))
	time.Sleep(time.Second * 3)
}
