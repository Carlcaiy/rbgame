package main

import (
	"fmt"
	"math/rand"
	"time"
)

func shift_up(vals []XXXX, curr int) {
	if curr == 0 {
		return
	}
	for curr > 0 {
		root := (curr - 1) / 2
		if vals[curr] < vals[root] {
			vals[curr], vals[root] = vals[root], vals[curr]
			curr = root
			continue
		}
		break
	}
}

func shift_down(vals []XXXX) {
	root := 0
	for {
		left := root*2 + 1
		right := root*2 + 2
		if left < len(vals) {
			son := root
			if vals[son] > vals[left] {
				son = left
			}
			if right < len(vals) && vals[son] > vals[right] {
				son = right
			}
			if son != root {
				vals[son], vals[root] = vals[root], vals[son]
				root = son
				continue
			}
		}
		break
	}
}

type XXXX = int

func main() {
	fmt.Println(time.Now().Unix(), time.Now().UnixMilli(), time.Now().UnixMicro(), time.Now().UnixNano())
	heap := make([]XXXX, 0, 100)
	for i := 0; i < 100; i++ {
		heap = append(heap, rand.Intn(100))
		shift_up(heap, len(heap)-1)
	}

	sli := make([]byte, 100)
	for i := range sli {
		sli[i] = byte(i)
	}
	nameLen := 20
	for pos := range sli {
		if len(sli) < pos+nameLen {
			continue
		}
		fmt.Println(string(sli[pos : pos+nameLen-1]))
	}

	peek := func() int {
		if len(heap) > 0 {
			return heap[0]
		}
		return -1
	}

	pop := func() {
		heap[0], heap[len(heap)-1] = heap[len(heap)-1], heap[0]
		heap = heap[:len(heap)-1]
		shift_down(heap)
	}

	for v := peek(); len(heap) > 0 && v > -1; pop() {
		print(heap)
		fmt.Println("------------------------------------")
	}
}

func print(heap []XXXX) {
	l := 0
	for i := 1; ; i *= 2 {
		r := l + i
		if l+i >= len(heap) {
			r = len(heap)
		}
		fmt.Println(heap[l:r])
		l = r
		if l >= len(heap) {
			break
		}
	}
}
