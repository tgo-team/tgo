package tgo

import "fmt"

type Handler func(ctx Context)

type TGO struct {
	opts        *Options // TGO启动参数
	servers     []Server // server集合
	acceptChan  chan Context
	pro         Protocol
	runExitChan chan int
	handler     Handler
	waitGroup   WaitGroupWrapper
	router      Router
}

func New(options *Options) *TGO {

	return &TGO{opts: options, servers: make([]Server, 0), runExitChan: make(chan int, 0), acceptChan: make(chan Context, 1024)}
}

// Start 开始TGO
func (t *TGO) Start() {

	for _, svr := range t.servers {
		err := svr.Start(t.serverContext(svr))
		if err != nil {
			panic(err)
		}
	}
	t.waitGroup.Wrap(t.msgLoop)
}

func (t *TGO) Run() {
	t.Start()

	<-t.runExitChan
}

// Stop 停止TGO
func (t *TGO) Stop() {
	for _, svr := range t.servers {
		err := svr.Stop()
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("TGO stopped")
}

// UseServer 指定server服务器
func (t *TGO) UseServer(server Server) {
	t.servers = append(t.servers, server)
}

// UseProtocol 指定协议
func (t *TGO) UseProtocol(p Protocol) {
	t.pro = p
}

func (t *TGO) UseRouter(router Router) {
	t.router = router
}

// UseHandler 处理者
func (t *TGO) UseHandler(handler Handler) {
	t.handler = handler
}

// GetProtocol 获取协议
func (t *TGO) GetProtocol() Protocol {
	return t.pro
}

func (t *TGO) serverContext(svr Server) *ServerContext {
	return NewServerContext(t, svr)
}

func (t *TGO) msgLoop() {
	for {
		select {
		case context := <-t.acceptChan:
			if t.handler != nil {
				t.handler(context)
			}
			// 匹配处理者
			t.matchHandler(context)

		}
	}
}

// 匹配处理者
func (t *TGO) matchHandler(context Context)  {
	if t.router!=nil {
		handler := t.router.MatchHandler(context)
		if handler!=nil {
			handler(context)
		}
	}
}
