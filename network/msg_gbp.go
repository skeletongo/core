// Gbp
package network

import (
	"errors"

	"code.google.com/p/goprotobuf/proto"
)

var ErrorTypeNotFit = errors.New("packet not proto.Message type")

var gpbEncDecoder = &GbpEncDecoder{}

type GbpEncDecoder struct {
}

func (d *GbpEncDecoder) Unmarshal(data []byte, msg interface{}) error {
	protoMsg, ok := msg.(proto.Message)
	if !ok {
		return ErrorTypeNotFit
	}
	return proto.Unmarshal(data, protoMsg)
}

func (d *GbpEncDecoder) Marshal(msg interface{}) ([]byte, error) {
	protoMsg, ok := msg.(proto.Message)
	if !ok {
		return nil, ErrorTypeNotFit
	}
	return proto.Marshal(protoMsg)
}

func init() {
	RegisterEncoding(TypeGPB, gpbEncDecoder)
}
