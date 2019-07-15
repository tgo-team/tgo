package mqtt

import (
	"github.com/tgo-team/tgo-core/src"
	"io"
)

type MqttPacket struct {
	*FixedHeader                  // 固定头部
	VariableHeader VariableHeader // 可变头部
	Payload        interface{}          // 数据
}

func NewMqttPacket(header *FixedHeader, variableHeader VariableHeader, payload interface{}) *MqttPacket {
	return &MqttPacket{VariableHeader: variableHeader, Payload: payload, FixedHeader: header}
}

type FixedHeader struct {
	PacketType      PacketType
	DUP             bool
	QoS             QoSLevel
	Retain          bool
	RemainingLength uint64 // 控制报文总长度等于固定报头的长度加上剩余长度
}

func NewFixedHeader(packetType PacketType, dup bool, qos QoSLevel, retain bool) *FixedHeader {
	return &FixedHeader{PacketType: packetType, DUP: dup, QoS: qos, Retain: retain}
}

type VariableHeader interface {
	GetProps() *Props
}

type Codec interface {
	Encode() ([]byte, error)
	Decode(data []byte) ([]byte,error)
}

type MQTTCodec struct {
}

// 解码消息
func (m *MQTTCodec) DecodePacket(connContext tgo.ConnContext) (interface{}, error) {
	conn := connContext.GetConn()
	 header,err := m.decodeFixedHeader(conn)
	 if err!=nil {
	 	return nil,err
	 }
	 body := make([]byte,header.RemainingLength)
	_, err = conn.Read(body)
	if err!=nil {
		return nil,err
	}
	var vHeader VariableHeader
	var payload interface{}
	switch header.PacketType {
	case CONNECT:
		vHeader,payload,err = m.decodeConnect(body)
		if err!=nil {
			return nil,err
		}
	}

	return NewMqttPacket(header,vHeader,payload), nil
}

// 编码消息
func (m *MQTTCodec) EncodePacket(packet interface{}) ([]byte, error) {
	var mqttObj = packet.(*MqttPacket)

	var bodyBytes []byte
	var err error
	switch mqttObj.PacketType {
	case CONNECT:
		bodyBytes,err = m.encodeConnect(mqttObj)
		if err!=nil {
			return nil,err
		}
	}

	enc := newEncoder()

	// FixedHeader
	headerBytes, err := m.encodeFixedHeader(mqttObj.FixedHeader, uint64(len(bodyBytes)))
	if err != nil {
		return nil, err
	}
	enc.WriteBytes(headerBytes)
	enc.WriteBytes(bodyBytes)

	return enc.Bytes(), nil
}

func (m *MQTTCodec) encodeFixedHeader(f *FixedHeader, remainingLength uint64) ([]byte, error) {
	header := []byte{byte(int(f.PacketType<<4) | encodeBool(f.DUP)<<3 | int(f.QoS)<<1 | encodeBool(f.Retain))}
	varHeader := encodeVariable(remainingLength)

	return append(header, varHeader...), nil
}

func (m *MQTTCodec) decodeFixedHeader(conn tgo.Conn) (*FixedHeader, error) {

	b := make([]byte, 1)
	_, err := io.ReadFull(conn, b)
	if err != nil {
		return nil, err
	}
	typeAndFlags := b[0]
	fh := &FixedHeader{}
	fh.PacketType = PacketType(typeAndFlags >> 4)
	fh.DUP = (typeAndFlags>>3)&0x01 > 0
	fh.QoS = QoSLevel((typeAndFlags >> 1) & 0x03)
	fh.Retain = typeAndFlags&0x01 > 0
	fh.RemainingLength = uint64(decodeLength(conn))
	return fh,nil
}
