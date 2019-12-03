package tgo

// PacketContext 包上下文
type PacketContext interface {
	// GetMsg 获取消息
	GetPacket() interface{}
}

// Packet 包接口
type Packet interface {
	String() string
}

// DefaultPacketContext 默认包上下文
type DefaultPacketContext struct {
	p interface{}
}

// NewDefaultPacketContext NewDefaultPacketContext
func NewDefaultPacketContext(p interface{}) *DefaultPacketContext {
	return &DefaultPacketContext{p: p}
}

// GetPacket 获取包对象
func (d *DefaultPacketContext) GetPacket() interface{} {
	return d.p
}
