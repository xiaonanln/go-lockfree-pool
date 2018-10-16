package fastpool

import (
	"github.com/xiaonanln/go-lockfree-queue"
	"sync"
)

type FastPool struct {
	q        *lfqueue.Queue
	fallback sync.Pool
}

func NewFastPool(capacity int, New func() interface{}) *FastPool {
	fp := &FastPool{
		q: lfqueue.NewQueue(capacity),
	}
	fp.fallback.New = New
	return fp
}

func (fp *FastPool) Put(x interface{}) {
	if ok := fp.q.Put(x); !ok {
		fp.fallback.Put(x)
	}
}

func (fp *FastPool) Get() interface{} {
	if x, ok := fp.q.Get(); ok {
		return x
	} else {
		return fp.fallback.Get()
	}
}
