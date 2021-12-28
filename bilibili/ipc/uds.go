package ipc

import "net"

/*!
Unix Domain Socket
*/

func UdsListen(address string) error {
	listener, err := net.Listen("unix", address)
	if err != nil {
		return err
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		go new(EncapInterface).IPC(conn)
	}
}
