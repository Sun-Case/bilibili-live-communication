package biliJsonConv

import (
	"encoding/json"
	jsoniter "github.com/json-iterator/go"
)

type roomRealTimeMessageUpdate struct {
	Cmd  string `json:"cmd"`
	Data struct {
		RoomId    json.Number `json:"roomid"`
		Fans      json.Number `json:"fans"`
		RedNotice json.Number `json:"red_notice"`
		FansClub  json.Number `json:"fans_club"`
	} `json:"data"`
}
type RoomRealTimeMessageUpdateStruct struct {
	RawData roomRealTimeMessageUpdate

	RoomId    string
	Fans      string
	RedNotice string
	FansClub  string
}

func RoomRealTimeMessageUpdate(b []byte) (RoomRealTimeMessageUpdateStruct, error) {
	var rrtmu roomRealTimeMessageUpdate
	var rrtmus RoomRealTimeMessageUpdateStruct
	var err error

	if err = jsoniter.Unmarshal(b, &rrtmu); err != nil {
		goto Label
	}

	rrtmus.RawData = rrtmu
	rrtmus.RoomId = rrtmu.Data.RoomId.String()
	rrtmus.Fans = rrtmu.Data.Fans.String()
	rrtmus.RedNotice = rrtmu.Data.RedNotice.String()
	rrtmus.FansClub = rrtmu.Data.FansClub.String()

Label:
	return rrtmus, err
}
