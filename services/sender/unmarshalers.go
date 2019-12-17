package sender

import (
	"encoding/json"

	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
)

var (
	_ MessageUnmarshaler = ProtoUnmarshaler{}
	_ MessageUnmarshaler = JSONUnmarshaler{}
)

type ProtoUnmarshaler struct{}

func (ProtoUnmarshaler) Unmarshal(b []byte, to interface{}) error {
	message, ok := to.(proto.Message)
	if !ok {
		return errors.New("bad struct type")
	}

	return proto.Unmarshal(b, message)
}

type JSONUnmarshaler struct{}

func (JSONUnmarshaler) Unmarshal(b []byte, to interface{}) error {
	message, ok := to.(proto.Message)
	if !ok {
		return errors.New("bad struct type")
	}

	return json.Unmarshal(b, message)
}
