package gateway

import "github.com/gamedevelope/go-push/src/common"

// PushJob 推送任务
type PushJob struct {
	pushType int    // 推送类型
	roomId   string // 房间ID
	// union {
	bizMsg *common.BizMessage // 未序列化的业务消息
	wsMsg  *common.WSMessage  //  已序列化的业务消息
	// }
}

// ConnMgr 连接管理器
type ConnMgr struct {
	buckets []*Bucket
	jobChan []chan *PushJob // 每个Bucket对应一个Job Queue

	dispatchChan chan *PushJob // 待分发消息队列
}

// 消息分发到Bucket
func (connMgr *ConnMgr) dispatchWorkerMain(dispatchWorkerIdx int) {
	var (
		bucketIdx int
		pushJob   *PushJob
		err       error
	)
	for {
		select {
		case pushJob = <-connMgr.dispatchChan:
			defaultServer.gStats.DispatchpendingDesc()

			// 序列化
			if pushJob.wsMsg, err = common.EncodeWSMessage(pushJob.bizMsg); err != nil {
				continue
			}
			// 分发给所有Bucket, 若Bucket拥塞则等待
			for bucketIdx, _ = range connMgr.buckets {
				defaultServer.gStats.PushjobpendingIncr()
				connMgr.jobChan[bucketIdx] <- pushJob
			}
		}
	}
}

// Job负责消息广播给客户端
func (connMgr *ConnMgr) jobWorkerMain(jobWorkerIdx int, bucketIdx int) {
	bucket := connMgr.buckets[bucketIdx]

	for {
		select {
		case pushJob := <-connMgr.jobChan[bucketIdx]: // 从Bucket的job queue取出一个任务
			defaultServer.gStats.PushjobpendingDesc()
			if pushJob.pushType == common.PushTypeAll {
				bucket.PushAll(pushJob.wsMsg)
			} else if pushJob.pushType == common.PushTypeRoom {
				bucket.PushRoom(pushJob.roomId, pushJob.wsMsg)
			}
		}
	}
}

func (connMgr *ConnMgr) GetBucket(wsConnection *WSConnection) (bucket *Bucket) {
	bucket = connMgr.buckets[wsConnection.connId%uint64(len(connMgr.buckets))]
	return
}

func (connMgr *ConnMgr) AddConn(wsConnection *WSConnection) {
	bucket := connMgr.GetBucket(wsConnection)
	bucket.AddConn(wsConnection)

	defaultServer.gStats.OnlineconnectionsIncr()
}

func (connMgr *ConnMgr) DelConn(wsConnection *WSConnection) {
	var (
		bucket *Bucket
	)

	bucket = connMgr.GetBucket(wsConnection)
	bucket.DelConn(wsConnection)

	defaultServer.gStats.OnlineconnectionsDesc()
}

func (connMgr *ConnMgr) JoinRoom(roomId string, wsConn *WSConnection) (err error) {
	var (
		bucket *Bucket
	)

	bucket = connMgr.GetBucket(wsConn)
	err = bucket.JoinRoom(roomId, wsConn)
	return
}

func (connMgr *ConnMgr) LeaveRoom(roomId string, wsConn *WSConnection) (err error) {
	var (
		bucket *Bucket
	)

	bucket = connMgr.GetBucket(wsConn)
	err = bucket.LeaveRoom(roomId, wsConn)
	return
}

// PushAll 向所有在线用户发送消息
func (connMgr *ConnMgr) PushAll(bizMsg *common.BizMessage) (err error) {
	var (
		pushJob *PushJob
	)

	pushJob = &PushJob{
		pushType: common.PushTypeAll,
		bizMsg:   bizMsg,
	}

	select {
	case connMgr.dispatchChan <- pushJob:
		defaultServer.gStats.DispatchpendingIncr()
	default:
		err = common.ErrDispatchChannelFull
		defaultServer.gStats.DispatchfailIncr()
	}
	return
}

// PushRoom 向指定房间发送消息
func (connMgr *ConnMgr) PushRoom(roomId string, bizMsg *common.BizMessage) (err error) {
	var (
		pushJob *PushJob
	)

	pushJob = &PushJob{
		pushType: common.PushTypeRoom,
		bizMsg:   bizMsg,
		roomId:   roomId,
	}

	select {
	case connMgr.dispatchChan <- pushJob:
		defaultServer.gStats.DispatchpendingIncr()
	default:
		err = common.ErrDispatchChannelFull
		defaultServer.gStats.DispatchfailIncr()
	}
	return
}
