package lutil

import (
	"fmt"
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

func TestKeyLock_manyKeysBounded(t *testing.T) {
	const capacity = 8
	kl := NewKeyLock(capacity)
	require.NotNil(t, kl)

	for i := 0; i < 500; i++ {
		key := fmt.Sprintf("k%d", i)
		kl.Lock(key)
		kl.Unlock(key)
		assert.Equal(t, 0, len(kl.locks), "locks must not retain unlocked keys")
		assert.LessOrEqual(t, kl.cache.Len(), capacity)
	}
	assert.Equal(t, 0, len(kl.locks))
	assert.LessOrEqual(t, kl.cache.Len(), capacity)
}

func TestKeyLock_cacheReuse(t *testing.T) {
	kl := NewKeyLock(4)
	require.NotNil(t, kl)

	kl.Lock("a")
	kl.Unlock("a")
	assert.Equal(t, 0, len(kl.locks))
	assert.True(t, kl.cache.Contains("a"))

	kl.Lock("a")
	assert.Equal(t, 1, len(kl.locks))
	assert.False(t, kl.cache.Contains("a"))
	kl.Unlock("a")
	assert.True(t, kl.cache.Contains("a"))
}

func TestKeyLock_capacityEviction(t *testing.T) {
	const capacity = 2
	kl := NewKeyLock(capacity)
	require.NotNil(t, kl)

	for _, key := range []string{"a", "b", "c"} {
		kl.Lock(key)
		kl.Unlock(key)
	}
	assert.Equal(t, 0, len(kl.locks))
	assert.Equal(t, capacity, kl.cache.Len())
	assert.False(t, kl.cache.Contains("a"), "oldest cached key should be evicted")
	assert.True(t, kl.cache.Contains("b"))
	assert.True(t, kl.cache.Contains("c"))

	// 被淘汰的 key 仍可通过新互斥锁正常工作。
	kl.Lock("a")
	kl.Unlock("a")
	assert.Equal(t, 0, len(kl.locks))
}
