package network

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

// PkgHead 包头
type PkgHead struct {
	Len     uint16 // 长度
	Seq     uint16 // 序号
	LogicNo uint32 // 逻辑号
}

// PkgData 数据包
type PkgData struct {
	Head PkgHead
	//
	Seq  uint16
	data []byte // 业务层数据
}

var (
	// 包头长度
	PkgHeadLen = binary.Size(&PkgHead{})
	MaxDataLen = math.MaxUint16
)

func Encode(buf *PkgData, data []byte) (packets [][]byte, err error) {
	dataLen := len(data)
	if dataLen > MaxDataLen {
		// 分包
		for _, v := range PkgCutFunc(data) {
			ds, err := Encode(buf, v)
			if err != nil {
				return nil, err
			}
			packets = append(packets, ds...)
		}
		return
	}

	buf.Seq++

	ioBuf := new(bytes.Buffer)
	if err = binary.Write(ioBuf, defaultEndian, dataLen); err != nil {
		return nil, err
	}
	if err = binary.Write(ioBuf, defaultEndian, buf.Seq); err != nil {
		return nil, err
	}
	if err = binary.Write(ioBuf, defaultEndian, buf.Head.LogicNo); err != nil {
		return nil, err
	}
	if _, err = ioBuf.Write(data); err != nil {
		return nil, err
	}
	return [][]byte{ioBuf.Bytes()}, err
}

func Decode(buf *PkgData, r io.Reader) (data []byte, err error) {
	if err = binary.Read(r, defaultEndian, &buf.Head); err != nil {
		return
	}

	if int(buf.Head.Len) > MaxDataLen {
		err = fmt.Errorf("PacketHeader len exceed MaxDataLen. get %v limit %v", buf.Head.Len, MaxDataLen)
		return
	}

	if buf.Head.Seq != buf.Seq+1 {
		err = fmt.Errorf("PacketHeader sno not matched. get %v want %v", buf.Head.Seq, buf.Seq+1)
		return
	}
	buf.Seq++

	data = buf.data[0:buf.Head.Len]
	_, err = io.ReadFull(r, data)
	return
}

// PkgCutFunc 消息分包
func PkgCutFunc(data []byte) (dataArr [][]byte) {

	return nil
}
