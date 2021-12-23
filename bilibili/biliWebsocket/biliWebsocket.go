package biliWebsocket

import (
	"bilibili-live-communication/bilibili/biliDataConv"
	"bilibili-live-communication/bilibili/websocketP"
	"github.com/gorilla/websocket"
	"log"
	"sync"
	"time"
)

/*!

biliWebsocket 底层使用 websocketP
	负责维持心跳处理, 超时处理, 原始二进制解码, JSON编码成二进制

BiliWebsocket.sendCh: 接收外部的数据（结构体），将数据编码成二进制，传到 BiliWebsocket.wbpSendCh
	该通道由外部设置, BiliWebsocket.Stop 不负责关闭
BiliWebsocket.receiveCh: 将 JSON数据(由 BiliWebsocket.wspReceiveCh通道 获取并解码而来) 传输到该通道
	该通道由外部设置, BiliWebsocket.Stop 不负责关闭

BiliWebsocket.wspSendCh: 将 二进制数据(由 BiliWebsocket.wbpSendCh通道 获取并编码而来) 传到该通道, 交由 WebsocketP 处理
	该通道由 BiliWebsocket.New 创建, 由 BiliWebsocket.Stop 关闭
BiliWebsocket.wspReceiveCh: 读取该通道的原始二进制数据, 将原始二进制数据解码成 JSON 数据, 传到 BiliWebsocket.receiveCh
	该通道由 BiliWebsocket.New 创建, 由 BiliWebsocket.Stop 关闭

BiliWebsocket.IsClosedCh: 外部通过该通道判断 websocket、相关协程是否结束
	该通道由 BiliWebsocket.Net 创建, 由 BiliWebsocket.Stop 关闭
*/

const debug = true

type BiliWebsocket struct {
	WebsocketP *websocketP.WebSocketP

	// User 与 BiliWebsocket 通信
	sendCh     *chan biliDataConv.EncodeArgs
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
		case b, ok := <-*blw.sendCh:
			if !ok {
				return
			}
			bin, err := biliDataConv.Encode(b.Version, b.Operation, b.Sequence, b.Body)
			if err != nil {
				return
			}
			select {
			case blw.wspSendCh <- bin:
			case <-blw.WebsocketP.IsClosedCh:
				return
			}
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

	defer blw.timeoutTicker.Stop()

	select {
	case <-blw.timeoutTicker.C:
	case <-blw.WebsocketP.IsClosedCh:
	}
}

func New(conn *websocket.Conn, receiveChannel *chan []byte, sendChannel *chan biliDataConv.EncodeArgs, onOpen func(blw *BiliWebsocket) biliDataConv.EncodeArgs) *BiliWebsocket {
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
	blw.timeoutTicker = time.NewTicker(time.Second * 30)

	blw.WebsocketP = websocketP.New(conn, &blw.wspReceiveCh, &blw.wspSendCh)

	go func() {
		select {
		case *blw.sendCh <- onOpen(blw):
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
