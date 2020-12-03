package network

import (
	"fmt"
	"reflect"
)

type Handler interface {
	Process(s ISession, msgID int, msg interface{}) error
}

type HandlerWrapper func(s ISession, msgID int, msg interface{}) error

func (hw HandlerWrapper) Process(s ISession, msgID int, msg interface{}) error {
	return hw(s, msgID, msg)
}

var messages = make(map[int]reflect.Type)

func CreateMessage(msgID int) interface{} {
	v, ok := messages[msgID]
	if !ok {
		return nil
	}
	return reflect.New(v.Elem()).Interface()
}

var handlers = make(map[int]Handler)

// SetHandler
func SetHandler(msgID int, msg interface{}, handler Handler) {
	if _, ok := messages[msgID]; ok {
		panic(fmt.Sprintf("message already exist, msgID: %d", msgID))
		return
	}

	msgType := reflect.TypeOf(msg)
	if msgType == nil || msgType.Kind() != reflect.Ptr {
		panic(fmt.Sprintf("message pointer required, msgID: %d", msgID))
		return
	}

	if handler == nil {
		panic(fmt.Sprintf("message handler is nil, msgID: %d", msgID))
		return
	}

	messages[msgID] = msgType
	handlers[msgID] = handler
}

// GetHandler
func GetHandler(msgID int) Handler {
	return handlers[msgID]
}
