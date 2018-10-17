package lfpool

import (
	"github.com/xiaonanln/go-lockfree-queue"
)

type Pool struct {
	q   *lfqueue.Queue
	new func() interface{}
}

func NewFastPool(capacity int, New func() interface{}) *Pool {
	fp := &Pool{
		q:   lfqueue.NewQueue(capacity),
		new: New,
	}
	return fp
}

func (fp *Pool) Put(x interface{}) {
	fp.q.Put(x)
}

func (fp *Pool) Get() interface{} {
	if x, ok := fp.q.Get(); ok {
		return x
	} else {
		return fp.new()
	}
}
