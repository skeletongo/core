package timer

import (
	"github.com/skeletongo/core/basic"
	"github.com/skeletongo/core/log"
)

// SendTimer 执行延时方法
// o 执行节点
// t 延时方法
func SendTimer(o *basic.Object, t *Timer) {
	if t == nil {
		return
	}
	if o == nil {
		log.Warnf("Timer error: no object")
		return
	}
	o.Send(basic.CommandWrapper(func(o *basic.Object) error {
		t.a.OnTimer(t.h, t.data)
		return nil
	}))
}
