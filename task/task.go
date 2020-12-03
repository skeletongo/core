// 对多线程的支持
// 主要有以下几个方法
// Start：创建一个协程去执行，执行结束后协程关闭
// StartByExecutor：在预创建的协程节点中执行
// StartByFixExecutor：创建一个协程去执行，协程一旦创建就不会关闭
package task

import "github.com/skeletongo/core/basic"

// defaultObject 回调方法默认执行节点
var defaultObject *basic.Object

// SetObject 设置回调方法默认执行节点
func SetObject(o *basic.Object) {
	defaultObject = o
}

type Callable interface {
	// Call 需要在另外的协程中执行的方法
	// o 协程节点，就是这个方法执行的节点
	// 方法返回值会传递给回调方法
	Call(o *basic.Object) (ret interface{})
}

type CallableWrapper func(o *basic.Object) (ret interface{})

func (cw CallableWrapper) Call(o *basic.Object) (ret interface{}) {
	return cw(o)
}

type CompleteNotify interface {
	// Done 回调方法
	// ret Callable 方法的返回值
	Done(ret interface{}, t *Task)
}

type CompleteNotifyWrapper func(ret interface{}, t *Task)

func (cnw CompleteNotifyWrapper) Done(ret interface{}, t *Task) {
	cnw(ret, t)
}

// Task 任务，需要在另外的协程中处理的方法，通常是一些耗时操作
type Task struct {
	O    *basic.Object  // 回调方法执行节点
	Name string         // 任务名称
	c    Callable       // 需要并发执行的方法
	cb   CompleteNotify // 回调方法
	ret  interface{}    // Callable 方法执行返回值
}

func (t *Task) run(o *basic.Object) {
	if t.c == nil {
		return
	}
	t.ret = t.c.Call(o)
	if t.cb == nil {
		return
	}
	// 在回调方法执行节点执行回调方法
	sendCallback(t.O, t)
}

// New 创建任务
// o 回调方法执行的节点
// c 需要并发执行的方法
// cb 回调方法
// name 任务名称
func New(o *basic.Object, c Callable, cb CompleteNotify, name ...string) *Task {
	ret := &Task{
		O:  o,
		c:  c,
		cb: cb,
	}
	if len(name) > 0 {
		ret.Name = name[0]
	}
	if o == nil {
		ret.O = defaultObject
	}
	return ret
}

// Start 创建一个协程去执行，执行结束后协程关闭
func (t *Task) Start() {
	go t.run(nil)
}

// StartByExecutor 在预创建的协程节点中执行
// name 任务名称，名称相同的任务会在同一个协程中串行执行
func (t *Task) StartByExecutor(name string) {
	sendToExecutor(t, name)
}

// StartByFixExecutor 创建一个协程去执行，协程一旦创建就不会关闭
// 如果已经有任务名称相同的协程了(已经使用相同的name调用过此方法)，就不会再创建新协程
// name 任务名称，名称相同的任务会在同一个协程中串行执行
func (t *Task) StartByFixExecutor(name string) {
	sendToFixExecutor(t, name)
}
