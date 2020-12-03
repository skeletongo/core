package network

import (
	"fmt"
	"testing"
)

type D struct {
	Name string
	Age  int
}

func TestMarshal(t *testing.T) {
	SetHandler(1, new(D), HandlerWrapper(func(a ISession, msgID int, msg interface{}) error {
		return nil
	}))

	data, err := Marshal(1, &D{
		Name: "Tom",
		Age:  20,
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	id, msg, err := Unmarshal(data)
	fmt.Printf("msgID:%v Msg:%v Err:%v\n", id, *msg.(*D), err)
}

func TestMarshalNoMsgID(t *testing.T) {
	data, err := MarshalNoMsgID(&D{
		Name: "Tom",
		Age:  20,
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	var msg D
	err = UnmarshalNoMsgID(data, &msg)
	fmt.Printf("Msg:%v Err:%v\n", msg, err)
}
