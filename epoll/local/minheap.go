package local

import (
	"errors"
)

type IMinHeap interface {
	Val() int64
	SetIndex(int)
	Index() int
}

type MinHeap struct {
	data []IMinHeap
	len  int
}

func NewMinHeap(s int) *MinHeap {
	return &MinHeap{
		data: make([]IMinHeap, s),
		len:  0,
	}
}

func (m *MinHeap) Len() int {
	return m.len
}

func (m *MinHeap) Top() IMinHeap {
	if m.len < 1 {
		return nil
	}
	return m.data[0]
}

func (m *MinHeap) Push(e IMinHeap) error {
	if len(m.data) <= m.len {
		return errors.New("heap full")
	}
	m.data[m.len] = e
	idx := m.len
	m.len++
	m.shiftUp(idx)

	for i := 0; i < m.len; i++ {
		m.data[i].SetIndex(i)
	}

	return nil
}

func (m *MinHeap) Pop() (d interface{}, err error) {
	if m.len <= 0 {
		return nil, errors.New("heap empty")
	}
	d = m.data[0]
	m.len--
	m.data[0], m.data[m.len] = m.data[m.len], m.data[0]
	m.data[m.len] = nil

	m.shiftDown(0)

	for i := 0; i < m.len; i++ {
		m.data[i].SetIndex(i)
	}
	return d, nil
}

func (m *MinHeap) Drop(i IMinHeap) error {
	idx := i.Index()
	if m.len <= idx {
		return errors.New("Drop heap error")
	}
	m.len--
	m.data[idx], m.data[m.len] = m.data[m.len], m.data[idx]
	m.data[m.len] = nil

	m.shiftDown(idx)

	for i := idx; i < m.len; i++ {
		m.data[i].SetIndex(i)
	}
	return nil
}

func (m *MinHeap) Build() {
	for i := m.len / 2; i >= 0; i-- {
		m.shiftDown(0)
	}
	for i := 0; i < m.len; i++ {
		m.data[i].SetIndex(i)
	}
}

func (m *MinHeap) shiftDown(idx int) {
	for idx*2+1 < m.len {
		temp := idx*2 + 1
		if idx*2+2 < m.len {
			if m.data[idx*2+1].Val() > m.data[idx*2+2].Val() {
				temp = idx*2 + 2
			} else {
				temp = idx*2 + 1
			}
		}
		if m.data[idx].Val() > m.data[temp].Val() {
			m.data[idx], m.data[temp] = m.data[temp], m.data[idx]
			idx = temp
		} else {
			break
		}
	}
}

func (m *MinHeap) shiftUp(idx int) {
	for idx > 0 {
		root := (idx - 1) / 2
		if m.data[idx].Val() < m.data[root].Val() {
			m.data[idx], m.data[root] = m.data[root], m.data[idx]
			if idx%2 == 0 && m.data[idx-1].Val() < m.data[root].Val() {
				m.data[idx-1], m.data[root] = m.data[root], m.data[idx-1]
			}
		} else {
			break
		}
		idx = root
	}
}
