package biliJsonConv

import (
	"bytes"
	jsoniter "github.com/json-iterator/go"
)

type comboSend struct {
	Cmd  string `json:"cmd"`
	Data struct {
		Action         string          `json:"action"`
		BatchComboID   string          `json:"batch_combo_id"`
		BatchComboNum  jsoniter.Number `json:"batch_combo_num"`
		ComboID        string          `json:"combo_id"`
		ComboNum       jsoniter.Number `json:"combo_num"`
		ComboTotalCoin jsoniter.Number `json:"combo_total_coin"`
		Dmscore        jsoniter.Number `json:"dmscore"`
		GiftID         jsoniter.Number `json:"gift_id"`
		GiftName       string          `json:"gift_name"`
		GiftNum        jsoniter.Number `json:"gift_num"`
		IsShow         jsoniter.Number `json:"is_show"`
		MedalInfo      struct {
			AnchorRoomid     jsoniter.Number `json:"anchor_roomid"`
			AnchorUname      string          `json:"anchor_uname"`
			GuardLevel       jsoniter.Number `json:"guard_level"`
			IconID           jsoniter.Number `json:"icon_id"`
			IsLighted        jsoniter.Number `json:"is_lighted"`
			MedalColor       jsoniter.Number `json:"medal_color"`
			MedalColorBorder jsoniter.Number `json:"medal_color_border"`
			MedalColorEnd    jsoniter.Number `json:"medal_color_end"`
			MedalColorStart  jsoniter.Number `json:"medal_color_start"`
			MedalLevel       jsoniter.Number `json:"medal_level"`
			MedalName        string          `json:"medal_name"`
			Special          string          `json:"special"`
			TargetID         jsoniter.Number `json:"target_id"`
		} `json:"medal_info"`
		NameColor  string          `json:"name_color"`
		RUname     string          `json:"r_uname"`
		Ruid       jsoniter.Number `json:"ruid"`
		SendMaster interface{}     `json:"send_master"`
		TotalNum   jsoniter.Number `json:"total_num"`
		UID        jsoniter.Number `json:"uid"`
		Uname      string          `json:"uname"`
	} `json:"data"`
}

type ComboSendStruct struct {
	RawData comboSend

	Action                    string
	BatchComboNum, ComboNum   string
	GiftID, GiftName, GiftNum string
	IsShow                    string
	RUname                    string
	Ruid                      string
	TotalNum                  string
	UID                       string
	Uname                     string
}

func ComboSend(b []byte) (ComboSendStruct, error) {
	var cs comboSend
	var css ComboSendStruct
	var err error

	decoder := jsoniter.NewDecoder(bytes.NewReader(b))
	if err = decoder.Decode(&cs); err != nil {
		goto Label
	}

	css.Action = cs.Data.Action
	css.BatchComboNum = cs.Data.BatchComboNum.String()
	css.ComboNum = cs.Data.ComboNum.String()
	css.GiftID = cs.Data.GiftID.String()
	css.GiftName = cs.Data.GiftName
	css.GiftNum = cs.Data.GiftNum.String()
	css.IsShow = cs.Data.IsShow.String()
	css.RUname = cs.Data.RUname
	css.Ruid = cs.Data.Ruid.String()
	css.TotalNum = cs.Data.TotalNum.String()
	css.Uname = cs.Data.Uname

Label:
	return css, err
}
