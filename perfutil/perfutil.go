package perfutil

import (
	"fmt"
	"strings"
	"time"
)

// PerfTime 独立计时器实例（支持多链路、并发安全）
type PerfTime struct {
	name    string
	current time.Time
}

// NewPerfTime 创建计时器实例（初始化时自动重置时间）
func NewPerfTime(name ...string) *PerfTime {
	return &PerfTime{
		name:    strings.Join(name, ":"),
		current: time.Now(),
	}
}

// Reset 重置计时器起点
func (p *PerfTime) Reset() {
	p.current = time.Now()
}

// Mark 记录当前步骤耗时，并更新计时器起点（指针接收者确保状态更新）
func (p *PerfTime) Mark(msg string) {
	now := time.Now()
	duration := now.Sub(p.current)
	fmt.Printf("[PERF] %s %s: %v\n", p.name, msg, duration)
	p.current = now
}
