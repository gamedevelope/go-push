package gateway

import (
	"encoding/json"
	"sync/atomic"
)

type Stats struct {
	// 反馈在线长连接的数量
	OnlineConnections int64 `json:"onlineConnections"`

	// 反馈客户端的推送压力
	SendMessageTotal int64 `json:"sendMessageTotal"`
	SendMessageFail  int64 `json:"sendMessageFail"`

	// 反馈ConnMgr消息分发模块的压力
	DispatchPending int64 `json:"dispatchPending"`
	PushJobPending  int64 `json:"pushJobPending"`
	DispatchFail    int64 `json:"dispatchFail"`

	// 返回出在线的房间总数, 有利于分析内存上涨的原因
	RoomCount int64 `json:"roomCount"`

	// Merger模块处理队列, 反馈出消息合并的压力情况
	MergerPending int64 `json:"mergerPending"`

	// Merger模块合并发送的消息总数与失败总数
	MergerRoomTotal int64 `json:"mergerRoomTotal"`
	MergerAllTotal  int64 `json:"mergerAllTotal"`
	MergerRoomFail  int64 `json:"mergerRoomFail"`
	MergerAllFail   int64 `json:"mergerAllFail"`
}

func (s *Stats) DispatchpendingIncr() {
	atomic.AddInt64(&s.DispatchPending, 1)
}

func (s *Stats) DispatchpendingDesc() {
	atomic.AddInt64(&s.DispatchPending, -1)
}

func (s *Stats) PushjobpendingIncr() {
	atomic.AddInt64(&s.PushJobPending, 1)
}

func (s *Stats) PushjobpendingDesc() {
	atomic.AddInt64(&s.PushJobPending, -1)
}

func (s *Stats) OnlineconnectionsIncr() {
	atomic.AddInt64(&s.OnlineConnections, 1)
}

func (s *Stats) OnlineconnectionsDesc() {
	atomic.AddInt64(&s.OnlineConnections, -1)
}

func (s *Stats) RoomCount_INCR() {
	atomic.AddInt64(&s.RoomCount, 1)
}

func (s *Stats) RoomcountDesc() {
	atomic.AddInt64(&s.RoomCount, -1)
}

func (s *Stats) MergerpendingIncr() {
	atomic.AddInt64(&s.MergerPending, 1)
}

func (s *Stats) MergerpendingDesc() {
	atomic.AddInt64(&s.MergerPending, -1)
}

func (s *Stats) MergerroomtotalIncr(batchSize int64) {
	atomic.AddInt64(&s.MergerRoomTotal, batchSize)
}

func (s *Stats) MergeralltotalIncr(batchSize int64) {
	atomic.AddInt64(&s.MergerAllTotal, batchSize)
}

func (s *Stats) MergerroomfailIncr(batchSize int64) {
	atomic.AddInt64(&s.MergerRoomFail, batchSize)
}

func (s *Stats) MergerallfailIncr(batchSize int64) {
	atomic.AddInt64(&s.MergerAllFail, batchSize)
}

func (s *Stats) DispatchfailIncr() {
	atomic.AddInt64(&s.DispatchFail, 1)
}

func (s *Stats) SendmessagefailIncr() {
	atomic.AddInt64(&s.SendMessageFail, 1)
}

func (s *Stats) SendmessagetotalIncr() {
	atomic.AddInt64(&s.SendMessageTotal, 1)
}

func (s *Stats) Dump() (data []byte, err error) {
	return json.Marshal(s)
}
