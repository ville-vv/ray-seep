package msg

import "context"

// 解包接口
type UnPacker interface {
	UnPack(data []byte, pkg *Package) (err error)
}

// 打包接口
type Packer interface {
	Pack(pkg *Package) (data []byte, err error)
}

// 消息包管理器
type PackerManager interface {
	UnPacker
	Packer
}

// 消息接收器
type Receiver interface {
	RecvMsg() (buf []byte, err error)
}

// 消息发送器
type Sender interface {
	SendMsg(data []byte) (err error)
}

type ResponseSender interface {
	Send(pg *Package) error
	SendCh() chan<- Package
}

type Request struct {
	Ctx  context.Context
	Body *Package
}

type RouterFunc func(req *Request, wt ResponseSender) error
