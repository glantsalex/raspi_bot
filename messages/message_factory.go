package messages

import (
	"sync/atomic"
	"time"

	. "raspi-bot/errors_wrapper"
	m "raspi-bot/messaging"
)

var messageId uint64 = 0

type MessageBuilder interface {
	Build(code m.OpCode, cid m.MessageId, fe IFacilityError, body interface{}) m.IMessage
}

var builder MessageBuilder

type msgBuilder struct{}

func (_ msgBuilder) Build(opcode m.OpCode, cid m.MessageId, fe IFacilityError, body interface{}) m.IMessage {
	var errCode = Success
	if fe != nil {
		errCode = fe.ErrorCode()
	}
	r := m.ResponseMessage{MessageHeader: createMessageHeader(opcode, cid, errCode), Body: body}
	return &r
}

func Builder() MessageBuilder {
	if builder == nil {
		builder = msgBuilder{}
	}
	return builder
}

func createMessageHeader(opCode m.OpCode, corrId m.MessageId, errCode FacilityErrCode) m.MessageHeader {
	return m.MessageHeader{
		MessageId:     getMessageId(),
		CorrelationId: corrId,
		OpCode:        opCode,
		ErrCode:       errCode,
		Version:       ProtocolVersion,
		TimeStamp:     m.MessageTimeStamp(time.Now().Unix()),
	}
}

func getMessageId() m.MessageId {
	atomic.AddUint64(&messageId, 1)
	return m.MessageId(atomic.LoadUint64(&messageId))
}
