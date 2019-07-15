package main

import "github.com/tgo-team/tgo/src"

func main()  {

	// 创建TGO
	tg := tgo.New(tgo.NewOptions())
	// 指定server
	tg.UseServer(tgo.NewServerTCP("0.0.0.0:6666"))
	// 指定包处理者
	tg.UseHandler(func(ctx tgo.Context) {

	})

	// 开启TGO
	tg.Run()
}
