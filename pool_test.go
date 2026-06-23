package lutil

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

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
	pool := NewPool(1, 1, DiscardPolicy)
	block := make(chan struct{})
	pool.Submit(func() { <-block })
	err := pool.SubmitErr(func() {})
	assert.Error(t, err)
	close(block)
	pool.Shutdown()
}

func TestCallerRunsPolicy(t *testing.T) {
	pool := NewPool(1, 1, CallerRunsPolicy)
	block := make(chan struct{})
	pool.Submit(func() { <-block })
	var ran int32
	pool.Submit(func() { atomic.StoreInt32(&ran, 1) })
	assert.Equal(t, int32(1), atomic.LoadInt32(&ran))
	close(block)
	pool.Shutdown()
}

func TestDiscardOldestPolicy(t *testing.T) {
	pool := NewPool(1, 1, DiscardOldestPolicy)
	block := make(chan struct{})
	pool.Submit(func() { <-block })
	done := make(chan struct{})
	pool.Submit(func() { close(done) })
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("timeout")
	}
	close(block)
	pool.Shutdown()
}
