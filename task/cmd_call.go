package task

import (
	"github.com/skeletongo/core/basic"
	"github.com/skeletongo/core/log"
)

// sendCall 给协程节点发送要执行的任务
// o 需要执行 Task.c 的节点
func sendCall(o *basic.Object, t *Task) {
	if t == nil {
		return
	}
	if o == nil {
		log.Errorf("Task [%s] sendCall error: object is nil", t.Name)
		return
	}
	o.Send(basic.CommandWrapper(func(o *basic.Object) error {
		t.run(o)
		return nil
	}))
}
