// 内部模块消息通道
package msg

type InternalMsgChan struct {
	msgCache chan Package
}

func (im *InternalMsgChan) Push(pg *Package) {
	im.msgCache <- *pg
}

func (InternalMsgChan) Watch() {
}
