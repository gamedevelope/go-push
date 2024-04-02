package gateway

import (
	"github.com/gamedevelope/go-push/src/common"
	"github.com/sirupsen/logrus"
	"sync"
)

type Bucket struct {
	rwMutex sync.RWMutex
	index   int                      // 我是第几个桶
	id2Conn map[uint64]*WSConnection // 连接列表(key=连接唯一ID)
	rooms   map[string]*Room         // 房间列表
}

func (b *Bucket) GetConn(uid uint64) (*WSConnection, bool) {
	b.rwMutex.RLock()
	defer b.rwMutex.RUnlock()

	conn, exist := b.id2Conn[uid]

	return conn, exist
}

func (b *Bucket) GetRoom(roomId string) (*Room, bool) {
	b.rwMutex.RLock()
	defer b.rwMutex.RUnlock()

	room, exist := b.rooms[roomId]

	return room, exist
}

func (b *Bucket) SetRoom(roomId string, room *Room) {
	b.rwMutex.Lock()
	defer b.rwMutex.Unlock()

	b.rooms[roomId] = room
	defaultServer.gStats.RoomCount_INCR()
}

func (b *Bucket) DelRoom(roomId string) {
	b.rwMutex.Lock()
	defer b.rwMutex.Unlock()

	delete(b.rooms, roomId)
	defaultServer.gStats.RoomcountDesc()
}

func InitBucket(bucketIdx int) (bucket *Bucket) {
	bucket = &Bucket{
		index:   bucketIdx,
		id2Conn: make(map[uint64]*WSConnection),
		rooms:   make(map[string]*Room),
	}
	return
}

func (b *Bucket) AddConn(wsConn *WSConnection) {
	b.rwMutex.Lock()
	defer b.rwMutex.Unlock()

	b.id2Conn[wsConn.connId] = wsConn
}

func (b *Bucket) DelConn(wsConn *WSConnection) {
	b.rwMutex.Lock()
	defer b.rwMutex.Unlock()

	delete(b.id2Conn, wsConn.connId)
}

func (b *Bucket) JoinRoom(roomId string, wsConn *WSConnection) (err error) {
	var (
		existed bool
		room    *Room
	)

	// 找到房间
	if room, existed = b.GetRoom(roomId); !existed {
		room = InitRoom(roomId)
		b.SetRoom(roomId, room)
	}
	// 加入房间
	err = room.Join(wsConn)
	return
}

func (b *Bucket) LeaveRoom(roomId string, wsConn *WSConnection) (err error) {
	var (
		existed bool
		room    *Room
	)

	// 找到房间
	if room, existed = b.GetRoom(roomId); !existed {
		err = common.ErrNotInRoom
		return
	}

	err = room.Leave(wsConn)

	// 房间为空, 则删除
	if room.Count() == 0 {
		b.DelRoom(roomId)
	}
	return
}

// PushOne 单一推送
func (b *Bucket) PushOne(uid uint64, message *common.WSMessage) {
	wsConn, exist := b.GetConn(uid)

	if exist {
		err := wsConn.SendMessage(message)
		if err != nil {
			logrus.Errorf(`send message error %v`, err)
			return
		}
	}
}

// PushAll 推送给Bucket内所有用户
func (b *Bucket) PushAll(wsMsg *common.WSMessage) {
	var (
		wsConn *WSConnection
	)

	// 锁Bucket
	b.rwMutex.RLock()
	defer b.rwMutex.RUnlock()

	// 全量非阻塞推送
	for _, wsConn = range b.id2Conn {
		err := wsConn.SendMessage(wsMsg)
		if err != nil {
			logrus.Errorf(`send message error %v`, err)
			return
		}
	}
}

// PushRoom 推送给某个房间的所有用户
func (b *Bucket) PushRoom(roomId string, wsMsg *common.WSMessage) {
	var (
		room    *Room
		existed bool
	)

	// 锁Bucket
	room, existed = b.GetRoom(roomId)

	// 房间不存在
	if !existed {
		return
	}

	// 向房间做推送
	room.Push(wsMsg)
}
