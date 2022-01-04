package biliWebsocket

import (
	"bilibili-live-communication/bilibili/biliBinConv"
	"bilibili-live-communication/bilibili/connEncapsulate"
	"log"
	"testing"
)

func TestBiliWebsocket_Init(t *testing.T) {
	var interf connEncapsulate.ReadWriteCloser
	interf, _ = NewBiliWebsocket("1")

	for {
		_, bin, err := interf.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		bcs, _ := biliBinConv.Decode(bin)
		for _, it := range bcs {
			log.Println(string(it.Body))
		}
	}
}
