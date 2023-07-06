## pb-broadcast-error__2022-04-04_00-00-00.log

广播服可能会有一些问题，向中心服发送消息的时候，发现连接已经关闭，但是仍然发送消息

``` md
panic: send on closed channel

goroutine 18 [running]:
cynking.com/pkg/network/tcp.(*Conn).doWrite(...)
	/mnt/d/servercode/go/pkg/network/tcp/conn.go:88
cynking.com/pkg/network/tcp.(*Conn).SendMessage(0x914585, {0x9b8e30, 0xc000138050})
	/mnt/d/servercode/go/pkg/network/tcp/conn.go:100 +0x45
cynking.com/game/internal/broadcast/modules/center.(*Module).register.func1.1()
	/mnt/d/servercode/go/game/internal/broadcast/modules/center/module.go:84 +0xb6
created by cynking.com/game/internal/broadcast/modules/center.(*Module).register.func1
	/mnt/d/servercode/go/game/internal/broadcast/modules/center/module.go:78 +0x5f
```# rbgame
