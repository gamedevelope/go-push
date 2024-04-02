package gateway

import (
	"github.com/gamedevelope/go-push/src/common"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"net"
	"net/http"
	"strconv"
	"time"
)

var (
	defaultServer *Server
)

type Server struct {
	cfg     *Config
	connMgr *ConnMgr

	wsServer   *WSServer
	wsUpgrader *websocket.Upgrader

	gStats   *Stats
	gMerger  *Merger
	gService *Service
}

func (s *Server) InitStats(c *Config) (err error) {
	s.gStats = &Stats{}
	return
}

func (s *Server) InitConnMgr(c *Config) (err error) {
	var (
		bucketIdx         int
		jobWorkerIdx      int
		dispatchWorkerIdx int
		cm                *ConnMgr
	)

	cm = &ConnMgr{
		buckets:      make([]*Bucket, c.BucketCount),
		jobChan:      make([]chan *PushJob, c.BucketCount),
		dispatchChan: make(chan *PushJob, c.DispatchChannelSize),
	}

	for bucketIdx, _ = range cm.buckets {
		cm.buckets[bucketIdx] = InitBucket(bucketIdx)                       // 初始化Bucket
		cm.jobChan[bucketIdx] = make(chan *PushJob, c.BucketJobChannelSize) // Bucket的Job队列
		// Bucket的Job worker
		for jobWorkerIdx = 0; jobWorkerIdx < c.BucketJobWorkerCount; jobWorkerIdx++ {
			go cm.jobWorkerMain(jobWorkerIdx, bucketIdx)
		}
	}
	// 初始化分发协程, 用于将消息扇出给各个Bucket
	for dispatchWorkerIdx = 0; dispatchWorkerIdx < c.DispatchWorkerCount; dispatchWorkerIdx++ {
		go cm.dispatchWorkerMain(dispatchWorkerIdx)
	}

	s.connMgr = cm
	return
}

func (s *Server) InitWSServer(c *Config) (err error) {
	var (
		mux      *http.ServeMux
		server   *http.Server
		listener net.Listener
	)

	// 路由
	mux = http.NewServeMux()
	mux.HandleFunc("/connect", s.handleConnect)

	// HTTP服务
	server = &http.Server{
		ReadTimeout:  time.Duration(c.WsReadTimeout) * time.Millisecond,
		WriteTimeout: time.Duration(c.WsWriteTimeout) * time.Millisecond,
		Handler:      mux,
	}

	// 监听端口
	if listener, err = net.Listen("tcp", ":"+strconv.Itoa(c.WsPort)); err != nil {
		return
	}

	// 赋值全局变量
	s.wsServer = &WSServer{
		server:    server,
		curConnId: uint64(time.Now().Unix()),
	}

	// 拉起服务
	go server.Serve(listener)

	return
}

func (s *Server) InitMerger(c *Config) (err error) {
	s.gMerger = &Merger{
		MergerWorkerCount: c.MergerWorkerCount,
		roomWorkers:       make([]*MergeWorker, c.MergerWorkerCount),
	}

	for workerIdx := 0; workerIdx < c.MergerWorkerCount; workerIdx++ {
		s.gMerger.roomWorkers[workerIdx] = initMergeWorker(common.PushTypeRoom, c)
	}
	s.gMerger.broadcastWorker = initMergeWorker(common.PushTypeAll, c)

	return
}

func (s *Server) InitService(c *Config) (err error) {
	// 路由
	mux := http.NewServeMux()
	mux.HandleFunc("/push/all", s.gService.handlePushAll)
	mux.HandleFunc("/push/room", s.gService.handlePushRoom)
	mux.HandleFunc("/stats", s.gService.handleStats)

	// HTTP/2 TLS服务
	s.gService = &Service{
		server: &http.Server{
			Addr:         ":" + strconv.Itoa(c.ServicePort),
			ReadTimeout:  time.Duration(c.ServiceReadTimeout) * time.Millisecond,
			WriteTimeout: time.Duration(c.ServiceWriteTimeout) * time.Millisecond,
			Handler:      mux,
		}}

	go s.gService.server.ListenAndServe()
	return
}

func NewServer(c *Config) *Server {
	var err error

	defaultServer = &Server{
		cfg: c,
	}

	// 统计
	if err = defaultServer.InitStats(c); err != nil {
		logrus.Panicf(`init stats %v`, err)
	}

	// 初始化连接管理器
	if err = defaultServer.InitConnMgr(c); err != nil {
		logrus.Panicf(`init conn mgr %v`, err)
	}

	// 初始化websocket服务器
	if err = defaultServer.InitWSServer(c); err != nil {
		logrus.Panicf(`init ws server %v`, err)
	}

	// 初始化merger合并层
	if err = defaultServer.InitMerger(c); err != nil {
		logrus.Panicf(`init merger %v`, err)
	}

	// 初始化service接口
	if err = defaultServer.InitService(c); err != nil {
		logrus.Panicf(`init service %v`, err)
	}

	return defaultServer
}
