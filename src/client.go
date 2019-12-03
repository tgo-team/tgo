package tgo

import (
	"github.com/tgo-team/tgo/src/log"
	"go.uber.org/zap"
	"net"
	"time"
)

// TCPClient TCPClient
type TCPClient struct {
	clientID string
	net.Conn
	sCtx *ServerContext
}

// NewTCPClient NewTCPClient
func NewTCPClient(clientID string, conn net.Conn, sCtx *ServerContext) *TCPClient {
	c := &TCPClient{}
	c.clientID = clientID
	c.Conn = conn
	c.sCtx = sCtx
	go c.msgLoop()
	return c
}

func (c *TCPClient) msgLoop() {
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

// SetDeadline 客户端死亡线
func (c *TCPClient) SetDeadline(t time.Time) {
	c.Conn.SetDeadline(t)
}

// GetID 客户端唯一ID
func (c *TCPClient) GetID() string {
	return c.clientID
}

// GetUID 客户端的uid
func (c *TCPClient) GetUID() string {
	return c.clientID
}

// Close 关闭客户端
func (c *TCPClient) Close() error {
	return c.Conn.Close()
}

// KeepAlive 客户端保活
func (c *TCPClient) KeepAlive() {

}
