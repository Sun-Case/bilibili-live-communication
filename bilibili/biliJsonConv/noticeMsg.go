package biliJsonConv

import (
	"bytes"
	jsoniter "github.com/json-iterator/go"
)

type noticeMsg struct {
	Cmd  string          `json:"cmd"`
	ID   jsoniter.Number `json:"id"`
	Name string          `json:"name"`
	Full struct {
		HeadIcon    string          `json:"head_icon"`
		TailIcon    string          `json:"tail_icon"`
		HeadIconFa  string          `json:"head_icon_fa"`
		TailIconFa  string          `json:"tail_icon_fa"`
		HeadIconFan jsoniter.Number `json:"head_icon_fan"`
		TailIconFan jsoniter.Number `json:"tail_icon_fan"`
		Background  string          `json:"background"`
		Color       string          `json:"color"`
		Highlight   string          `json:"highlight"`
		Time        jsoniter.Number `json:"time"`
	} `json:"full"`
	Half struct {
		HeadIcon   string          `json:"head_icon"`
		TailIcon   string          `json:"tail_icon"`
		Background string          `json:"background"`
		Color      string          `json:"color"`
		Highlight  string          `json:"highlight"`
		Time       jsoniter.Number `json:"time"`
	} `json:"half"`
	Side struct {
		HeadIcon   string `json:"head_icon"`
		Background string `json:"background"`
		Color      string `json:"color"`
		Highlight  string `json:"highlight"`
		Border     string `json:"border"`
	} `json:"side"`
	Roomid     jsoniter.Number `json:"roomid"`
	RealRoomid jsoniter.Number `json:"real_roomid"`
	MsgCommon  string          `json:"msg_common"`
	MsgSelf    string          `json:"msg_self"`
	LinkURL    string          `json:"link_url"`
	MsgType    jsoniter.Number `json:"msg_type"`
	ShieldUID  jsoniter.Number `json:"shield_uid"`
	BusinessID string          `json:"business_id"`
	Scatter    struct {
		Min jsoniter.Number `json:"min"`
		Max jsoniter.Number `json:"max"`
	} `json:"scatter"`
	MarqueeID  string          `json:"marquee_id"`
	NoticeType jsoniter.Number `json:"notice_type"`
}

type NoticeMsgStruct struct {
	RawData noticeMsg

	Name       string
	RoomId     string
	RealRoomId string
	MsgCommon  string
	MsgSelf    string
	LinkUrl    string
	MsgType    string
}

func NoticeMsg(b []byte) (NoticeMsgStruct, error) {
	var nm noticeMsg
	var nms NoticeMsgStruct
	var err error

	decoder := jsoniter.NewDecoder(bytes.NewReader(b))
	if err = decoder.Decode(&nm); err != nil {
		goto Label
	}

	nms.RawData = nm

	nms.Name = nm.Name
	nms.RoomId = nm.Roomid.String()
	nms.RealRoomId = nm.RealRoomid.String()
	nms.MsgCommon = nm.MsgCommon
	nms.MsgSelf = nm.MsgSelf
	nms.LinkUrl = nm.LinkURL
	nms.MsgType = nm.MsgType.String()

Label:
	return nms, err
}
