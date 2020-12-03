package network

import (
	"net"
)

type TCPConn struct {
	*SessionState
	conn net.Conn
}

func NewTCPConn(conn net.Conn, s *SessionState) *TCPConn {
	t := &TCPConn{
		SessionState: s,
		conn:         conn,
	}
	return t
}

func (t *TCPConn) LocalAddr() net.Addr {
	return t.conn.LocalAddr()
}

func (t *TCPConn) RemoteAddr() net.Addr {
	return t.conn.RemoteAddr()
}

func (t *TCPConn) Read(b []byte) (int, error) {
	return t.conn.Read(b)
}

// read goroutine
func (t *TCPConn) ReadMsg() ([]byte, error) {
	return Decode(t.readBuf, t)
}

// write goroutine
func (t *TCPConn) WriteMsg(data []byte) error {
	ps, err := Encode(t.writeBuf, data)
	if err != nil {
		return err
	}

	for _, v := range ps {
		if _, err = t.conn.Write(v); err != nil {
			return err
		}
	}
	return nil
}

func (t *TCPConn) Close() {
	t.conn.Close()
}
