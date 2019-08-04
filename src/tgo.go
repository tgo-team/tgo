package tgo

import (
	"fmt"
	"github.com/tgo-team/tgo/src/log"
	"go.uber.org/zap"
	"sync/atomic"
)

type Handler func(ctx Context)

type TGO struct {
	opts             *Options // TGO启动参数
	servers          []Server // server集合
	AcceptChan       chan Context
	pro              Protocol
	runExitChan      chan int
	handler          Handler
	waitGroup        WaitGroupWrapper
	router           Router
	clientId         uint64
	handlerWaitGroup WaitGroupWrapper
	ConnManager      ConnManager
}

func New(options *Options) *TGO {

	return &TGO{opts: options, servers: make([]Server, 0), runExitChan: make(chan int, 0), AcceptChan: make(chan Context, 1024), ConnManager: NewDefaultConnManager()}
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
	t.debug("退出")
}

// UseServer 指定server服务器
func (t *TGO) UseServer(server Server) {
	t.servers = append(t.servers, server)
}

func (t *TGO) ClearServers() {
	t.servers = make([]Server, 0)
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

// GenClientId 生成客户端ID
func (t *TGO) GenClientId() uint64 {
	return atomic.AddUint64(&t.clientId, 1)
}

func (t *TGO) serverContext(svr Server) *ServerContext {
	return NewServerContext(t, svr)
}

func (t *TGO) msgLoop() {
	for {
		select {
		case context := <-t.AcceptChan:
			if t.handler != nil {
				t.handlerWaitGroup.Wrap(func() {
					t.handler(context)
				})
			}
			// 匹配处理者
			t.matchHandler(context)
		}
	}
}

// 匹配处理者
func (t *TGO) matchHandler(context Context) {
	if t.router != nil {
		handler := t.router.MatchHandler(context)
		if handler != nil {
			t.handlerWaitGroup.Wrap(func() {
				handler(context)
			})
		}
	}
}

func (t *TGO) debug(msg string, fields ...zap.Field) {
	log.Debug(fmt.Sprintf("【TGO】%s", msg), fields...)
}

func (t *TGO) info(msg string, fields ...zap.Field) {
	log.Info(fmt.Sprintf("【TGO】%s", msg), fields...)
}
