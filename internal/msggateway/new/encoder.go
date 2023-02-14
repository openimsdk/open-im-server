package new

import (
	"bytes"
	"encoding/gob"
	"open_im_sdk/pkg/utils"
)

type Encoder interface {
	Encode(data interface{}) ([]byte, error)
	Decode(encodeData []byte, decodeData interface{}) error
}

type GobEncoder struct {
}

func NewGobEncoder() *GobEncoder {
	return &GobEncoder{}
}
func (g *GobEncoder) Encode(data interface{}) ([]byte, error) {
	buff := bytes.Buffer{}
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(data)
	if err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}
func (g *GobEncoder) Decode(encodeData []byte, decodeData interface{}) error {
	buff := bytes.NewBuffer(encodeData)
	dec := gob.NewDecoder(buff)
	err := dec.Decode(decodeData)
	if err != nil {
		return utils.Wrap(err, "")
	}
	return nil
}
