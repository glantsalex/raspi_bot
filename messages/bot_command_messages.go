package messages

import (
	"fmt"
	m "raspi-bot/messaging"
)

const(
	DcMotorForward  = 1
	DcMotorBackward = 2
	DcMotorLeft		= 3
	DcMotorRight	= 4
)


type RbBITMessage struct{
	BaseMessage
}


type DcMotorCmdMessage struct{
	BaseMessage
	PwmLevel uint8 `json:"pwm_level"`
	MotorCmd byte `json:"motor_cmd"`
}

type DcMotorCmdResponse struct{
	m.MessageHeader
}


func ( m *DcMotorCmdMessage) String() string {
	return fmt.Sprintf("DcMotorCmdMessage. Cmd: %d Pwm level: %d ", m.MotorCmd, m.PwmLevel )
}