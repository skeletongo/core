package network

import (
	"github.com/skeletongo/core/log"
)

type action struct {
	s     ISession
	msgID int
	msg   interface{}
}

func (a *action) do() {
	h := GetHandler(a.msgID)
	if h == nil {
		log.Errorf("%v not register handler", a.msgID)
		return
	}
	if err := h.Process(a.s, a.msgID, a.msg); err != nil {
		log.Errorf("%v process error: %v", a.msgID, err)
	}
}
