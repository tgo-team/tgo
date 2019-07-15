# tgo-core

## 简介

一款简单，高效,扩展性极强的现代通讯服务器，适用于即时通讯，物联网通讯，AI智能等等

## 启动TGO

```
// 创建TGO
tg := tgo.New(tgo.NewOptions())

...

// 运行TGO
tg.Run()
```

## 自定义通信服务

```
tg.UseServer(tgo.NewServerTCP())
```
### 自定义数据协议

```
tg.UseProtocol(tgo.NewProtocolMQTT())
```


## 设置包处理者

```
tg.UseHandler(func(ctx tgo.Context) {

})
```

## 自定义路由

```
tg.UseRouter(router)
```