package ipc

import (
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

		log.Println(string(message))
	}
}
