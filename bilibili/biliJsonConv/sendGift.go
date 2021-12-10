package biliJsonConv

import (
	"bytes"
	jsoniter "github.com/json-iterator/go"
)

type sendGift struct {
	Cmd  string `json:"cmd"`
	Data struct {
		Action            string          `json:"action"`
		BatchComboId      string          `json:"batch_combo_id"`
		BatchComboSend    interface{}     `json:"batch_combo_send"`
		BeatId            string          `json:"beatId"`
		BizSource         string          `json:"biz_source"`
		BlindGift         interface{}     `json:"blind_gift"`
		BroadcastId       jsoniter.Number `json:"broadcast_id"`
		CoinType          string          `json:"coin_type"`
		ComboResourceId   jsoniter.Number `json:"combo_resource_id"`
		ComboSend         interface{}     `json:"combo_send"`
		ComboStayTime     jsoniter.Number `json:"combo_stay_time"`
		ComboTotalCoin    jsoniter.Number `json:"combo_total_coin"`
		CritProb          jsoniter.Number `json:"crit_prob"`
		Demarcation       jsoniter.Number `json:"demarcation"`
		DiscountPrice     jsoniter.Number `json:"discount_price"`
		DmScore           jsoniter.Number `json:"dmscore"`
		Draw              jsoniter.Number `json:"draw"`
		Effect            jsoniter.Number `json:"effect"`
		EffectBlock       jsoniter.Number `json:"effect_block"`
		Face              string          `json:"face"`
		FloatScResourceId jsoniter.Number `json:"float_sc_resource_id"`
		GiftId            jsoniter.Number `json:"giftId"`
		GiftName          string          `json:"giftName"`
		GiftType          jsoniter.Number `json:"giftType"`
		Gold              jsoniter.Number `json:"gold"`
		GuardLevel        jsoniter.Number `json:"guard_level"`
		IsFirst           bool            `json:"is_first"`
		IsSpecialBatch    jsoniter.Number `json:"is_special_batch"`
		Magnification     jsoniter.Number `json:"magnification"`
		MedalInfo         struct {
			AnchorRoomId     jsoniter.Number `json:"anchor_roomid"`
			AnchorUname      string          `json:"anchor_uname"`
			GuardLevel       jsoniter.Number `json:"guard_level"`
			IconId           jsoniter.Number `json:"icon_id"`
			IsLighted        jsoniter.Number `json:"is_lighted"`
			MedalColor       jsoniter.Number `json:"medal_color"`
			MedalColorBorder jsoniter.Number `json:"medal_color_border"`
			MedalColorEnd    jsoniter.Number `json:"medal_color_end"`
			MedalColorStart  jsoniter.Number `json:"medal_color_start"`
			MedalLevel       jsoniter.Number `json:"medal_level"`
			MedalName        string          `json:"medal_name"`
			Special          string          `json:"special"`
			TargetId         jsoniter.Number `json:"target_id"`
		} `json:"medal_info"`
		NameColor         string          `json:"name_color"`
		Num               jsoniter.Number `json:"num"`
		OriginalGiftName  string          `json:"original_gift_name"`
		Price             jsoniter.Number `json:"price"`
		RCost             jsoniter.Number `json:"rcost"`
		Remain            jsoniter.Number `json:"remain"`
		Rnd               string          `json:"rnd"`
		SendMaster        interface{}     `json:"send_master"`
		Silver            jsoniter.Number `json:"silver"`
		Super             jsoniter.Number `json:"super"`
		SuperBatchGiftNum jsoniter.Number `json:"super_batch_gift_num"`
		SuperGiftNum      jsoniter.Number `json:"super_gift_num"`
		SvgaBlock         jsoniter.Number `json:"svga_block"`
		TagImage          string          `json:"tag_image"`
		Tid               string          `json:"tid"`
		Timestamp         jsoniter.Number `json:"timestamp"`
		TopList           interface{}     `json:"top_list"`
		TotalCoin         jsoniter.Number `json:"total_coin"`
		Uid               jsoniter.Number `json:"uid"`
		Uname             string          `json:"uname"`
	} `json:"data"`
}
type SendGiftStruct struct {
	RawData sendGift

	Timestamp  string
	Uname      string
	Uid        string
	Tid        string
	Num        string
	MedalName  string
	MedalLevel string
	IsLighted  string
	GiftName   string
	GiftId     string
	GiftType   string
	Action     string
}

func SendGift(b []byte) (SendGiftStruct, error) {
	var sg sendGift
	var sgs SendGiftStruct
	var err error

	decoder := jsoniter.NewDecoder(bytes.NewReader(b))
	if err = decoder.Decode(&sg); err != nil {
		goto Label
	}

	sgs.RawData = sg

	sgs.Timestamp = sg.Data.Timestamp.String()
	sgs.Uname = sg.Data.Uname
	sgs.Uid = sg.Data.Uid.String()
	sgs.Tid = sg.Data.Tid
	sgs.Num = sg.Data.SuperGiftNum.String()
	sgs.MedalName = sg.Data.MedalInfo.MedalName
	sgs.MedalLevel = sg.Data.MedalInfo.MedalLevel.String()
	sgs.IsLighted = sg.Data.MedalInfo.IsLighted.String()
	sgs.GiftName = sg.Data.GiftName
	sgs.GiftId = sg.Data.GiftId.String()
	sgs.GiftType = sg.Data.GiftType.String()
	sgs.Action = sg.Data.Action

Label:
	return sgs, err
}
