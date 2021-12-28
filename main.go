package main

import (
	"bilibili-live-communication/bilibili/ipc"
	"os"
)

func main() {
	go ipc.TcpListener("0.0.0.0:22333")
	os.Remove("/tmp/golang.sock")
	go ipc.UdsListen("/tmp/golang.sock")
	select {}
}
