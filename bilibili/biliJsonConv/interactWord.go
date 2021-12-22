package biliJsonConv

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	jsoniter "github.com/json-iterator/go"
)

type interactWord struct {
	Cmd  string                 `json:"cmd"`
	Data map[string]interface{} `json:"data"`
}

type InteractWordStruct struct {
	RawData interactWord

	RoomId      string
	Score       string
	Timestamp   string
	TriggerTime string
	UId         string
	UName       string
	UNameColor  string
}

// InteractWord
// 进入直播间
func InteractWord(b []byte) (InteractWordStruct, error) {
	var iw interactWord
	var iwS InteractWordStruct
	var globalErr error

	decoder := jsoniter.NewDecoder(bytes.NewReader(b))
	decoder.UseNumber()

	if globalErr = decoder.Decode(&iw); globalErr != nil {
		globalErr = errors.New(fmt.Sprintf("InteractWord json decode fail: %s", globalErr.Error()))
		goto Label
	}

	for _, k := range []struct {
		Key   string
		Value *string
	}{{Key: "roomid", Value: &iwS.RoomId},
		{Key: "score", Value: &iwS.Score},
		{Key: "timestamp", Value: &iwS.Timestamp},
		{Key: "trigger_time", Value: &iwS.TriggerTime},
		{Key: "uid", Value: &iwS.UId},
		{Key: "uname", Value: &iwS.UName},
		{Key: "uname_color", Value: &iwS.UNameColor}} {
		if s, ok := iw.Data[k.Key]; !ok {
			globalErr = errors.New(fmt.Sprintf("InteractWord not found: %s", k.Key))
			goto Label
		} else {
			switch ta := s.(type) {
			case json.Number:
				*k.Value = ta.String()
			case string:
				*k.Value = ta
			default:
				globalErr = errors.New(fmt.Sprintf("InteractWord type assertion fail: %s", k.Key))
				goto Label
			}
		}
	}

Label:
	return iwS, globalErr
}
