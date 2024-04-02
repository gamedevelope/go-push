package gateway

import (
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"net/http"
	"sync/atomic"
)

// WSServer WebSocket服务端
type WSServer struct {
	server    *http.Server
	curConnId uint64
}

var (
	wsUpgrader = websocket.Upgrader{
		// 允许所有CORS跨域请求
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func (s *Server) handleConnect(resp http.ResponseWriter, req *http.Request) {
	// 检查请求是否合法

	// WebSocket握手
	wsSocket, err := wsUpgrader.Upgrade(resp, req, nil)
	if err != nil {
		return
	}

	// 连接唯一标识
	connId := atomic.AddUint64(&defaultServer.wsServer.curConnId, 1)

	// 初始化WebSocket的读写协程
	wsConn := s.InitWSConnection(connId, wsSocket)

	logrus.Infof(`收到新链接 %v`, connId)

	// 开始处理websocket消息
	wsConn.WSHandle()
}
