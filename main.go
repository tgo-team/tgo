package main

import (
	"github.com/tgo-team/tgo/src"
)

func main() {

	// 创建TGO
	tg := tgo.New()
	// 指定server
	tg.UseServer(tgo.NewServerTCP(tgo.TcpHandshake(func(packet interface{}, conn tgo.StatefulConn) bool {
		return true
	}), tgo.TcpAddr("0.0.0.0:0")))
	// 指定包处理者
	tg.UseHandler(func(ctx tgo.Context) {

	})

	// 开启TGO
	tg.Run()
}
