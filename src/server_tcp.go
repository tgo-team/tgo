package tgo

import (
	"fmt"
	"github.com/tgo-team/tgo/src/log"
	"go.uber.org/zap"
	"net"
	"runtime"
	"strings"
)

type NewContextFnc func(packetCtx PacketContext,sContext *ServerContext, conn Conn, statefulConn StatefulConn) Context
type ServerTCP struct {
	tcpListener   net.Listener
	ctx           *ServerContext
	exitChan      chan int
	waitGroup     WaitGroupWrapper
	pro           Protocol
	addr          string
	RealAddr      string // 真实连接地址
	newContextFnc NewContextFnc
}

func NewServerTCP(addr string, newContextFnc NewContextFnc) *ServerTCP {
	return &ServerTCP{addr: addr, exitChan: make(chan int, 0), newContextFnc: newContextFnc}
}

func (s *ServerTCP) Start(context *ServerContext) error {
	s.ctx = context
	s.pro = context.GetProtocol()
	var err error
	s.tcpListener, err = net.Listen("tcp", s.addr)
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

	//var tgoConn Conn
	//var statefulConn StatefulConn
	//var packet interface{}
	//var err error
	//if s.handshakeFnc != nil {
	//	packet, statefulConn, err = s.handshakeFnc(cn,s.ctx, 10*time.Second)
	//	if err != nil {
	//		log.Debug("握手失败！", zap.Error(err))
	//		cn.Close()
	//		return
	//	}
	//	tgoConn = statefulConn
	//} else {
	//	tgoConn = NewStatefulConn(cn, fmt.Sprintf("%d",s.ctx.T.GenClientId()), nil)
	//	packet, err = s.pro.DecodePacket(tgoConn)
	//	if err != nil {
	//		fmt.Println("解码消息失败！-> ", err.Error())
	//		cn.Close()
	//		return
	//	}
	//}
	//if tgoConn == nil  {
	//	log.Debug("握手失败！")
	//	cn.Close()
	//	return
	//}
	statefulConn := NewStatefulConn(cn, fmt.Sprintf("%d", s.ctx.T.GenClientId()), nil)
	packet, err := s.pro.DecodePacket(cn)
	if err != nil {
		log.Error("解码连接消息失败！", zap.Error(err))
		return
	}
	pCtx := NewDefaultPacketContext(packet)
	var context Context
	if s.newContextFnc != nil {
		context = s.newContextFnc(pCtx,s.ctx, statefulConn, statefulConn)
	} else {
		context = NewDefaultContext(pCtx, s.ctx, statefulConn, s.pro, statefulConn)
	}
	s.ctx.Accept(context)
}

func (s *ServerTCP) debug(msg string, fields ...zap.Field) {
	log.Debug(fmt.Sprintf("【TCP Server】%s", msg), fields...)
}

func (s *ServerTCP) info(msg string, fields ...zap.Field) {
	log.Info(fmt.Sprintf("【TCP Server】%s", msg), fields...)
}
