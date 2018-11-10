package handlers

import (
	. "raspi-bot/messages"
)

type BaseHandler struct{

}


var msgBuilder MessageBuilder


func init()  {
	msgBuilder = Builder()
}
