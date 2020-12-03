package task

import (
	"fmt"

	"github.com/skeletongo/core/basic"

	"github.com/stathat/consistent"
)

var defaultMaster *master

// worker 协程节点
type worker struct {
	*basic.Object
}

type master struct {
	// 预创建协程的序号
	i int
	c *consistent.Consistent
	// 所有协程节点
	workers map[string]*worker
}

func newMaster(n int) *master {
	m := &master{
		c:       consistent.New(),
		workers: make(map[string]*worker),
	}

	for i := 0; i < n; i++ {
		m.addWorker()
	}
	return m
}

func (m *master) addWorkerByName(name string) *worker {
	w := new(worker)
	w.Object = basic.NewObject(m.i, name, Config.Worker.Options, nil)
	w.Object.Run()
	w.Data = w
	Obj.AddChild(w.Object)
	m.workers[w.Name] = w
	m.i++
	return w
}

func (m *master) addWorker() *worker {
	name := fmt.Sprintf("worker_%d", m.i)
	m.c.Add(name)
	return m.addWorkerByName(name)
}

func (m *master) getWorkerByName(name string) *worker {
	if w, ok := m.workers[name]; ok {
		return w
	}
	return nil
}

func (m *master) getWorker(name string) *worker {
	workName, err := m.c.Get(name)
	if err != nil {
		return nil
	}
	return m.getWorkerByName(workName)
}
