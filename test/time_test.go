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

func TestRound(t *testing.T) {
	now := time.Now()
	RegisterTime := now.Unix()
	fmt.Println(RegisterTime)
	for i := 0; i < 10; i++ {
		xx := GetNextDayUnix(RegisterTime, int32(i))
		dx := GetDayIndex(RegisterTime, xx)
		fmt.Println(i, xx, dx)
	}

}

// 获取第n天00:00:00的时间戳
func GetNextDayUnix(curr int64, day int32) int64 {
	date := time.Unix(curr, 0)
	date = time.Date(date.Year(), date.Month(), date.Day()+int(day), 0, 0, 0, 0, time.Local)
	return date.Unix()
}

// 计算当前是第几天
func GetDayIndex(registerTime int64, now int64) int64 {
	var (
		begin time.Time = time.Unix(registerTime, 0)
		end   time.Time = time.Unix(now, 0)
	)
	begin = time.Date(begin.Year(), begin.Month(), begin.Day(), 0, 0, 0, 0, time.Local)
	end = time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, time.Local)
	return (end.Unix()-begin.Unix())/86400 + 1
}
