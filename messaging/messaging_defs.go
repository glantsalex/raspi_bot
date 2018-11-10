package messaging

import (
. "raspi-bot/errors_wrapper"
)

type OpCode				uint16
type MessageId			uint64
type Version 			uint8
type MessageTimeStamp	uint64
type ErrorCode			int

type IMessage interface{
	Opcode() 	OpCode
	MessageID() MessageId
	ErrorCode() FacilityErrCode
}

type MessageHeader struct{
	OpCode 			OpCode			`json:"op"`
	MessageId		MessageId		`json:"id"`
	CorrelationId	MessageId		`json:"cid"`
	ErrCode			FacilityErrCode `json:"ok"`
	TimeStamp 		MessageTimeStamp`json:"ts"`
	RequestSendTime int64			`json:"-"`
	Version			Version			`json:"ver"`
}

func (h MessageHeader ) Opcode() OpCode{
	return h.OpCode
}

func (h MessageHeader ) ErrorCode() FacilityErrCode{
	return h.ErrCode
}

func (h MessageHeader ) MessageID() MessageId{
	return h.MessageId
}

type IHandler interface {
	Handle( message IMessage ) (IMessage, IFacilityError )
}

type ICodec interface {
	Encode( message IMessage) ( []byte, error)
	Decode( buffer []byte ) (IMessage, error)
}

type ReplyMarker interface {
	MarkAsReply()
}

type ResponseMessage struct {
	MessageHeader
	Body interface{} `json:"body,omitempty"`
}

func ( m ResponseMessage ) MarkAsReply(){}