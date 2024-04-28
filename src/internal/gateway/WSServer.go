package gateway

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

// WSServer WebSocket服务端
type WSServer struct {
	server    *http.Server
	curConnId uint64
}

type User struct {
	Id                    int    `json:"id"`
	Nickname              string `json:"nickname"`
	Status                int    `json:"status"`
	Avatar                string `json:"avatar"`
	Gender                int    `json:"gender"`
	IpLoc                 string `json:"ip_loc"`
	Balance               int    `json:"balance"`
	IsAdmin               bool   `json:"is_admin"`
	SellerVerify          int    `json:"seller_verify"`
	CreatedOn             int    `json:"created_on"`
	Follows               int    `json:"follows"`
	Followings            int    `json:"followings"`
	TweetCount            int    `json:"tweet_count"`
	LikeCount             int    `json:"like_count"`
	SellCount             int    `json:"sell_count"`
	SeekCount             int    `json:"seek_count"`
	CommentCount          int    `json:"comment_count"`
	OfficialCertification int    `json:"official_certification"`
	ActivePoint           int    `json:"active_point"`
	ActiveLevel           int    `json:"active_level"`
	BuyPoint              int    `json:"buy_point"`
	BuyLevel              int    `json:"buy_level"`
	BuyLevelStatus        int    `json:"buy_level_status"`
	SellPoint             int    `json:"sell_point"`
	SellLevel             int    `json:"sell_level"`
	SellLevelStatus       int    `json:"sell_level_status"`
	CreditPoint           int    `json:"credit_point"`
	CreditLevel           int    `json:"credit_level"`
	CreditLevel1          int    `json:"creditLevel"`
	SaleLockAt            int    `json:"sale_lock_at"`
	SaleFoul              int    `json:"sale_foul"`
}

type UserInfoResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data User   `json:"data"`
}

func getUserInfo(bearer string) (*User, error) {
	// 从远程接口中获取认证信息
	logrus.Infof(`从远程接口中获取认证信息 %v`, gServer.cfg.AuthInterface)
	client := &http.Client{}

	// get user info from remote interface with header
	req, _ := http.NewRequest("GET", gServer.cfg.AuthInterface, nil)
	req.Header.Set("Authorization", `Bearer `+bearer)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(`%v`, resp.StatusCode)
	}

	// parse user info
	var userResp UserInfoResponse
	if err = json.NewDecoder(resp.Body).Decode(&userResp); err != nil {
		return nil, err
	}

	return &userResp.Data, nil
}

func (s *Server) handleConnect(resp http.ResponseWriter, req *http.Request) {
	// WebSocket握手
	wsSocket, err := s.wsUpgrader.Upgrade(resp, req, nil)
	if err != nil {
		return
	}

	//authToken := req.URL.Query().Get(`authorization`)

	uid := req.URL.Query().Get(`uid`)

	// 将uid 转成uint64
	connId, err := strconv.ParseUint(uid, 10, 64)
	logrus.Infof(`新的连接 %v`, connId)
	//
	//// 获取用户信息
	//userInfo, err := getUserInfo(authToken)
	//if err != nil {
	//	logrus.Error(err)
	//	return
	//}
	//logrus.Infof(`%+v`, userInfo)
	//
	//// 连接唯一标识
	//connId := uint64(userInfo.Id)

	//connId := atomic.AddUint64(&gServer.wsServer.curConnId, 1)

	// 初始化WebSocket的读写协程
	wsConn := s.InitWSConnection(connId, wsSocket)

	logrus.Infof(`收到新链接 %v`, connId)

	// 开始处理websocket消息
	wsConn.WSHandle()
}
