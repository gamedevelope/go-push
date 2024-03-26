package gateway

import (
	"github.com/gamedevelope/go-push/src/common"
	"github.com/sirupsen/logrus"
	"sync"
)

// Room 房间
type Room struct {
	rwMutex  sync.RWMutex
	roomId   string
	connHash map[uint64]*WSConnection
}

func InitRoom(roomId string) (room *Room) {
	room = &Room{
		roomId:   roomId,
		connHash: make(map[uint64]*WSConnection),
	}
	return
}

func (room *Room) Join(wsConn *WSConnection) (err error) {
	room.rwMutex.Lock()
	defer room.rwMutex.Unlock()

	if _, existed := room.connHash[wsConn.connId]; existed {
		err = common.ErrJoinRoomTwice
		return
	}

	room.connHash[wsConn.connId] = wsConn
	return
}

func (room *Room) Leave(wsConn *WSConnection) (err error) {
	room.rwMutex.Lock()
	defer room.rwMutex.Unlock()

	if _, existed := room.connHash[wsConn.connId]; !existed {
		err = common.ErrNotInRoom
		return
	}

	delete(room.connHash, wsConn.connId)
	return
}

func (room *Room) Count() int {
	room.rwMutex.RLock()
	defer room.rwMutex.RUnlock()

	return len(room.connHash)
}

func (room *Room) Push(wsMsg *common.WSMessage) {
	room.rwMutex.RLock()
	defer room.rwMutex.RUnlock()

	for _, wsConn := range room.connHash {
		err := wsConn.SendMessage(wsMsg)
		if err != nil {
			logrus.Errorf(`send message error %v`, err)
			return
		}
	}
}
