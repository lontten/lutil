package lutil

import (
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// occupyPool 让 1 个 worker 阻塞，并填满 queueSize 容量的队列，返回解除阻塞的函数。
func occupyPool(t *testing.T, pool *Pool, queueSize int) (unblock func()) {
	t.Helper()
	block := make(chan struct{})
	started := make(chan struct{})
	pool.Submit(func() {
		close(started)
		<-block
	})
	select {
	case <-started:
	case <-time.After(time.Second):
		t.Fatal("worker did not start")
	}
	for i := 0; i < queueSize; i++ {
		if err := pool.SubmitErr(func() {}); err != nil {
			t.Fatalf("failed to fill queue at %d: %v", i, err)
		}
	}
	return func() { close(block) }
}

func TestPoolSubmitAndShutdown(t *testing.T) {
	pool := NewPool(2, 4, nil)
	var n int32
	done := make(chan struct{})
	pool.Submit(func() {
		atomic.AddInt32(&n, 1)
		close(done)
	})
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("task timeout")
	}
	pool.Shutdown()
	assert.Equal(t, int32(1), atomic.LoadInt32(&n))
}

func TestPoolSubmitErr_full(t *testing.T) {
	const queueSize = 1
	pool := NewPool(1, queueSize, DiscardPolicy)
	unblock := occupyPool(t, pool, queueSize)

	err := pool.SubmitErr(func() {})
	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrQueueFull))

	unblock()
	pool.Shutdown()
}

func TestPoolSubmitErr_closed(t *testing.T) {
	pool := NewPool(1, 1, DiscardPolicy)
	pool.Shutdown()
	err := pool.SubmitErr(func() {})
	assert.True(t, errors.Is(err, ErrPoolClosed))
	pool.Submit(func() {}) // 关闭后不得 panic
	pool.Shutdown()         // 重复关闭安全
}

func TestCallerRunsPolicy(t *testing.T) {
	const queueSize = 1
	pool := NewPool(1, queueSize, CallerRunsPolicy)
	unblock := occupyPool(t, pool, queueSize)

	var ran int32
	pool.Submit(func() { atomic.StoreInt32(&ran, 1) })
	assert.Equal(t, int32(1), atomic.LoadInt32(&ran))

	unblock()
	pool.Shutdown()
}

func TestDiscardOldestPolicy(t *testing.T) {
	const queueSize = 1
	pool := NewPool(1, queueSize, DiscardOldestPolicy)
	block := make(chan struct{})
	started := make(chan struct{})
	pool.Submit(func() {
		close(started)
		<-block
	})
	<-started

	var oldRan, newRan int32
	require.NoError(t, pool.SubmitErr(func() { atomic.StoreInt32(&oldRan, 1) }))
	pool.Submit(func() { atomic.StoreInt32(&newRan, 1) }) // 应丢弃最老任务

	close(block)
	pool.Shutdown()
	assert.Equal(t, int32(0), atomic.LoadInt32(&oldRan))
	assert.Equal(t, int32(1), atomic.LoadInt32(&newRan))
}

func TestDiscardPolicy(t *testing.T) {
	const queueSize = 1
	pool := NewPool(1, queueSize, DiscardPolicy)
	unblock := occupyPool(t, pool, queueSize)

	var ran int32
	pool.Submit(func() { atomic.StoreInt32(&ran, 1) })
	assert.Equal(t, int32(0), atomic.LoadInt32(&ran))

	unblock()
	pool.Shutdown()
}

func TestWorkerRecoversFromPanic(t *testing.T) {
	pool := NewPool(1, 2, DiscardPolicy)
	pool.Submit(func() { panic("boom") })

	done := make(chan struct{})
	pool.Submit(func() { close(done) })
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("worker did not recover after panic")
	}
	pool.Shutdown()
}
