package biliJsonConv

import (
	"bytes"
	jsoniter "github.com/json-iterator/go"
)

type guardBuy struct {
	Cmd  string `json:"cmd"`
	Data struct {
		UID        jsoniter.Number `json:"uid"`
		Username   string          `json:"username"`
		GuardLevel jsoniter.Number `json:"guard_level"`
		Num        jsoniter.Number `json:"num"`
		Price      jsoniter.Number `json:"price"`
		GiftID     jsoniter.Number `json:"gift_id"`
		GiftName   string          `json:"gift_name"`
		StartTime  jsoniter.Number `json:"start_time"`
		EndTime    jsoniter.Number `json:"end_time"`
	} `json:"data"`
}

type GuardBuyStruct struct {
	RawData guardBuy

	Uid        string `json:"UID"`
	Username   string `json:"Username"`
	GuardLevel string `json:"GuardLevel"`
	Num        string `json:"Num"`
	Price      string `json:"Price"`
	GiftId     string `json:"GiftID"`
	GiftName   string `json:"GiftName"`
	StartTime  string `json:"StartTime"`
	EndTime    string `json:"EndTime"`
}

func GuardBuy(b []byte) (GuardBuyStruct, error) {
	var gb guardBuy
	var gbs GuardBuyStruct
	var err error

	decoder := jsoniter.NewDecoder(bytes.NewReader(b))
	decoder.UseNumber()
	if err = decoder.Decode(&gb); err != nil {
		goto Label
	}

	gbs.RawData = gb
	gbs.Uid = gb.Data.UID.String()
	gbs.Username = gb.Data.Username
	gbs.GuardLevel = gb.Data.GuardLevel.String()
	gbs.Num = gb.Data.Num.String()
	gbs.Price = gb.Data.Price.String()
	gbs.GiftId = gb.Data.GiftID.String()
	gbs.GiftName = gb.Data.GiftName
	gbs.StartTime = gb.Data.StartTime.String()
	gbs.EndTime = gb.Data.EndTime.String()

Label:
	return gbs, err
}
