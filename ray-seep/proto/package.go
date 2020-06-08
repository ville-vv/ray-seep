// @File     : package
// @Author   : Ville
// @Time     : 19-9-24 上午11:25
// msg
package proto

import (
	"encoding/binary"
	"errors"
	jsoniter "github.com/json-iterator/go"
	"hash/crc32"
	"strconv"
)

// 协议版本
const version = uint8(1)

type Frame struct {
	Cmd  uint32
	sum  uint32
	Body []byte
}

func (f *Frame) Pack() []byte {
	l := len(f.Body) + 40
	buf := make([]byte, l)
	f.sum = crc32.ChecksumIEEE(f.Body)
	binary.LittleEndian.PutUint32(buf[0:4], f.Cmd)
	binary.LittleEndian.PutUint32(buf[4:8], f.sum)
	buf[8] = version
	buf = append(buf[0:9], f.Body...)
	return buf
}

func (f *Frame) UnPack(data []byte) error {
	if len(data) < 9 {
		return nil
	}
	f.Cmd = binary.LittleEndian.Uint32(data[0:4])
	f.sum = binary.LittleEndian.Uint32(data[4:8])
	f.Body = append(f.Body, data[9:]...)
	// 检测版本
	if version != data[8] {
		return errors.New("frame version error")
	}
	// 校验和
	if f.sum != crc32.ChecksumIEEE(f.Body) {
		return errors.New("frame check sum error")
	}
	return nil
}

type Package struct {
	Cmd  int32  `json:"cmd"`
	Body []byte `json:"body"`
}

// 新建一个包 body 是 []byte 类型
func New(cmd int32, body []byte) *Package {
	return &Package{Cmd: cmd, Body: body}
}

// NewPackage 新建一个 pack, body 支持 int, string, []byte 类型，其他类型转换成json数据
func NewPackage(cmd int32, body interface{}) *Package {
	switch body.(type) {
	case int:
		return &Package{Cmd: cmd, Body: []byte(strconv.Itoa(body.(int)))}
	case string:
		return &Package{Cmd: cmd, Body: []byte(body.(string))}
	case []byte:
		return &Package{Cmd: cmd, Body: body.([]byte)}
	default:
		dt, _ := jsoniter.Marshal(body)
		return &Package{Cmd: cmd, Body: dt}
	}
}

// NewWithObj 新建一个pack body 支持 map, struct 类型，并且直接转换成 json 格式
func NewWithObj(cmd int32, body interface{}) *Package {
	dt, _ := jsoniter.Marshal(body)
	return &Package{Cmd: cmd, Body: dt}
}

// unPack 解包
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

// pack 打包
func Pack(pkg *Package) (data []byte, err error) {
	frame := Frame{}
	frame.Body, err = jsoniter.Marshal(pkg)
	if err != nil {
		return nil, err
	}
	return frame.Pack(), nil
}
