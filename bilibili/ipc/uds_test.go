package ipc

import (
	"github.com/gorilla/websocket"
	"log"
	"net"
	"os"
	"testing"
	"time"
)

func TestUdsListen(t *testing.T) {
	_ = os.Remove("/tmp/golang.sock")
	go UdsListen("/tmp/golang.sock")
	time.Sleep(time.Second)

	dial, err := net.Dial("unix", "/tmp/golang.sock")
	if err != nil {
		log.Panicln(err)
	}
	conn := NetConnEncapsulate(dial)

	err = conn.WriteMessage(websocket.BinaryMessage, []byte("{\"room_id\": \"1\"}"))
	if err != nil {
		log.Panicln(err)
	}
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Panicln(err)
		}
		log.Println(string(message))
	}
}
