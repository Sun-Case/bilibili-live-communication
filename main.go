package main

import "bilibili-live-communication/bilibili/ipc"

func main() {
	go ipc.TcpListen("0.0.0.0:22333")
	go ipc.UdsListen("/tmp/bili.sock")
	go ipc.WsListen("/", "0.0.0.0:22334")

	select {}
}
