package lutil

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKeyLock(t *testing.T) {
	req := require.New(t)
	kl := NewKeyLock(8)
	req.NotNil(kl)

	var count int32
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			kl.Lock("same")
			atomic.AddInt32(&count, 1)
			time.Sleep(5 * time.Millisecond)
			kl.Unlock("same")
		}()
	}
	wg.Wait()
	assert.Equal(t, int32(10), count)
}

func TestKeyLock_unlockPanic(t *testing.T) {
	kl := NewKeyLock(4)
	assert.Panics(t, func() { kl.Unlock("missing") })
}
