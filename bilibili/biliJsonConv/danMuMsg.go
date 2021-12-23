package biliJsonConv

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	jsoniter "github.com/json-iterator/go"
)

// danMuMsg
// 弹幕
type danMuMsg struct {
	Cmd  string        `json:"cmd"`
	Info []interface{} `json:"info"`
}

type DanMuMsgStruct struct {
	RawData danMuMsg

	FontSize   string // 字号
	UnixMs     string // 时间戳: 毫秒
	UnixS      string // 时间戳: 秒
	Color      string // 弹幕颜色
	Msg        string // 弹幕内容
	Name       string // 用户名
	UId        string // 用户ID
	CT         string
	Brand      string // 牌子
	BrandLevel string // 牌子等级
	ShowBrand  bool   // 显示牌子
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

	// 牌子
	if v1, ok := dmm.Info[3].([]interface{}); ok && len(v1) != 0 {
		// BrandLevel
		if bl, ok := v1[0].(json.Number); ok {
			dmmt.BrandLevel = bl.String()
		} else {
			globalErr = errors.New("DanMuMsg get BrandLevel fail")
			goto Label
		}
		// Brand
		if dmmt.Brand, ok = v1[1].(string); !ok {
			globalErr = errors.New("DanMuMsg get Brand fail")
			goto Label
		}
		// ShowBrand
		if sb, ok := v1[11].(json.Number); ok {
			if sb.String() == "1" {
				dmmt.ShowBrand = true
			} else {
				dmmt.ShowBrand = false
			}
		}
	} else if !ok {
		globalErr = errors.New("DanMuMsg.Info[3] type assertion fail")
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
