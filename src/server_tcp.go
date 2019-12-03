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

// TCPOptions 配置
type TCPOptions struct {
	Addr              string        // tcp连接地址
	lg                log.Log       // 日志
	handshake         HandshakeFnc  // 握手
	HeartbeatInterval time.Duration //心跳间隔
	router            Router
	protocol          Protocol
	createClientFnc   CreateClientFnc // 创建客户端的方法
}

// NewTCPOptions 创建一个默认配置
func NewTCPOptions() *TCPOptions {
	return &TCPOptions{
		Addr:              "0.0.0.0:6666",
		HeartbeatInterval: time.Second * 30,
	}
}

// TCPOption 配置项
type TCPOption func(opt *TCPOptions)

// TCPAddr tcp连接地址
func TCPAddr(addr string) TCPOption {
	return func(opt *TCPOptions) {
		opt.Addr = addr
	}
}

// TCPProtocol 协议接口
func TCPProtocol(protocol Protocol) TCPOption {
	return func(opt *TCPOptions) {
		opt.protocol = protocol
	}
}

// TCPRouter tcp路由
func TCPRouter(router Router) TCPOption {
	return func(opt *TCPOptions) {
		opt.router = router
	}
}

// TCPLog 日志
func TCPLog(lg log.Log) TCPOption {
	return func(opt *TCPOptions) {
		opt.lg = lg
	}
}

// TCPHandshake 握手函数
func TCPHandshake(handshake HandshakeFnc) TCPOption {
	return func(opt *TCPOptions) {
		opt.handshake = handshake
	}
}

// TCPHeartbeatInterval 心跳
func TCPHeartbeatInterval(heartbeatInterval time.Duration) TCPOption {
	return func(opt *TCPOptions) {
		opt.HeartbeatInterval = heartbeatInterval
	}
}

// TCPCreateClient 创建tcp客户端
func TCPCreateClient(createClientFnc CreateClientFnc) TCPOption {
	return func(opt *TCPOptions) {
		opt.createClientFnc = createClientFnc
	}
}

// NewContextFnc NewContextFnc
type NewContextFnc func(packetCtx PacketContext, sContext *ServerContext, conn Conn, client Client) Context

// CreateClientFnc 创建客户端函数
type CreateClientFnc func(clientId string, packet interface{}, conn net.Conn, sCtx *ServerContext) (Client, error)

// ServerTCP ServerTCP
type ServerTCP struct {
	tcpListener net.Listener
	ctx         *ServerContext
	exitChan    chan int
	waitGroup   WaitGroupWrapper
	realAddr    string // 真实连接地址
	opts        *TCPOptions
}

// HandshakeFnc 握手函数
type HandshakeFnc func(packet interface{}, conn net.Conn, ctx *ServerContext) (string, error)

// NewServerTCP NewServerTCP
func NewServerTCP(opts ...TCPOption) *ServerTCP {
	options := NewTCPOptions()
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
	return &ServerTCP{opts: options, exitChan: make(chan int)}
}

// Start 开始运行
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

// GetProtocol 获取协议
func (s *ServerTCP) GetProtocol() Protocol {
	return s.opts.protocol
}

// GetRealAddr 获取真实的tcp的地址
func (s *ServerTCP) GetRealAddr() string {
	return s.realAddr
}

// Stop 停止
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
	var clientID string
	if s.opts.handshake != nil {
		clientID, err = s.opts.handshake(packet, cn, s.ctx)
		if err != nil {
			s.Debug("握手失败！", zap.Error(err))
			cn.Close()
			return
		}
	}
	pCtx := NewDefaultPacketContext(packet)
	context := NewDefaultContext(pCtx, s.ctx, s.GetProtocol(), nil)
	if s.opts.createClientFnc != nil {
		context.client, err = s.opts.createClientFnc(clientID, packet, cn, s.ctx)
		if err != nil {
			s.Error("客户端创建失败！", zap.Error(err))
			return
		}
	} else {
		// 创建一个tcp客户端
		context.client = NewTCPClient(clientID, cn, s.ctx)
	}

	s.ctx.Accept(context)
}

// GetRouter 获取路由
func (s *ServerTCP) GetRouter() Router {

	return s.opts.router
}

// Debug Debug
func (s *ServerTCP) Debug(msg string, fields ...zap.Field) {
	s.opts.lg.Debug(msg, fields...)
}

// Info Info
func (s *ServerTCP) Info(msg string, fields ...zap.Field) {
	s.opts.lg.Info(msg, fields...)
}

// Error Error
func (s *ServerTCP) Error(msg string, fields ...zap.Field) {
	s.opts.lg.Error(msg, fields...)
}
