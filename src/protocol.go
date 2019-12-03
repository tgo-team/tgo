package tgo

import "io"

// Protocol Protocol
type Protocol interface {
	// DecodePacket 解码消息
	DecodePacket(conn io.Reader) (interface{}, error)
	// EncodePacket 编码消息
	EncodePacket(packet interface{}) ([]byte, error)
}
