package biliJsonConv

import (
	"encoding/json"
	jsoniter "github.com/json-iterator/go"
)

type stopLiveRoomList struct {
	Cmd  string `json:"cmd"`
	Data struct {
		RoomIdList []json.Number `json:"room_id_list"`
	} `json:"data"`
}

type StopLiveRoomListStruct struct {
	RawData stopLiveRoomList

	RoomIdList []string
}

func StopLiveRoomList(b []byte) (StopLiveRoomListStruct, error) {
	var slrl stopLiveRoomList
	var slrls StopLiveRoomListStruct
	var err error

	if err = jsoniter.Unmarshal(b, &slrl); err != nil {
		goto Label
	}

	slrls.RawData = slrl
	slrls.RoomIdList = make([]string, len(slrl.Data.RoomIdList))

	for i := range slrl.Data.RoomIdList {
		slrls.RoomIdList[i] = slrl.Data.RoomIdList[i].String()
	}

Label:
	return slrls, err
}
