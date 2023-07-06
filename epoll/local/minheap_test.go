package local

import (
	"fmt"
	"testing"
)

type XXX struct {
	val   int
	index int
}

func (x XXX) Val() int64 {
	return int64(x.val)
}

func (x *XXX) SetIndex(i int) {
	x.index = i
}

func (x XXX) Index() int {
	return x.index
}

func NewXXX(v int) *XXX {
	return &XXX{
		val: v,
	}
}

func TestHeap(t *testing.T) {
	h := NewMinHeap(50)
	h.Push(NewXXX(5))
	fmt.Println(h.data[:h.len])
	h.Push(NewXXX(8))
	fmt.Println(h.data[:h.len])
	h.Push(NewXXX(6))
	fmt.Println(h.data[:h.len])
	h.Push(NewXXX(3))
	fmt.Println(h.data[:h.len])
	h.Push(NewXXX(10))
	fmt.Println(h.data[:h.len])
	h.Push(NewXXX(1))
	fmt.Println(h.data[:h.len])
	h.Push(NewXXX(2))
	fmt.Println(h.data[:h.len])
	a, _ := h.Pop()
	fmt.Println(a, h.data[:h.len])
	a, _ = h.Pop()
	fmt.Println(a, h.data[:h.len])
	a, _ = h.Pop()
	fmt.Println(a, h.data[:h.len])
	a, _ = h.Pop()
	fmt.Println(a, h.data[:h.len])
	a, _ = h.Pop()
	fmt.Println(a, h.data[:h.len])
	a, _ = h.Pop()
	fmt.Println(a, h.data[:h.len])
	a, _ = h.Pop()
	fmt.Println(a, h.data[:h.len])
	a, _ = h.Pop()
	fmt.Println(a, h.data[:h.len])
}
