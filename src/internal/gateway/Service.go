package gateway

import (
	"encoding/json"
	"net/http"
)

type Service struct {
	server *http.Server
}

// 全量推送POST msg={}
func (s *Service) handlePushAll(resp http.ResponseWriter, req *http.Request) {
	var (
		err    error
		items  string
		msgArr []json.RawMessage
		msgIdx int
	)
	if err = req.ParseForm(); err != nil {
		return
	}

	items = req.PostForm.Get("items")
	if err = json.Unmarshal([]byte(items), &msgArr); err != nil {
		return
	}

	for msgIdx, _ = range msgArr {
		defaultServer.gMerger.PushAll(&msgArr[msgIdx])
	}
}

// 房间推送POST room=xxx&msg
func (s *Service) handlePushRoom(resp http.ResponseWriter, req *http.Request) {
	var (
		err    error
		room   string
		items  string
		msgArr []json.RawMessage
		msgIdx int
	)
	if err = req.ParseForm(); err != nil {
		return
	}

	room = req.PostForm.Get("room")
	items = req.PostForm.Get("items")

	if err = json.Unmarshal([]byte(items), &msgArr); err != nil {
		return
	}

	for msgIdx, _ = range msgArr {
		defaultServer.gMerger.PushRoom(room, &msgArr[msgIdx])
	}
}

// 统计
func (s *Service) handleStats(resp http.ResponseWriter, req *http.Request) {
	var (
		data []byte
		err  error
	)

	if data, err = defaultServer.gStats.Dump(); err != nil {
		return
	}

	resp.Write(data)
}
