package messages

import (
	"encoding/json"
	"errors"
	"fmt"

	m "raspi-bot/messaging"

	"github.com/cihub/seelog"
)

type JsonCodec struct{}

func NewJsonCodec() m.ICodec {
	return JsonCodec{}
}

func (_ JsonCodec) Encode(m m.IMessage) ([]byte, error) {

	return json.Marshal(m)
}

func (_ JsonCodec) Decode(buffer []byte) (m.IMessage, error) {

	header := &m.MessageHeader{}
	if err := json.Unmarshal(buffer, header); err != nil {
		seelog.Errorf("JSON decoder: %s", err)
		return nil, err
	}

	if _, ok := factories[header.OpCode]; !ok {
		return nil, errors.New(fmt.Sprintf("No message factory function for opcode: %d", header.OpCode))
	}

	msg := factories[header.OpCode]()

	return msg, json.Unmarshal(buffer, msg)
}
