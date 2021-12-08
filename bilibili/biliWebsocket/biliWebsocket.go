package biliWebsocket

import (
	"bilibili-live-communication/bilibili/biliDataConv"
	"bilibili-live-communication/bilibili/websocketP"
	"github.com/gorilla/websocket"
	"log"
	"sync"
	"time"
)

const debug = true

type BiliWebsocket struct {
	WebsocketP *websocketP.WebSocketP

	// User 与 BiliWebsocket 通信
	sendCh     *chan []byte
	receiveCh  *chan []byte
	IsClosedCh chan struct{}

	// BiliWebsocket 与 websocketP 通信
	wspSendCh    chan []byte
	wspReceiveCh chan []byte

	timeoutTicker *time.Ticker

	wg   sync.WaitGroup
	once sync.Once
}

// receive
// 接收，并反序列化数据
func (blw *BiliWebsocket) receive() {
	defer blw.wg.Done()
	defer blw.Stop()

	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	for {
		select {
		case b, ok := <-blw.wspReceiveCh:
			if !ok {
				return
			}

			bcs, err := biliDataConv.Decode(b)
			if err != nil {
				log.Println(err)
			} else {
				for i := range bcs {
					select {
					case *blw.receiveCh <- bcs[i].Body:
					case <-blw.WebsocketP.IsClosedCh:
						return
					}
				}
			}

			blw.timeoutTicker.Reset(time.Second * 30)
		case <-blw.WebsocketP.IsClosedCh:
			return
		}
	}
}

// send
// 发送
func (blw *BiliWebsocket) send() {
	defer blw.wg.Done()
	defer blw.Stop()

	for {
		select {
		case blw.wspSendCh <- <-*blw.sendCh:
			blw.timeoutTicker.Reset(time.Second * 30)
		case <-blw.WebsocketP.IsClosedCh:
			return
		}
	}
}

// Stop
// 停止
func (blw *BiliWebsocket) Stop() {
	go blw.once.Do(blw.stop)
}
func (blw *BiliWebsocket) stop() {
	if debug {
		log.Println("biliWebsocket StopFunc: Running")
		defer func() {
			log.Println("biliWebsocket StopFunc: Exit")
		}()
	}

	blw.WebsocketP.Stop()
	if debug {
		log.Println("biliWebsocket StopFunc: WebsocketP stop success")
	}

	blw.wg.Wait()
	if debug {
		log.Println("biliWebsocket StopFunc: other goroutines exit")
	}

	close(blw.wspReceiveCh)
	close(blw.wspSendCh)
	close(blw.IsClosedCh)
	log.Println("biliWebsocket StopFunc: all channel close")
}

// heartbeat
// 心跳包
func (blw *BiliWebsocket) heartbeat() {
	defer blw.wg.Done()
	defer blw.Stop()

	if debug {
		log.Println("biliWebsocket HeartbeatFunc: Running")
		defer func() {
			log.Println("biliWebsocket HeartbeatFunc: Exit")
		}()
	}
	ticker := time.NewTicker(time.Second * 15)
	defer ticker.Stop()

	for {
		var b []byte
		var err error
		select {
		case <-ticker.C:
			log.Println("biliWebsocket HeartbeatFunc: heartbeat ticker")
			if b, err = biliDataConv.Encode(0x01, 0x00000002, 0x00000001, []byte("object Object")); err != nil {
				return
			}

			select {
			case blw.wspSendCh <- b:
			case <-blw.WebsocketP.IsClosedCh:
				return
			}
		case <-blw.WebsocketP.IsClosedCh:
			return
		}
	}
}

// timeout
// 超时
func (blw *BiliWebsocket) timeout() {
	defer blw.wg.Done()
	defer blw.Stop()

	if debug {
		log.Println("biliWebsocket TimeoutFunc: Running")
		defer func() {
			log.Println("biliWebsocket TimeoutFunc: Exit")
		}()
	}

	blw.timeoutTicker = time.NewTicker(time.Second * 30)
	defer blw.timeoutTicker.Stop()

	select {
	case <-blw.timeoutTicker.C:
	case <-blw.WebsocketP.IsClosedCh:
	}
}

func New(conn *websocket.Conn, receiveChannel, sendChannel *chan []byte, onOpen func(blw *BiliWebsocket) []byte) *BiliWebsocket {
	if debug {
		log.Println("biliWebSocket NewFunc: Running")
		defer func() {
			log.Println("biliWebSocket NewFunc: Exit")
		}()
	}

	blw := new(BiliWebsocket)

	blw.receiveCh = receiveChannel
	blw.sendCh = sendChannel
	blw.wspSendCh = make(chan []byte)
	blw.wspReceiveCh = make(chan []byte)
	blw.IsClosedCh = make(chan struct{})

	blw.WebsocketP = websocketP.New(conn, &blw.wspReceiveCh, &blw.wspSendCh)

	go func() {
		select {
		case blw.wspSendCh <- onOpen(blw):
			blw.wg.Add(2)
			// 心跳包
			go blw.heartbeat()
			// 超时
			go blw.timeout()
		case <-blw.WebsocketP.IsClosedCh:
		}
	}()

	blw.wg.Add(2)
	// 发送
	go blw.send()
	// 接收
	go blw.receive()

	return blw
}
