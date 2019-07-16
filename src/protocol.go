package tgo

import "io"

type Protocol interface {
	// 解码消息
	DecodePacket(conn io.Reader) (interface{},error)
	// 编码消息
	EncodePacket(packet interface{}) ([]byte,error)
}