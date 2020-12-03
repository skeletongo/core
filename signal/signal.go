package signal

import (
	"github.com/skeletongo/core/module"
	"os"
	"os/signal"
	"sync"

	"github.com/skeletongo/core/log"
)

var handles = new(sync.Map)

func Register(sig os.Signal, f func()) {
	handles.Store(sig, f)
}

func Run() {
	c := make(chan os.Signal, 10)
	signal.Notify(c)

	for {
		select {
		case sig := <-c:
			log.Infof("Core receive signal: %v", sig)
			if h, ok := handles.Load(sig); ok {
				h.(func())()
			}
		}
	}
}

func init() {
	Register(os.Interrupt, func() {
		module.Stop()
	})
	Register(os.Kill, func() {
		module.Stop()
	})
}