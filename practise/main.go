package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"time"
)

const WIDTH = 5

var arr []uint8 = []uint8{
	1, 1, 1, 1, 1,
	2, 3, 3, 3, 3,
	4, 4, 4, 4, 4,
}

func main() {
	go func() {
		err := http.ListenAndServe(":9909", nil)
		if err != nil {
			panic(err)
		}
	}()
	arr = arr[5:6]
	fmt.Println(arr)
	arr = arr[0:5]
	fmt.Println(arr)
	arr32 := []int32{
		1, 12, 44, 55, 66,
	}
	fmt.Println(string(arr32[0]))
	for {
		time.Sleep(time.Second)
		fmt.Println(slitostring(arr32))
	}
}

func slitostring(sli interface{}) string {
	str := ""
	switch get := sli.(type) {
	case []uint8:
		for i, v := range get {
			if i == 0 {
				str += fmt.Sprintf("%d", v)
			} else {
				str += fmt.Sprintf(",%d", v)
			}
		}
	case []int32:
		for i, v := range get {
			if i == 0 {
				str += fmt.Sprintf("%d", v)
			} else {
				str += fmt.Sprintf(",%d", v)
			}
		}
	default:
		fmt.Println("slitostring unsupported")
	}
	return str
}
