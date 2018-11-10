package ws

import (
m "raspi-bot/messaging"
. "net/http"
)


type restEntry struct {
	method  string
	path    string
	handler Handler
}

var clientConnected 	= make( chan *WsClient )

var clientDisconnected  = make( chan *WsClient )

var wsHandlers = make(map[m.OpCode]m.IHandler)

var restHandlers = make([]*restEntry, 0)

func RegisterWsHandler(opcode m.OpCode, handler m.IHandler) {
	wsHandlers[opcode] = handler
}

