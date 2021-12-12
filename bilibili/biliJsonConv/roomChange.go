package biliJsonConv

import (
	"bytes"
	jsoniter "github.com/json-iterator/go"
)

type roomChange struct {
	Cmd  string `json:"cmd"`
	Data struct {
		Title          string          `json:"title"`
		AreaID         jsoniter.Number `json:"area_id"`
		ParentAreaID   jsoniter.Number `json:"parent_area_id"`
		AreaName       string          `json:"area_name"`
		ParentAreaName string          `json:"parent_area_name"`
		LiveKey        string          `json:"live_key"`
		SubSessionKey  string          `json:"sub_session_key"`
	} `json:"data"`
}

type RoomChangeStruct struct {
	RawData roomChange

	Title, AreaId, ParentAreaId, AreaName, ParentAreaName, LiveKey, SubSessionKey string
}

func RoomChange(b []byte) (RoomChangeStruct, error) {
	var rc roomChange
	var rcs RoomChangeStruct
	var err error

	decoder := jsoniter.NewDecoder(bytes.NewReader(b))
	if err = decoder.Decode(&rc); err != nil {
		goto Label
	}

	rcs.RawData = rc
	rcs.Title = rc.Data.Title
	rcs.AreaId = rc.Data.AreaID.String()
	rcs.ParentAreaId = rc.Data.ParentAreaID.String()
	rcs.AreaName = rc.Data.AreaName
	rcs.ParentAreaName = rc.Data.ParentAreaName
	rcs.LiveKey = rc.Data.LiveKey
	rcs.SubSessionKey = rc.Data.SubSessionKey

Label:
	return rcs, err
}
