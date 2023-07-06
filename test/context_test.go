package test

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestContext(t *testing.T) {
	ctx, cancle := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("退出监听错误协程")
				return
			default:
				fmt.Println("逻辑处理中...")
			}
		}
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("退出监听re-balance协程")
				return
			default:
				fmt.Println("逻辑处理中...")
			}
		}
	}()

	time.Sleep(time.Millisecond)
	// 调用cancelFunc, 结束消费
	cancle()
}
