package mqtt

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/tgo-team/tgo-core/src"
	"testing"
)

func TestConnectEncodeAndDecode(t *testing.T) {
	codec := &MQTTCodec{}
	vHeader := &ConnectVariableHeader{
		ProtocolName:    "MQTT",
		ProtocolVersion: 1,
		FlagUsername:    true,
		FlagPassword:    true,
		FlagWill:        true,
		Props: &ConnectProps{
			SessionExpiryInterval:100,
			UserProperty: map[string]string{
				"test":"111",
			},
		},
	}
	payload := &ConnectPayload{
		ClientId:    "dddd",
		Username:    "admin",
		Password:    []byte("admin"),
		WillTopic:   "test",
		WillPayload: "d12",
	}
	data, err := codec.EncodePacket(NewConnect(vHeader, payload))
	if err != nil {
		panic(err)
	}

	fmt.Println(data)

	cData := bytes.NewBuffer(data)
	packet, err := codec.DecodePacket(tgo.NewDefaultConnContext(cData))
	assert.NoError(t, err)

	vHeader = packet.(*MqttPacket).VariableHeader.(*ConnectVariableHeader)
	payload = packet.(*MqttPacket).Payload.(*ConnectPayload)

	assert.Equal(t, vHeader.ProtocolName, "MQTT")
	assert.Equal(t, vHeader.ProtocolVersion, uint8(1))
	assert.Equal(t, payload.Username, "admin")
	assert.Equal(t, payload.Password, []byte("admin"))
	assert.Equal(t, payload.ClientId, "dddd")
	assert.Equal(t, payload.WillPayload, "d12")
	assert.Equal(t,vHeader.Props.SessionExpiryInterval,uint32(100))

	assert.Equal(t,vHeader.Props.UserProperty["test"],"111")

	//time.Sleep(time.Second*5)
}
