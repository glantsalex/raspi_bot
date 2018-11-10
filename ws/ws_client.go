package ws

import (
	"errors"
	"strings"
	"sync/atomic"
	"time"

	. "github.com/cihub/seelog"
	"github.com/gorilla/websocket"

	m "raspi-bot/messaging"
)

const (
	maxMessageSize = 1024 * 1024 * 5
)

type WsClient struct {
	id                       string
	port                     string
	addr                     string
	logPrefix                string
	conn                     *websocket.Conn
	codec                    m.ICodec
	send                     chan []byte
	stop                     chan bool
	closeFlag                uint32
	unregisterListener       chan *WsClient
}

func NewWsClient(id string, conn *websocket.Conn, lp string) *WsClient {
	addr := conn.RemoteAddr().String()
	idx := strings.Index(addr, ":")
	port := addr[idx+1:]
	return &WsClient{
		port:      port,
		conn:      conn,
		addr:      conn.RemoteAddr().String(),
		send:      make(chan []byte, 1000),
		stop:      make(chan bool, 1),
		id:        id,
		logPrefix: lp,
	}
}

func (c *WsClient) Port() string {
	return c.port
}

func (c *WsClient) SetCodec(codec m.ICodec) {
	c.codec = codec
}

func (c *WsClient) WriteResponse(r m.IMessage) error {

	if atomic.LoadUint32(&c.closeFlag) == 1 {
		err := errors.New("write to closed connection")
		return Errorf("WriteResponse: %s peer: %s", err, c.addr)
	}

	if buffer, err := c.codec.Encode(r); err != nil {
		return Errorf("WriteResponse: %s peer %s", err, c.addr)
	} else {
		c.send <- buffer
		Debugf("Outgoing message queue size: %d for peer: %s", len(c.send), c.addr)
	}
	return nil
}

func (c *WsClient) Id() string {
	return c.id
}

func (c *WsClient) SetId(id string) {
	c.id = id
}

func (c *WsClient) Send(r m.IMessage) error {
	return c.WriteResponse(r)
}

func (c *WsClient) Run() *WsClient {
	go c.writePump()
	go c.readPump()
	return c
}


func (c *WsClient) writePump() {
LOOP:
	for {
		select {
		case <-c.stop:
			Debugf("%s Connection stop signaled when write.", c.logPrefix)
			break LOOP
		case msg := <-c.send:
			if atomic.LoadUint32(&c.closeFlag) == 1 {
				break LOOP
			}
			c.conn.SetWriteDeadline(time.Now().Add(3 * time.Second))
			if err := c.conn.WriteMessage(websocket.TextMessage, msg); err == nil {
				Debugf("%s Message sent to client:  %s", c.logPrefix, string(msg))
			} else {
				Errorf("%s WRITE Client error:  %s. Connection closed.", c.logPrefix, err)
			}
			msg = nil
		}
	}
}

func (c *WsClient) readPump() {

	if atomic.LoadUint32(&c.closeFlag) == 1 {
		return
	}

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Time{})
	for {
		_, rawMessage, err := c.conn.ReadMessage()

		if err != nil {
			atomic.StoreUint32(&c.closeFlag, 1)
			Errorf("%s READ Client: %s.", c.logPrefix, err)
			c.stop <- true
			c.closeConn()
			break
		}
		Debugf("%s Message received: %s", c.logPrefix, string(rawMessage))
		if im, err := c.codec.Decode(rawMessage); err != nil {
			Errorf("%s Error decoding message. Msg: %s", c.logPrefix, string(rawMessage))
		} else {
			if handler, ok := wsHandlers[im.Opcode()]; ok {
				go c.handleMessage(handler, im)
			} else {
				Errorf("%s No handler registered for OPCODE: %d ", c.logPrefix, im.Opcode())
			}
		}
	}
}

func (c *WsClient) handleMessage(handler m.IHandler, msg m.IMessage) {

	start := time.Now()

	defer func() {
		timeSpentInHandler := time.Since(start).Seconds()
		Debugf("%s Time in handler for msg_id %d: %f sec.", c.logPrefix, msg.MessageID(), timeSpentInHandler)
		if timeSpentInHandler > 3.0 {
			Errorf("%s handling time exceeded %f secs for msg_id %d. peer: %s", c.logPrefix, msg.MessageID(), timeSpentInHandler)
		}
	}()

	response, fe := handler.Handle(msg)

	if fe != nil {
		Errorf("%s %s", c.logPrefix, fe.Error())
	}
	c.WriteResponse(response)
}

func (c *WsClient) drainOutgoingMessageChannel() {
	if size := len(c.send); size > 0 {
		for i := 0; i <= size; i++ {
			<-c.send
		}
	}
	close(c.send)
	c.send = nil
}

func (c *WsClient) closeConn() {
	if c.conn != nil {
		c.conn.Close()
		clientDisconnected <- c
		c.drainOutgoingMessageChannel()
		c.codec = nil
		c.conn = nil
		Debugf(" %s Connection disposed.", c.logPrefix)
	}
}



