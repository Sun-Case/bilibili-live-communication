package biliWebsocket

import (
	"bilibili-live-communication/bilibili/biliAPI"
	"bilibili-live-communication/bilibili/biliBinConv"
	"bilibili-live-communication/bilibili/connEncapsulate"
	"fmt"
	"log"
	"testing"
)

func TestBiliWebsocket_Init(t *testing.T) {
	roomId, _ := biliAPI.GetRoomInfo("1")
	wsList, token, _ := biliAPI.GetDanmuInfo(roomId)
	log.Println(wsList)
	log.Println(token)

	var bw BiliWebsocket

	log.Println(bw.Init(wsList[0], func(_ *BiliWebsocket) []byte {
		data := fmt.Sprintf("{\"uid\":0,\"roomid\":%s,\"protover\":3,\"platform\":\"web\",\"type\":2,\"key\":\"%s\"}", roomId, token)
		bin, _ := biliBinConv.Encode(0x0001, 0x00000007, 0x00000001, []byte(data))
		return bin
	}))
	log.Println(bw.Start())

	var interf connEncapsulate.ReadWriteCloser = &bw

	for {
		_, bin, err := interf.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		bcs, _ := biliBinConv.Decode(bin)
		for _, it := range bcs {
			log.Println(string(it.Body))
		}
	}
}
