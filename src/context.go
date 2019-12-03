package tgo

import (
	"github.com/tgo-team/tgo/src/log"
	"go.uber.org/zap"
)

// Context Context
type Context interface {
	// 获取包
	GetPacket() interface{}
	// GetClient 获取客户端
	GetClient() Client
	// WritePacket 写包
	WritePacket(packet interface{})
	// GetServerContext 获取Server上下文
	GetServerContext() *ServerContext
}

// DefaultContext DefaultContext
type DefaultContext struct {
	packetCtx PacketContext
	client    Client
	pro       Protocol
	sCtx      *ServerContext
}

// NewDefaultContext 这里为了防止大量类型转换 两个参数可以选一个 无状态连接conn，一个有状态连接 statefulConn
func NewDefaultContext(packetCtx PacketContext, sCtx *ServerContext, pro Protocol, client Client) *DefaultContext {
	return &DefaultContext{
		packetCtx: packetCtx,
		client:    client,
		pro:       pro,
		sCtx:      sCtx,
	}
}

// GetServerContext Server的上下文
func (d *DefaultContext) GetServerContext() *ServerContext {
	return d.sCtx
}

// GetPacket 获取当前请求的包
func (d *DefaultContext) GetPacket() interface{} {

	return d.packetCtx.GetPacket()
}

// GetClient 获取客户端对象
func (d *DefaultContext) GetClient() Client {
	return d.client
}

// WritePacket 写入包
func (d *DefaultContext) WritePacket(packet interface{}) {
	packetBytes, err := d.pro.EncodePacket(packet)
	if err != nil {
		log.Warn("EncodePacket失败！", zap.Error(err))
	}
	_, err = d.GetClient().Write(packetBytes)
	if err != nil {
		log.Warn("Write失败！", zap.Error(err))
	}
}
