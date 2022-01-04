package ipc

import (
	"bilibili-live-communication/bilibili/biliBinConv"
	"github.com/gorilla/websocket"
	"log"
	"net"
	"testing"
	"time"
)

func TestTcpListen(t *testing.T) {
	go TcpListen("0.0.0.0:22333")

	time.Sleep(time.Second)

	dial, err := net.Dial("tcp", "127.0.0.1:22333")
	if err != nil {
		log.Println(err)
		return
	}
	conn := NetConnEncapsulate(dial)
	// 发送房间ID
	err = conn.WriteMessage(websocket.BinaryMessage, []byte("{\"room_id\": \"1\"}"))
	if err != nil {
		log.Panicln(err)
		return
	}
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		d, err := biliBinConv.Decode(message)
		if err != nil {
			log.Println(err)
			return
		}

		for _, l := range d {
			log.Println(string(l.Body))
		}
	}
}
