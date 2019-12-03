package tgo

// Server Server
type Server interface {
	Start(context *ServerContext) error
	Stop() error
	GetRouter() Router
	GetProtocol() Protocol
}

// ServerContext ServerContext
type ServerContext struct {
	T   *TGO
	svr Server
}

// NewServerContext NewServerContext
func NewServerContext(tg *TGO, svr Server) *ServerContext {
	return &ServerContext{T: tg, svr: svr}
}

// Accept 接收请求
func (sc *ServerContext) Accept(ctx Context) {
	sc.T.AcceptChan <- ctx
}

// GetProtocol 获取协议
func (sc *ServerContext) GetProtocol() Protocol {

	return sc.svr.GetProtocol()
}
