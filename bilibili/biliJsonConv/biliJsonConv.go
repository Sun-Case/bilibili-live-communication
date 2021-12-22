package biliJsonConv

import (
	jsoniter "github.com/json-iterator/go"
)

type CmdType int

const (
	CodeType CmdType = iota
	DanMuMsgType
	InteractWordType
	OnlineRankCountType
	StopLiveRoomListType
	RoomRealTimeMessageUpdateType
	OnlineRankV2Type
	HotRankChangedV2Type
	PreparingType
	SendGiftType
	EntryEffectType
	LiveInteractiveGameType
	ComboSendType
	NoticeMsgType
	OnlineRankTop3Type
	LiveType
	RoomChangeType
	WidgetBannerType
	GuardBuyType

	UnknownType
	ErrorType
)

type BaseKey struct {
	Code jsoniter.Number `json:"code"`
	Cmd  string          `json:"cmd"`
}

// GetCmdType
// 获取数据类型
func GetCmdType(b []byte) CmdType {
	var bk BaseKey
	if err := jsoniter.Unmarshal(b, &bk); err != nil {
		return ErrorType
	}

	if bk.Cmd == "" {
		return CodeType
	}

	switch bk.Cmd {
	case "DANMU_MSG":
		return DanMuMsgType
	case "INTERACT_WORD":
		return InteractWordType
	case "ONLINE_RANK_COUNT":
		return OnlineRankCountType
	case "STOP_LIVE_ROOM_LIST":
		return StopLiveRoomListType
	case "ROOM_REAL_TIME_MESSAGE_UPDATE":
		return RoomRealTimeMessageUpdateType
	case "ONLINE_RANK_V2":
		return OnlineRankV2Type
	case "HOT_RANK_CHANGED_V2":
		return HotRankChangedV2Type
	case "PREPARING":
		return PreparingType
	case "SEND_GIFT":
		return SendGiftType
	case "ENTRY_EFFECT":
		return EntryEffectType
	case "LIVE_INTERACTIVE_GAME":
		return LiveInteractiveGameType
	case "COMBO_SEND":
		return ComboSendType
	case "NOTICE_MSG":
		return NoticeMsgType
	case "ONLINE_RANK_TOP3":
		return OnlineRankTop3Type
	case "LIVE":
		return LiveType
	case "ROOM_CHANGE":
		return RoomChangeType
	case "WIDGET_BANNER":
		return WidgetBannerType
	case "GUARD_BUY":
		return GuardBuyType
	default:
		return UnknownType
	}
}
