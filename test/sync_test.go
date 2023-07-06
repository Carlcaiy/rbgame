package test

import (
	"fmt"
	"math/rand"
	"net"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

var createNum int32

func createBuffer() interface{} {
	atomic.AddInt32(&createNum, 1)
	buffer := make([]byte, 1024)
	fmt.Println()
	return buffer
}

/*
sync.Pool最重要的是减少GC负担，并且是并发安全的
特征如下:
池里的元素随时可能释放掉，释放策略完全由runtime内部管理
Get获取到的对象可能是刚创建的，也可能是之前创建好cache的，使用者无法区分
池里的元素个数你无法知道
*/
func TestPool(t *testing.T) {
	pool := &sync.Pool{New: createBuffer}
	workerPool := 1024 * 1024
	var wg sync.WaitGroup
	wg.Add(workerPool)

	for i := 0; i < workerPool; i++ {
		go func() {
			defer wg.Done()
			buffer := pool.Get()
			// _ := buffer.([]byte)
			defer pool.Put(buffer)
		}()
	}
	wg.Wait()
	fmt.Printf("%d buffer objects were create\n", createNum)
	time.Sleep(5 * time.Second)
}

var cond sync.Cond
var num int64
var ato int64

func consumer(in <-chan int, index int) {
	for {
		cond.L.Lock()
		// for len(in) == 0 {
		for num <= 0 {
			// fmt.Println(index, "len == 0")
			cond.Wait()
			// fmt.Println(index, "len == 0 recv sig")
		}
		// time.Sleep(time.Microsecond)
		// num := <-in
		num--
		atomic.AddInt64(&ato, -1)
		if num != ato {
			fmt.Println("*********************")
		}
		// fmt.Println("消费者：", index, ato, num)
		cond.L.Unlock()
		cond.Signal()
	}
}

func producer(out chan<- int, index int) {
	for {
		cond.L.Lock()
		// for len(out) == 10 {
		for num >= 40 {
			// fmt.Println(index, "len == 10")
			cond.Wait()
			// fmt.Println(index, "len == 10 recv sig")
		}
		// num := rand.Intn(800)
		// time.Sleep(time.Microsecond)
		// out <- num
		num++
		atomic.AddInt64(&ato, +1)
		if num != ato {
			fmt.Println("*********************")
		}
		// fmt.Println("生产者：", index, ato, num)

		cond.L.Unlock()
		cond.Signal()
	}
}

func TestCond(t *testing.T) {
	ch := make(chan int, 10)
	rand.Seed(time.Now().UnixMilli())
	cond.L = new(sync.Mutex)

	for i := 1; i <= 4; i++ {
		go consumer(ch, i)
	}

	time.Sleep(time.Millisecond)
	for i := 1; i <= 4; i++ {
		go producer(ch, i)
	}

	quit := make(chan []struct{})
	<-quit
}

func TestEAgain(t *testing.T) {
	ln, _ := net.ListenTCP("tcp", &net.TCPAddr{Port: 8080})
	for {
		conn, _ := ln.Accept()
		for {
			xx := []byte{}
			conn.Read(xx)
		}
	}
}
