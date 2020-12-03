package basic

import (
	"time"
)

// Options 节点配置
type Options struct {
	// Interval 定时任务的执行时间间隔
	Interval time.Duration
}

// State 节点状态
type State struct {
	QueueLen   uint64 // 待处理消息数量
	EnqueueNum uint64 // 收到的消息总数
	DoneNum    uint64 // 已处理的消息数
}
