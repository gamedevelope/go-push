package gateway

import (
	"encoding/json"
	"github.com/gamedevelope/go-push/src/common"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"strings"
)

type Service struct {
	server *http.Server
}

type SendRangeEnum string

const (
	SendRangeAll  SendRangeEnum = "all"
	SendRangeRoom SendRangeEnum = "room"
	SendRangeOne  SendRangeEnum = "one"
)

type SendReq struct {
	Range    string          `json:"range"`
	UniqueId string          `json:"unique_id"`
	Message  json.RawMessage `json:"message"`
}

// handleSend ...
func (s *Service) handleSend(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	resp.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	resp.Header().Set("content-type", "application/json")             //返回数据格式是json

	if strings.EqualFold(req.Method, "OPTIONS") {
		return
	}

	var sr SendReq
	if err := json.NewDecoder(req.Body).Decode(&sr); err != nil {
		logrus.Errorf(`%v`, err)
		return
	}

	logrus.Infof(`%+v`, string(sr.Message))

	switch SendRangeEnum(sr.Range) {
	case SendRangeAll:
		if err := gServer.gMerger.PushAll(&sr.Message); err != nil {
			logrus.Errorf(`%v`, err)
			return
		}
	case SendRangeRoom:
		if err := gServer.gMerger.PushRoom(sr.UniqueId, &sr.Message); err != nil {
			logrus.Errorf(`%v`, err)
			return
		}
	case SendRangeOne:
		bizMessage := common.BizMessageData{
			List: []*json.RawMessage{&sr.Message},
		}

		buf, err := json.Marshal(bizMessage)
		if err != nil {
			return
		}

		bizMsg := &common.BizMessage{
			Type: common.MESSAGE,
			Data: buf,
		}

		// unique_id 转成 uint64
		if connId, err := strconv.ParseUint(sr.UniqueId, 10, 64); err != nil {
			logrus.Errorf(`%v`, err)
			return
		} else if err = gServer.connMgr.PushOne(connId, bizMsg); err != nil {
			logrus.Errorf(`%v`, err)
			return
		}
	default:
		return
	}
}

// 统计
func (s *Service) handleStats(resp http.ResponseWriter, req *http.Request) {
	data, err := gServer.gStats.Dump()
	if err != nil {
		return
	}

	resp.Write(data)
}
