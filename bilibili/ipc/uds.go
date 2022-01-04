package ipc

import (
	"log"
	"net"
)

func UdsListen(address string) {
	listener, err := net.Listen("unix", address)
	if err != nil {
		log.Panicln(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Panicln(err)
		}

		go IPC(NetConnEncapsulate(conn))
	}
}
