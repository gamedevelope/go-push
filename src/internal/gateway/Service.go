package gateway

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

type Service struct {
	server *http.Server
}

var (
	gService *Service
)

// 全量推送POST msg={}
func handlePushAll(resp http.ResponseWriter, req *http.Request) {
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
		GMerger.PushAll(&msgArr[msgIdx])
	}
}

// 房间推送POST room=xxx&msg
func handlePushRoom(resp http.ResponseWriter, req *http.Request) {
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
		GMerger.PushRoom(room, &msgArr[msgIdx])
	}
}

// 统计
func handleStats(resp http.ResponseWriter, req *http.Request) {
	var (
		data []byte
		err  error
	)

	if data, err = GStats.Dump(); err != nil {
		return
	}

	resp.Write(data)
}

func InitService() (err error) {
	var (
		mux    *http.ServeMux
		server *http.Server
	)

	// 路由
	mux = http.NewServeMux()
	mux.HandleFunc("/push/all", handlePushAll)
	mux.HandleFunc("/push/room", handlePushRoom)
	mux.HandleFunc("/stats", handleStats)

	// HTTP/2 TLS服务
	server = &http.Server{
		Addr:         ":" + strconv.Itoa(GConfig.ServicePort),
		ReadTimeout:  time.Duration(GConfig.ServiceReadTimeout) * time.Millisecond,
		WriteTimeout: time.Duration(GConfig.ServiceWriteTimeout) * time.Millisecond,
		Handler:      mux,
	}

	// 赋值全局变量
	gService = &Service{
		server: server,
	}

	go gService.server.ListenAndServe()
	return
}
