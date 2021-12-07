package biliWebsocket

import (
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

const debug = true

type BiliLiveWebsocket struct {
	// 各种事件
	onEvent struct {
		onOpen func(*BiliLiveWebsocket) []byte
	}

	// 协程等待，用于等待相关协程的退出
	wg sync.WaitGroup

	// WebSocket 链接
	conn *websocket.Conn

	// send通道，可用户初始化（用户不初始化则自动初始化），将需要发送的数据传递到该通道，由 sendFunc协程 读取并发送数据
	sendCh *chan []byte
	// receive通道，可用户初始化（用户不初始化则自动初始化），将websocket接收到的数据传递到该通道
	receiveCh *chan []byte
	// IsClosed通道，供外界判断的通道
	IsClosedCh chan struct{}
	// needClose通道
	needCloseCh chan struct{}

	// 只执行一次 Stop 函数
	once sync.Once
}

//
// NewBiliLiveWebsocket
//  @Description: 初始化并启动相关协程
//  @return *BiliLiveWebsocket
//
func NewBiliLiveWebsocket(conn *websocket.Conn, receiveChannel *chan []byte, sendChannel *chan []byte, onOpen func(*BiliLiveWebsocket) []byte) *BiliLiveWebsocket {
	if debug {
		log.Println("new func: run")
		defer func() {
			log.Println("new func: exit")
		}()
	}

	blw := new(BiliLiveWebsocket)

	blw.onEvent = struct {
		onOpen func(*BiliLiveWebsocket) []byte
	}{onOpen: onOpen}

	blw.conn = conn
	blw.sendCh = sendChannel
	blw.receiveCh = receiveChannel
	blw.IsClosedCh = make(chan struct{})
	blw.needCloseCh = make(chan struct{})

	go blw.run()

	return blw
}

func (blw *BiliLiveWebsocket) run() {
	if debug {
		log.Println("run func: run")
		defer func() {
			log.Println("run func: exit")
		}()
	}

	blw.wg.Add(3)

	go blw.onOpenFunc()
	go blw.sendFunc()

	if blw.receiveCh == nil {
		go blw.nonReceiveFunc()
	} else {
		go blw.receiveFunc()
	}
}

func (blw *BiliLiveWebsocket) stop() {
	go blw.__stop__()
}

func (blw *BiliLiveWebsocket) __stop__() {
	if debug {
		log.Println("stop func: run")
		defer func() {
			log.Println("stop func: exit")
		}()
	}

	close(blw.needCloseCh)
	if debug {
		log.Println("stop func: closed 'needCloseCh'")
	}

	_ = blw.conn.Close()
	if debug {
		log.Println("stop func: closed 'conn'")
	}

	close(blw.IsClosedCh)
	if debug {
		log.Println("stop func: closed 'IsCloseCh'")
	}

	if debug {
		log.Println("stop func: waiting other goroutines")
	}
	blw.wg.Wait()
	if debug {
		log.Println("stop func: all other goroutines exit")
	}
}

func (blw *BiliLiveWebsocket) sendFunc() {
	if debug {
		log.Println("send func: run")
		defer func() {
			log.Println("send func: exit")
		}()
	}
	defer func() {
		blw.wg.Done()
		blw.once.Do(blw.stop)
	}()

	for {
		select {
		case b, ok := <-*blw.sendCh:
			if !ok {
				return
			}
			if err := blw.conn.WriteMessage(websocket.BinaryMessage, b); err != nil {
				return
			}
		case <-blw.needCloseCh:
			return
		}
	}
}

func (blw *BiliLiveWebsocket) receiveFunc() {
	if debug {
		log.Println("receive func: run")
		defer func() {
			log.Println("receive func: exit")
		}()
	}

	defer func() {
		blw.wg.Done()
		blw.once.Do(blw.stop)
	}()
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	for {
		_, p, err := blw.conn.ReadMessage()
		if err != nil {
			return
		}

		select {
		case *blw.receiveCh <- p:
			break
		case <-blw.needCloseCh:
			return
		}
	}
}

func (blw *BiliLiveWebsocket) nonReceiveFunc() {
	if debug {
		log.Println("non receive func: run")
		defer func() {
			log.Println("non receive func: exit")
		}()
	}

	defer func() {
		blw.wg.Done()
	}()

	for {
		if _, p, err := blw.conn.ReadMessage(); err != nil {
			return
		} else {
			log.Println(p)
		}
	}
}

func (blw *BiliLiveWebsocket) onOpenFunc() {
	if debug {
		log.Println("onOpen func: run")
		defer func() {
			log.Println("onOpen func: exit")
		}()
	}

	defer func() {
		blw.wg.Done()
	}()

	if blw.onEvent.onOpen == nil {
		return
	}

	select {
	case *blw.sendCh <- blw.onEvent.onOpen(blw):
	case <-blw.needCloseCh:
	}
}
