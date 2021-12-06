package biliAPI

import (
	"bytes"
	"errors"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"io"
	"log"
	"net/http"
	"strconv"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// DanMuReport
// 举报直播弹幕
func DanMuReport(id, roomId, tUid, msg, reason, ts, sign, reasonId, token, dmType, csrf, visitId string) error {
	u := "https://api.live.bilibili.com/xlive/web-ucenter/v1/dMReport/Report"

	var data = struct {
		Id        string `json:"id"`         // 未知，为0
		RoomId    string `json:"roomid"`     // 直播间真实ID
		TUId      string `json:"tuid"`       // 发送者UID
		Msg       string `json:"msg"`        // 内容
		Reason    string `json:"reason"`     // 举报理由(辱骂引战/...)就是点举报后，要选择的理由
		Ts        string `json:"ts"`         // 秒级时间戳，注意是和 ct 同一个 Object 的那个 ts
		Sign      string `json:"sign"`       // 弹幕标志，即 ct
		ReasonId  string `json:"reason_id"`  // 似乎是理由对应的序号，选择举报理由的时候，有个下拉框，第一个理由对应1
		Token     string `json:"token"`      // 未知，为空
		DmType    string `json:"dm_type"`    // 字面上理解为弹幕类型，似乎与 HTML的 data-type 属性一致
		CsrfToken string `json:"csrf_token"` // Cookies 中的 bili_jct
		Csrf      string `json:"csrf"`       // 同上
		VisitId   string `json:"visit_id"`   // 未知，为空
	}{Id: id, RoomId: roomId, TUId: tUid, Msg: msg, Reason: reason, Ts: ts, Sign: sign, ReasonId: reasonId, Token: token, DmType: dmType, CsrfToken: csrf, Csrf: csrf, VisitId: visitId}

	buf, err := jsoniter.Marshal(&data)
	if err != nil {
		return err
	}

	cli := http.Client{}
	req, err := http.NewRequest("POST", u, bytes.NewReader(buf))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json;charset=utf-8")

	_, err = cli.Do(req)
	if err != nil {
		return err
	}

	return nil
}

// GetRoomInfo
// 获取直播真实ID
func GetRoomInfo(id string) (string, error) {
	u := fmt.Sprintf("https://api.live.bilibili.com/room/v1/Room/room_init?id=%s", id)

	resp, err := http.Get(u)
	if err != nil {
		return "", err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(resp.Body)

	var jsonData struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			RoomId      uint64 `json:"room_id"`
			ShortId     uint64 `json:"short_id"`
			UId         uint64 `json:"uid"`
			NeedP2p     int    `json:"need_p2p"`
			IsHidden    bool   `json:"is_hidden"`
			IsLocked    bool   `json:"is_locked"`
			IsPortrait  bool   `json:"is_portrait"`
			LiveStatus  int    `json:"live_status"`
			HiddenTill  int    `json:"hidden_till"`
			LockTill    int    `json:"lock_till"`
			Encrypted   bool   `json:"encrypted"`
			PwdVerified bool   `json:"pwd_verified"`
			LiveTime    uint64 `json:"live_time"`
			RoomShield  int    `json:"room_shield"`
			IsSp        int    `json:"is_sp"`
			SpecialType int    `json:"special_type"`
		} `json:"data"`
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(b, &jsonData)
	if err != nil {
		return "", err
	}

	return strconv.FormatUint(jsonData.Data.RoomId, 10), nil
}

// GetDanmuInfo
// 获取弹幕服务器地址及Token
func GetDanmuInfo(id string) ([]string, string, error) {
	const u string = "https://api.live.bilibili.com/xlive/web-room/v1/index/getDanmuInfo?id=%s&type=0"
	fU := fmt.Sprintf(u, id)

	if resp, err := http.Get(fU); err != nil {
		return nil, "", err
	} else {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Println(err)
			}
		}(resp.Body)

		var jsonData struct {
			Code int64 `json:"code"`
			Data struct {
				BusinessId int64  `json:"business_id"`
				Group      string `json:"group"`
				HostList   []struct {
					Host    string `json:"host"`
					Port    uint32 `json:"port"`
					WsPort  uint32 `json:"ws_port"`
					WssPort uint32 `json:"wss_port"`
				} `json:"host_list"`
				MaxDelay         int64   `json:"max_delay"`
				RefreshRate      int64   `json:"refresh_rate"`
				RefreshRowFactor float32 `json:"refresh_row_factor"`
				Token            string  `json:"token"`
			} `json:"data"`
			Message string `json:"message"`
			Ttl     int64  `json:"ttl"`
		}

		if b, err := io.ReadAll(resp.Body); err != nil {
			return nil, "", err
		} else {
			err := json.Unmarshal(b, &jsonData)
			if err != nil {
				return nil, "", err
			}

			if jsonData.Code != 0 {
				return nil, "", errors.New(fmt.Sprintf("Code fail: %d", jsonData.Code))
			}

			hosts := make([]string, 0)
			for _, v := range jsonData.Data.HostList {
				hosts = append(hosts, v.Host)
			}

			return hosts, jsonData.Data.Token, nil
		}
	}
}
