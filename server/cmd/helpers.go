package main

import (
	m "raspi-bot/messages"
	ms "raspi-bot/messaging"
	h "raspi-bot/server/handlers"
)

var wsHandlers = []struct {
	opcode  ms.OpCode
	handler ms.IHandler
}{
	{m.RbDcMotorCmdOpCode, &h.DcMotorCommandHandler{}},
}
