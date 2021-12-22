package biliJsonConv

import (
	"encoding/json"
	jsoniter "github.com/json-iterator/go"
)

type onlineRankCount struct {
	Cmd  string `json:"cmd"`
	Data struct {
		Count json.Number `json:"count"`
	} `json:"data"`
}
type OnlineRankCountStruct struct {
	RawData onlineRankCount

	Count string
}

func OnlineRankCount(b []byte) (OnlineRankCountStruct, error) {
	var orc onlineRankCount
	var orcs OnlineRankCountStruct
	var err error
	if err = jsoniter.Unmarshal(b, &orc); err != nil {
		goto Label
	}

	orcs.RawData = orc
	orcs.Count = orc.Data.Count.String()

Label:
	return orcs, err
}
