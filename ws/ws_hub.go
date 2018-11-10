package ws

import (
	m "raspi-bot/messaging"
	ms "raspi-bot/messages"
	"github.com/cihub/seelog"
)

type pushPayload map[string]interface{}

type WsHub struct {
	codec           m.ICodec
	msgBuilder      ms.MessageBuilder
}

func NewWsHub(codec m.ICodec ) *WsHub {

	wsHub := &WsHub{
		codec:           codec,
		msgBuilder:      ms.Builder(),
	}
	return wsHub
}

func (h *WsHub) RunAsync() {

	go func() {
		for {
			select {

			case client := <-clientConnected:
				client.SetCodec(h.codec)
				client.Run()
				seelog.Debugf("New client connection registered. Remote address is: %s", client.conn.RemoteAddr())

			case client := <-clientDisconnected:
				h.handleAsServerDisconnected(client.Id(), ms.AsServerDisconnected)

				seelog.Debugf("client unregistered: %s", client.addr)
			}
		}
	}()
}

func (h *WsHub) handleAsServerDisconnected(asUid string, state byte) {
	//msg := ms.BuildAsServerDisconnectedMsg(asUid)
	//wsHandlers[msg.Opcode()].Handle(msg)
}
