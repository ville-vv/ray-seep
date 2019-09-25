// @File     : msg
// @Author   : Ville
// @Time     : 19-9-24 下午2:59
// msg
package msg

import (
	jsoniter "github.com/json-iterator/go"
)

type Message struct {
	Cmd  string      `json:"cmd"`
	Body interface{} `json:"body"`
}

func UnPack(data []byte, pkg *Message) (err error) {
	frame := Frame{}
	if err = frame.UnPack(data); err != nil {
		return
	}
	if len(frame.Body) > 0 {
		if err = jsoniter.Unmarshal(frame.Body, pkg); err != nil {
			return
		}
	}
	return
}

func Pack(pkg *Message) (data []byte, err error) {
	frame := Frame{}
	frame.Body, err = jsoniter.Marshal(pkg)
	if err != nil {
		return nil, err
	}
	return frame.Pack(), nil
}
