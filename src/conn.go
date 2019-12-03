package tgo

import (
	"time"
)

// Conn Conn
type Conn interface {
	Read(b []byte) (n int, err error)
	Write(b []byte) (n int, err error)
	Close() error
}

// Client 客户端接口
type Client interface {
	Read(b []byte) (n int, err error)
	Write(b []byte) (n int, err error)
	SetDeadline(t time.Time)
	GetID() string
	GetUID() string
	Close() error
	KeepAlive()
}

// ConnContext 连接上下文
type ConnContext interface {
	// GetConn 获取连接
	GetConn() Conn
}

// DefaultConnContext DefaultConnContext
type DefaultConnContext struct {
	conn Conn
}

// NewDefaultConnContext NewDefaultConnContext
func NewDefaultConnContext(conn Conn) *DefaultConnContext {

	return &DefaultConnContext{conn: conn}
}

// GetConn 获取连接
func (c *DefaultConnContext) GetConn() Conn {
	return c.conn
}
