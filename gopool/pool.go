package gopool

import "time"

type Pool struct {
	work chan func()
	sem  chan struct{}
}

func New(size int) *Pool {
	return &Pool{
		work: make(chan func()),
		sem:  make(chan struct{}, size),
	}
}

func (p *Pool) Schedule(task func()) error {
	select {
	case p.work <- task:
	case p.sem <- struct{}{}:
		go p.worker(task)
	}
	return nil
}

func (p *Pool) ScheduleTimeout(timeout time.Duration, task func()) error {
	ch := time.After(timeout)
	select {
	case p.work <- task:
	case p.sem <- struct{}{}:
		go p.worker(task)
	case <-ch:
	}
	return nil
}

func (p *Pool) worker(task func()) {
	// defer func() { <-p.sem }
	for {
		task()
		task = <-p.work
	}
}
