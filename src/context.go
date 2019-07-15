package tgo

type Context interface {
	GetPacket() interface{}
}

type DefaultContext struct {
	packetCtx PacketContext
	connCtx    ConnContext
}

func NewDefaultContext(packetCtx PacketContext, connCtx ConnContext) *DefaultContext {
	return &DefaultContext{
		packetCtx: packetCtx,
		connCtx:    connCtx,
	}
}

func (d *DefaultContext) GetPacket() interface{} {

	return d.packetCtx.GetPacket()
}