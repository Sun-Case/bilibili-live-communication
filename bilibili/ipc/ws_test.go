package ipc

import "testing"

func TestWsListen(t *testing.T) {
	WsListen("/", "0.0.0.0:22334")
}
