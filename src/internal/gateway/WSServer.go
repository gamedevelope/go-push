package gateway

import (
	"github.com/gorilla/websocket"
	"net"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"
)

// WSServer WebSocket服务端
type WSServer struct {
	server    *http.Server
	curConnId uint64
}

var (
	GWsserver *WSServer

	wsUpgrader = websocket.Upgrader{
		// 允许所有CORS跨域请求
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func handleConnect(resp http.ResponseWriter, req *http.Request) {
	var (
		err      error
		wsSocket *websocket.Conn
		connId   uint64
		wsConn   *WSConnection
	)

	// WebSocket握手
	if wsSocket, err = wsUpgrader.Upgrade(resp, req, nil); err != nil {
		return
	}

	// 连接唯一标识
	connId = atomic.AddUint64(&GWsserver.curConnId, 1)

	// 初始化WebSocket的读写协程
	wsConn = InitWSConnection(connId, wsSocket)

	// 开始处理websocket消息
	wsConn.WSHandle()
}

func InitWSServer() (err error) {
	var (
		mux      *http.ServeMux
		server   *http.Server
		listener net.Listener
	)

	// 路由
	mux = http.NewServeMux()
	mux.HandleFunc("/connect", handleConnect)

	// HTTP服务
	server = &http.Server{
		ReadTimeout:  time.Duration(GConfig.WsReadTimeout) * time.Millisecond,
		WriteTimeout: time.Duration(GConfig.WsWriteTimeout) * time.Millisecond,
		Handler:      mux,
	}

	// 监听端口
	if listener, err = net.Listen("tcp", ":"+strconv.Itoa(GConfig.WsPort)); err != nil {
		return
	}

	// 赋值全局变量
	GWsserver = &WSServer{
		server:    server,
		curConnId: uint64(time.Now().Unix()),
	}

	// 拉起服务
	go server.Serve(listener)

	return
}
