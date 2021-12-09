package biliJsonConv

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	jsoniter "github.com/json-iterator/go"
)

type CmdType int

const (
	CodeType CmdType = iota
	DanMuMsgType
	InteractWordType
	OnlineRankCountType
	StopLiveRoomListType
	RoomRealTimeMessageUpdateType
	OnlineRankV2Type

	UnknownType
	ErrorType
)

type BaseKey struct {
	Code int    `json:"code"`
	Cmd  string `json:"cmd"`
}

// GetCmdType
// 获取数据类型
func GetCmdType(b []byte) CmdType {
	var bk BaseKey
	if err := jsoniter.Unmarshal(b, &bk); err != nil {
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
	case "ONLINE_RANK_COUNT":
		return OnlineRankCountType
	case "STOP_LIVE_ROOM_LIST":
		return StopLiveRoomListType
	case "ROOM_REAL_TIME_MESSAGE_UPDATE":
		return RoomRealTimeMessageUpdateType
	case "ONLINE_RANK_V2":
		return OnlineRankV2Type
	default:
		return UnknownType
	}
}

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

	decoder := jsoniter.NewDecoder(bytes.NewReader(b))
	decoder.UseNumber()

	if globalErr = decoder.Decode(&dmm); globalErr != nil {
		globalErr = errors.New(fmt.Sprintf("DanMuMsg json decode fail: %s", globalErr.Error()))
		goto Label
	}

	if v1, ok := dmm.Info[0].([]interface{}); ok {
		// FontSize
		if jn, ok := v1[2].(json.Number); ok {
			dmmt.FontSize = jn.String()
		} else {
			globalErr = errors.New("DanMuMsg get FontSize fail")
			goto Label
		}

		// Color
		if jn, ok := v1[3].(json.Number); ok {
			dmmt.Color = jn.String()
		} else {
			globalErr = errors.New("DanMuMsg get Color fail")
			goto Label
		}

		// UnixMs
		if jn, ok := v1[4].(json.Number); ok {
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
		if jn, ok := v1[0].(json.Number); ok {
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
			if jn, ok := us.(json.Number); ok {
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
