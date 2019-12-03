package tgo

import (
	"bufio"
	"io"
	"net"
	"os"
)

// NewTestTGO NewTestTGO
func NewTestTGO() (*TGO, net.Addr) {
	fileName := "test.socket"
	os.Remove(fileName)
	// 创建TGO
	tg := New()
	s := NewServerUnix(fileName)
	// 指定server
	tg.UseServer(s)
	return tg, s.addr
}

// TestPro TestPro
type TestPro struct {
}

// DecodePacket 解码消息
func (t *TestPro) DecodePacket(conn io.Reader) (interface{}, error) {
	testBytes := make([]byte, 1024)
	cn, err := bufio.NewReader(conn).Read(testBytes)
	if err != nil {
		return nil, err
	}
	return string(testBytes[:cn]), nil
}

// EncodePacket 编码消息
func (t *TestPro) EncodePacket(packet interface{}) ([]byte, error) {
	return []byte(packet.(string)), nil
}
