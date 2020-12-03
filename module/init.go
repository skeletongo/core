package module

import (
	"time"

	"github.com/skeletongo/core/basic"
	"github.com/skeletongo/core/pkg"
)

// Obj 模块功能节点
var Obj *basic.Object

// Config 节点配置
var Config = new(Configuration)

type Configuration struct {
	Options *basic.Options
}

func (c *Configuration) Name() string {
	return "module"
}

func (c *Configuration) Init() error {
	if c.Options.Interval <= 0 {
		c.Options.Interval = time.Millisecond * 10
	} else {
		c.Options.Interval = time.Millisecond * c.Options.Interval
	}
	Obj = basic.NewObject(basic.ModuleID, "module", c.Options, new(sink))
	Obj.Run()
	return nil
}

func (c *Configuration) Close() error {
	return nil
}

func init() {
	pkg.RegisterPackage(Config)
}
