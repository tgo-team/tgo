package tgo

type Protocol interface {
	// 解码消息
	DecodePacket(connContext ConnContext) (interface{},error)
	// 编码消息
	EncodePacket(packet interface{}) ([]byte,error)
}