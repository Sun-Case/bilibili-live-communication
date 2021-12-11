package biliJsonConv

import (
	"bytes"
	jsoniter "github.com/json-iterator/go"
)

type live struct {
	Cmd             string          `json:"cmd"`
	LiveKey         string          `json:"live_key"`
	VoiceBackground string          `json:"voice_background"`
	SubSessionKey   string          `json:"sub_session_key"`
	LivePlatform    string          `json:"live_platform"`
	LiveModel       jsoniter.Number `json:"live_model"`
	Roomid          jsoniter.Number `json:"roomid"`
}

type LiveStruct struct {
	RawData live

	LiveKey, VoiceBackground, SubSessionKey, LivePlatform, LiveMode, RoomId string
}

func Live(b []byte) (LiveStruct, error) {
	var l live
	var ls LiveStruct
	var err error

	decoder := jsoniter.NewDecoder(bytes.NewReader(b))
	if err = decoder.Decode(&l); err != nil {
		goto Label
	}

	ls.RawData = l

	ls.LiveKey = l.LiveKey
	ls.VoiceBackground = l.VoiceBackground
	ls.SubSessionKey = l.SubSessionKey
	ls.LivePlatform = l.LivePlatform
	ls.LiveMode = l.LiveModel.String()
	ls.RoomId = l.Roomid.String()

Label:
	return ls, err
}
