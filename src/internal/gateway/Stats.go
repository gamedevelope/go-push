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

var (
	GStats *Stats
)

func InitStats() (err error) {
	GStats = &Stats{}
	return
}

func DispatchpendingIncr() {
	atomic.AddInt64(&GStats.DispatchPending, 1)
}

func DispatchpendingDesc() {
	atomic.AddInt64(&GStats.DispatchPending, -1)
}

func PushjobpendingIncr() {
	atomic.AddInt64(&GStats.PushJobPending, 1)
}

func PushjobpendingDesc() {
	atomic.AddInt64(&GStats.PushJobPending, -1)
}

func OnlineconnectionsIncr() {
	atomic.AddInt64(&GStats.OnlineConnections, 1)
}

func OnlineconnectionsDesc() {
	atomic.AddInt64(&GStats.OnlineConnections, -1)
}

func RoomCount_INCR() {
	atomic.AddInt64(&GStats.RoomCount, 1)
}

func RoomcountDesc() {
	atomic.AddInt64(&GStats.RoomCount, -1)
}

func MergerpendingIncr() {
	atomic.AddInt64(&GStats.MergerPending, 1)
}

func MergerpendingDesc() {
	atomic.AddInt64(&GStats.MergerPending, -1)
}

func MergerroomtotalIncr(batchSize int64) {
	atomic.AddInt64(&GStats.MergerRoomTotal, batchSize)
}

func MergeralltotalIncr(batchSize int64) {
	atomic.AddInt64(&GStats.MergerAllTotal, batchSize)
}

func MergerroomfailIncr(batchSize int64) {
	atomic.AddInt64(&GStats.MergerRoomFail, batchSize)
}

func MergerallfailIncr(batchSize int64) {
	atomic.AddInt64(&GStats.MergerAllFail, batchSize)
}

func DispatchfailIncr() {
	atomic.AddInt64(&GStats.DispatchFail, 1)
}

func SendmessagefailIncr() {
	atomic.AddInt64(&GStats.SendMessageFail, 1)
}

func SendmessagetotalIncr() {
	atomic.AddInt64(&GStats.SendMessageTotal, 1)
}

func (stats *Stats) Dump() (data []byte, err error) {
	return json.Marshal(GStats)
}
