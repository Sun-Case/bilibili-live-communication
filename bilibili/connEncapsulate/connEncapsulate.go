package connEncapsulate

import (
	"net"
	"time"
)

type ReadWriteCloser interface {
	ReadMessage() (int, []byte, error)
	WriteMessage(int, []byte) error
	Close() error
}

type Conn interface {
	//Read([]byte) (int, error)
	//Write([]byte) (int, error)

	ReadWriteCloser
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	SetReadDeadline(time.Time) error
	SetDeadline(time.Time) error
	SetWriteDeadline(time.Time) error
}
