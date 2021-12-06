package biliWebsocket

import (
	"github.com/gorilla/websocket"
	"log"
	"sync"
	"time"
)

const debug = true

type BiliLiveWebsocket struct {
	// 消息通道，用于通知各协程关闭
	channelClose struct {
		heartbeat chan struct{} // 心跳协程
		timeout   chan struct{} // 超时协程
		receive   chan struct{} // 接收协程
		send      chan struct{} // 发送协程
	}

	// 各种事件
	onEvent struct {
		onOpen      func(*BiliLiveWebsocket) []byte
		onHeartbeat func(*BiliLiveWebsocket) []byte
	}

	// 协程等待，用于等待相关协程的关闭
	wg sync.WaitGroup

	// WebSocket 链接
	conn *websocket.Conn

	// 发送通道，自动初始化，将需要发送的数据传递到该通道，由 发送协程 读取并发送数据
	sendChannel chan []byte
	// 接收通道，由用户初始化，到时候将websocket接收到的数据传递到该通道
	receiveChannel *chan []byte

	// 超时定时器，由 接收协程 重置，超时则触发停止函数
	timeoutTicker *time.Timer

	// 只执行一次 Stop 函数
	onceStop sync.Once
}

//
// NewBiliLiveWebsocket
//  @Description: 初始化并启动相关协程
//  @return *BiliLiveWebsocket
//
func NewBiliLiveWebsocket(conn *websocket.Conn, receiveChannel *chan []byte, onOpen, onHeartbeat func(*BiliLiveWebsocket) []byte) *BiliLiveWebsocket {
	blw := new(BiliLiveWebsocket)

	blw.channelClose = struct {
		heartbeat chan struct{}
		timeout   chan struct{}
		receive   chan struct{}
		send      chan struct{}
	}{heartbeat: make(chan struct{}, 1), timeout: make(chan struct{}, 1), receive: make(chan struct{}, 1), send: make(chan struct{}, 1)}

	blw.onEvent = struct {
		onOpen      func(*BiliLiveWebsocket) []byte
		onHeartbeat func(*BiliLiveWebsocket) []byte
	}{onOpen: onOpen, onHeartbeat: onHeartbeat}

	blw.conn = conn
	blw.sendChannel = make(chan []byte)
	blw.receiveChannel = receiveChannel
	blw.timeoutTicker = time.NewTimer(time.Second * 60)

	go blw.run()

	return blw
}

//
// run
//  @Description: 运行
//  @receiver blw
//
func (blw *BiliLiveWebsocket) run() {
	blw.wg.Add(4)

	// 发送协程
	go blw.sendDataFunc()
	// 心跳
	go blw.heartbeatFunc()
	// 超时
	go blw.timeoutFunc()
	// 接收协程
	go blw.receiveFunc()

	if blw.onEvent.onOpen != nil {
		blw.sendChannel <- blw.onEvent.onOpen(blw)
	}
}

func (blw *BiliLiveWebsocket) receiveFunc() {
	if debug {
		log.Println("接收协程运行")
		defer func() {
			log.Println("接收协程结束")
		}()
	}
	defer blw.wg.Done()
	defer blw.Stop()

	for {
		_, message, err := blw.conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		//log.Println("重置超时定时器")
		blw.timeoutTicker.Reset(time.Second * 60)
		*blw.receiveChannel <- message
	}
}

func (blw *BiliLiveWebsocket) timeoutFunc() {
	if debug {
		log.Println("超时协程运行")
		defer func() {
			log.Println("超时协程结束")
		}()
	}
	defer blw.wg.Done()
	defer blw.Stop()

	for {
		select {
		case <-blw.timeoutTicker.C:
			return
		case <-blw.channelClose.timeout:
			return
		}
	}
}

func (blw *BiliLiveWebsocket) sendDataFunc() {
	if debug {
		log.Println("发送协程运行")
		defer func() {
			log.Println("发送协程结束")
		}()
	}
	defer blw.wg.Done()
	defer blw.Stop()

	for {
		select {
		case b := <-blw.sendChannel:
			log.Println("发送数据")
			err := blw.conn.WriteMessage(websocket.BinaryMessage, b)
			if err != nil {
				log.Println(err)
				return
			}
			break
		case <-blw.channelClose.send:
			return
		}
	}
}

func (blw *BiliLiveWebsocket) heartbeatFunc() {
	if debug {
		log.Println("心跳协程运行")
		defer func() {
			log.Println("心跳协程结束")
		}()
	}
	defer blw.wg.Done()

	ticker := time.NewTicker(time.Second * 30)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if blw.onEvent.onHeartbeat != nil {
				log.Println("发送心跳包")
				blw.sendChannel <- blw.onEvent.onHeartbeat(blw)
			}
			break
		case <-blw.channelClose.heartbeat:
			return
		}
	}
}

func (blw *BiliLiveWebsocket) Stop() {
	blw.onceStop.Do(func() {
		go func() {
			//blw.channelClose.heartbeat <- struct{}{}
			//blw.channelClose.timeout <- struct{}{}
			//blw.channelClose.receive <- struct{}{}
			//blw.channelClose.send <- struct{}{}
			noWaitChannel(&blw.channelClose.heartbeat)
			noWaitChannel(&blw.channelClose.timeout)
			noWaitChannel(&blw.channelClose.receive)
			noWaitChannel(&blw.channelClose.send)
			// 有一个 Bug，就是 心跳协程 正在通过通道传输数据时，如果这时 发送协程 结束了，那就会阻塞，无法结束
			close(blw.sendChannel) // 所以这里把通道关闭，虽然可能会报错，不过问题不大

			err := blw.conn.Close()
			if err != nil {
				log.Println()
				return
			}

			log.Println("等待协程结束")
			blw.wg.Wait()

			close(*blw.receiveChannel)
		}()
	})
}

func noWaitChannel(ch *chan struct{}) {
	select {
	case *ch <- struct{}{}:
		return
	default:
		return
	}
}
