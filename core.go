package core

import (
	"github.com/skeletongo/core/basic"
	"github.com/skeletongo/core/log"
	"github.com/skeletongo/core/module"
	"github.com/skeletongo/core/pkg"
	"github.com/skeletongo/core/signal"
	"github.com/skeletongo/core/task"
	"github.com/skeletongo/core/timer"
)

func Run(config string) {
	log.Infof("Core %v starting up", Version)
	defer log.Flush()

	// 初始化
	pkg.Load(config)
	defer pkg.Close()

	// 关联节点
	timer.SetObject(module.Obj)
	task.SetObject(module.Obj)
	basic.Root.AddChild(module.Obj)
	basic.Root.AddChild(task.Obj)

	// 启动 module
	module.Start()

	// 信号监听
	go signal.Run()

	// 等待所有节点关闭
	basic.WG.Wait()
}
