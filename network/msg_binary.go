package network

import (
	"bytes"
	"encoding/binary"
)

var binaryEncDecoder = &BinaryEncDecoder{}

type BinaryEncDecoder struct {
}

func (d *BinaryEncDecoder) Unmarshal(data []byte, msg interface{}) error {
	return binary.Read(bytes.NewReader(data), defaultEndian, msg)
}

func (d *BinaryEncDecoder) Marshal(msg interface{}) ([]byte, error) {
	writer := bytes.NewBuffer(nil)
	err := binary.Write(writer, defaultEndian, msg)
	return writer.Bytes(), err
}

func init() {
	RegisterEncoding(TypeBinary, binaryEncDecoder)
}
