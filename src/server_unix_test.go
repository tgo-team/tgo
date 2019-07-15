package tgo

import (
	"fmt"
	"net"
	"testing"
	"time"
)

func TestServerUnix(t *testing.T)  {
	tg,addr := NewTestTGO()
	c, err := net.DialUnix("unix", nil, addr.(*net.UnixAddr))
	if err!=nil {
		panic(err)
	}

	tg.UseHandler(func(ctx Context) {
		fmt.Println("dddddzz",ctx.GetPacket())
	})

	tg.Start()

	c.Write([]byte("zdds"))

	time.Sleep(time.Millisecond*200)


}
