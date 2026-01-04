package queue

import (
	"context"
	"sync"
)

type Task func(context.Context) error

type Queue struct {
	mu     sync.Mutex
	closed bool
	ch     chan Task
}

func New(size int) *Queue {
	if size < 1 {
		size = 1
	}
	return &Queue{ch: make(chan Task, size)}
}

func (q *Queue) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case task, ok := <-q.ch:
			if !ok {
				return
			}
			if task != nil {
				_ = task(ctx)
			}
		}
	}
}

func (q *Queue) Enqueue(task Task) bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.closed {
		return false
	}
	q.ch <- task
	return true
}

func (q *Queue) Close() {
	q.mu.Lock()
	if q.closed {
		q.mu.Unlock()
		return
	}
	q.closed = true
	close(q.ch)
	q.mu.Unlock()
}
