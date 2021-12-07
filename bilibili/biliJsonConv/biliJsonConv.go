package biliJsonConv

import (
	"bytes"
	json_ "encoding/json"
	"errors"
	"fmt"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type CmdType int

const (
	CodeType CmdType = iota
	DanMuMsgType
	InteractWordType

	UnknownType
	ErrorType
)

type BaseKey struct {
	Code int    `json:"code"`
	Cmd  string `json:"cmd"`
}

func GetCmdType(b []byte) CmdType {
	var bk BaseKey
	if err := json.Unmarshal(b, &bk); err != nil {
		return ErrorType
	}

	if bk.Cmd == "" {
		return CodeType
	}

	switch bk.Cmd {
	case "DANMU_MSG":
		return DanMuMsgType
	case "INTERACT_WORD":
		return InteractWordType
	default:
		return UnknownType
	}
}

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

func InteractWord(b []byte) (InteractWordStruct, error) {
	var iw interactWord
	var iwS InteractWordStruct
	var globalErr error

	decoder := json.NewDecoder(bytes.NewReader(b))
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
			case json_.Number:
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

// danMuMsg
// 弹幕
type danMuMsg struct {
	Cmd  string        `json:"cmd"`
	Info []interface{} `json:"info"`
}

type DanMuMsgStruct struct {
	RawData danMuMsg

	FontSize string // 字号
	UnixMs   string // 时间戳: 毫秒
	UnixS    string // 时间戳: 秒
	Color    string // 弹幕颜色
	Msg      string // 弹幕内容
	Name     string // 用户名
	UId      string // 用户ID
	CT       string
}

func DanMuMsg(b []byte) (DanMuMsgStruct, error) {
	var dmm danMuMsg
	var dmmt DanMuMsgStruct
	var globalErr error

	decoder := json.NewDecoder(bytes.NewReader(b))
	decoder.UseNumber()

	if globalErr = decoder.Decode(&dmm); globalErr != nil {
		globalErr = errors.New(fmt.Sprintf("DanMuMsg json decode fail: %s", globalErr.Error()))
		goto Label
	}

	if v1, ok := dmm.Info[0].([]interface{}); ok {
		// FontSize
		if jn, ok := v1[2].(json_.Number); ok {
			dmmt.FontSize = jn.String()
		} else {
			globalErr = errors.New("DanMuMsg get FontSize fail")
			goto Label
		}

		// Color
		if jn, ok := v1[3].(json_.Number); ok {
			dmmt.Color = jn.String()
		} else {
			globalErr = errors.New("DanMuMsg get Color fail")
			goto Label
		}

		// UnixMs
		if jn, ok := v1[4].(json_.Number); ok {
			dmmt.UnixMs = jn.String()
		} else {
			globalErr = errors.New("DanMuMsg get UnixMs fail")
			goto Label
		}

		// UnixS
		// 不知道什么原因，有些数据不对，有些数据是负数，所以不在这里获取
		//if jn, ok := v1[5].(json_.Number); ok {
		//	dmmt.UnixS = jn.String()
		//} else {
		//	globalErr = errors.New("DanMuMsg get UnixS fail")
		//	goto Label
		//}
	} else {
		globalErr = errors.New("DanMuMsg.Info[0] type assertion fail")
		goto Label
	}

	// Msg
	switch v1 := dmm.Info[1].(type) {
	case string:
		dmmt.Msg = v1
	default:
		globalErr = errors.New("DanMuMsg get Msg fail")
		goto Label
	}

	if v1, ok := dmm.Info[2].([]interface{}); ok {
		// UId
		if jn, ok := v1[0].(json_.Number); ok {
			dmmt.UId = jn.String()
		} else {
			globalErr = errors.New("DanMuMsg get UId fail")
			goto Label
		}

		// Name
		if dmmt.Name, ok = v1[1].(string); !ok {
			globalErr = errors.New("DanMuMsg get Name fail")
			goto Label
		}
	} else {
		globalErr = errors.New("DanMuMsg.Info[2] type assertion fail")
		goto Label
	}

	if v1, ok := dmm.Info[9].(map[string]interface{}); ok {
		// CT
		if ct, ok := v1["ct"]; ok {
			if dmmt.CT, ok = ct.(string); !ok {
				globalErr = errors.New("DanMuMsg CT type assertion fail")
				goto Label
			}
		} else {
			globalErr = errors.New("DanMuMsg get CT fail")
			goto Label
		}

		// UnixS
		if us, ok := v1["ts"]; ok {
			if jn, ok := us.(json_.Number); ok {
				dmmt.UnixS = jn.String()
			} else {
				globalErr = errors.New("DanMuMsg UnixS type assertion fail")
				goto Label
			}
		} else {
			globalErr = errors.New("DanMuMsg get TS fail")
			goto Label
		}
	} else {
		globalErr = errors.New("DanMuMsg.Info[9] type assertion fail")
		goto Label
	}

	dmmt.RawData = dmm

Label:
	return dmmt, globalErr
}
