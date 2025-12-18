package lutil

import (
	"github.com/hashicorp/golang-lru/v2"
	"sync"
)

// KeyLock 提供基于键的互斥锁功能
type KeyLock struct {
	mu    sync.RWMutex
	locks map[string]*sync.Mutex
	cache *lru.Cache[string, *sync.Mutex]
}

func NewKeyLock(size int) *KeyLock {
	l, _ := lru.New[string, *sync.Mutex](size)
	if l == nil {
		return nil
	}
	return &KeyLock{
		mu:    sync.RWMutex{},
		locks: make(map[string]*sync.Mutex),
		cache: l,
	}
}

// Lock 获取指定键的锁
func (kl *KeyLock) Lock(key string) {
	kl.mu.RLock()
	mtx, ok := kl.locks[key]
	if ok {
		kl.mu.RUnlock()
		mtx.Lock()
		return
	}
	kl.mu.RUnlock()
	kl.mu.Lock()

	mtx, ok = kl.cache.Get(key)
	if ok {
		mtx.Lock()
		kl.cache.Remove(key)
		kl.locks[key] = mtx
		kl.mu.Unlock()
		return
	}

	mtx = &sync.Mutex{}
	mtx.Lock()
	kl.locks[key] = mtx
	kl.mu.Unlock()
}

// Unlock 释放指定键的锁
func (kl *KeyLock) Unlock(key string) {
	kl.mu.Lock()
	defer kl.mu.Unlock()

	mtx, ok := kl.locks[key]
	if !ok {
		panic("unlock of unlocked mutex")
	}
	mtx.Unlock()
	kl.cache.Add(key, mtx)
}
