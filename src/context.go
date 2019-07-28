package tgo

type Context interface {
	// 获取包
	GetPacket() interface{}
	// GetConn 获取无状态连接
	GetConn() Conn
	// GetStatefulConn 获取有状态连接
	GetStatefulConn() StatefulConn
	// WritePacket 写包
	WritePacket(packet interface{})
}

type DefaultContext struct {
	packetCtx PacketContext
	conn    Conn
	statefulConn StatefulConn
	pro Protocol
}

// NewDefaultContext 这里为了防止大量类型转换 两个参数可以选一个 无状态连接conn，一个有状态连接 statefulConn
func NewDefaultContext(packetCtx PacketContext, conn Conn,pro Protocol,statefulConn StatefulConn) *DefaultContext {
	return &DefaultContext{
		packetCtx: packetCtx,
		conn:    conn,
		statefulConn:statefulConn,
		pro:pro,
	}
}

func (d *DefaultContext) GetPacket() interface{} {

	return d.packetCtx.GetPacket()
}

func (d *DefaultContext) GetConn() Conn  {
	if d.conn == nil {
		return d.statefulConn
	}
	return d.conn
}

func (d *DefaultContext) GetStatefulConn() StatefulConn  {
	return d.statefulConn
}

func (d *DefaultContext) WritePacket(packet interface{})   {
	packetBytes,err := d.pro.EncodePacket(packet)
	if err!=nil {
		panic(err)
	}
	_,err =  d.GetConn().Write(packetBytes)
	if err!=nil {
		panic(err)
	}
}