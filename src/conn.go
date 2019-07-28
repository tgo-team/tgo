package tgo

import (
	"net"
	"time"
)

type Conn interface {
	Read(b []byte) (n int, err error)
	Write(b []byte) (n int, err error)
	Close() error
}


type StatefulConn interface {
	Read(b []byte) (n int, err error)
	Write(b []byte) (n int, err error)
	SetDeadline(t time.Time)
	GetClientId() uint64
	Close() error
	GetProps() map[string]interface{}
}

type DefaultStatefulConn struct {
	conn net.Conn
	clientId uint64
	props map[string]interface{}
}

func NewStatefulConn(conn net.Conn,clientId uint64,props map[string]interface{}) *DefaultStatefulConn {
	return &DefaultStatefulConn{conn: conn,clientId:clientId,props:props}
}

func (c *DefaultStatefulConn) Read(b []byte) (n int, err error) {
	return c.conn.Read(b)
}

func (c *DefaultStatefulConn) Write(b []byte) (n int, err error) {
	return c.conn.Write(b)
}

func (c *DefaultStatefulConn) SetDeadline(t time.Time)  {
	c.conn.SetDeadline(t)
}
func (c *DefaultStatefulConn) GetClientId() uint64  {
	return c.clientId
}

func (c *DefaultStatefulConn) Close() error  {
	return c.conn.Close()
}


type ConnContext interface {
	// GetConn 获取连接
	GetConn() Conn
}

type DefaultConnContext struct {
	conn Conn
}

func NewDefaultConnContext(conn Conn) *DefaultConnContext {

	return &DefaultConnContext{conn: conn}
}

func (c *DefaultConnContext) GetConn() Conn {
	return c.conn
}
