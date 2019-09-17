package main

import (
	"github.com/tgo-team/tgo/src"
	"net"
)

func main() {

	// 创建TGO
	tg := tgo.New()
	// 指定server
	tg.UseServer(tgo.NewServerTCP(tgo.TcpHandshake(func(packet interface{}, conn net.Conn,ctx *tgo.ServerContext) (error,uint64) {
		return nil,1
	}), tgo.TcpAddr("0.0.0.0:0")))
	// 指定包处理者
	tg.UseHandler(func(ctx tgo.Context) {

	})

	// 开启TGO
	tg.Run()
}
