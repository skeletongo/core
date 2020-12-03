package network

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/skeletongo/core/log"
)

type TCPServer struct {
	network   *Network
	SC        *SessionConfig
	conns     map[net.Conn]struct{}
	actions   chan *action // 消息队列，所有链接接收到的消息都进入这个队列
	closeSign chan struct{}
	m         sync.Mutex
	ln        net.Listener
	wgLn      sync.WaitGroup
	wgConns   sync.WaitGroup
}

func NewTCPServer(n *Network, cfg *SessionConfig) *TCPServer {
	return &TCPServer{
		network:   n,
		SC:        cfg,
		conns:     make(map[net.Conn]struct{}),
		actions:   make(chan *action, cfg.MaxDone),
		closeSign: make(chan struct{}),
	}
}

func (t *TCPServer) Start() error {
	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", t.SC.Ip, t.SC.Port))
	if err != nil {
		log.Error("TCPServer Listen error:", err)
		return err
	}

	t.ln = ln

	go func() {
		// 监听端口
		t.wgLn.Add(1)
		defer t.wgLn.Done()

		var tempDelay time.Duration
		for {
			conn, err := t.ln.Accept()
			if err != nil {
				if ne, ok := err.(net.Error); ok && ne.Temporary() {
					if tempDelay == 0 {
						tempDelay = 5 * time.Millisecond
					} else {
						tempDelay *= 2
					}
					if max := 1 * time.Second; tempDelay > max {
						tempDelay = max
					}
					log.Infof("accept error: %v; retrying in %v", err, tempDelay)
					time.Sleep(tempDelay)
					continue
				}
				return
			}
			tempDelay = 0

			t.m.Lock()
			if len(t.conns) > t.SC.MaxConn {
				t.m.Unlock()
				conn.Close()
				log.Warn("too many connections")
				continue
			}
			t.conns[conn] = struct{}{}
			t.m.Unlock()

			c := conn.(*net.TCPConn)
			if t.SC.IsInnerLink {
				var timeZero time.Time
				c.SetReadDeadline(timeZero)
				c.SetWriteDeadline(timeZero)
			} else {
				now := time.Now()
				if t.SC.ReadTimeout != 0 {
					c.SetReadDeadline(now.Add(t.SC.ReadTimeout))
				}
				if t.SC.WriteTimeout != 0 {
					c.SetWriteDeadline(now.Add(t.SC.WriteTimeout))
				}
			}
			c.SetLinger(t.SC.Linger)
			c.SetNoDelay(t.SC.NoDelay)
			c.SetReadBuffer(t.SC.ReadBuffer)
			c.SetWriteBuffer(t.SC.WriteBuffer)
			c.SetKeepAlive(t.SC.KeepAlive)
			if t.SC.KeepAlive {
				c.SetKeepAlivePeriod(t.SC.KeepAlivePeriod)
			}

			t.wgConns.Add(1)

			state := NewSessionState(t.SC)
			tcpConn := NewTCPConn(conn, state)
			s := &Session{
				SessionState: state,
				conn:         tcpConn,
				actions:      t.actions,
			}
			go s.WriteMsg()
			go func() {
				s.ReadMsg()

				t.m.Lock()
				delete(t.conns, conn)
				t.m.Unlock()
				s.Close()

				t.wgConns.Done()
			}()
		}
	}()
	return nil
}

func (t *TCPServer) Update() {
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

func (t *TCPServer) Shutdown() {
	t.ln.Close()
	t.wgLn.Wait()

	t.m.Lock()
	for conn := range t.conns {
		conn.Close()
	}
	t.conns = nil
	t.m.Unlock()
	t.wgConns.Wait()
	close(t.closeSign)
}
