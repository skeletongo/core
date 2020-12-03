package network

import "fmt"

type ParsePacketTypeErr struct {
	EncodeType int16
	MsgID      int16
	Err        error
}

func (e *ParsePacketTypeErr) Error() string {
	return fmt.Sprintf("cannot parse proto type:%v msgID:%v err:%v", e.EncodeType, e.MsgID, e.Err)
}

func NewParsePacketTypeErr(et, msgID int16, err error) *ParsePacketTypeErr {
	return &ParsePacketTypeErr{EncodeType: et, MsgID: msgID, Err: err}
}
