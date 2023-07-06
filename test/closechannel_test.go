package test

import (
	"fmt"
	"testing"
	"time"
)

func TestClose(t *testing.T) {
	quit := make(chan struct{})
	go func() {
		// 模拟工作
		fmt.Println("工作中...")
		time.Sleep(3 * time.Second)
		// 关闭退出信号
		close(quit)
	}()

	// 阻塞，等待退出信号被关闭
	ans, ok := <-quit
	ans1, ok1 := <-quit
	ans2, ok2 := <-quit
	ans3, ok3 := <-quit
	fmt.Println("已收到退出信号，退出中...", ans, ok, ans1, ok1, ans2, ok2, ans3, ok3)
}

func TestClosechs(t *testing.T) {
	quit := make(chan int, 10)
	for i := 0; i < 10; i++ {
		quit <- i
	}
	// 阻塞，等待退出信号被关闭
	ans, ok := <-quit
	ans1, ok1 := <-quit
	close(quit)
	ans2, ok2 := <-quit
	ans3, ok3 := <-quit
	for x := range quit {
		fmt.Println(x)
	}
	fmt.Println("已收到退出信号，退出中...", ans, ok, ans1, ok1, ans2, ok2, ans3, ok3)
	x := time.NewTimer(time.Second)
	select {
	case a, ok := <-quit:
		fmt.Println("quit", a, ok)
	case <-x.C:
		fmt.Println("time in")
	}

}
