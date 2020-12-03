package task

import (
	"time"

	"github.com/skeletongo/core/basic"
	"github.com/skeletongo/core/pkg"
)

// Obj 协程管理节点，所有新协程的创建都是由这个节点完成
var Obj *basic.Object

// Config 配置
var Config = new(Configuration)

type WorkerConfig struct {
	Options   *basic.Options // 协程节点配置
	WorkerCnt int            // 预创建的协程数量
}

type Configuration struct {
	Options *basic.Options // 协程管理节点配置
	Worker  *WorkerConfig  // 协程节点配置
}

func (c *Configuration) Name() string {
	return "task"
}

func (c *Configuration) Init() error {
	if c.Options.Interval <= 0 {
		c.Options.Interval = time.Millisecond * 10
	} else {
		c.Options.Interval = time.Millisecond * c.Options.Interval
	}
	if c.Worker.WorkerCnt <= 0 {
		c.Worker.WorkerCnt = 4
	}
	Obj = basic.NewObject(basic.TaskID, "task", c.Options, nil)
	Obj.Run()
	// 预创建协程节点，并连接到 Obj 节点，作为子节点
	defaultMaster = newMaster(c.Worker.WorkerCnt)
	return nil
}

func (c *Configuration) Close() error {
	return nil
}

func init() {
	pkg.RegisterPackage(Config)
}
