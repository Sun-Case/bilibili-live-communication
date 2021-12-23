package websocketP

import (
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

/*
对 websocket 进行封装，以适用 Bilibili Live Websocket 通信

WebsocketP.sendCh: 从该通道获取外部传输进来的二进制数据, 该通道由外部设置, WebsocketP.Stop 不负责关闭
	传输的数据为已编码的二进制数据

WebsocketP.receiveCh: 将弹幕服务器的二进制数据传输到该通道, 该通道由外部设置, WebsocketP.Stop 不负责关闭
	传输的数据为已编码压缩的二进制数据，也可能是多个子二进制数据拼接后再编码压缩而生成的二进制数据

WebsocketP.IsClosedCh: 外部判断 websocket 是否断开、协程是否全部关闭
	在 websocket 断开后, 调用 WebsocketP.Stop 关闭相关协程, 全部协程关闭后, 再关闭 WebsocketP.IsClosedCh

*/

const debug = true

type WebSocketP struct {
	// 用于等待其他协程
	wg sync.WaitGroup

	// WebSocket 连接
	conn *websocket.Conn

	// send通道，可用户初始化（不初始化则为nil），用户将需要发送的数据传递到该通道，由 sendFunc协程 读取并发送数据
	sendCh *chan []byte
	// receive通道，可用户初始化（用户不初始化则为nil），将websocket接收到的数据传递到该通道
	receiveCh *chan []byte
	// IsClosed通道，供外界判断所有协程是否关闭
	IsClosedCh chan struct{}
	// needClose通道，通知相关协程退出
	needCloseCh chan struct{}

	// 只执行一次 stop函数
	once sync.Once
}

// Stop
// 停止, 可执行多次
func (wsp *WebSocketP) Stop() {
	wsp.once.Do(func() {
		go wsp.stop()
	})
}

// stop
// 停止函数，只能执行一次，否则会 panic
func (wsp *WebSocketP) stop() {
	if debug {
		log.Println("websocketP StopFunc: Running")
		defer func() {
			log.Println("websocketP StopFunc: Exit")
		}()
	}

	close(wsp.needCloseCh)
	if debug {
		log.Println("websocketP StopFunc: closed \"needCloseCh\"")
	}

	_ = wsp.conn.Close()
	if debug {
		log.Println("websocketP StopFunc: closed \"conn\"")
	}

	close(wsp.IsClosedCh)
	if debug {
		log.Println("websocketP StopFunc: closed \"IsClosedCh\"")
	}

	if debug {
		log.Println("websocketP StopFunc: Waiting other goroutines")
	}
	wsp.wg.Wait()
	if debug {
		log.Println("websocketP StopFunc: All other goroutines exit")
	}
}

// nonReceiveFunc
// 丢弃接收
func (wsp *WebSocketP) nonReceiveFunc() {
	if debug {
		log.Println("websocketP NonReceiveFunc: Running")
		defer func() {
			log.Println("websocketP NonReceiveFunc: Exit")
		}()
	}
	defer wsp.wg.Done()

	for {
		if _, _, err := wsp.conn.ReadMessage(); err != nil {
			return
		}
	}
}

// receiveFunc
// 接收
func (wsp *WebSocketP) receiveFunc() {
	if debug {
		log.Println("websocketP ReceiveFunc: Running")
		defer func() {
			log.Println("websocketP ReceiveFunc: Exit")
		}()
	}

	defer func() {
		wsp.wg.Done()
		wsp.Stop()
	}()
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	for {
		_, p, err := wsp.conn.ReadMessage()
		if err != nil {
			return
		}

		select {
		case *wsp.receiveCh <- p:
			break
		case <-wsp.needCloseCh:
			return
		}
	}
}

// sendFunc
// 发送
func (wsp *WebSocketP) sendFunc() {
	if debug {
		log.Println("websocketP SendFunc: Running")
		defer func() {
			log.Println("websocketP SendFunc: Exit")
		}()
	}
	defer func() {
		wsp.wg.Done()
		wsp.Stop()
	}()

	for {
		select {
		case b, ok := <-*wsp.sendCh:
			if !ok {
				return
			}
			if err := wsp.conn.WriteMessage(websocket.BinaryMessage, b); err != nil {
				return
			}
		case <-wsp.needCloseCh:
			return
		}
	}
}

// run
// 启动
func (wsp *WebSocketP) run() {
	if debug {
		log.Println("websocketP RunFunc: Running")
		defer func() {
			log.Println("websocketP RunFunc: Exit")
		}()
	}

	// 收发协程共两个
	wsp.wg.Add(2)

	go wsp.sendFunc()
	if wsp.receiveCh == nil {
		go wsp.nonReceiveFunc()
	} else {
		go wsp.receiveFunc()
	}
}

// New
// 新建
func New(conn *websocket.Conn, receiveChannel, sendChannel *chan []byte) *WebSocketP {
	if debug {
		log.Println("NewFunc: Running")
		defer func() {
			log.Println("NewFunc: Exit")
		}()
	}

	wsP := new(WebSocketP)

	// 初始化成员
	wsP.conn = conn
	wsP.sendCh = sendChannel
	wsP.receiveCh = receiveChannel
	wsP.IsClosedCh = make(chan struct{})
	wsP.needCloseCh = make(chan struct{})

	// 启动收发协程
	go wsP.run()

	return wsP
}
