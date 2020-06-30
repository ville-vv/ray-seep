package msg

import (
	jsoniter "github.com/json-iterator/go"
)

type MessagePusher struct {
	ResponseSender
}

func (p *MessagePusher) PushInJson(cmd int32, obj interface{}) (err error) {
	body, err := jsoniter.Marshal(obj)
	if err != nil {
		return err
	}
	p.SendCh() <- Package{Cmd: cmd, Body: body}
	return
}

func (p *MessagePusher) PushInByte(cmd int32, data []byte) (err error) {
	p.SendCh() <- Package{Cmd: cmd, Body: data[:]}
	return
}
