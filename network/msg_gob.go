package network

import (
	"bytes"
	"encoding/gob"
)

var gobEncDecoder = &GobEncDecoder{}

type GobEncDecoder struct {
}

func (d *GobEncDecoder) Unmarshal(data []byte, msg interface{}) error {
	return gob.NewDecoder(bytes.NewBuffer(data)).Decode(msg)
}

func (d *GobEncDecoder) Marshal(msg interface{}) ([]byte, error) {
	data := new(bytes.Buffer)
	err := gob.NewEncoder(data).Encode(msg)
	return data.Bytes(), err
}

func init() {
	RegisterEncoding(TypeGob, gobEncDecoder)
}
