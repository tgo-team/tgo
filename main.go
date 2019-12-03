package main

import (
	"net"

	tgo "github.com/tgo-team/tgo/src"
)

func main() {

	// 创建TGO
	tg := tgo.New()
	// 指定server
	tg.UseServer(tgo.NewServerTCP(tgo.TCPHandshake(func(packet interface{}, conn net.Conn, ctx *tgo.ServerContext) (s string, e error) {
		return "1", nil
	}), tgo.TCPAddr("0.0.0.0:0")))
	// 指定包处理者
	tg.UseHandler(func(ctx tgo.Context) {

	})

	// 开启TGO
	tg.Run()
}
