package handlers

import (
	. "raspi-bot/errors_wrapper"
	. "raspi-bot/messages"
	. "raspi-bot/messaging"

	"github.com/cihub/seelog"
)

type DcMotorCommandHandler struct {
	*BaseHandler
}

func (h DcMotorCommandHandler) Handle(im IMessage) (response IMessage, fe IFacilityError) {

	msg := im.( DcMotorCmdMessage )

	seelog.Debug("DcMotorCommandHandler handler invoked. msg: %s", msg  )
	response = msgBuilder.Build(RbDcMotorCmdReplyOpCode, im.MessageID(), fe, nil )

	return
}
