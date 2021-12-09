package biliJsonConv

import jsoniter "github.com/json-iterator/go"

type preparing struct {
	Cmd    string `json:"cmd"`
	RoomId string `json:"roomid"`
}
type PreparingStruct struct {
	RawData preparing

	RoomId string
}

func Preparing(b []byte) (PreparingStruct, error) {
	var p preparing
	var ps PreparingStruct
	var err error

	if err = jsoniter.Unmarshal(b, &p); err != nil {
		goto Label
	}

	ps.RawData = p
	ps.RoomId = p.RoomId

Label:
	return ps, err
}
