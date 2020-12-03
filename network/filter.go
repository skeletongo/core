package network

import (
	"container/list"
)

const (
	Opened uint = iota
	Closed
	Idle
	Received
	Send
	Max
)

type Filter interface {
	Name() string
	InterestOps() uint
	Opened(s *Session) bool                                                     //run in main goroutine
	Closed(s *Session) bool                                                     //run in main goroutine
	Idle(s *Session) bool                                                       //run in main goroutine
	Received(s *Session, packetid int, logicNo uint32, packet interface{}) bool //run in session receive goroutine
	Send(s *Session, packetid int, logicNo uint32, data []byte) bool            //run in session send goroutine
}

type FilterCreator func() Filter

var filterCreators = make(map[string]FilterCreator)

type FilterChain struct {
	filters            *list.List
	filtersInterestOps [Max]*list.List
}

func NewFilterChain() *FilterChain {
	sfc := &FilterChain{
		filters: list.New(),
	}
	for i := uint(0); i < Max; i++ {
		sfc.filtersInterestOps[i] = list.New()
	}
	return sfc
}

func (sfc *FilterChain) AddFirst(sf Filter) {
	sfc.filters.PushFront(sf)
	ops := sf.InterestOps()
	for i := uint(0); i < Max; i++ {
		if ops&(1<<i) != 0 {
			sfc.filtersInterestOps[i].PushFront(sf)
		}
	}
}

func (sfc *FilterChain) AddLast(sf Filter) {
	sfc.filters.PushBack(sf)
	ops := sf.InterestOps()
	for i := uint(0); i < Max; i++ {
		if ops&(1<<i) != 0 {
			sfc.filtersInterestOps[i].PushBack(sf)
		}
	}
}

func (sfc *FilterChain) GetFilter(name string) Filter {
	for e := sfc.filters.Front(); e != nil; e = e.Next() {
		sf := e.Value.(Filter)
		if sf != nil && sf.Name() == name {
			return sf
		}
	}
	return nil
}

func (sfc *FilterChain) OnSessionOpened(s *Session) bool {
	for e := sfc.filtersInterestOps[Opened].Front(); e != nil; e = e.Next() {
		sf := e.Value.(Filter)
		if sf != nil {
			if !sf.Opened(s) {
				return false
			}
		}
	}
	return true
}

func (sfc *FilterChain) OnSessionClosed(s *Session) bool {
	for e := sfc.filtersInterestOps[Closed].Front(); e != nil; e = e.Next() {
		sf := e.Value.(Filter)
		if sf != nil {
			if !sf.Closed(s) {
				return false
			}
		}
	}
	return true
}

func (sfc *FilterChain) OnSessionIdle(s *Session) bool {
	for e := sfc.filtersInterestOps[Idle].Front(); e != nil; e = e.Next() {
		sf := e.Value.(Filter)
		if sf != nil {
			if !sf.Idle(s) {
				return false
			}
		}
	}
	return true
}

func (sfc *FilterChain) OnPacketReceived(s *Session, packetid int, logicNo uint32, packet interface{}) bool {
	for e := sfc.filtersInterestOps[Received].Front(); e != nil; e = e.Next() {
		sf := e.Value.(Filter)
		if sf != nil {
			if !sf.Received(s, packetid, logicNo, packet) {
				return false
			}
		}
	}
	return true
}

func (sfc *FilterChain) OnPacketSent(s *Session, packetid int, logicNo uint32, data []byte) bool {
	for e := sfc.filtersInterestOps[Send].Front(); e != nil; e = e.Next() {
		sf := e.Value.(Filter)
		if sf != nil {
			if !sf.Send(s, packetid, logicNo, data) {
				return false
			}
		}
	}
	return true
}

func RegisterFilterCreator(name string, sfc FilterCreator) {
	if sfc == nil {
		return
	}
	if _, exist := filterCreators[name]; exist {
		panic("repeat register Filter:" + name)
	}
	filterCreators[name] = sfc
}

func GetFilterCreator(name string) FilterCreator {
	if sfc, exist := filterCreators[name]; exist {
		return sfc
	}
	return nil
}
