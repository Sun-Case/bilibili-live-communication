package connEncapsulate

import (
	"github.com/gorilla/websocket"
	"net"
	"time"
)

// websocket.Conn
type wsConn struct {
	conn *websocket.Conn
}

func WsConn(conn *websocket.Conn) Conn {
	var wC = wsConn{conn: conn}
	return &wC
}

func (wc *wsConn) SetDeadline(t time.Time) error {
	if err := wc.SetReadDeadline(t); err != nil {
		return err
	}
	return wc.SetWriteDeadline(t)
}

func (wc *wsConn) SetWriteDeadline(t time.Time) error {
	return wc.conn.SetWriteDeadline(t)
}

func (wc *wsConn) SetReadDeadline(t time.Time) error {
	return wc.conn.SetReadDeadline(t)
}

func (wc *wsConn) Close() error {
	return wc.conn.Close()
}

func (wc *wsConn) RemoteAddr() net.Addr {
	return wc.conn.RemoteAddr()
}

func (wc *wsConn) LocalAddr() net.Addr {
	return wc.conn.LocalAddr()
}

func (wc *wsConn) WriteMessage(messageType int, data []byte) error {
	return wc.conn.WriteMessage(messageType, data)
}

func (wc *wsConn) ReadMessage() (int, []byte, error) {
	return wc.conn.ReadMessage()
}
