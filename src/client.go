package tgo

import (
	"github.com/tgo-team/tgo/src/log"
	"go.uber.org/zap"
	"net"
	"time"
)

type TCPClient struct {
	clientId uint64
	net.Conn
	sCtx *ServerContext
}

func NewTCPClient(clientId uint64,conn net.Conn,sCtx *ServerContext) *TCPClient  {
	c := &TCPClient{}
	c.clientId = clientId
	c.Conn = conn
	c.sCtx = sCtx
	go c.msgLoop()
	return c
}


func (c *TCPClient) msgLoop()  {
	for {
		packet, err := c.sCtx.GetProtocol().DecodePacket(c)
		if err != nil {
			log.Debug("连接关闭", zap.Error(err))
			c.Close()
			return
		}
		pCtx := NewDefaultPacketContext(packet)
		c.sCtx.Accept(NewDefaultContext(pCtx, c.sCtx, c.sCtx.GetProtocol(), c))
	}
}


func (c *TCPClient) Read(b []byte) (n int, err error) {
	return c.Conn.Read(b)
}

func (c *TCPClient) Write(b []byte) (n int, err error) {
	return c.Conn.Write(b)
}

func (c *TCPClient) SetDeadline(t time.Time)  {
	c.Conn.SetDeadline(t)
}
func (c *TCPClient) GetId() uint64  {
	return c.clientId
}

func (c *TCPClient) Close() error  {
	return c.Conn.Close()
}

func (c *TCPClient)  KeepAlive()  {

}

