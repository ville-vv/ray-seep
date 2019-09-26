// @File     : package
// @Author   : Ville
// @Time     : 19-9-24 上午11:25
// msg
package pkg

import (
	"encoding/binary"
	"errors"
	jsoniter "github.com/json-iterator/go"
	"hash/crc32"
	"reflect"
)

// 协议版本
const version = uint8(1)

type Frame struct {
	Cmd  uint32
	Sum  uint32
	Body []byte
}

func (f *Frame) Pack() []byte {
	l := len(f.Body) + 40
	buf := make([]byte, l)
	f.Sum = crc32.ChecksumIEEE(f.Body)
	binary.LittleEndian.PutUint32(buf[0:4], f.Cmd)
	binary.LittleEndian.PutUint32(buf[4:8], f.Sum)
	buf[8] = version
	buf = append(buf[0:9], f.Body...)
	return buf
}

func (f *Frame) UnPack(data []byte) error {
	if len(data) < 9 {
		return nil
	}
	f.Cmd = binary.LittleEndian.Uint32(data[0:4])
	f.Sum = binary.LittleEndian.Uint32(data[4:8])
	f.Body = append(f.Body, data[9:]...)
	// 检测版本
	if version != data[8] {
		return errors.New("frame version error")
	}
	// 校验和
	if f.Sum != crc32.ChecksumIEEE(f.Body) {
		return errors.New("frame check sum error")
	}
	return nil
}

type Package struct {
	Cmd  string      `json:"cmd"`
	Body interface{} `json:"body"`
}

func NewPackage(cmd string, body interface{}) *Package {
	if cmd == "" {
		cmd = reflect.TypeOf(body).Elem().Name()
	}
	return &Package{Cmd: cmd, Body: body}
}

func UnPack(data []byte, pkg *Package) (err error) {
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

func Pack(pkg *Package) (data []byte, err error) {
	frame := Frame{}
	frame.Body, err = jsoniter.Marshal(pkg)
	if err != nil {
		return nil, err
	}
	return frame.Pack(), nil
}
