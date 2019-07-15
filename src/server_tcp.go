package tgo

import (
	"fmt"
	"net"
	"runtime"
	"strings"
)

type ServerTCP struct {
	tcpListener net.Listener
	ctx         *ServerContext
	exitChan         chan int
	waitGroup        WaitGroupWrapper
	pro Protocol
	addr string
}

func NewServerTCP(addr string) *ServerTCP  {
	return &ServerTCP{addr:addr}
}

func (s *ServerTCP) Start(context *ServerContext) error {
	s.ctx = context
	s.pro = context.GetProtocol()
	var err error
	s.tcpListener, err = net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	s.waitGroup.Wrap(s.connLoop)
	return nil
}

func (s *ServerTCP) Stop() error {

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
	fmt.Println("退出Server")
}

func (s *ServerTCP) handleConn(cn net.Conn)  {
	cCtx := NewDefaultConnContext(cn)
	packet,err := s.pro.DecodePacket(cCtx)
	if err!=nil {
		fmt.Println("解码消息失败！-> ",err.Error())
		return
	}
	pCtx := NewDefaultPacketContext(packet)
	s.ctx.Accept(NewDefaultContext(pCtx,cCtx))
}


