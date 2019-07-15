package tgo

type Server interface {
	Start(context *ServerContext) error
	Stop() error
}

type ServerContext struct {
	tg *TGO
	svr Server
}

func NewServerContext(tg *TGO,svr Server) *ServerContext {
	return &ServerContext{tg: tg,svr:svr}
}

func (sc *ServerContext) Accept(ctx Context)  {
	sc.tg.acceptChan <- ctx
}

func (sc *ServerContext) GetProtocol() Protocol {

	return sc.tg.pro
}