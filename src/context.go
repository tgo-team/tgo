package tgo

type Context interface {
	GetPacket() interface{}
}

type DefaultContext struct {
	packetCtx PacketContext
	conn    Conn
}

func NewDefaultContext(packetCtx PacketContext, conn Conn) *DefaultContext {
	return &DefaultContext{
		packetCtx: packetCtx,
		conn:    conn,
	}
}

func (d *DefaultContext) GetPacket() interface{} {

	return d.packetCtx.GetPacket()
}