package network

import (
	"github.com/skeletongo/core/log"
	"net"
	"sync"
)

type ISession interface {
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	GetSessionConfig() *SessionConfig
	SetAttribute(key, value interface{})
	RemoveAttribute(key interface{})
	GetAttribute(key interface{}) interface{}
	Send(msgID int, msg interface{}) error
	SendEx(msgID int, logicNo uint32, msg interface{}) error
	Close()
}

type pack struct {
	logicNo uint32
	b       []byte
}

type SessionState struct {
	ID       int
	SID      int64
	Auth     bool
	SC       *SessionConfig
	userData map[interface{}]interface{}
	packPool sync.Pool
	send     chan *pack
	writeBuf *PkgData
	readBuf  *PkgData
}

func NewSessionState(sc *SessionConfig) *SessionState {
	ret := &SessionState{
		ID:       sc.GetSeq(),
		SC:       sc,
		userData: make(map[interface{}]interface{}),
		send:     make(chan *pack, sc.MaxSend),
		writeBuf: sc.pkgDataPool.Get().(*PkgData),
		readBuf:  sc.pkgDataPool.Get().(*PkgData),
	}
	ret.writeBuf.Seq = 0
	ret.readBuf.Seq = 0
	ret.packPool.New = func() interface{} {
		return new(pack)
	}
	return ret
}

func (s *SessionState) GetSessionConfig() *SessionConfig {
	return s.SC
}

func (s *SessionState) SetAttribute(key, value interface{}) {
	s.userData[key] = value
}

func (s *SessionState) RemoveAttribute(key interface{}) {
	delete(s.userData, key)
}

func (s *SessionState) GetAttribute(key interface{}) interface{} {
	return s.userData[key]
}

func (s *SessionState) Send(msgID int, msg interface{}) error {
	return s.SendEx(msgID, 0, msg)
}

func (s *SessionState) SendEx(msgID int, logicNo uint32, msg interface{}) error {
	b, err := Marshal(msgID, msg)
	if err != nil {
		return err
	}

	p := s.packPool.Get().(*pack)
	p.logicNo = logicNo
	p.b = b

	select {
	case s.send <- p:
	default:
	}
	return nil
}

// Session implement ISession
type Session struct {
	*SessionState
	conn    Conn
	actions chan *action
}

func (a *Session) LocalAddr() net.Addr {
	return a.conn.LocalAddr()
}

func (a *Session) RemoteAddr() net.Addr {
	return a.conn.RemoteAddr()
}

// read goroutine
func (a *Session) ReadMsg() {
	for {
		data, err := a.conn.ReadMsg()
		if err != nil {
			log.Errorf("read message error: %v", err)
			break
		}

		msgID, msg, err := Unmarshal(data)
		if err != nil {
			log.Errorf("encoding.Unmarshal error: %v", err)
			break
		}

		ac := a.SC.actionPool.Get().(*action)
		ac.s = a
		ac.msgID = msgID
		ac.msg = msg
		a.actions <- ac
	}

	a.SC.pkgDataPool.Put(a.readBuf)
	a.conn.Close()
}

// write goroutine
func (a *Session) WriteMsg() {
	for v := range a.send {
		if v == nil {
			break
		}
		a.writeBuf.Head.LogicNo = v.logicNo
		if err := a.conn.WriteMsg(v.b); err != nil {
			log.Errorf("send message error: %v", err)
			break
		}

		a.packPool.Put(v)
	}

	a.SC.pkgDataPool.Put(a.writeBuf)
	a.conn.Close()
}

// goroutine safe
func (a *Session) Close() {
	a.conn.Close()
	select {
	case a.send <- nil:
	default:
	}
}
