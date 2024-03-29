package gateway

import (
	"github.com/gamedevelope/go-push/src/common"
	"github.com/gorilla/websocket"
	"sync"
	"time"
)

type WSConnection struct {
	mutex             sync.Mutex
	connId            uint64
	wsSocket          *websocket.Conn
	inChan            chan *common.WSMessage
	outChan           chan *common.WSMessage
	closeChan         chan byte
	isClosed          bool
	lastHeartbeatTime time.Time       // 最近一次心跳时间
	rooms             map[string]bool // 加入了哪些房间
}

// 读websocket
func (wsConnection *WSConnection) readLoop() {
	var (
		msgType int
		msgData []byte
		message *common.WSMessage
		err     error
	)
	for {
		if msgType, msgData, err = wsConnection.wsSocket.ReadMessage(); err != nil {
			goto ERR
		}

		message = common.BuildWSMessage(msgType, msgData)

		select {
		case wsConnection.inChan <- message:
		case <-wsConnection.closeChan:
			goto CLOSED
		}
	}

ERR:
	wsConnection.Close()
CLOSED:
}

// 写websocket
func (wsConnection *WSConnection) writeLoop() {
	var (
		message *common.WSMessage
		err     error
	)
	for {
		select {
		case message = <-wsConnection.outChan:
			if err = wsConnection.wsSocket.WriteMessage(message.MsgType, message.MsgData); err != nil {
				goto ERR
			}
		case <-wsConnection.closeChan:
			goto CLOSED
		}
	}
ERR:
	wsConnection.Close()
CLOSED:
}

func InitWSConnection(connId uint64, wsSocket *websocket.Conn) (wsConnection *WSConnection) {
	wsConnection = &WSConnection{
		wsSocket:          wsSocket,
		connId:            connId,
		inChan:            make(chan *common.WSMessage, GConfig.WsInChannelSize),
		outChan:           make(chan *common.WSMessage, GConfig.WsOutChannelSize),
		closeChan:         make(chan byte),
		lastHeartbeatTime: time.Now(),
		rooms:             make(map[string]bool),
	}

	go wsConnection.readLoop()
	go wsConnection.writeLoop()

	return
}

// SendMessage 发送消息
func (wsConnection *WSConnection) SendMessage(message *common.WSMessage) (err error) {
	select {
	case wsConnection.outChan <- message:
		SendmessagetotalIncr()
	case <-wsConnection.closeChan:
		err = common.ErrConnectionLoss
	default: // 写操作不会阻塞, 因为channel已经预留给websocket一定的缓冲空间
		err = common.ErrSendMessageFull
		SendmessagefailIncr()
	}
	return
}

// ReadMessage 读取消息
func (wsConnection *WSConnection) ReadMessage() (message *common.WSMessage, err error) {
	select {
	case message = <-wsConnection.inChan:
	case <-wsConnection.closeChan:
		err = common.ErrConnectionLoss
	}
	return
}

// Close 关闭连接
func (wsConnection *WSConnection) Close() {
	wsConnection.wsSocket.Close()

	wsConnection.mutex.Lock()
	defer wsConnection.mutex.Unlock()

	if !wsConnection.isClosed {
		wsConnection.isClosed = true
		close(wsConnection.closeChan)
	}
}

// IsAlive 检查心跳（不需要太频繁）
func (wsConnection *WSConnection) IsAlive() bool {
	var (
		now = time.Now()
	)

	wsConnection.mutex.Lock()
	defer wsConnection.mutex.Unlock()

	// 连接已关闭 或者 太久没有心跳
	if wsConnection.isClosed || now.Sub(wsConnection.lastHeartbeatTime) > time.Duration(GConfig.WsHeartbeatInterval)*time.Second {
		return false
	}
	return true
}

// KeepAlive 更新心跳
func (wsConnection *WSConnection) KeepAlive() {
	var (
		now = time.Now()
	)

	wsConnection.mutex.Lock()
	defer wsConnection.mutex.Unlock()

	wsConnection.lastHeartbeatTime = now
}
