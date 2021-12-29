package connEncapsulate

import (
	"net"
	"time"
)

type Conn interface {
	//Read([]byte) (int, error)
	//Write([]byte) (int, error)

	ReadMessage() (int, []byte, error)
	WriteMessage(int, []byte) error
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	Close() error
	SetReadDeadline(time.Time) error
	SetDeadline(time.Time) error
	SetWriteDeadline(time.Time) error
}
