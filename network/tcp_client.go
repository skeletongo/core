package network

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/skeletongo/core/log"
)

type TCPClient struct {
	network   *Network
	SC        *SessionConfig
	conns     map[net.Conn]struct{}
	actions   chan *action // 消息队列，所有链接接收到的消息都进入这个队列
	closeSign chan struct{}
	closing   bool
	m         sync.Mutex
	wg        sync.WaitGroup
}

func NewTCPClient(n *Network, sc *SessionConfig) *TCPClient {
	return &TCPClient{
		network:   n,
		SC:        sc,
		conns:     make(map[net.Conn]struct{}),
		actions:   make(chan *action, sc.MaxDone),
		closeSign: make(chan struct{}),
	}
}

func (t *TCPClient) Start() error {
	for i := 0; i < t.SC.ConnNum; i++ {
		t.wg.Add(1)
		go t.connect(t.SC)
	}
	return nil
}

func (t *TCPClient) dial(addr string) net.Conn {
	for {
		conn, err := net.Dial("tcp", addr)
		if err == nil || t.closing {
			return conn
		}

		log.Infof("connect to %v error: %v", addr, err)
		time.Sleep(t.SC.ReconnectInterval)
	}
}

func (t *TCPClient) connect(sc *SessionConfig) {
	defer t.wg.Done()

reconnect:
	addr := fmt.Sprintf("%v:%v", sc.Ip, sc.Port)
	conn := t.dial(addr)
	if conn == nil {
		return
	}

	t.m.Lock()
	if t.closing {
		t.m.Unlock()
		conn.Close()
		return
	}
	t.conns[conn] = struct{}{}
	t.m.Unlock()

	state := NewSessionState(sc)
	tcpConn := NewTCPConn(conn, state)

	s := &Session{
		SessionState: state,
		conn:         tcpConn,
		actions:      t.actions,
	}
	go s.WriteMsg()
	s.ReadMsg()

	t.m.Lock()
	delete(t.conns, conn)
	t.m.Unlock()
	s.Close()

	if t.SC.AutoReconnect {
		time.Sleep(t.SC.ReconnectInterval)
		goto reconnect
	}
}

func (t *TCPClient) Update() {
	for {
		select {
		case <-t.closeSign:
			t.network.ServiceClosed(t.SC)
		case v := <-t.actions:
			v.do()
			t.SC.actionPool.Put(v)
		default:
			return
		}
	}
}

func (t *TCPClient) Shutdown() {
	t.m.Lock()
	t.closing = true
	for conn := range t.conns {
		conn.Close()
	}
	t.conns = nil
	t.m.Unlock()
	t.wg.Wait()
	close(t.closeSign)
}
