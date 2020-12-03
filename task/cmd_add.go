package task

import (
	"errors"

	"github.com/skeletongo/core/basic"
)

var ErrCannotFindWorker = errors.New("Cannot find worker ")

// sendToExecutor 给预创建的协程节点发送待执行的任务
func sendToExecutor(t *Task, name string) {
	if t == nil {
		return
	}
	Obj.Send(basic.CommandWrapper(func(o *basic.Object) error {
		w := defaultMaster.getWorker(name)
		if w == nil {
			return ErrCannotFindWorker
		}

		sendCall(w.Object, t)
		return nil
	}))
}

// sendToFixExecutor 给指定的一个协程节点发送待执行的任务
func sendToFixExecutor(t *Task, name string) {
	if t == nil {
		return
	}
	Obj.Send(basic.CommandWrapper(func(o *basic.Object) error {
		w := defaultMaster.getWorkerByName(name)
		if w == nil {
			// 创建新的协程节点
			w = defaultMaster.addWorkerByName(name)
		}

		sendCall(w.Object, t)
		return nil
	}))
}
