// Package lutil 提供协程池与按键互斥锁等基础工具。
package lutil

import (
	"sync"

	"github.com/hashicorp/golang-lru/v2"
)

// keyMutex 带引用计数的互斥锁；refs 统计已进入 Lock、尚未完成 Unlock 的调用数。
type keyMutex struct {
	mu   sync.Mutex
	refs int
}

// KeyLock 提供基于键的互斥锁功能
type KeyLock struct {
	mu    sync.Mutex
	locks map[string]*keyMutex
	cache *lru.Cache[string, *keyMutex]
}

// NewKeyLock 创建带 LRU 缓存的按键互斥锁，size 为缓存容量。
func NewKeyLock(size int) *KeyLock {
	l, _ := lru.New[string, *keyMutex](size)
	if l == nil {
		return nil
	}
	return &KeyLock{
		locks: make(map[string]*keyMutex),
		cache: l,
	}
}

// Lock 获取指定键的锁
func (kl *KeyLock) Lock(key string) {
	kl.mu.Lock()
	km, ok := kl.locks[key]
	if !ok {
		if cached, hit := kl.cache.Get(key); hit {
			km = cached
			kl.cache.Remove(key)
		} else {
			km = &keyMutex{}
		}
		kl.locks[key] = km
	}
	km.refs++
	kl.mu.Unlock()
	km.mu.Lock()
}

// Unlock 释放指定键的锁。无等待者时从 locks 移除并放入 LRU，避免 locks 随键无限增长。
func (kl *KeyLock) Unlock(key string) {
	kl.mu.Lock()
	defer kl.mu.Unlock()

	km, ok := kl.locks[key]
	if !ok {
		panic("unlock of unlocked mutex")
	}
	km.mu.Unlock()
	km.refs--
	if km.refs == 0 {
		delete(kl.locks, key)
		kl.cache.Add(key, km)
	}
}
