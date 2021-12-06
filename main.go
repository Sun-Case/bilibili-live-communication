package main

import (
	"bilibili-live-communication/bilibili/biliAPI"
	"bilibili-live-communication/bilibili/biliDataConv"
	"bilibili-live-communication/bilibili/biliJsonConv"
	"bilibili-live-communication/bilibili/biliWebsocket"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
)

func main() {
	rid := ""

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

	//time.Sleep(time.Second * 5)

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

	receiveChannel := make(chan []byte, 1)
	blw := biliWebsocket.NewBiliLiveWebsocket(c, &receiveChannel, func(liveWebsocket *biliWebsocket.BiliLiveWebsocket) []byte {
		var convStruct biliDataConv.BiliConvStruct

		// 注意，如果获取弹幕地址及Token时，没有携带账号的Cookie，则下面的uid必须为0，否则鉴权失败，直接断开连接
		data := fmt.Sprintf("{\"uid\":0,\"roomid\":%s,\"protover\":3,\"platform\":\"web\",\"type\":2,\"key\":\"%s\"}", rid, token)
		log.Println(data)
		if err := convStruct.Encode(0x01, 0x00000007, 0x00000001, []byte(data)); err != nil {
			log.Panicln(err)
		}
		return convStruct.Bin
	}, func(liveWebsocket *biliWebsocket.BiliLiveWebsocket) []byte {
		var convStruct biliDataConv.BiliConvStruct

		if err := convStruct.Encode(0x01, 0x00000002, 0x00000001, []byte("[object Object]")); err != nil {
			log.Panicln(err)
		}
		return convStruct.Bin
	})

	for {
		select {
		case b, ok := <-receiveChannel:
			if !ok {
				blw.Stop()
				return
			}
			if slice, err := biliDataConv.Decode(b); err == nil {
				for _, v := range slice {
					switch biliJsonConv.GetCmdType(v.Body) {
					case biliJsonConv.DanMuMsgType:
						log.Println(biliJsonConv.DanMuMsg(v.Body))
					case biliJsonConv.ErrorType:
						log.Println("Error Type:", string(v.Body))
					case biliJsonConv.UnknownType:
						fallthrough
					default:
						//log.Println(string(v.Body))
					}
				}
			} else {
				log.Println(err)
			}
		}
	}
}
