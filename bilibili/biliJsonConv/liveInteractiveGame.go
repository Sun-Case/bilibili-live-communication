package biliJsonConv

import (
	"bytes"
	jsoniter "github.com/json-iterator/go"
)

type liveInteractiveGame struct {
	Cmd  string `json:"cmd"`
	Data struct {
		Type           jsoniter.Number `json:"type"`
		Uid            jsoniter.Number `json:"uid"`
		Uname          string          `json:"uname"`
		Uface          string          `json:"uface"`
		GiftId         jsoniter.Number `json:"gift_id"`
		GiftName       string          `json:"gift_name"`
		GiftNum        jsoniter.Number `json:"gift_num"`
		Price          jsoniter.Number `json:"price"`
		Paid           bool            `json:"paid"`
		Msg            string          `json:"msg"`
		FansMedalLevel jsoniter.Number `json:"fans_medal_level"`
		GuardLevel     jsoniter.Number `json:"guard_level"`
		Timestamp      jsoniter.Number `json:"timestamp"`
		AnchorLottery  interface{}     `json:"anchor_lottery"`
		PkInfo         interface{}     `json:"pk_info"`
		AnchorInfo     struct {
			Uid   jsoniter.Number `json:"uid"`
			Uname string          `json:"uname"`
			Uface string          `json:"uface"`
		} `json:"anchor_info"`
	} `json:"data"`
}
type LiveInteractiveGameStruct struct {
	RawData liveInteractiveGame

	Uid            string
	Uname          string
	GiftId         string
	GiftName       string
	GiftNum        string
	FansMedalLevel string
	Timestamp      string
	AnchorInfo     struct {
		Uid   string
		Uname string
	}
}

func LiveInteractiveGame(b []byte) (LiveInteractiveGameStruct, error) {
	var lig liveInteractiveGame
	var ligs LiveInteractiveGameStruct
	var err error

	decoder := jsoniter.NewDecoder(bytes.NewReader(b))
	if err = decoder.Decode(&lig); err != nil {
		goto Label
	}

	ligs.RawData = lig

	ligs.Uid = lig.Data.Uid.String()
	ligs.Uname = lig.Data.Uname
	ligs.GiftId = lig.Data.GiftId.String()
	ligs.GiftName = lig.Data.GiftName
	ligs.GiftNum = lig.Data.GiftNum.String()
	ligs.FansMedalLevel = lig.Data.FansMedalLevel.String()
	ligs.Timestamp = lig.Data.Timestamp.String()
	ligs.AnchorInfo.Uid = lig.Data.AnchorInfo.Uid.String()
	ligs.AnchorInfo.Uname = lig.Data.AnchorInfo.Uname

Label:
	return ligs, err
}
