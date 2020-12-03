package network

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/golang/protobuf/proto"
)

const (
	TypeNil = iota
	TypeGPB
	TypeBinary
	TypeGob
	TypeMax
)

// MsgHead 应用层协议包头
type MsgHead struct {
	// 编码类型
	EncodeType int16
	// 消息号
	MsgID int16
}

// LenOfMsgHead 应用层协议包头长度
var LenOfMsgHead = binary.Size(&MsgHead{})

// EncDecoder 应用层数据编解码方式
type EncDecoder interface {
	Unmarshal(buf []byte, data interface{}) error
	Marshal(data interface{}) ([]byte, error)
}

var defaultEndian binary.ByteOrder = binary.LittleEndian

// SetEndian 设置大小端序
func SetEndian(endian binary.ByteOrder) {
	defaultEndian = endian
}

var encodingType [TypeMax]EncDecoder

// RegisterEncoding 注册编码方式
func RegisterEncoding(typeName int, ed EncDecoder) {
	if encodingType[typeName] != nil {
		panic(fmt.Sprintf("repeated registe EncDecoder %d", typeName))
	}
	encodingType[typeName] = ed
}

// typeTest 消息编码类型判断
func typeTest(msg interface{}) int {
	switch msg.(type) {
	case proto.Message:
		return TypeGPB
	case []byte:
		return TypeBinary
	default:
		return TypeGob
	}
}

// Marshal 消息编码
func Marshal(msgID int, msg interface{}) ([]byte, error) {
	et := typeTest(msg)
	if et < TypeNil || et >= TypeMax {
		return nil, fmt.Errorf("MarshalMessage unkown data type:%v", et)
	}

	data, err := encodingType[et].Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("MarshalMessage MsgID:%v Error:%v", msg, err.Error())
	}

	head := &MsgHead{
		EncodeType: int16(et),
		MsgID:      int16(msgID),
	}
	w := bytes.NewBuffer(nil)
	if err = binary.Write(w, defaultEndian, head); err != nil {
		return nil, err
	}
	if err = binary.Write(w, defaultEndian, data); err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

// Unmarshal 消息解码
func Unmarshal(data []byte) (msgID int, msg interface{}, err error) {
	var head MsgHead
	if err := binary.Read(bytes.NewReader(data), defaultEndian, &head); err != nil {
		return int(head.MsgID), nil, err
	}

	if head.EncodeType < TypeNil || head.EncodeType >= TypeMax {
		return int(head.MsgID), nil,
			NewParsePacketTypeErr(head.EncodeType, head.MsgID, fmt.Errorf("EncodeType:%d unregiste", head.EncodeType))
	}

	msg = CreateMessage(int(head.MsgID))
	if msg == nil {
		return int(head.MsgID), nil,
			NewParsePacketTypeErr(head.EncodeType, head.MsgID, fmt.Errorf("MsgID:%d unregiste", head.MsgID))
	}

	return int(head.MsgID), msg, encodingType[head.EncodeType].Unmarshal(data[LenOfMsgHead:], msg)
}

func MarshalNoMsgID(msg interface{}) (data []byte, err error) {
	et := typeTest(msg)
	if et < TypeNil || et >= TypeMax {
		return nil, fmt.Errorf("MarshalNoMsgId unkown data type:%v", et)
	}

	data, err = encodingType[et].Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("MarshalNoMsgId MsgID:%v Error:%v", msg, err.Error())
	}

	head := &MsgHead{
		EncodeType: int16(et),
	}
	w := bytes.NewBuffer(nil)
	if err = binary.Write(w, defaultEndian, head); err != nil {
		return nil, err
	}
	if err = binary.Write(w, defaultEndian, data); err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

func UnmarshalNoMsgID(data []byte, msg interface{}) error {
	var head MsgHead
	if err := binary.Read(bytes.NewReader(data), defaultEndian, &head); err != nil {
		return err
	}

	if head.EncodeType < TypeNil || head.EncodeType >= TypeMax {
		return NewParsePacketTypeErr(head.EncodeType, head.MsgID, fmt.Errorf("EncodeType:%d unregiste", head.EncodeType))
	}

	if err := encodingType[head.EncodeType].Unmarshal(data[LenOfMsgHead:], msg); err != nil {
		return NewParsePacketTypeErr(head.EncodeType, head.MsgID, err)
	}
	return nil
}
