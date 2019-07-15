package tgo

import (
	"bufio"
	"net"
	"os"
)

func NewTestTGO() (*TGO,net.Addr) {
	fileName := "test.socket"
	os.Remove(fileName)
	// 创建TGO
	tg := New(NewOptions())
	s := NewServerUnix(fileName)
	// 指定server
	tg.UseServer(s)
	tg.UseProtocol(&TestPro{})
	return tg,s.addr
}


type TestPro struct {

}

// 解码消息
func (t *TestPro) DecodePacket(connContext ConnContext) (interface{},error) {
	 testBytes := make([]byte,1024)
	cn,err := bufio.NewReader(connContext.GetConn()).Read(testBytes)
	if err!=nil {
		return nil,err
	}
	return string(testBytes[:cn]),nil
}
// 编码消息
func (t *TestPro) EncodePacket(packet interface{}) ([]byte,error) {
	return []byte(packet.(string)),nil
}