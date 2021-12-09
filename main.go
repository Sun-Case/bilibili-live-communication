package main

import (
	"bilibili-live-communication/bilibili/biliAPI"
	"bilibili-live-communication/bilibili/biliDataConv"
	"bilibili-live-communication/bilibili/biliJsonConv"
	"bilibili-live-communication/bilibili/biliWebsocket"
	"fmt"
	"github.com/gorilla/websocket"
	jsoniter "github.com/json-iterator/go"
	"log"
	"os"
	"strconv"
	"time"
)

var KV map[string]string

func main() {
	fb, _ := os.ReadFile("config.json")

	_ = jsoniter.Unmarshal(fb, &KV)

	rid := KV["room_id"]

	// 获取真正直播间ID
	rid, err := biliAPI.GetRoomInfo(rid)
	if err != nil {
		log.Panicln(err)
	}

	// 获取弹幕服务器地址及token
	hosts, token, err := biliAPI.GetDanmuInfo(rid)
	if err != nil {
		log.Panicln(err)
	}

	wsHost := fmt.Sprintf("wss://%s/sub", hosts[0])
	log.Println(wsHost)
	c, _, err := websocket.DefaultDialer.Dial(wsHost, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer func(c *websocket.Conn) {
		err := c.Close()
		if err != nil {
			log.Println(err)
		}
	}(c)

	sendCannel := make(chan []byte)
	receiveChannel := make(chan []byte)
	blw := biliWebsocket.New(c, &receiveChannel, &sendCannel, func(liveWebsocket *biliWebsocket.BiliWebsocket) []byte {
		var b []byte
		var err error

		// 注意，如果获取弹幕地址及Token时，没有携带账号的Cookie，则下面的uid必须为0，否则鉴权失败，直接断开连接
		data := fmt.Sprintf("{\"uid\":0,\"roomid\":%s,\"protover\":3,\"platform\":\"web\",\"type\":2,\"key\":\"%s\"}", rid, token)
		log.Println(data)
		if b, err = biliDataConv.Encode(0x01, 0x00000007, 0x00000001, []byte(data)); err != nil {
			log.Panicln(err)
		}
		return b
	})

	for {
		select {
		case <-blw.IsClosedCh:
			return
		case b, ok := <-receiveChannel:
			if !ok {
				return
			}
			switch biliJsonConv.GetCmdType(b) {
			case biliJsonConv.DanMuMsgType:
				if b, err := biliJsonConv.DanMuMsg(b); err == nil {
					log.Printf("%s: %s", b.Name, b.Msg)
				} else {
					log.Println(err)
				}
			case biliJsonConv.InteractWordType:
				if b, err := biliJsonConv.InteractWord(b); err == nil {
					log.Println(b.UName, "进入直播间")
				} else {
					log.Println(err)
				}
			case biliJsonConv.OnlineRankCountType:
				if b, err := biliJsonConv.OnlineRankCount(b); err == nil {
					log.Println("在线排名统计:", b.Count)
				} else {
					log.Println(err)
				}
			case biliJsonConv.StopLiveRoomListType:
				if b, err := biliJsonConv.StopLiveRoomList(b); err == nil {
					log.Println("下播直播间ID:", b.RoomIdList)
				}
			case biliJsonConv.RoomRealTimeMessageUpdateType:
				if b, err := biliJsonConv.RoomRealTimeMessageUpdate(b); err == nil {
					log.Printf("直播间: %s, 粉丝数: %s, 粉丝团数: %s\n", b.RoomId, b.Fans, b.FansClub)
				} else {
					log.Println(err)
				}
			case biliJsonConv.ErrorType:
				log.Println("Error Type:", string(b))
			case biliJsonConv.UnknownType:
				fallthrough
			default:
				log.Println(string(b))
				save(b)
			}
		}
	}
}

func save(b []byte) {
	var cmd struct {
		Cmd string `json:"cmd"`
	}
	if err := jsoniter.Unmarshal(b, &cmd); err != nil {
		return
	} else {
		if err := os.MkdirAll(KV["save_path"]+cmd.Cmd, 0755); err != nil {
			return
		} else {
			_ = os.WriteFile(KV["save_path"]+cmd.Cmd+"/"+strconv.FormatUint(uint64(time.Now().UnixNano()), 10)+".json", b, 0644)
		}
	}
}
