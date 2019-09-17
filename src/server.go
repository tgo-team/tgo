package tgo

type Server interface {
	Start(context *ServerContext) error
	Stop() error
	GetRouter() Router
	GetProtocol() Protocol
}

type ServerContext struct {
	T *TGO
	svr Server
}

func NewServerContext(tg *TGO,svr Server) *ServerContext {
	return &ServerContext{T: tg,svr:svr}
}

func (sc *ServerContext) Accept(ctx Context)  {
	sc.T.AcceptChan <- ctx
}

func (sc *ServerContext) GetProtocol() Protocol {

	return sc.svr.GetProtocol()
}