package test

import (
	"fmt"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"testing"
	"time"
)

func randSeed() {
	rand.Seed(time.Now().UnixNano())
	time.AfterFunc(time.Duration(1+rand.Intn(5))*time.Hour, randSeed)
}

func TestTick(t *testing.T) {
	go func() {
		http.ListenAndServe("172.20.11.69:8080", nil)
	}()

	time.Sleep(time.Second)
	go randSeed()
	ticker := time.NewTicker(time.Second)
	go func() {
		for range ticker.C {
			for i := 0; i < 100; i++ {
				fmt.Println(i)
			}
		}
		fmt.Println("stop")
	}()
	ticker.Stop()
	time.Sleep(time.Second * 10)
}
