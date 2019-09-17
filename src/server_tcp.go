package tgo

import (
	"fmt"
	"github.com/tgo-team/tgo/src/log"
	"go.uber.org/zap"
	"net"
	"runtime"
	"strings"
	"time"
)

type TcpOptions struct {
	Addr              string        // tcp连接地址
	lg                log.Log       // 日志
	handshake         HandshakeFnc  // 握手
	HeartbeatInterval time.Duration //心跳间隔
	router Router
	protocol Protocol
	createClientFnc CreateClientFnc // 创建客户端的方法
}

func NewTcpOptions() *TcpOptions {
	return &TcpOptions{
		Addr:              "0.0.0.0:6666",
		HeartbeatInterval: time.Second * 30,
	}
}

type TcpOption func(opt *TcpOptions)

func TcpAddr(addr string) TcpOption {
	return func(opt *TcpOptions) {
		opt.Addr = addr
	}
}

func TcpProtocol(protocol Protocol) TcpOption {
	return func(opt *TcpOptions) {
		opt.protocol = protocol
	}
}



func TcpRouter(router Router) TcpOption {
	return func(opt *TcpOptions) {
		opt.router = router
	}
}

func TcpLog(lg log.Log) TcpOption {
	return func(opt *TcpOptions) {
		opt.lg = lg
	}
}
func TcpHandshake(handshake HandshakeFnc) TcpOption {
	return func(opt *TcpOptions) {
		opt.handshake = handshake
	}
}

func TcpHeartbeatInterval(heartbeatInterval time.Duration) TcpOption {
	return func(opt *TcpOptions) {
		opt.HeartbeatInterval = heartbeatInterval
	}
}

func TcpCreateClient(createClientFnc CreateClientFnc) TcpOption {
	return func(opt *TcpOptions) {
		opt.createClientFnc = createClientFnc
	}
}

type NewContextFnc func(packetCtx PacketContext, sContext *ServerContext, conn Conn, client Client) Context
type CreateClientFnc func(clientId uint64,packet interface{},conn net.Conn,sCtx *ServerContext) (Client,error)
type ServerTCP struct {
	tcpListener net.Listener
	ctx         *ServerContext
	exitChan    chan int
	waitGroup   WaitGroupWrapper
	realAddr    string // 真实连接地址
	opts        *TcpOptions
}

type HandshakeFnc func(packet interface{}, conn net.Conn,ctx *ServerContext) (error,uint64)

func NewServerTCP(opts ...TcpOption) *ServerTCP {
	options := NewTcpOptions()
	if opts != nil {
		for _, opt := range opts {
			opt(options)
		}
	}
	var nLog log.Log
	if options.lg == nil {
		nLog = log.NewTLog("ServerTCP")
	} else {
		nLog = options.lg
	}
	options.lg = nLog
	return &ServerTCP{opts: options,exitChan:make(chan int)}
}

func (s *ServerTCP) Start(context *ServerContext) error {
	s.ctx = context
	var err error
	s.tcpListener, err = net.Listen("tcp", s.opts.Addr)
	if err != nil {
		return err
	}
	s.realAddr = s.tcpListener.Addr().String()
	s.Info("启动 ", zap.String("addr", s.realAddr))
	s.waitGroup.Wrap(s.connLoop)
	return nil
}

func (s *ServerTCP) GetProtocol() Protocol  {
	return s.opts.protocol
}

func (s *ServerTCP) GetRealAddr() string  {
	return s.realAddr
}

func (s *ServerTCP) Stop() error {
	err := s.tcpListener.Close()
	if err != nil {
		return err
	}
	s.exitChan <- 1
	s.Debug("退出")
	return nil
}

func (s *ServerTCP) connLoop() {
	for {
		select {
		case <-s.exitChan:
			goto exit
		default:
			cn, err := s.tcpListener.Accept()
			if err != nil {
				if nerr, ok := err.(net.Error); ok && nerr.Temporary() {
					fmt.Println("temporary Accept() failure - ", err)
					runtime.Gosched()
					continue
				}
				// theres no direct way to detect this error because it is not exposed
				if !strings.Contains(err.Error(), "use of closed network connection") {
					fmt.Println("listener.Accept() - ", err)
				}
				break
			}
			s.waitGroup.Wrap(func() {
				s.handleConn(cn)
			})
		}
	}
exit:
	s.Debug("TCP Server 停止监听")
}

func (s *ServerTCP) handleConn(cn net.Conn) {
	cn.SetDeadline(time.Now().Add(2 * time.Second)) // 指定时间内没有握手成功，就关闭连接
	packet, err := s.GetProtocol().DecodePacket(cn)
	if err != nil {
		log.Error("解码连接消息失败！", zap.Error(err))
		return
	}
	var clientId uint64
	if s.opts.handshake!=nil {
		err,clientId = s.opts.handshake(packet, cn,s.ctx)
		if err!=nil {
			s.Debug("握手失败！",zap.Error(err))
			cn.Close()
			return
		}
	}
	pCtx := NewDefaultPacketContext(packet)
	context := NewDefaultContext(pCtx, s.ctx, s.GetProtocol(), nil)
	if s.opts.createClientFnc!=nil {
		context.client,err = s.opts.createClientFnc(clientId,packet,cn,s.ctx)
		if err!=nil {
			s.Error("客户端创建失败！",zap.Error(err))
			return
		}
	}else {
		// 创建一个tcp客户端
		context.client = NewTCPClient(clientId,cn,s.ctx)
	}

	s.ctx.Accept(context)
}

func (s *ServerTCP) GetRouter() Router {

	return s.opts.router
}

func (s *ServerTCP) Debug(msg string, fields ...zap.Field) {
	s.opts.lg.Debug(msg, fields...)
}

func (s *ServerTCP) Info(msg string, fields ...zap.Field) {
	s.opts.lg.Info(msg, fields...)
}

func (s *ServerTCP) Error(msg string, fields ...zap.Field) {
	s.opts.lg.Error(msg, fields...)
}
