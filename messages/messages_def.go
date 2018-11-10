package messages

import (
	m "raspi-bot/messaging"
)

const (
	ProtocolVersion = 1
	ShiftValueForRepliesOpCodes = 1000
)

const (
	RbBITOpCode							= m.OpCode(0)
	RbDcMotorCmdOpCode					= m.OpCode(1)	//
)

const(
	RbDcMotorCmdReplyOpCode = m.OpCode( RbDcMotorCmdOpCode + ShiftValueForRepliesOpCodes)
)

type factory func() m.IMessage

var factories =  map[m.OpCode]factory {
	RbBITOpCode 			: func() m.IMessage { return &RbBITMessage{}},
	RbDcMotorCmdOpCode  	: func() m.IMessage { return &DcMotorCmdMessage{}},
	RbDcMotorCmdReplyOpCode	: func() m.IMessage { return &DcMotorCmdResponse{}},
}


type BaseMessage struct {
	m.MessageHeader
}

type RbCommandResponse struct {
	BaseMessage
	Status byte `json:"status"`
}

const (
	AsServerConnected 		byte = 0
	AsServerDisconnected 	byte = 1
)

type AsServerConnectionStateChanged struct {
	m.MessageHeader
	StateInfo struct{
		AsServerUid 	string
		ConnectionState byte
	}
}




