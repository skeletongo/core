package basic

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/skeletongo/core/container/queue"
	"github.com/skeletongo/core/utils"
)

// WG 用来等待所有节点都关闭
var WG = &sync.WaitGroup{}

// Object 节点
type Object struct {
	sync.Mutex
	// ID 节点ID
	ID int
	// Name 节点名称
	Name string
	// Data 节点保存的数据
	Data interface{}
	// Opt 节点配置
	Opt *Options
	// Closed 节点是否已经关闭
	Closed bool
	// CloseSign 是否已经给当前节点发送过关闭消息
	CloseSign bool
	// closing 节点正在关闭
	closing bool
	// doneNum 已处理消息数量
	doneNum uint64
	// sendNum 收到的消息总数
	sendNum uint64
	// ack 当节点正在关闭时， ack 的值等于当前节点直接子节点的数量，当收到其中一个子节点已经关闭的消息后 ack 减一
	// 当收到所有直接子节点的已关闭消息后 ack 为零；当 ack 为零时当前节点才能关闭；
	// 另外判定当前节点是否已经关闭还有一些其它条件，见 checkAck 方法
	ack uint64
	// child 记录当前节点的直接子节点; key:子节点ID,value:子节点
	child sync.Map
	// owner 父节点
	owner *Object
	// q 消息队列
	q queue.Queue
	// sign 收到新消息的信号
	// 作用：当消息队列为空时，阻塞当前节点所在的协程，当收到新消息后不再阻塞
	sign chan struct{}
	// ticker 定时器，用来定时处理定时任务
	ticker *time.Ticker
	// sinker .
	sinker Sinker
}

// NewObject 创建节点
// id 节点ID
// name 节点名称
// opt 节点配置
// sinker
func NewObject(id int, name string, opt *Options, sinker Sinker) *Object {
	if opt == nil {
		panic("NewObject error: required Options")
	}
	o := &Object{
		ID:     id,
		Name:   name,
		Opt:    opt,
		sinker: sinker,
		sign:   make(chan struct{}, 1),
		q:      queue.NewSyncQueue(),
	}
	return o
}

// FullName 完整名称
func (o *Object) FullName() string {
	name := o.Name
	parent := o.owner
	for parent != nil {
		name = parent.Name + "/" + name
		parent = parent.owner
	}
	return "/" + name
}

func (o *Object) safeDone(cmd Command) {
	defer utils.DumpStackIfPanic("Object::Command::Done")

	defer atomic.AddUint64(&o.doneNum, 1)
	err := cmd.Done(o)
	if err != nil {
		panic(err)
	}
}

func (o *Object) safeStart() {
	defer utils.DumpStackIfPanic("Object::OnStart")

	if o.sinker != nil {
		o.sinker.OnStart()
	}
}

func (o *Object) safeTick() {
	defer utils.DumpStackIfPanic("Object::OnTick")

	if o.sinker != nil {
		o.sinker.OnTick()
	}
}

func (o *Object) safeStop() {
	defer utils.DumpStackIfPanic("Object::OnStop")

	if o.sinker != nil {
		o.sinker.OnStop()
	}
}

// State 获取节点状态
func (o *Object) State() *State {
	return &State{
		QueueLen:   uint64(o.q.Len()),
		EnqueueNum: atomic.LoadUint64(&o.sendNum),
		DoneNum:    atomic.LoadUint64(&o.doneNum),
	}
}

// GetStates 获取节点状态包含所有子节点的状态
func (o *Object) GetStates() map[string]*State {
	stats := make(map[string]*State)
	stats[o.FullName()] = o.State()
	o.child.Range(func(key, value interface{}) bool {
		if c, ok := value.(*Object); ok && c != nil {
			stats[c.FullName()] = c.State()
			for k, v := range c.GetStates() {
				stats[k] = v
			}
		}
		return true
	})
	return stats
}

// checkAck 判定节点是否可以关闭
// 关闭条件：所有子节点已经关闭，所有收到的消息已经处理
func (o *Object) checkAck() bool {
	if !o.closing || o.ack > 0 || o.sendNum > o.doneNum {
		return false
	}
	sendAck(o.owner)
	o.Lock()
	o.Closed = true
	o.Unlock()
	return true
}

func (o *Object) run() {
	defer WG.Done()
	// 定时器
	if o.Opt.Interval > 0 && o.sinker != nil {
		o.ticker = time.NewTicker(o.Opt.Interval)
	}
	// 队列，定时任务
	for !o.checkAck() {
		if o.q.Len() <= 0 {
			if o.ticker == nil {
				<-o.sign
				continue
			}
			select {
			case <-o.sign:
			case <-o.ticker.C:
				o.safeTick()
			}
		} else {
			cmd, ok := o.q.Dequeue().(Command)
			if !ok {
				continue
			}
			o.safeDone(cmd)
			if o.ticker != nil {
				select {
				case <-o.ticker.C:
					o.safeTick()
				default:
				}
			}
		}
	}
}

// Run 启动节点
// 创建一个协程来处理消息队列中的消息和定时任务
func (o *Object) Run() {
	WG.Add(1)
	o.safeStart()
	go o.run()
}

// Send 给当前节点发送消息
// 此方法为非阻塞方法，消息为异步处理，消息先进入消息队列等待处理
func (o *Object) Send(c Command) {
	atomic.AddUint64(&o.sendNum, 1)
	o.q.Enqueue(c)
	select {
	case o.sign <- struct{}{}:
	default:
	}
}

// AddChild 添加一个子节点
func (o *Object) AddChild(c *Object) {
	if c == nil {
		return
	}

	o.Lock()
	if o.Closed || o.CloseSign {
		o.Unlock()
		c.Close()
		return
	}
	o.Unlock()

	if c.owner != nil {
		panic("AddChild error: An object can have only one parent node")
	}
	c.owner = o

	// 通知父节点子节点已添加
	sendAddChild(o, c)
}

// IsClosed 是否已经关闭
func (o *Object) IsClosed() bool {
	o.Lock()
	if o.Closed {
		o.Unlock()
		return true
	}
	o.Unlock()
	return false
}

// Close 关闭节点
// 关闭一个节点需要它的所有子节点都关闭
func (o *Object) Close() {
	// 判定节点是否已经关闭
	o.Lock()
	if o.Closed {
		o.Unlock()
		return
	}
	// 节点只接收一次关闭消息
	if o.CloseSign {
		o.Unlock()
		return
	}
	o.CloseSign = true
	o.Unlock()

	if o.owner != nil {
		sendReqClose(o.owner, o)
	} else {
		sendClose(o)
	}
}
