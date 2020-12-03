package network

var _nil = &NilEncDecoder{}

type NilEncDecoder struct {
}

func (d *NilEncDecoder) Unmarshal(data []byte, msg interface{}) error {
	return nil
}

func (d *NilEncDecoder) Marshal(msg interface{}) ([]byte, error) {
	return nil, nil
}

func init() {
	RegisterEncoding(TypeNil, _nil)
}
