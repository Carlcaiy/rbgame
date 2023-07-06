package local

import (
	"fmt"
	"time"
)

type Timer struct {
	*MinHeap
}

type TimerEvent struct {
	Loop        bool
	duration    time.Duration
	triggerTime time.Time
	event       func()
	index       int
}

func NewTimer(n int) *Timer {
	return &Timer{
		MinHeap: NewMinHeap(n),
	}
}

func (t *TimerEvent) Val() int64 {
	return t.triggerTime.UnixNano()
}

func (t *TimerEvent) Index() int {
	return t.index
}

func (t *TimerEvent) SetIndex(i int) {
	t.index = i
}

func (t *Timer) Push(e *TimerEvent) {
	t.MinHeap.Push(e)
}

func (t *Timer) Get() *TimerEvent {
	return t.MinHeap.Top().(*TimerEvent)
}

func (t *Timer) Len() int {
	return t.MinHeap.Len()
}

func (t *Timer) Pop() *TimerEvent {
	if data, err := t.MinHeap.Pop(); err != nil {
		fmt.Println(err)
		return nil
	} else {
		return data.(*TimerEvent)
	}
}

func (t *Timer) FrameCheck() {
	for t.Len() > 0 && t.Get().Val() <= time.Now().UnixNano() {
		if e := t.Pop(); e != nil {
			e.event()
			if e.Loop {
				e.triggerTime = e.triggerTime.Add(e.duration)
				t.Push(e)
			}
		}
	}
}

func (t *Timer) Stop(e *TimerEvent) {
	t.Drop(e)
}

func (t *Timer) Start(e *TimerEvent) {
	if e.index > 0 {
		fmt.Println("timer hava started, use reset")
		return
	}
	e.triggerTime = time.Now().Add(e.duration)
	t.Push(e)
}

func (t *Timer) Reset(e *TimerEvent) {
	t.Drop(e)
	e.triggerTime = time.Now().Add(e.duration)
	t.Push(e)
}

func NewLoopEvent(dur time.Duration, f func()) *TimerEvent {
	return &TimerEvent{
		duration: dur,
		Loop:     true,
		event:    f,
		index:    -1,
	}
}

func NewDelayEvent(dur time.Duration, f func()) *TimerEvent {
	return &TimerEvent{
		duration: dur,
		Loop:     false,
		event:    f,
		index:    -1,
	}
}
