package ipc

import (
	"log"
	"net"
)

func TcpListen(address string) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Panicln(err)
	}
	for {
		conn, err := listener.Accept()
		log.Println("new conn:", conn.RemoteAddr())
		if err != nil {
			log.Panicln(err)
		}

		go IPC(NetConnEncapsulate(conn))
	}
}
