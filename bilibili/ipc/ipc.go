package ipc

import (
	"bilibili-live-communication/bilibili/biliBinConv"
	"bilibili-live-communication/bilibili/biliWebsocket"
	"bilibili-live-communication/bilibili/connEncapsulate"
	"github.com/gorilla/websocket"
	jsoniter "github.com/json-iterator/go"
	"log"
)

func IPC(rwc connEncapsulate.ReadWriteCloser) {
	defer func(rwc connEncapsulate.ReadWriteCloser) {
		err := rwc.Close()
		if err != nil {
			log.Println(err)
		}
	}(rwc)

	// 获取roomid
	_, roomIdJson, err := rwc.ReadMessage()
	if err != nil {
		log.Println(err)
		return
	}

	roomId := jsoniter.Get(roomIdJson, "room_id").ToString()

	// 连接弹幕服务器
	wsConn, err := biliWebsocket.NewBiliWebsocket(roomId)
	if err != nil {
		log.Println(err)
		return
	}
	defer func(wsConn connEncapsulate.ReadWriteCloser) {
		err := wsConn.Close()
		if err != nil {
			log.Println(err)
		}
	}(wsConn)

	for {
		_, message, err := wsConn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		decodeMessage, err := biliBinConv.Decode(message)
		if err != nil {
			log.Println(err)
			return
		}
		for index := range decodeMessage {
			err := rwc.WriteMessage(websocket.TextMessage, decodeMessage[index].Body)
			if err != nil {
				log.Println(err)
				return
			}
		}
	}
}
