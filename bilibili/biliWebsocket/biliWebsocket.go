package biliWebsocket

import (
	"bilibili-live-communication/bilibili/biliAPI"
	"bilibili-live-communication/bilibili/biliBinConv"
	"bilibili-live-communication/bilibili/connEncapsulate"
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"net/url"
	"time"
)

func NewBiliWebsocket(roomId string) (connEncapsulate.ReadWriteCloser, error) {
	roomId_, err := biliAPI.GetRoomInfo(roomId)
	if err != nil {
		return nil, err
	}
	wsList, token, err := biliAPI.GetDanmuInfo(roomId_)
	if err != nil {
		return nil, err
	}
	if len(wsList) == 0 {
		return nil, errors.New("not found websocket url")
	}

	var BW = new(BiliWebsocket)
	err = BW.Init(wsList[0], func(self *BiliWebsocket) []byte {
		return join_room(roomId_, token)
	})
	if err != nil {
		return nil, err
	}

	err = BW.Start()
	if err != nil {
		return nil, err
	}

	return BW, nil
}

func join_room(roomId, token string) []byte {
	data := fmt.Sprintf("{\"uid\":0,\"roomid\":%s,\"protover\":3,\"platform\":\"web\",\"type\":2,\"key\":\"%s\"}", roomId, token)
	bin, _ := biliBinConv.Encode(0x0001, 0x00000007, 0x00000001, []byte(data))
	return bin
}

type BiliWebsocket struct {
	conn    *websocket.Conn
	context struct {
		ctx    context.Context
		cancel context.CancelFunc
	}

	// url
	wsUrl *url.URL

	onOpen func(self *BiliWebsocket) []byte
}

func (self *BiliWebsocket) heartbeat() {
	ticker := time.NewTicker(time.Second * 15)
	defer ticker.Stop()

	for {
		select {
		case <-self.context.ctx.Done():
			return
		case <-ticker.C:
			data, _ := biliBinConv.Encode(0x0001, 0x00000002, 0x00000001, []byte("[object Object]"))
			if self.WriteMessage(websocket.BinaryMessage, data) != nil {
				return
			}
		}
	}
}

func (self *BiliWebsocket) Init(wsUrl string, onOpen func(self *BiliWebsocket) []byte) error {
	var err error

	self.onOpen = onOpen

	if self.wsUrl, err = url.Parse(fmt.Sprintf("wss://%s/sub", wsUrl)); err != nil {
		return err
	}

	self.context.ctx, self.context.cancel = context.WithCancel(context.Background())

	return nil
}

func (self *BiliWebsocket) Start() error {
	var err error
	if self.conn, _, err = websocket.DefaultDialer.DialContext(self.context.ctx, self.wsUrl.String(), nil); err != nil {
		return err
	}

	if err = self.setReadDeadline(30); err != nil {
		return err
	}

	if err = self.WriteMessage(websocket.BinaryMessage, self.onOpen(self)); err != nil {
		return err
	}

	go self.heartbeat()

	return nil
}

func (self *BiliWebsocket) Close() error {
	return self.Stop()
}

func (self *BiliWebsocket) Stop() error {
	self.context.cancel()
	return self.conn.Close()
}

func (self *BiliWebsocket) setReadDeadline(second time.Duration) error {
	return self.conn.SetReadDeadline(time.Now().Add(time.Second * second))
}

func (self *BiliWebsocket) ReadMessage() (int, []byte, error) {
	mT, p, err := self.conn.ReadMessage()
	if err == nil {
		_ = self.setReadDeadline(30)
	}
	return mT, p, err
}

func (self *BiliWebsocket) WriteMessage(messageType int, data []byte) error {
	return self.conn.WriteMessage(messageType, data)
}
