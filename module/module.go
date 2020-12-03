package module

import (
	"container/list"
	"strings"
	"time"

	"github.com/skeletongo/core/basic"
	"github.com/skeletongo/core/log"
	"github.com/skeletongo/core/timer"
	"github.com/skeletongo/core/utils"
)

// Module 扩展模块，实现自己的功能
type Module interface {
	// Name 模块名称
	Name() string
	// Init 模块初始化方法
	Init()
	// Update 模块更新方法
	Update()
	// Close 模块关闭方法
	Close()
}

type module struct {
	// 最后一次更新时间
	lastTime time.Time
	// 更新时间间隔
	interval time.Duration
	// 优先级，越小优先级越高
	priority int
	// implement Module
	mi Module
}

func (m *module) safeInit() {
	defer utils.DumpStackIfPanic("Module.safeInit")
	m.mi.Init()
}

func (m *module) safeUpdate(t time.Time) {
	defer utils.DumpStackIfPanic("Module.safeUpdate")
	if m.interval == 0 || t.Sub(m.lastTime) >= m.interval {
		m.lastTime = t
		m.mi.Update()
	}
}

func (m *module) safeClose() {
	defer utils.DumpStackIfPanic("Module.safeClose")
	m.mi.Close()
}

// 模块管理器状态
const (
	StateInvalid = iota // 停止
	StateInit           // 初始化
	StateUpdate         // 运行中
	StateClose          // 开始关闭
	StateClosing        // 关闭中
	StateClosed         // 已关闭
)

// moduleMgr 模块管理器
type moduleMgr struct {
	// state 模块管理器状态
	state int
	// mods 所有模块
	mods *list.List
	// modSign 接收模块关闭信号
	modSign chan string
	// t 定时输出还有哪些模块没有关闭
	t <-chan time.Time
}

func (m *moduleMgr) onTick() {
	switch m.state {
	case StateInit:
		m.init()
	case StateUpdate:
		m.update()
	case StateClose:
		m.close()
	case StateClosing:
		m.closing()
	case StateClosed:
		m.closed()
	}
}

func (m *moduleMgr) init() {
	log.Infof("module init...")
	for e := m.mods.Front(); e != nil; e = e.Next() {
		mod := e.Value.(*module)
		log.Infof("module [%16s] init...", mod.mi.Name())
		mod.safeInit()
		log.Infof("module [%16s] init[ok]", mod.mi.Name())
	}
	log.Infof("module init[ok]")

	m.state = StateUpdate
}

func (m *moduleMgr) update() {
	nowTime := time.Now()
	for e := m.mods.Front(); e != nil; e = e.Next() {
		e.Value.(*module).safeUpdate(nowTime)
	}
}

func (m *moduleMgr) close() {
	m.modSign = make(chan string, m.mods.Len())

	// 停止所有定时任务
	timer.StopAll()
	log.Infof("timer close")

	log.Infof("module close...")
	for e := m.mods.Back(); e != nil; e = e.Prev() {
		mod := e.Value.(*module)
		log.Infof("module [%16s] close...", mod.mi.Name())
		mod.safeClose()
		log.Infof("module [%16s] close[ok]", mod.mi.Name())
	}
	log.Infof("module close[ok]")

	m.state = StateClosing

	m.t = time.Tick(time.Second)
}

func (m *moduleMgr) closing() {
	for {
		select {
		case name := <-m.modSign:
			for e := m.mods.Front(); e != nil; e = e.Next() {
				if e.Value.(*module).mi.Name() == name {
					m.mods.Remove(e)
					break
				}
			}
		case <-m.t:
			if m.mods.Len() > 0 {
				var names []string
				for e := m.mods.Front(); e != nil; e = e.Next() {
					names = append(names, e.Value.(*module).mi.Name())
				}
				log.Info("module closing ", strings.Join(names, "|"))
			}
		default:
			if m.mods.Len() == 0 {
				m.state = StateClosed
			} else {
				m.update()
			}
			return
		}
	}
}

func (m *moduleMgr) closed() {
	// 关闭根节点
	basic.Root.Close()

	m.state = StateInvalid
}

func newModuleMgr() *moduleMgr {
	ret := &moduleMgr{
		state: StateInvalid,
		mods:  list.New(),
	}
	return ret
}

var defaultModuleMgr = newModuleMgr()

func Closed(m Module) {
	defaultModuleMgr.modSign <- m.Name()
}

// Register 模块注册
// interval 间隔时长；如果值为0表示以最短间隔时间执行update,取值范围大于等于0
// priority 优先级；值越小越优先处理
func Register(m Module, interval time.Duration, priority int) {
	mod := &module{
		lastTime: time.Now(),
		interval: interval,
		priority: priority,
		mi:       m,
	}
	for e := defaultModuleMgr.mods.Front(); e != nil; e = e.Next() {
		if me, ok := e.Value.(*module); ok {
			if priority < me.priority {
				defaultModuleMgr.mods.InsertBefore(mod, e)
				return
			}
		}
	}
	defaultModuleMgr.mods.PushBack(mod)
}

// Start 启动模块
func Start() {
	Obj.Send(basic.CommandWrapper(func(o *basic.Object) error {
		defaultModuleMgr.state = StateInit
		return nil
	}))
}

// Stop 停止所有模块
func Stop() {
	Obj.Send(basic.CommandWrapper(func(o *basic.Object) error {
		defaultModuleMgr.state = StateClose
		return nil
	}))
}
