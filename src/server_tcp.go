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
	newContextFnc     NewContextFnc // 自定义context
	lg                log.Log       // 日志
	handshake         HandshakeFnc  // 握手
	HeartbeatInterval time.Duration //心跳间隔
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

func TcpNewContextFnc(newContextFnc NewContextFnc) TcpOption {
	return func(opt *TcpOptions) {
		opt.newContextFnc = newContextFnc
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

type NewContextFnc func(packetCtx PacketContext, sContext *ServerContext, conn Conn, statefulConn StatefulConn) Context
type ServerTCP struct {
	tcpListener net.Listener
	ctx         *ServerContext
	exitChan    chan int
	waitGroup   WaitGroupWrapper
	pro         Protocol
	RealAddr    string // 真实连接地址
	opts        *TcpOptions
}

type HandshakeFnc func(packet interface{}, conn StatefulConn) bool

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
	return &ServerTCP{opts: options}
}

func (s *ServerTCP) Start(context *ServerContext) error {
	s.ctx = context
	s.pro = context.GetProtocol()
	var err error
	s.tcpListener, err = net.Listen("tcp", s.opts.Addr)
	if err != nil {
		return err
	}
	s.RealAddr = s.tcpListener.Addr().String()
	s.info("启动 ", zap.String("addr", s.tcpListener.Addr().String()))
	s.waitGroup.Wrap(s.connLoop)
	return nil
}

func (s *ServerTCP) Stop() error {
	err := s.tcpListener.Close()
	if err != nil {
		return err
	}
	s.exitChan <- 1
	s.debug("退出")
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
	s.debug("TCP Server 停止监听")
}

func (s *ServerTCP) handleConn(cn net.Conn) {
	cn.SetDeadline(time.Now().Add(2 * time.Second)) // 指定时间内没有握手成功，就关闭连接
	statefulConn := NewStatefulConn(cn, fmt.Sprintf("%d", s.ctx.T.GenClientId()), nil)

	packet, err := s.pro.DecodePacket(cn)
	if err != nil {
		log.Error("解码连接消息失败！", zap.Error(err))
		return
	}
	var isHandshake = s.opts.handshake(packet, statefulConn)
	if !isHandshake {
		s.debug("握手失败！")
		statefulConn.Close()
		return
	}
	cn.SetDeadline(time.Now().Add(s.opts.HeartbeatInterval*2)) // 握手成功后将死亡时间设置为心跳的2倍
	pCtx := NewDefaultPacketContext(packet)
	var context Context
	if s.opts.newContextFnc != nil {
		context = s.opts.newContextFnc(pCtx, s.ctx, statefulConn, statefulConn)
	} else {
		context = NewDefaultContext(pCtx, s.ctx, statefulConn, s.pro, statefulConn)
	}
	s.ctx.Accept(context)
}

func (s *ServerTCP) debug(msg string, fields ...zap.Field) {
	s.opts.lg.Debug(msg, fields...)
}

func (s *ServerTCP) info(msg string, fields ...zap.Field) {
	s.opts.lg.Info(msg, fields...)
}
