package tgo

import (
	"context"
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
	"time"
)

func TestRouter(t *testing.T)  {
	tg,addr := NewTestTGO()
	c, err := net.DialUnix("unix", nil, addr.(*net.UnixAddr))
	if err!=nil {
		panic(err)
	}
	ct,cancelFnc := context.WithTimeout(context.Background(),time.Millisecond*300)

	result :=""
	router := NewDefaultRouter(func(ctx Context, handlerMap map[interface{}]Handler) Handler {
		return  handlerMap[ctx.GetPacket()]
	})
	router.Route("hello", func(ctx Context) {
		result = ctx.GetPacket().(string)
		cancelFnc()
	})
	tg.UseRouter(router)

	tg.Start()
	c.Write([]byte("hello"))
	<-ct.Done()
	assert.Equal(t,"hello",result)

}
