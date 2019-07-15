package mqtt

import "github.com/pkg/errors"

func NewConnect(variableHeader *ConnectVariableHeader, payload *ConnectPayload) *MqttPacket {

	return NewMqttPacket(NewFixedHeader(CONNECT, false, QoS0, false), variableHeader, payload)
}

type ConnectProps struct {
	SessionExpiryInterval      uint32
	AuthenticationMethod       string
	AuthenticationData         []byte
	RequestProblemInformation  bool
	RequestResponseInformation bool
	ReceiveMaximum             uint16
	TopicAliasMaximum          uint16
	UserProperty               map[string]string
	MaximumPacketSize          uint32
}

func (c *ConnectProps) ToProps() *Props {
	return &Props{
		SessionExpiryInterval:      c.SessionExpiryInterval,
		AuthenticationMethod:       c.AuthenticationMethod,
		AuthenticationData:         c.AuthenticationData,
		RequestProblemInformation:  c.RequestProblemInformation,
		RequestResponseInformation: c.RequestResponseInformation,
		ReceiveMaximum:             c.ReceiveMaximum,
		TopicAliasMaximum:          c.TopicAliasMaximum,
		UserProperty:               c.UserProperty,
		MaximumPacketSize:          c.MaximumPacketSize,
	}
}

type WillProps struct {
	PayloadFormatIndicator uint8
	MessageExpiryInterval  uint32
	ContentType            string
	ResponseTopic          string
	CorrelationData        []byte
	WillDelayInterval      uint32
	UserProperty           map[string]string
}

func (w *WillProps) ToProps() *Props {
	return &Props{
		PayloadFormatIndicator: w.PayloadFormatIndicator,
		MessageExpiryInterval:  w.MessageExpiryInterval,
		ContentType:            w.ContentType,
		ResponseTopic:          w.ResponseTopic,
		CorrelationData:        w.CorrelationData,
		WillDelayInterval:      w.WillDelayInterval,
		UserProperty:           w.UserProperty,
	}
}

type ConnectVariableHeader struct {
	// Protocol defintions
	ProtocolName    string
	ProtocolVersion uint8

	// Connection flags
	FlagUsername bool
	FlagPassword bool
	WillRetain   bool
	WillQoS      QoSLevel
	FlagWill     bool
	CleanStart   bool
	KeepAlive    uint16
	// Connection properties
	Props *ConnectProps
}


func (c *ConnectVariableHeader) GetProps() *Props {
	return c.Props.ToProps()
}

type ConnectPayload struct {
	// Payloads
	ClientId    string
	WillProps   *WillProps
	WillTopic   string
	WillPayload string
	Username    string
	Password    []byte
}



func (m *MQTTCodec) encodeConnect(packet *MqttPacket) ([]byte,error) {

	// ConnectVariableHeader
	vHeader := packet.VariableHeader.(*ConnectVariableHeader)
	enc := newEncoder()
	enc.WriteString(vHeader.ProtocolName)
	enc.WriteUint8(vHeader.ProtocolVersion)
	flag := encodeBool(vHeader.FlagUsername)<<7 | encodeBool(vHeader.FlagPassword)<<6 | encodeBool(vHeader.WillRetain)<<5 | int(vHeader.WillQoS)<<3 | encodeBool(vHeader.FlagWill)<<2 | encodeBool(vHeader.CleanStart)<<1
	enc.WriteInt(flag)
	enc.WriteUint16(vHeader.KeepAlive)
	if vHeader.Props != nil {
		enc.WriteProperty(vHeader.Props.ToProps())
	} else {
		enc.WriteUint8(0)
	}


	// ConnectPayload
	payload := packet.Payload.(*ConnectPayload)
	enc.WriteString(payload.ClientId)
	if payload.WillProps != nil {
		enc.WriteProperty(payload.WillProps.ToProps())
	} else {
		enc.WriteUint8(0)
	}
	if payload.WillTopic != "" {
		enc.WriteString(payload.WillTopic)
	}
	if payload.WillPayload != "" {
		enc.WriteString(payload.WillPayload)
	}
	if payload.Username != "" {
		enc.WriteString(payload.Username)
	}
	if payload.Password != nil && len(payload.Password) > 0 {
		enc.WriteBinary(payload.Password)
	}

	return enc.Bytes(),nil
}

func (m *MQTTCodec) decodeConnect(data []byte) (VariableHeader,*ConnectPayload,error) {
	dec := newDecoder(data)

	vHeader := &ConnectVariableHeader{}

	var err error
	if vHeader.ProtocolName, err = dec.String(); err != nil {
		return  nil,nil,errors.Wrap(err, "failed to decode as string")
	}
	if vHeader.ProtocolVersion, err = dec.Uint8(); err != nil {
		return  nil,nil,errors.Wrap(err, "failed to decode as int")
	}
	var flag int
	if flag, err = dec.Int(); err != nil {
		return nil,nil,errors.Wrap(err, "failed to decode as int")
	}
	vHeader.FlagUsername = ((flag >> 7) & 0x01) > 0
	vHeader.FlagPassword = ((flag >> 6) & 0x01) > 0
	vHeader.WillRetain = ((flag >> 5) & 0x01) > 0
	vHeader.WillQoS = QoSLevel((flag >> 3) & 0x03)
	vHeader.FlagWill = ((flag >> 2) & 0x01) > 0
	vHeader.CleanStart = ((flag >> 1) & 0x01) > 0
	if vHeader.KeepAlive, err = dec.Uint16(); err != nil {
		return nil,nil,errors.Wrap(err, "failed to decode as uint16")
	}
	// Connection variable properties enables on v5
	if prop, err := dec.Property(); err != nil {
		return nil,nil,errors.Wrap(err, "failed to decode property")
	} else if prop != nil {
		vHeader.Props = prop.ToConnect()
	}

	payload := &ConnectPayload{}

	if payload.ClientId, err = dec.String(); err != nil {
		return nil,nil, errors.Wrap(err, "failed to decode as string")
	}
	// Will properties enables on v5
	if prop, err := dec.Property(); err != nil {
		return nil,nil, errors.Wrap(err, "failed to decode will property")
	} else if prop != nil {
		payload.WillProps = prop.ToWill()
	}
	if vHeader.FlagWill {
		if payload.WillTopic, err = dec.String(); err != nil {
			return nil,nil, errors.Wrap(err, "failed to decode will topic")
		}
		if payload.WillPayload, err = dec.String(); err != nil {
			return nil,nil, errors.Wrap(err, "failed to decode will payload")
		}
	}
	if vHeader.FlagUsername {
		if payload.Username, err = dec.String(); err != nil {
			return nil,nil, errors.Wrap(err, "failed to decode username")
		}
	}
	if vHeader.FlagPassword {
		if payload.Password, err = dec.Binary(); err != nil {
			return nil,nil, errors.Wrap(err, "failed to decode password")
		}
	}
	return vHeader,payload,nil

}