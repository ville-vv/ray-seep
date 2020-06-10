package msg

import (
	"encoding/binary"
	"errors"
	jsoniter "github.com/json-iterator/go"
	"hash/crc32"
	"ray-seep/ray-seep/common/util"
)

// 协议版本
const version = uint8(1)

// 每次通信时的帧数据
type frame struct {
	head    uint32 //
	version uint8
	sum     uint32
	Body    []byte
}

func newFrame() frame {
	return frame{version: version}
}

func (f *frame) pack() []byte {
	l := len(f.Body) + 40
	buf := make([]byte, l)
	f.sum = crc32.ChecksumIEEE(f.Body)
	binary.LittleEndian.PutUint32(buf[0:4], f.head)
	binary.LittleEndian.PutUint32(buf[4:8], f.sum)
	buf[8] = f.version
	buf = append(buf[0:9], f.Body...)
	return buf
}

func (f *frame) unPack(data []byte) error {
	if len(data) < 9 {
		return nil
	}
	f.head = binary.LittleEndian.Uint32(data[0:4])
	f.sum = binary.LittleEndian.Uint32(data[4:8])
	f.Body = append(f.Body, data[9:]...)
	// 检测版本
	if f.version != data[8] {
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

type packerManagerJson struct {
}

// unPack 解包
func (*packerManagerJson) UnPack(data []byte, pkg *Package) (err error) {
	frame := newFrame()
	if err = frame.unPack(data); err != nil {
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
func (*packerManagerJson) Pack(pkg *Package) (data []byte, err error) {
	frame := newFrame()
	frame.Body, err = jsoniter.Marshal(pkg)
	if err != nil {
		return nil, err
	}
	return frame.pack(), nil
}

type packerManager01 struct {
}

// unPack 解包
func (*packerManager01) UnPack(data []byte, pkg *Package) (err error) {
	frame := newFrame()
	if err = frame.unPack(data); err != nil {
		return
	}
	pkg.Cmd, err = util.BytesToInt32(frame.Body[:4])
	if err != nil {
		return
	}
	pkg.Body = make([]byte, 0, len(frame.Body)-4)
	pkg.Body = append(pkg.Body, frame.Body[4:]...)
	return
}

// pack 打包
func (*packerManager01) Pack(pkg *Package) (data []byte, err error) {
	frame := newFrame()
	cmd := util.Int32ToBytes(pkg.Cmd)
	frame.Body = make([]byte, 0, len(cmd)+len(pkg.Body))
	frame.Body = append(frame.Body, cmd...)
	frame.Body = append(frame.Body, pkg.Body...)
	return frame.pack(), nil
}
