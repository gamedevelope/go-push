package common

import "errors"

var (
	ErrConnectionLoss           = errors.New("ERR_CONNECTION_LOSS")
	ErrSendMessageFull          = errors.New("ERR_SEND_MESSAGE_FULL")
	ErrJoinRoomTwice            = errors.New("ERR_JOIN_ROOM_TWICE")
	ErrNotInRoom                = errors.New("ERR_NOT_IN_ROOM")
	ErrRoomIdInvalid            = errors.New("ERR_ROOM_ID_INVALID")
	ErrDispatchChannelFull      = errors.New("ERR_DISPATCH_CHANNEL_FULL")
	ErrMergeChannelFull         = errors.New("ERR_MERGE_CHANNEL_FULL")
	ErrCertInvalid              = errors.New("ERR_CERT_INVALID")
	ErrLogicDispatchChannelFull = errors.New("ERR_LOGIC_DISPATCH_CHANNEL_FULL")
)
