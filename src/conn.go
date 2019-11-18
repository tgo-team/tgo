package tgo

import (
	"time"
)

type Conn interface {
	Read(b []byte) (n int, err error)
	Write(b []byte) (n int, err error)
	Close() error
}


type Client interface {
	Read(b []byte) (n int, err error)
	Write(b []byte) (n int, err error)
	SetDeadline(t time.Time)
	GetId() string
	GetUid() string
	Close() error
    KeepAlive()
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
