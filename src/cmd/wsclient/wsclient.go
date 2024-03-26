package wsclient

import (
	"encoding/json"
	"fmt"
	"github.com/gamedevelope/go-push/src/cli"
	"github.com/gamedevelope/go-push/src/common"
	"github.com/gamedevelope/go-push/src/internal/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"log"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type websocketClientManager struct {
	conn        *websocket.Conn
	addr        *string
	path        string
	sendMsgChan chan string
	recvMsgChan chan string
	isAlive     bool
	timeout     int
}

// NewWsClientManager 构造函数
func NewWsClientManager(addrIp, addrPort, path string, timeout int) *websocketClientManager {
	addrString := addrIp + ":" + addrPort
	var sendChan = make(chan string, 10) //定义channel大小，需要及时处理消费，否则会阻塞
	var recvChan = make(chan string, 10) //定义channel大小，需要及时处理消费，否则会阻塞
	var conn *websocket.Conn
	return &websocketClientManager{
		addr:        &addrString,
		path:        path,
		conn:        conn,
		sendMsgChan: sendChan,
		recvMsgChan: recvChan,
		isAlive:     false,
		timeout:     timeout,
	}
}

// 链接服务端
func (wsc *websocketClientManager) dail() {
	var err error
	u := url.URL{Scheme: "ws", Host: *wsc.addr, Path: wsc.path}
	fmt.Println("connecting to:", u.String())
	wsc.conn, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	wsc.isAlive = true
	log.Printf("connecting to %s 链接成功！！！", u.String())

}

// 发送消息到服务端
func (wsc *websocketClientManager) sendMsgThread() {
	go func() {
		for {
			msg := <-wsc.sendMsgChan
			fmt.Println("发送消息:", msg)
			// websocket.TextMessage类型
			err := wsc.conn.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				fmt.Println("write:", err)
				continue
			}
		}
	}()
}

// 读取服务端消息
func (wsc *websocketClientManager) readMsgThread() {
	go func() {
		for {
			if wsc.conn != nil {
				_, message, err := wsc.conn.ReadMessage()
				if err != nil {
					log.Println("readErr:", err)
					wsc.isAlive = false
					// 出现错误，退出读取，尝试重连
					break
				}
				// 需要读取数据，不然会阻塞
				wsc.recvMsgChan <- string(message)

			}
		}
	}()
}

// 开启服务并重连
func (wsc *websocketClientManager) start() {
	for {
		if wsc.isAlive == false {
			wsc.dail()
			wsc.sendMsgThread()
			wsc.readMsgThread()
			wsc.Msg()  //构造假消息
			wsc.Recv() //接收处理服务端返回到消息
		}
		time.Sleep(time.Second * time.Duration(wsc.timeout))
	}
}

// Msg 模拟websocket心跳包，假数据
func (wsc *websocketClientManager) Msg() {
	joinData := common.BizJoinData{Room: `room1`}
	bs, _ := json.Marshal(joinData)
	msg := common.BizMessage{
		Type: "JOIN",
		Data: bs,
	}

	bs, _ = json.Marshal(msg)

	wsc.sendMsgChan <- string(bs)

	go func() {
		for {
			msg = common.BizMessage{
				Type: "PING",
				Data: nil,
			}
			bs, _ = json.Marshal(msg)

			wsc.sendMsgChan <- string(bs)
			time.Sleep(time.Second * 10)
		}
	}()
}

// Recv 接收处理服务端返回到消息
func (wsc *websocketClientManager) Recv() {
	go func() {
		for {
			msg, ok := <-wsc.recvMsgChan
			if ok {
				fmt.Println("收到消息：", msg)
			}
		}
	}()
}

var (
	appPath string
	appMode string
)

func init() {
	serveCmd := &cobra.Command{
		Use:   "client",
		Short: "开启同步服务",
		Long:  "开启同步服务，同步链上事件并进行处理",
		Run:   clientRun,
	}

	logrus.Infof(`cli init`)
	serveCmd.Flags().StringVarP(&appPath, "app_path", `P`, ``, "application path")
	serveCmd.Flags().StringVarP(&appMode, "app_mode", `M`, ``, "application mode: local/testnet/prod")

	cli.Register(serveCmd)
}

func clientRun(cmd *cobra.Command, args []string) {
	cf := appPath + `/env/app.` + appMode + `.json`

	if err := config.Parse(cf); err != nil {
		logrus.Panicf(`error on parse file %v`, err)
	} else {
		logrus.Infof(`成功解析配置文件 %v`, cf)
	}

	wsc := NewWsClientManager("127.0.0.1", strconv.Itoa(config.AppConf.GatewayConf.WsPort), "/connect", 10)
	wsc.start()

	var w1 sync.WaitGroup
	w1.Add(1)
	w1.Wait()
}
