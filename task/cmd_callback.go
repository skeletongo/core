package task

import (
	"github.com/skeletongo/core/basic"
	"github.com/skeletongo/core/log"
)

// sendCallback 执行回调函数
func sendCallback(o *basic.Object, t *Task) {
	if t == nil {
		return
	}
	if o == nil {
		log.Errorf("Task [%s] sendCallback error: object is nil", t.Name)
	}
	t.O.Send(basic.CommandWrapper(func(o *basic.Object) error {
		t.cb.Done(t.ret, t)
		return nil
	}))
}
