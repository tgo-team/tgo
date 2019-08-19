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
	// GetServerContext 获取Server上下文
	GetServerContext() *ServerContext
}

type DefaultContext struct {
	packetCtx PacketContext
	conn    Conn
	statefulConn StatefulConn
	pro Protocol
	sCtx *ServerContext
}

// NewDefaultContext 这里为了防止大量类型转换 两个参数可以选一个 无状态连接conn，一个有状态连接 statefulConn
func NewDefaultContext(packetCtx PacketContext,sCtx *ServerContext, conn Conn,pro Protocol,statefulConn StatefulConn) *DefaultContext {
	return &DefaultContext{
		packetCtx: packetCtx,
		conn:    conn,
		statefulConn:statefulConn,
		pro:pro,
		sCtx: sCtx,
	}
}

// GetServerContext Server的上下文
func (d *DefaultContext) GetServerContext() *ServerContext  {
	return  d.sCtx
}

// GetPacket 获取当前请求的包
func (d *DefaultContext) GetPacket() interface{} {

	return d.packetCtx.GetPacket()
}

// GetConn 当前请求的连接
func (d *DefaultContext) GetConn() Conn  {
	if d.conn == nil {
		return d.statefulConn
	}
	return d.conn
}

// GetStatefulConn 当前请求的有状态连接
func (d *DefaultContext) GetStatefulConn() StatefulConn  {
	return d.statefulConn
}

// WritePacket 写入包
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