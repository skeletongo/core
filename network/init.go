package network

import (
	"github.com/skeletongo/core/log"
)

var Config = new(Configuration)

type Configuration struct {
	SrvInfo  ServerInfo
	Services []*SessionConfig
}

func (c *Configuration) Name() string {
	return "network"
}

func (c *Configuration) Init() error {
	for _, v := range c.SrvInfo.Banner {
		log.Info(v)
	}
	// 服务初始化配置
	for i := 0; i < len(c.Services); i++ {
		c.Services[i].Init()
	}
	return nil
}

func (c *Configuration) Close() error {
	return nil
}

func init() {
	//pkg.RegisterPackage(Config)
	//module.Register(New(), 0, math.MaxInt32)
}
