package ipc

import "net"

func TcpListener(address string) error {
	listener, err := net.Listen("tcp", address)
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
