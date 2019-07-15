package tgo

type PacketContext interface {
	// GetMsg 获取消息
	 GetPacket() interface{}
}

type Packet interface {
	String() string
}

type DefaultPacketContext struct {
	p interface{}
}

func NewDefaultPacketContext(p interface{}) *DefaultPacketContext {
	return &DefaultPacketContext{p:p}
}
func (d *DefaultPacketContext) GetPacket() interface{} {
	return  d.p
}