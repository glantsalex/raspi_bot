package ws

import (
	"fmt"
	"innovi-as-proxy/utils"
	"net/http"
	. "raspi-bot/server/config"
	"time"

	"github.com/cihub/seelog"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WsServer struct {
	host       string
	port       string
	wsEndpoint string
	httpServer *http.Server
}

func NewWsServer(c *Config) *WsServer {

	defer func() {
		seelog.Infof("HTTP server created.\nREAD timeout: %d ms\nWRITE timeout:%d ms", 3000,3000)
	}()

	return &WsServer{
		port:       c.Port(),
		host:       c.Host(),
		wsEndpoint: c.WsApiRoot(),
		httpServer: &http.Server{
			Addr:         fmt.Sprintf("%s:%s", c.Host(), c.Port()),
			ReadTimeout:  time.Millisecond * time.Duration(3000),
			WriteTimeout: time.Millisecond * time.Duration(3000),
			Handler:      nil,
		},
	}
}

func (srv *WsServer) Run() {

	routes := mux.NewRouter()
	routes.StrictSlash(true)

	routes.HandleFunc(srv.wsEndpoint, srv.listenForWSConnections)

	routes.NotFoundHandler = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		seelog.Errorf("REST unknown route request: %s from peer: %s", r.URL.String(), utils.ResolveRemoteIpFromHttpRequest(r))
	})

	http.Handle("/", routes)

	seelog.Infof("Waiting for client connections at  %s:%s...\n", srv.host, srv.port)

	if err := srv.httpServer.ListenAndServe(); err != nil {
		seelog.Errorf("ListenAndServe: %s", err)
	}
}

func (srv *WsServer) Stop() {
	srv.httpServer.Close()
}

func (srv *WsServer) listenForWSConnections(w http.ResponseWriter, r *http.Request) {

	var err error

	seelog.Infof("connection request arrived from %s.", r.RemoteAddr)

	lp := ""

	conn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		seelog.Errorf(lp+"Error upgrading connection to WS: %s", err)
		return
	}
	clientConnected <- NewWsClient(uuid.New().String(), conn, lp)
}
