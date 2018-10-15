package fastpool

import (
	"os"
	"runtime/pprof"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

const (
	concurrencyCount = 5000
	testDuration     = time.Second * 10
)

type testObject struct {
	b [1024]byte
}

func newTestObject() interface{} {
	return &testObject{}
}

type pool interface {
	Put(x interface{})
	Get() (x interface{})
}

func TestSyncPool(t *testing.T) {
	p := &sync.Pool{
		New: newTestObject,
	}
	testPool(t, "sync.Pool.pprof", p)
}

func TestFastPool(t *testing.T) {
	p := NewFastPool(1024, newTestObject)
	testPool(t, "FastPool.pprof", p)
}

func testPool(t *testing.T, profileName string, p pool) {
	w, err := os.OpenFile(profileName, os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		t.Fatal(err)
	}
	pprof.StartCPUProfile(w)
	defer pprof.StopCPUProfile()
	var waitDone sync.WaitGroup
	startTestTime := time.Now()
	stopTestTime := startTestTime.Add(testDuration)
	waitDone.Add(concurrencyCount)
	var putCount, getCount uint64
	for i := 0; i < concurrencyCount; i++ {
		go func() {
			defer waitDone.Done()
			for time.Now().Before(stopTestTime) {
				x := p.Get()
				atomic.AddUint64(&getCount, 1)
				time.Sleep(time.Millisecond)
				p.Put(x)
				atomic.AddUint64(&putCount, 1)
			}
		}()
	}
	waitDone.Wait()
	t.Logf("Put %d, Get %d", putCount, getCount)
}
