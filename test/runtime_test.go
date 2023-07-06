package test

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"sync"
	"testing"
	"time"
)

const runtime_num = 10000

func cal() {
	for i := 0; i < runtime_num; i++ {
		// fmt.Println("cal", i)
		runtime.Gosched()
	}
}

func TestGosched(t *testing.T) {
	runtime.GOMAXPROCS(2)
	start := time.Now().UnixNano()
	go cal()
	for i := 0; i < runtime_num; i++ {
		// fmt.Println("main", i)
		runtime.Gosched()
	}
	end := time.Now().UnixNano()
	fmt.Printf("tatal %vns per %vns \n", end-start, (end-start)/runtime_num)

	debug.SetGCPercent(-1)
}

func createLargeNumGoroutine(runtime_num int, wg *sync.WaitGroup) {
	wg.Add(runtime_num)
	for i := 0; i < runtime_num; i++ {
		go func() {
			defer wg.Done()
		}()
	}
}

func TestRuntime(te *testing.T) {
	// 只设置一个 Processor 保证 Go 程串行执行
	runtime.GOMAXPROCS(1)
	// 关闭GC改为手动执行
	debug.SetGCPercent(-1)

	var wg sync.WaitGroup
	createLargeNumGoroutine(10000, &wg)
	wg.Wait()
	t := time.Now()
	runtime.GC() // 手动GC
	cost := time.Since(t)
	fmt.Printf("GC cost %v when goroutine runtime_num is %v\n", cost, 10000)

	createLargeNumGoroutine(100000, &wg)
	wg.Wait()
	t = time.Now()
	runtime.GC() // 手动GC
	cost = time.Since(t)
	fmt.Printf("GC cost %v when goroutine runtime_num is %v\n", cost, 100000)

	createLargeNumGoroutine(1000000, &wg)
	wg.Wait()
	t = time.Now()
	runtime.GC() // 手动GC
	cost = time.Since(t)
	fmt.Printf("GC cost %v when goroutine runtime_num is %v\n", cost, 100000)
}
