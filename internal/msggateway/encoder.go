package msggateway

import (
	"bytes"
	"encoding/gob"
	"encoding/json"

	"github.com/openimsdk/tools/errs"
)

type Encoder interface {
	Encode(data any) ([]byte, error)
	Decode(encodeData []byte, decodeData any) error
}

type GobEncoder struct{}

func NewGobEncoder() Encoder {
	return GobEncoder{}
}

func (g GobEncoder) Encode(data any) ([]byte, error) {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	if err := enc.Encode(data); err != nil {
		return nil, errs.WrapMsg(err, "GobEncoder.Encode failed", "action", "encode")
	}
	return buff.Bytes(), nil
}

func (g GobEncoder) Decode(encodeData []byte, decodeData any) error {
	buff := bytes.NewBuffer(encodeData)
	dec := gob.NewDecoder(buff)
	if err := dec.Decode(decodeData); err != nil {
		return errs.WrapMsg(err, "GobEncoder.Decode failed", "action", "decode")
	}
	return nil
}

type JsonEncoder struct{}

func NewJsonEncoder() Encoder {
	return JsonEncoder{}
}

func (g JsonEncoder) Encode(data any) ([]byte, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return nil, errs.New("JsonEncoder.Encode failed", "action", "encode")
	}
	return b, nil
}

func (g JsonEncoder) Decode(encodeData []byte, decodeData any) error {
	err := json.Unmarshal(encodeData, decodeData)
	if err != nil {
		return errs.New("JsonEncoder.Decode failed", "action", "decode")
	}
	return nil
}
