package fastpool

import "github.com/xiaonanln/go-lockfree-queue"

type FastPool struct {
	q   *lfqueue.Queue
	New func() interface{}
}

func NewFastPool(capacity int, New func() interface{}) *FastPool {
	fp := &FastPool{
		q:   lfqueue.NewQueue(capacity),
		New: New,
	}
	return fp
}

func (fp *FastPool) Put(x interface{}) {
	fp.q.Put(x)
}

func (fp *FastPool) Get() interface{} {
	if x, ok := fp.q.Get(); ok {
		return x
	} else {
		return fp.New()
	}
}
