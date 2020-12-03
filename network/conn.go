package network

import (
	"net"
)

type Conn interface {
	ReadMsg() ([]byte, error)
	WriteMsg([]byte) error
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	Close()
}
