package gateway

import (
	"github.com/sirupsen/logrus"
	"net/http"
)

func (s *Server) handleConnect(resp http.ResponseWriter, req *http.Request) {
	authToken := req.URL.Query().Get(`token`)
	if authToken == `` {
		http.Error(resp, `token is required`, http.StatusBadRequest)
		return
	}

	uid, err := s.auth.GetUid(s.cfg.AuthInterface, authToken)
	if err != nil {
		http.Error(resp, `token is invalid`, http.StatusUnauthorized)
		return
	}

	// WebSocket握手
	wsSocket, err := s.wsUpgrader.Upgrade(resp, req, nil)
	if err != nil {
		return
	}

	// 初始化WebSocket的读写协程
	wsConn := s.InitWSConnection(uid, wsSocket)

	logrus.Infof(`收到新链接 uid %v`, uid)

	// 开始处理websocket消息
	wsConn.WSHandle()
}
