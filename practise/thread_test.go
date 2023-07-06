package main

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
)

// 下面的代码会输出什么？并说明原因。

func TestPack(t *testing.T) {
	runtime.GOMAXPROCS(1)
	wg := sync.WaitGroup{}
	wg.Add(20)
	for i := 0; i < 10; i++ {
		go func() {
			fmt.Println("ai: ", i)
			wg.Done()
		}()
	}
	for i := 0; i < 10; i++ {
		go func(i int) {
			fmt.Println("bi: ", i)
			wg.Done()
		}(i)
	}
	wg.Wait()
}
