package basic

import (
	"fmt"
	"testing"
	"time"
)

// 测试时先把根节点关掉，因为 WG 是公用的，关闭根节点重置 WG
func init() {
	Root.Close()
	WG.Wait()
}

func ExampleObject_Send() {
	var n []int
	obj := NewObject(0, "test", &Options{Interval: 0}, nil)
	obj.Run()
	obj.Send(CommandWrapper(func(o *Object) error {
		n = append(n, 1)
		return nil
	}))
	obj.Send(CommandWrapper(func(o *Object) error {
		fmt.Println(n)
		return nil
	}))
	obj.Send(CommandWrapper(func(o *Object) error {
		n = append(n, 2)
		return nil
	}))
	obj.Send(CommandWrapper(func(o *Object) error {
		n = append(n, 3)
		return nil
	}))
	obj.Send(CommandWrapper(func(o *Object) error {
		fmt.Println(n)
		return nil
	}))
	obj.Close()
	WG.Wait()
	// Output:
	// [1]
	// [1 2 3]
}

type testSinker struct {
	name string
	n    int
}

func (t *testSinker) OnStart() {
	fmt.Println("OnStart", t.name)
}

func (t *testSinker) OnTick() {
	if t.n < 1 {
		fmt.Println(t.n, t.name)
		t.n++
	}
}

func (t *testSinker) OnStop() {
	fmt.Println("OnStop", t.name)
}

func ExampleObject_Close() {
	a := NewObject(0, "a", &Options{Interval: time.Millisecond}, &testSinker{name: "a"})
	a.Run()
	b := NewObject(1, "b", &Options{Interval: time.Millisecond * 4}, &testSinker{name: "b"})
	b.Run()
	c := NewObject(2, "c", &Options{Interval: time.Millisecond * 8}, &testSinker{name: "c"})
	c.Run()

	time.Sleep(time.Millisecond * 10)

	a.AddChild(b)
	a.AddChild(c)
	c.Close()
	time.Sleep(time.Millisecond * 2)
	a.Close()
	time.Sleep(time.Millisecond * 2)
	// 重复关闭
	c.Close()
	a.Close()
	time.Sleep(time.Millisecond * 5)
	WG.Wait()
	// Output:
	// OnStart a
	// OnStart b
	// OnStart c
	// 0 a
	// 0 b
	// 0 c
	// OnStop c
	// OnStop a
	// OnStop b
}

type s1 struct {
}

func (s *s1) OnStart() {
}

func (s *s1) OnTick() {
}

func (s *s1) OnStop() {
}

type s2 struct {
}

func (s *s2) OnStart() {
}

func (s *s2) OnTick() {
}

func (s *s2) OnStop() {
	time.Sleep(time.Second * 2)
}

func TestObject_Close(t *testing.T) {
	o1 := NewObject(1, "a", &Options{Interval: time.Second}, &s1{})
	o2 := NewObject(2, "b", &Options{Interval: time.Second}, &s1{})
	o3 := NewObject(3, "c", &Options{Interval: time.Second}, &s2{})

	o1.Run()
	o2.Run()
	o3.Run()

	o1.AddChild(o2)
	o1.AddChild(o3)

	o1.Close()
	time.Sleep(time.Second)
	if o1.IsClosed() || !o2.IsClosed() || o3.IsClosed() {
		t.Error()
	}
	WG.Wait()
	if !o1.IsClosed() || !o2.IsClosed() || !o3.IsClosed() {
		t.Error()
	}
}

var RunCh = make(chan int, 11)

type runSinker struct {
}

func (r *runSinker) OnStart() {
	RunCh <- 1
}

func (r *runSinker) OnTick() {
	RunCh <- 6
}

func (r *runSinker) OnStop() {
	RunCh <- -1
}

func TestObject_Run(t *testing.T) {
	obj := NewObject(0, "test", &Options{Interval: time.Millisecond * 50}, new(runSinker))
	obj.Run()

	for i := 2; i < 6; i++ {
		go func(n int) {
			obj.Send(CommandWrapper(func(o *Object) error {
				RunCh <- n
				return nil
			}))
		}(i)
	}

	time.Sleep(time.Millisecond * 10)

	obj.Send(CommandWrapper(func(o *Object) error {
		time.Sleep(time.Millisecond * 40)
		return nil
	}))

	for i := 7; i < 10; i++ {
		go func(n int) {
			obj.Send(CommandWrapper(func(o *Object) error {
				RunCh <- n
				time.Sleep(time.Millisecond * 20)
				return nil
			}))
		}(i)
	}

	time.Sleep(time.Millisecond * 100)
	obj.Close()
	WG.Wait()

	res := make([]int, 11)
	for i := 0; i < 11; i++ {
		v := <-RunCh
		fmt.Printf("%d ", v)
		res[i] = v
	}

	switch {
	case res[0] != 1 || res[10] != -1 || res[5] != 6 || res[9] != 6:
		t.Error()
	case res[1]+res[2]+res[3]+res[4] != 2+3+4+5:
		t.Error()
	case res[6]+res[7]+res[8] != 7+8+9:
		t.Error()
	}
}
