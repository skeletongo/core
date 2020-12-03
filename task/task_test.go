package task_test

import (
	"fmt"
	"github.com/skeletongo/core/basic"
	"github.com/skeletongo/core/task"
	"testing"
	"time"
)

// task 初始化
func init() {
	task.Config.Options = new(basic.Options)
	task.Config.Worker = &task.WorkerConfig{
		Options:   new(basic.Options),
		WorkerCnt: 5,
	}
	task.Config.Init()
}

func ExampleTask_Start() {
	ch := make(chan struct{})
	task.New(basic.Root, task.CallableWrapper(func(o *basic.Object) interface{} {
		fmt.Println("1")
		return "2"
	}), task.CompleteNotifyWrapper(func(ret interface{}, t *task.Task) {
		fmt.Println(ret)
		fmt.Println(3)
		ch <- struct{}{}
	})).Start()
	<-ch
	// output:
	// 1
	// 2
	// 3
}

func TestTask_StartByExecutor(t *testing.T) {
	ch := make(chan string, 3)
	A, B := "a", "b"

	task.New(task.Obj, task.CallableWrapper(func(o *basic.Object) interface{} {
		t.Logf("do task name: %s, object name: %s\n", A, o.Name)
		return o.Name
	}), task.CompleteNotifyWrapper(func(ret interface{}, t *task.Task) {
		ch <- fmt.Sprint(ret, A)
	})).StartByExecutor(A)

	task.New(task.Obj, task.CallableWrapper(func(o *basic.Object) interface{} {
		t.Logf("do task name: %s, object name: %s\n", A, o.Name)
		return o.Name
	}), task.CompleteNotifyWrapper(func(ret interface{}, t *task.Task) {
		ch <- fmt.Sprint(ret, A)
	})).StartByExecutor(A)

	task.New(task.Obj, task.CallableWrapper(func(o *basic.Object) interface{} {
		t.Logf("do task name: %s, object name: %s\n", B, o.Name)
		return o.Name
	}), task.CompleteNotifyWrapper(func(ret interface{}, t *task.Task) {
		ch <- fmt.Sprint(ret, B)
	})).StartByExecutor(B)

	var names []string
	for i := 0; i < 3; i++ {
		select {
		case v :=<- ch:
			names = append(names, v)
		case <- time.Tick(time.Second):
			t.Error("1")
			return
		}
	}
	t.Log(names)

	if len(names) != 3 {
		t.Error("2")
		return
	}

	if names[0] == names[1] && names[0] != names[2] {
		return
	}
	if names[1] == names[2] && names[0] != names[1] {
		return
	}
	if names[0] == names[2] && names[0] != names[1] {
		return
	}
	t.Error("3")
}

func ExampleTask_StartByFixExecutor() {

}
