package tgo

import (
	"fmt"
	"net"
	"runtime"
	"strings"
)

type ServerUnix struct {
	listener *net.UnixListener
	ctx         *ServerContext
	exitChan         chan int
	waitGroup        WaitGroupWrapper
	pro Protocol
	fileName string
	addr *net.UnixAddr
}

func NewServerUnix(fileName string) *ServerUnix {
	s := &ServerUnix{}
	var err error
	s.addr, err = net.ResolveUnixAddr("unix", fileName)
	if err!=nil {
		panic(err)
	}
	s.listener, err = net.ListenUnix("unix", s.addr)
	if err!=nil {
		panic(err)
	}
	return s
}


func (s *ServerUnix) Start(context *ServerContext) error {
	s.pro = context.tg.pro
	s.ctx = context
	s.waitGroup.Wrap(s.connLoop)
	return nil
}

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

func (s *ServerUnix) handleConn(cn net.Conn)  {

	packet,err := s.pro.DecodePacket(cn)
	if err!=nil {
		fmt.Println("解码消息失败！-> ",err.Error())
		return
	}
	s.ctx.Accept(NewDefaultPacketContext(packet))
}
