package lutil

import (
	"errors"
	"runtime"
	"sync"
	"sync/atomic"
)

var (
	// ErrQueueFull 表示任务队列已满（SubmitErr 的 Abort 语义）。
	ErrQueueFull = errors.New("task queue is full")
	// ErrPoolClosed 表示协程池已关闭。
	ErrPoolClosed = errors.New("pool is closed")
)

// Task 定义任务类型
type Task func()

// RejectPolicy 定义拒绝策略类型（仅 Submit 在队列满时使用；Abort 请用 SubmitErr）
type RejectPolicy func(task Task, pool *Pool)

// Pool 协程池结构体
type Pool struct {
	maxWorkers   int          // 最大工作协程数
	tasks        chan Task    // 任务队列
	rejectPolicy RejectPolicy // 拒绝策略
	wg           sync.WaitGroup
	mu           sync.Mutex
	closed       atomic.Bool
}

// NewPool 创建一个新的协程池
func NewPool(maxWorkers int, queueSize int, rejectPolicy RejectPolicy) *Pool {
	if rejectPolicy == nil {
		rejectPolicy = CallerRunsPolicy
	}
	p := &Pool{
		maxWorkers:   maxWorkers,
		tasks:        make(chan Task, queueSize),
		rejectPolicy: rejectPolicy,
	}
	p.wg.Add(maxWorkers)

	for i := 0; i < p.maxWorkers; i++ {
		go p.worker()
	}
	return p
}

// worker 工作协程
func (p *Pool) worker() {
	defer p.wg.Done()
	for task := range p.tasks {
		p.runTask(task)
	}
}

func (p *Pool) runTask(task Task) {
	defer func() {
		_ = recover()
	}()
	task()
}

// Submit 提交任务；队列满时走拒绝策略。池已关闭则直接返回。
func (p *Pool) Submit(task Task) {
	if !p.trySend(task) {
		if p.closed.Load() {
			return
		}
		p.rejectPolicy(task, p)
	}
}

// SubmitErr 提交任务；等价于 AbortPolicy：队列满则拒绝任务并返回 ErrQueueFull，
// 不执行任务、不走 rejectPolicy。池已关闭则返回 ErrPoolClosed。
func (p *Pool) SubmitErr(task Task) error {
	if p.closed.Load() {
		return ErrPoolClosed
	}
	if p.trySend(task) {
		return nil
	}
	if p.closed.Load() {
		return ErrPoolClosed
	}
	return ErrQueueFull
}

// trySend 在池未关闭时尝试非阻塞入队。成功返回 true。
func (p *Pool) trySend(task Task) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.closed.Load() {
		return false
	}
	select {
	case p.tasks <- task:
		return true
	default:
		return false
	}
}

// Shutdown 关闭协程池；可安全重复调用。
func (p *Pool) Shutdown() {
	p.mu.Lock()
	if p.closed.Load() {
		p.mu.Unlock()
		return
	}
	p.closed.Store(true)
	close(p.tasks)
	p.mu.Unlock()
	p.wg.Wait()
}

// CallerRunsPolicy 由提交任务的 Goroutine 自己执行任务
func CallerRunsPolicy(task Task, pool *Pool) {
	if pool.closed.Load() {
		return
	}
	pool.wg.Add(1)
	defer pool.wg.Done()
	pool.runTask(task)
}

// DiscardPolicy 直接丢弃任务
func DiscardPolicy(task Task, pool *Pool) {
}

// DiscardOldestPolicy 丢弃队列中最老的任务，然后重新提交新任务
func DiscardOldestPolicy(task Task, pool *Pool) {
	for {
		pool.mu.Lock()
		if pool.closed.Load() {
			pool.mu.Unlock()
			return
		}
		select {
		case <-pool.tasks: // 丢弃最老的任务
		default:
		}
		select {
		case pool.tasks <- task:
			pool.mu.Unlock()
			return
		default:
			pool.mu.Unlock()
			runtime.Gosched()
		}
	}
}
