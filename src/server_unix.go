package tgo

import (
	"fmt"
	"net"
	"runtime"
	"strings"
)

// ServerUnix ServerUnix
type ServerUnix struct {
	listener  *net.UnixListener
	ctx       *ServerContext
	exitChan  chan int
	waitGroup WaitGroupWrapper
	fileName  string
	addr      *net.UnixAddr
}

// NewServerUnix NewServerUnix
func NewServerUnix(fileName string) *ServerUnix {
	s := &ServerUnix{}
	var err error
	s.addr, err = net.ResolveUnixAddr("unix", fileName)
	if err != nil {
		panic(err)
	}
	s.listener, err = net.ListenUnix("unix", s.addr)
	if err != nil {
		panic(err)
	}
	return s
}

// Start Start
func (s *ServerUnix) Start(context *ServerContext) error {
	s.ctx = context
	s.waitGroup.Wrap(s.connLoop)
	return nil
}

// Stop Stop
func (s *ServerUnix) Stop() error {
	err := s.listener.Close()
	s.waitGroup.Wait()
	fmt.Println("ServerUnix stopped")
	return err
}

func (s *ServerUnix) connLoop() {
	for {
		select {
		case <-s.exitChan:
			goto exit
		default:
			cn, err := s.listener.AcceptUnix()
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

func (s *ServerUnix) handleConn(cn net.Conn) {

	packet, err := s.GetProtocol().DecodePacket(cn)
	if err != nil {
		fmt.Println("解码消息失败！-> ", err.Error())
		return
	}
	pCtx := NewDefaultPacketContext(packet)
	s.ctx.Accept(NewDefaultContext(pCtx, s.ctx, s.GetProtocol(), nil))
}

// GetRouter GetRouter
func (s *ServerUnix) GetRouter() Router {
	return nil
}

// GetProtocol GetProtocol
func (s *ServerUnix) GetProtocol() Protocol {

	return nil
}
