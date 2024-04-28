package common

import (
	"encoding/json"
	"github.com/gorilla/websocket"
)

// 推送类型
const (
	PushTypeRoom = 1 // 推送房间
	PushTypeAll  = 2 // 推送在线
	PushTypeOne  = 3 // 推送单个
)

// WSMessage websocket的Message对象
type WSMessage struct {
	MsgType int
	MsgData []byte
}

type BizMessageType string

const (
	PING    BizMessageType = "PING"
	PONG    BizMessageType = "PONG"
	JOIN    BizMessageType = "JOIN"
	LEAVE   BizMessageType = "LEAVE"
	PUSH    BizMessageType = "PUSH"
	MESSAGE BizMessageType = "MESSAGE"
)

// BizMessage 业务消息的固定格式(type+data)
type BizMessage struct {
	Type BizMessageType  `json:"type"` // type消息类型: PING, PONG, JOIN, LEAVE, PUSH
	Data json.RawMessage `json:"data"` // data数据字段
}

// Data数据类型

// BizMessageData PUSH
type BizMessageData struct {
	List []*json.RawMessage `json:"list"`
}

// BizPingData PING
type BizPingData struct{}

// BizPongData PONG
type BizPongData struct{}

// BizJoinData JOIN
type BizJoinData struct {
	Room string `json:"room"`
}

// BizLeaveData LEAVE
type BizLeaveData struct {
	Room string `json:"room"`
}

func BuildWSMessage(msgType int, msgData []byte) (wsMessage *WSMessage) {
	return &WSMessage{
		MsgType: msgType,
		MsgData: msgData,
	}
}

func EncodeWSMessage(bizMessage *BizMessage) (wsMessage *WSMessage, err error) {
	var (
		buf []byte
	)
	if buf, err = json.Marshal(*bizMessage); err != nil {
		return
	}
	wsMessage = &WSMessage{websocket.TextMessage, buf}
	return
}

// DecodeBizMessage 解析{"type": "PING", "data": {...}}的包
func DecodeBizMessage(buf []byte) (bizMessage *BizMessage, err error) {
	var (
		bizMsgObj BizMessage
	)

	if err = json.Unmarshal(buf, &bizMsgObj); err != nil {
		return
	}

	bizMessage = &bizMsgObj
	return
}
